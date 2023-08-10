package services

import (
	"errors"
	"fmt"
	magnetarapi "github.com/c12s/magnetar/pkg/api"
	"github.com/c12s/star/internal/domain"
	"log"
)

type RegistrationService struct {
	api        *magnetarapi.AsyncRegistrationClient
	nodeIdRepo domain.NodeIdRepo
}

func NewRegistrationService(api *magnetarapi.AsyncRegistrationClient, nodeIdRepo domain.NodeIdRepo) *RegistrationService {
	return &RegistrationService{
		api:        api,
		nodeIdRepo: nodeIdRepo,
	}
}

func (rs *RegistrationService) Register(maxRetries int8) error {
	req := rs.buildReq()
	for attemptsLeft := maxRetries; attemptsLeft > 0; attemptsLeft-- {
		var errChan chan error
		err := rs.tryRegister(req, errChan)
		if err == nil {
			err = <-errChan
			if err == nil {
				return nil
			}
		}
		log.Println(err)
	}
	return errors.New("max registration attempts exceeded")
}

func (rs *RegistrationService) tryRegister(req *magnetarapi.RegistrationReq, errChan chan<- error) error {
	err := rs.api.Register(req, func(resp *magnetarapi.RegistrationResp) {
		err := rs.nodeIdRepo.Put(domain.NodeId{Value: resp.NodeId})
		errChan <- err
	})
	return err
}

func (rs *RegistrationService) buildReq() *magnetarapi.RegistrationReq {
	builder := magnetarapi.NewRegistrationReqBuilder()
	cpuCores, err := cpuCores()
	if err == nil {
		builder = builder.AddFloat64Label("cpuCores", cpuCores)
	}
	for coreId := 0; coreId < int(cpuCores); coreId++ {
		mhz, err := coreMhz(string(rune(coreId)))
		if err == nil {
			builder = builder.AddFloat64Label(fmt.Sprintf("core%dmhz", coreId), mhz)
		}
		vendorId, err := coreVendorId(string(rune(coreId)))
		if err == nil {
			builder = builder.AddStringLabel(fmt.Sprintf("core%dvendorId", coreId), vendorId)
		}
		model, err := coreModel(string(rune(coreId)))
		if err == nil {
			builder = builder.AddStringLabel(fmt.Sprintf("core%dmodel", coreId), model)
		}
		cacheKB, err := coreCacheKB(string(rune(coreId)))
		if err == nil {
			builder = builder.AddFloat64Label(fmt.Sprintf("core%dcacheKB", coreId), cacheKB)
		}
	}
	fsType, err := FsType()
	if err == nil {
		builder = builder.AddStringLabel("fsType", fsType)
	}
	diskTotalGB, err := diskTotalGB()
	if err == nil {
		builder = builder.AddFloat64Label("diskTotalGB", diskTotalGB)
	}
	diskFreeGB, err := diskFreeGB()
	if err == nil {
		builder = builder.AddFloat64Label("diskFreeGB", diskFreeGB)
	}
	kernelArch, err := kernelArch()
	if err == nil {
		builder = builder.AddStringLabel("kernelArch", kernelArch)
	}
	kernelVersion, err := kernelVersion()
	if err == nil {
		builder = builder.AddStringLabel("kernelVersion", kernelVersion)
	}
	platform, err := platform()
	if err == nil {
		builder = builder.AddStringLabel("platform", platform)
	}
	platformFamily, err := platformFamily()
	if err == nil {
		builder = builder.AddStringLabel("platformFamily", platformFamily)
	}
	platformVersion, err := platformVersion()
	if err == nil {
		builder = builder.AddStringLabel("platformVersion", platformVersion)
	}
	memoryTotalGB, err := memoryTotalGB()
	if err == nil {
		builder = builder.AddFloat64Label("memoryTotalGB", memoryTotalGB)
	}
	return builder.Request()
}

func (rs *RegistrationService) Registered() bool {
	if _, err := rs.nodeIdRepo.Get(); err != nil {
		return false
	}
	return true
}
