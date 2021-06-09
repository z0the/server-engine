package handler

import (
	"rpg/internal/service"
	"rpg/pkg/hubber"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	sync.Mutex
	logger       *logrus.Logger
	services     *service.Service
	handleChan   chan hubber.IRequest
	clients      map[int64]hubber.IClient
	clientsRoom  map[int64]int64
	messages     <-chan hubber.IResponse
	lastClientID int64
}

func NewHandler(logger *logrus.Logger, msgPipe <-chan hubber.IResponse, services *service.Service) *Handler {
	hdl := &Handler{
		logger:      logger,
		services:    services,
		handleChan:  make(chan hubber.IRequest),
		clients:     make(map[int64]hubber.IClient),
		clientsRoom: make(map[int64]int64),
		messages:    msgPipe,
	}
	go hdl.listenMessages()
	return hdl
}

func (h *Handler) Register(client hubber.IClient) int64 {
	h.Lock()
	defer h.Unlock()
	h.lastClientID++
	h.clients[h.lastClientID] = client
	h.logger.Infof("Register client #%d", h.lastClientID)
	return h.lastClientID
}

func (h *Handler) Unregister(id int64) {
	h.Lock()
	defer h.Unlock()
	if _, ok := h.clients[id]; ok {
		delete(h.clients, id)
		h.logger.Infof("Unregister client #%d", id)
	}
}

func (h *Handler) Handle(req hubber.IRequest) {
	go h.handleRequest(req)
}

func (h *Handler) handleRequest(req hubber.IRequest) {
	if strings.Contains(req.GetAction(), "game") {
		if err := h.services.SendToGame(h.clientsRoom[req.SenderID()], req); err != nil {
			h.handleError(req.SenderID(), err.Error())
		}
		return
	}
	switch req.GetAction() {
	case "GAME_EVENT":
	case "joinRoom":
		h.joinRoom(req)
	}
}

func (h *Handler) listenMessages() {
	for msg := range h.messages {
		go h.handleMessage(msg)
	}
}

func (h *Handler) handleMessage(msg hubber.IResponse) {
	if client, ok := h.clients[msg.ReceiverID()]; ok {
		client.Send(msg)
		return
	}
	h.logger.Warn("Can't send message to client")
}

func (h *Handler) proxyToRoom(req hubber.IRequest) {
	if roomID, ok := h.clientsRoom[req.SenderID()]; ok {
		if err := h.services.SendToGame(roomID, req); err != nil {
			h.handleError(req.SenderID(), err.Error())
		}
	}
}

func (h *Handler) handleError(receiverID int64, description string) {
	h.logger.Error(description)
	type Error struct {
		description string
	}
	err := &Error{description: description}
	var resp hubber.Response
	resp.SetReceiverID(receiverID)
	resp.Action = "Error"
	resp.WriteData(err)
	h.handleMessage(&resp)
}
