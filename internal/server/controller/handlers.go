package controller

import (
	"rpg/internal/server/api"
)

func (h *Controller) registrationHandler(req *api.BaseRequest) {
	payload := new(api.RegistrationPayload)
	err := req.DecodeJSONPayload(payload)
	if err != nil {
		h.lg.WithField("err", err).Error("failed to decode the payload")
		return
	}
	_, err = h.services.Auth.RegisterNewUser(payload.Login, payload.Password)
	h.clientConns[req.GetConnUID()].AsyncSend(nil)
}
