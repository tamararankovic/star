package apis

import (
	"github.com/c12s/star/domain"
	"github.com/nats-io/nats.go"
	"time"
)

type natsRegistrationAPI struct {
	conn                   *nats.Conn
	registrationReqSubject string
	reqTimeout             time.Duration
}

func NewNatsRegistrationAPI(conn *nats.Conn, registrationReqSubject string, reqTimeoutMilliseconds int64) domain.RegistrationAPI {
	return natsRegistrationAPI{
		conn:                   conn,
		registrationReqSubject: registrationReqSubject,
		reqTimeout:             time.Duration(reqTimeoutMilliseconds) * time.Millisecond,
	}
}

func (nra natsRegistrationAPI) Register(request domain.RegistrationReq) (*domain.RegistrationResp, error) {
	marshalledReq, err := MarshalReq(request)
	if err != nil {
		return nil, err
	}

	natsMsg, err := nra.conn.Request(nra.registrationReqSubject, marshalledReq, nra.reqTimeout)
	if err != nil {
		return nil, err
	}

	return UnmarshalResp(natsMsg.Data)
}
