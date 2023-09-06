package twiliofx

import (
	"github.com/twilio/twilio-go"
	"go.uber.org/fx"
)

var Module = fx.Module("twiliofx", fx.Provide(New))

type Params struct {
	fx.In

	Config *Config `optional:"true"`
}

type Result struct {
	fx.Out

	RestClient *twilio.RestClient
}

type Config struct {
	Username   string
	Password   string
	AccountSid string
}

func New(p Params) (Result, error) {
	var cli *twilio.RestClient

	if p.Config != nil {
		cli = twilio.NewRestClientWithParams(twilio.ClientParams{
			Username:   p.Config.Username,
			Password:   p.Config.Password,
			AccountSid: p.Config.AccountSid,
		})
	} else {
		cli = twilio.NewRestClient()
	}

	return Result{
		RestClient: cli,
	}, nil
}
