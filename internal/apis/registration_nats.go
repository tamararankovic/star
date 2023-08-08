package apis

import (
	"github.com/c12s/magnetar/pkg/magnetar"
	"github.com/c12s/star/internal/domain"
	"github.com/nats-io/nats.go"
	"time"
)

type natsRegistrationAPI struct {
	conn                   *nats.Conn
	registrationReqSubject string
	reqTimeout             time.Duration
	marshaller             magnetar.Marshaller
}

func NewNatsRegistrationAPI(conn *nats.Conn, registrationReqSubject string, reqTimeoutMilliseconds int64, marshaller magnetar.Marshaller) domain.RegistrationAPI {
	return natsRegistrationAPI{
		conn:                   conn,
		registrationReqSubject: registrationReqSubject,
		reqTimeout:             time.Duration(reqTimeoutMilliseconds) * time.Millisecond,
		marshaller:             marshaller,
	}
}

func (n natsRegistrationAPI) Register(request magnetar.RegistrationReq) (*magnetar.RegistrationResp, error) {
	marshalledReq, err := n.marshaller.MarshalRegistrationReq(request)
	if err != nil {
		return nil, err
	}

	natsMsg, err := n.conn.Request(n.registrationReqSubject, marshalledReq, n.reqTimeout)
	if err != nil {
		return nil, err
	}

	return n.marshaller.UnmarshalRegistrationResp(natsMsg.Data)
}
