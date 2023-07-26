package services

import (
	"github.com/c12s/star/domain"
	"log"
)

type RegistrationService struct {
	api domain.RegistrationAPI
}

func NewRegistrationService(api domain.RegistrationAPI) *RegistrationService {
	return &RegistrationService{api: api}
}

func (rs *RegistrationService) Register() {
	resp, err := rs.api.Register(domain.RegistrationReq{})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(resp.NodeId)
	// todo: save node id to a file
}

func (rs *RegistrationService) Registered() bool {
	// todo: check if node id has been assigned
	return false
}
