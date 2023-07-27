package apis

import (
	"errors"
	"github.com/c12s/star/domain"
	"github.com/golang/protobuf/proto"
)

func MarshalReq(req domain.RegistrationReq) ([]byte, error) {
	protoReq, err := RegistrationReq{}.fromDomain(req)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(protoReq)
}

func UnmarshalResp(resp []byte) (*domain.RegistrationResp, error) {
	protoResp := RegistrationResp{}
	err := proto.Unmarshal(resp, &protoResp)
	if err != nil {
		return nil, err
	}
	return protoResp.toDomain()
}

func (x RegistrationReq) fromDomain(req domain.RegistrationReq) (*RegistrationReq, error) {
	var protoLabels []*Label
	for _, label := range req.Labels {
		protoLabel, err := Label{}.fromDomain(label)
		if err != nil {
			return nil, err
		}
		protoLabels = append(protoLabels, protoLabel)
	}
	return &RegistrationReq{
		Labels: protoLabels,
	}, nil
}

func (x RegistrationResp) toDomain() (*domain.RegistrationResp, error) {
	return &domain.RegistrationResp{
		NodeId: x.NodeId,
	}, nil
}

func (x Label) fromDomain(label domain.Label) (*Label, error) {
	value, err := Value{}.fromDomain(label.Value())
	if err != nil {
		return nil, err
	}
	return &Label{
		Key:   label.Key(),
		Value: value,
	}, nil
}

func (x Value) fromDomain(value interface{}) (*Value, error) {
	var marshalled []byte
	var valueType Value_ValueTYpe
	var err error
	switch value.(type) {
	case bool:
		marshalled, err = proto.Marshal(&BoolValue{Value: value.(bool)})
		valueType = Value_Bool
	case float64:
		marshalled, err = proto.Marshal(&Float64Value{Value: value.(float64)})
		valueType = Value_Float64
	case string:
		marshalled, err = proto.Marshal(&StringValue{Value: value.(string)})
		valueType = Value_String
	default:
		err = errors.New("unsupported data type")
	}
	return &Value{
		Marshalled: marshalled,
		Type:       valueType,
	}, err
}
