package startup

import (
	"github.com/c12s/star/apis"
	"github.com/c12s/star/configs"
	"github.com/c12s/star/services"
)

func StartApp(config *configs.Config) error {
	natsConn, err := NewNatsConn(config.NatsAddress())
	if err != nil {
		return err
	}
	registrationAPI := apis.NewNatsRegistrationAPI(natsConn, config.RegistrationSubject(), config.RegistrationReqTimeoutMilliseconds())
	rs := services.NewRegistrationService(registrationAPI)
	if !rs.Registered() {
		rs.Register()
	}
	return nil
}
