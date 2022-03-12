package controller

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"

	"rpg/internal/server/api"
	"rpg/internal/server/service"
	"rpg/pkg/hubber"

	"github.com/sirupsen/logrus"
)

type Controller struct {
	sync.Mutex
	lg                 *logrus.Logger
	services           *service.Services
	handleRequestsChan chan hubber.Message
	handlers           map[string]func(baseReq *api.BaseRequest)
	clientConns        map[string]hubber.ClientConnection
	clientsRoom        map[string]string
}

func NewController(logger *logrus.Logger, services *service.Services) *Controller {
	hdl := &Controller{
		lg:                 logger,
		services:           services,
		handleRequestsChan: make(chan hubber.Message),
		clientConns:        make(map[string]hubber.ClientConnection),
		clientsRoom:        make(map[string]string),
	}

	hdl.handlers = map[string]func(baseReq *api.BaseRequest){
		api.RegistrationReqType.String(): hdl.registrationHandler,
	}

	go hdl.listenClientsRequests()

	return hdl
}

func (h *Controller) HandleClientConnection(client hubber.ClientConnection) {
	h.Lock()
	defer h.Unlock()

	clientUID := uuid.NewString()
	h.clientConns[clientUID] = client
	h.lg.WithField("uid", clientUID).Info("handle new client")
	client.Run(clientUID, h.handleRequestsChan)
}

func (h *Controller) clearClientConn(clientUID string) {
	h.Lock()
	defer h.Unlock()

	if _, ok := h.clientConns[clientUID]; ok {
		delete(h.clientConns, clientUID)
		h.lg.WithField("uid", clientUID).Info("clear client")
	}
}

func (h *Controller) handleRequestMsg(msg hubber.Message) {
	req := new(api.BaseRequest)
	err := json.Unmarshal(msg.GetRawData(), req)
	req.SetConnUID(msg.GetConnUID())
	if err != nil {
		h.lg.WithField("err", err).Error("failed to unmarshal request")
		return
	}

	err = req.Validate()
	if err != nil {
		h.lg.WithField("err", err).Error("validation failed")
		return
	}

	h.handlers[req.RequestTypeCode()](req)
}

func (h *Controller) listenClientsRequests() {
	for msg := range h.handleRequestsChan {
		go h.handleRequestMsg(msg)
	}
}

// func (h *Controller) sendMessageToClient(msg hubber.IResponse) {
// 	if client, ok := h.clientConns[msg.ReceiverID()]; ok {
// 		client.Send(msg)
// 		return
// 	}
// 	h.lg.Warn("Can't send message to client")
// }

// func (h *Controller) proxyToRoom(req hubber.IRequest) {
// 	if roomID, ok := h.clientsRoom[req.SenderID()]; ok {
// 		if err := h.services.SendToGame(roomID, req); err != nil {
// 			h.handleError(req.SenderID(), err.Error())
// 		}
// 	}
// }

// func (h *Controller) handleError(receiverID int64, description string) {
// 	h.lg.Error(description)
// 	type Error struct {
// 		description string
// 	}
// 	err := &Error{description: description}
// 	var resp hubber.Response
// 	resp.SetReceiverID(receiverID)
// 	resp.Action = "Error"
// 	resp.WriteData(err)
// 	h.sendMessageToClient(&resp)
// }
