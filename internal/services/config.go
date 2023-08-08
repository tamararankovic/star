package services

import (
	"context"
	oort "github.com/c12s/oort/pkg/proto"
	"github.com/c12s/star/internal/domain"
)

type ConfigService struct {
	repo      domain.ConfigRepo
	evaluator oort.OortEvaluatorClient
}

func NewConfigService(repo domain.ConfigRepo, evaluator oort.OortEvaluatorClient) (*ConfigService, error) {
	return &ConfigService{
		repo:      repo,
		evaluator: evaluator,
	}, nil
}

func (c *ConfigService) Put(req domain.PutConfigGroupReq) (*domain.PutConfigGroupResp, error) {
	//resp, err := c.evaluator.Authorize(context.TODO(), &oort.AuthorizationReq{
	//	Subject: &oort.Resource{
	//		Id:   req.SubId,
	//		Kind: req.SubKind,
	//	},
	//	Object: &oort.Resource{
	//		Id:   req.Group.Namespace,
	//		Kind: "namespace",
	//	},
	//	PermissionName: "config.put",
	//})
	//if err != nil {
	//	return nil, err
	//}
	//if !resp.Allowed {
	//	return nil, domain.ErrUnauthorized()
	//}

	err := c.repo.Put(req.Group)
	if err != nil {
		return nil, err
	}
	return &domain.PutConfigGroupResp{}, nil
}

func (c *ConfigService) Get(req domain.GetConfigGroupReq) (*domain.GetConfigGroupResp, error) {
	cg := &domain.ConfigGroup{
		Name:      req.GroupName,
		Namespace: req.Namespace,
	}

	resp, err := c.evaluator.Authorize(context.TODO(), &oort.AuthorizationReq{
		Subject: &oort.Resource{
			Id:   req.SubId,
			Kind: req.SubKind,
		},
		Object: &oort.Resource{
			Id:   cg.Id(),
			Kind: "config",
		},
		PermissionName: "config.get",
	})
	if err != nil {
		return nil, err
	}
	if !resp.Allowed {
		return nil, domain.ErrUnauthorized()
	}

	cg, err = c.repo.Get(cg.Id())
	if err != nil {
		return nil, err
	}
	return &domain.GetConfigGroupResp{
		Group: *cg,
	}, nil
}
