package async_controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"go.uber.org/zap"

	"rpg/internal/server/async_controller/asyncapi"
	"rpg/internal/server/auth"
	"rpg/internal/server/eventbus"
	"rpg/internal/server/matchmaker"
	"rpg/pkg/hubber"
)

func NewController(
	logger *zap.SugaredLogger,
	bus eventbus.Bus,
	auth auth.Service,
	matchMaker matchmaker.Service,
) *Controller {
	hdl := &Controller{
		lg:                 logger,
		bus:                bus,
		auth:               auth,
		matchMaker:         matchMaker,
		clientMessagesChan: make(chan hubber.Message, 100),
		conns:              make(map[string]*ClientConnWrapper),
	}

	hdl.handlers = map[string]Handler{
		asyncapi.LoginMsgType: hdl.loginHandler,
	}

	go hdl.listenClientsRequests()

	return hdl
}

type Controller struct {
	sync.RWMutex
	lg                 *zap.SugaredLogger
	bus                eventbus.Bus
	auth               auth.Service
	matchMaker         matchmaker.Service
	clientMessagesChan chan hubber.Message
	handlers           map[string]Handler
	conns              map[string]*ClientConnWrapper
}

func (h *Controller) HandleClientConnection(clientConn hubber.ConnectionWrapper) {
	h.Lock()
	defer h.Unlock()

	connUID := uuid.NewString()
	h.conns[connUID] = &ClientConnWrapper{
		ConnectionWrapper: clientConn,
	}
	h.lg.Infow("handle new clientConn", "uid", connUID)
	clientConn.StartReading(connUID, h.clientMessagesChan)

	h.bus.SubscribeWithToPrefix(connUID, h.busRepeater)

	h.sendMsgToClient(
		connUID,
		asyncapi.ServerConnectMsgType,
		asyncapi.ServerConnectPayloadOUT{ConnectionUID: connUID},
	)
}

func (h *Controller) busRepeater(msg eventbus.BusMsg) {
	fmt.Println("Repeater")

	_, isConnExists := h.conns[msg.RecipientConnUID]
	if !isConnExists {
		h.lg.Warnw("trying to send to not exists conn", "connUID", msg.RecipientConnUID)
		return
	}
	h.sendMsgToClient(msg.RecipientConnUID, asyncapi.MsgType(msg.MsgType), msg.Payload)
}

func (h *Controller) clearClientConn(clientUID string) {
	h.Lock()
	defer h.Unlock()

	if _, ok := h.conns[clientUID]; ok {
		delete(h.conns, clientUID)
		h.lg.Infow("clear client", "uid", clientUID)

	}
}

func (h *Controller) handleMsgFromClient(rawMsg hubber.Message) {
	msg := asyncapi.BaseMessage{}
	err := json.Unmarshal(rawMsg.GetRawData(), &msg)
	if err != nil {
		h.lg.Errorw("failed to unmarshal request", "err", err)
		return
	}

	err = msg.Validate()
	if err != nil {
		h.lg.Errorw("validation failed", "err", err)
		return
	}

	payload, err := decodePayload(msg)
	if err != nil {
		h.lg.Errorw("validation failed", "err", err)
		return
	}

	msg.SetConnUID(rawMsg.GetConnUID())

	// result, err := h.handlers[msg.RequestTypeCode()](msg)
	// if err != nil {
	// 	h.sendMsgToClient(
	// 		rawMsg.GetConnUID(),
	// 		asyncapi.ErrorMsgType,
	// 		asyncapi.ErrorPayloadOut{
	// 			Description: err.Error(),
	// 		},
	// 	)
	// 	return
	// }
	// if result == nil {
	// 	return
	// }

	h.bus.PublishWithFromPrefix(
		msg.GetConnUID(), eventbus.BusMsg{
			RecipientConnUID: msg.GetConnUID(),
			MsgType:          msg.MsgType.String(),
			Payload:          payload,
		},
	)
	// h.sendMsgToClient(
	// 	rawMsg.GetConnUID(),
	// 	msg.MsgType,
	// 	result,
	// )
}

func decodePayload(reqMsg asyncapi.BaseMessage) (any, error) {
	switch reqMsg.MsgType {
	case asyncapi.MoveMsgType:
		var payload asyncapi.MovePayloadIN
		err := reqMsg.DecodeJSONPayload(&payload)
		if err != nil {
			return nil, err
		}
		return payload, nil
	default:
		return nil, errors.New("unknown msg type")
	}
}

func (h *Controller) listenClientsRequests() {
	for msg := range h.clientMessagesChan {
		go h.handleMsgFromClient(msg)
	}
}

func (h *Controller) sendMsgToClient(recipientConnID string, msgType asyncapi.MsgType, payload any) {
	recipient, exists := h.conns[recipientConnID]
	if !exists {
		h.lg.Error("Recipient client connection does not exist")
		return
	}

	err := msgType.Validate()
	if err != nil {
		h.lg.Errorw("failed to validate request type from server message", "err", err)
		return
	}

	msg, err := asyncapi.NewBaseMessage(msgType, payload)
	if err != nil {
		h.lg.Errorw("failed to create baseMessage from server message", "err", err)
		return
	}

	rawMsg, err := json.Marshal(msg)
	if err != nil {
		h.lg.Errorw("failed to marshal msg", "err", err)
		return
	}

	recipient.Send(rawMsg)
}
