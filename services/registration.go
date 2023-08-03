package services

import (
	"fmt"
	"github.com/c12s/magnetar/pkg/magnetar"
	"github.com/c12s/star/domain"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
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
	labels := Labels()
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

func Labels() []magnetar.Label {
	labels := make([]magnetar.Label, 0)
	cpuCores, err := cpu.Counts(true)
	if err == nil {
		labels = append(labels, magnetar.NewFloat64Label("cpuCores", float64(cpuCores)))
	}
	info, err := cpu.Info()
	if err == nil {
		for _, core := range info {
			labelKey := fmt.Sprintf("core%smhz", core.CoreID)
			labels = append(labels, magnetar.NewFloat64Label(labelKey, core.Mhz))
			labelKey = fmt.Sprintf("core%svendorId", core.CoreID)
			labels = append(labels, magnetar.NewStringLabel(labelKey, core.VendorID))
			labelKey = fmt.Sprintf("core%smodel", core.CoreID)
			labels = append(labels, magnetar.NewStringLabel(labelKey, core.ModelName))
			labelKey = fmt.Sprintf("core%scacheKB", core.CoreID)
			labels = append(labels, magnetar.NewFloat64Label(labelKey, float64(core.CacheSize/1000)))
		}
	}
	usage, err := disk.Usage("/")
	if err == nil {
		labels = append(labels, magnetar.NewStringLabel("fsType", usage.Fstype))
		labels = append(labels, magnetar.NewFloat64Label("diskTotalGB", float64(usage.Total/1000000000)))
		labels = append(labels, magnetar.NewFloat64Label("diskFreeGB", float64(usage.Free/1000000000)))
	}
	kernelArch, err := host.KernelArch()
	if err == nil {
		labels = append(labels, magnetar.NewStringLabel("kernelArch", kernelArch))
	}
	kernelVersion, err := host.KernelVersion()
	if err == nil {
		labels = append(labels, magnetar.NewStringLabel("kernelVersion", kernelVersion))
	}
	platform, platformFamily, platformVersion, err := host.PlatformInformation()
	if err == nil {
		labels = append(labels, magnetar.NewStringLabel("platform", platform))
		labels = append(labels, magnetar.NewStringLabel("platformFamily", platformFamily))
		labels = append(labels, magnetar.NewStringLabel("platformVersion", platformVersion))
	}
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		labels = append(labels, magnetar.NewFloat64Label("memoryTotalGB", float64(memInfo.Total/1000000000)))
	}
	return labels
}
