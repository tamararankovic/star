package domain

import "github.com/c12s/magnetar/pkg/magnetar"

type RegistrationAPI interface {
	Register(request magnetar.RegistrationReq) (*magnetar.RegistrationResp, error)
}
