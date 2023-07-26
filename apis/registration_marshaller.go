package apis

import (
	"github.com/c12s/star/domain"
	"github.com/golang/protobuf/proto"
)

func MarshalReq(req domain.RegistrationReq) ([]byte, error) {
	protoReq := RegistrationReq{}.fromDomain(req)
	return proto.Marshal(protoReq)
}

func UnmarshalResp(resp []byte) (*domain.RegistrationResp, error) {
	protoResp := RegistrationResp{}
	err := proto.Unmarshal(resp, &protoResp)
	if err != nil {
		return nil, err
	}
	return protoResp.toDomain(), nil
}

func (rr RegistrationReq) fromDomain(req domain.RegistrationReq) *RegistrationReq {
	return &RegistrationReq{}
}

func (rr RegistrationResp) toDomain() *domain.RegistrationResp {
	return &domain.RegistrationResp{
		NodeId: rr.NodeId,
	}
}
