package domain

type RegistrationAPI interface {
	Register(request RegistrationReq) (*RegistrationResp, error)
}

type RegistrationReq struct {
	Labels []Label
}

type RegistrationResp struct {
	NodeId string
}
