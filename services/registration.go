package services

import (
	"github.com/c12s/magnetar/pkg/magnetar"
	"github.com/c12s/star/domain"
	"log"
)

type RegistrationService struct {
	api        domain.RegistrationAPI
	nodeIdRepo domain.NodeIdRepo
	maxRetries int8
}

func NewRegistrationService(api domain.RegistrationAPI, nodeIdRepo domain.NodeIdRepo, maxRetries int8) *RegistrationService {
	return &RegistrationService{
		api:        api,
		nodeIdRepo: nodeIdRepo,
		maxRetries: maxRetries,
	}
}

func (rs *RegistrationService) Register() error {
	var err error
	for attemptsLeft := rs.maxRetries; attemptsLeft > 0; attemptsLeft-- {
		err = rs.tryRegister()
		if err == nil {
			break
		}
		log.Println(err)
	}
	return err
}

func (rs *RegistrationService) tryRegister() error {
	labels := labels()
	resp, err := rs.api.Register(magnetar.RegistrationReq{
		Labels: labels,
	})
	if err != nil {
		return err
	}
	log.Println(resp.NodeId)
	return rs.nodeIdRepo.Put(domain.NodeId{Value: resp.NodeId})
}

func (rs *RegistrationService) Registered() bool {
	if _, err := rs.nodeIdRepo.Get(); err != nil {
		return false
	}
	return true
}

func labels() []magnetar.Label {
	// todo: replace dummy with real labels
	label := magnetar.NewBoolLabel("bkey", true)
	label2 := magnetar.NewStringLabel("skey", "abc")
	return []magnetar.Label{label, label2}
}
