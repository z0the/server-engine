package handler

import (
	"rpg/internal/service"
	"rpg/pkg/hubber"
)

func (h *Handler) joinRoom(req hubber.IRequest) {
	var data service.RoomData
	req.ParseData(&data)
	roomID, err := h.services.JoinRoom(req.SenderID(), data)
	if err != nil {
		h.handleError(req.SenderID(), err.Error())
		return
	}
	h.logger.Infof("User #%d has joined  room #%d", req.SenderID(), roomID)
	h.clientsRoom[req.SenderID()] = roomID
	resp := &hubber.Response{}
	resp.SetReceiverID(req.SenderID())
	resp.Action = "joined"
	resp.WriteData(&struct {
		ID int64
	}{
		ID: req.SenderID(),
	})
	h.handleMessage(resp)
}
