package openfga

import (
	"context"
	"github.com/open-policy-agent/opa/plugins"
	"github.com/open-policy-agent/opa/util"
	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/credentials"
	"sync"
)

const PluginName = "openfga"

type Config struct {
	ApiHost   string `json:"apiHost"`
	ApiScheme string `json:"apiScheme"`
	StoreId   string `json:"storeId"`
	ApiToken  string `json:"apiToken"`
}

type OpenFGAPlugin struct {
	manager *plugins.Manager
	mtx     sync.Mutex
	config  Config
	client  *openfga.APIClient
}

var instance *OpenFGAPlugin = nil

func GetOpenFGAClient() *openfga.APIClient {

	if instance == nil {
		return nil
	}

	instance.mtx.Lock()
	defer instance.mtx.Unlock()

	return instance.client
}

func (o *OpenFGAPlugin) Start(ctx context.Context) error {

	configuration, err := openfga.NewConfiguration(openfga.Configuration{
		ApiScheme: o.config.ApiScheme,
		ApiHost:   o.config.ApiHost,
		StoreId:   o.config.StoreId,
		Credentials: &credentials.Credentials{
			Method: credentials.CredentialsMethodApiToken,
			Config: &credentials.Config{
				ApiToken: o.config.ApiToken,
			},
		},
	})

	o.client = openfga.NewAPIClient(configuration)

	// HACK to expose plugin instance to be able to access the openfga client from the custom openfga check_permission builtin
	instance = o

	return err

}

func (o *OpenFGAPlugin) Stop(ctx context.Context) {
}

func (o *OpenFGAPlugin) Reconfigure(ctx context.Context, config any) {

	o.mtx.Lock()
	defer o.mtx.Unlock()

	if o.config.ApiHost != config.(Config).ApiHost {
		o.Stop(ctx)
		if err := o.Start(ctx); err != nil {
			o.manager.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateErr})
		}
	}
	o.config = config.(Config)
}

type Factory struct{}

func (Factory) New(m *plugins.Manager, config any) plugins.Plugin {

	m.UpdatePluginStatus(PluginName, &plugins.Status{State: plugins.StateNotReady})

	return &OpenFGAPlugin{
		manager: m,
		config:  config.(Config),
	}
}

func (Factory) Validate(_ *plugins.Manager, config []byte) (any, error) {
	parsedConfig := Config{}
	err := util.Unmarshal(config, &parsedConfig)
	return parsedConfig, err
}
