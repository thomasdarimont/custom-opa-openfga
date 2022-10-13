package plugins

import (
	"github.com/open-policy-agent/opa/runtime"
	"github.com/thomasdarimont/custom-opa/custom-opa-openfga/plugins/openfga"
)

func Register() {
	runtime.RegisterPlugin(openfga.PluginName, openfga.Factory{})
}
