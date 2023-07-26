package services

import (
	"fmt"
	"github.com/c12s/star/domain"
	"log"
)

type registrationService struct {
	api domain.RegistrationAPI
}

func NewRegistrationService(api domain.RegistrationAPI) *registrationService {
	return &registrationService{api: api}
}

func (rs *registrationService) Register() {
	resp, err := rs.api.Register(domain.RegistrationReq{})
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(resp.NodeId)
	// todo: save node id to a file
}

func (rs *registrationService) Registered() bool {
	// todo: check if node id has been assigned
	return false
}
