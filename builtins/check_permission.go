package builtins

import (
	"errors"
	"fmt"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
	openfga "github.com/openfga/go-sdk"
	opaopenfga "github.com/thomasdarimont/custom-opa/custom-opa-openfga/plugins/openfga"
)

var checkPermissionBuiltinDecl = &rego.Function{
	Name: "openfga.check_permission",
	Decl: types.NewFunction(
		types.Args(types.S, types.S, types.S), // subject, permission, resource
		types.B),                              // Returns a boolean
}

// Use a custom cache key type to avoid collisions with other builtins caching data!!
type checkPermissionCacheKeyType string

// checkPermissionBuiltinImpl TODO
func checkPermissionBuiltinImpl(bctx rego.BuiltinContext, subjectTerm, permissionTerm, resourceIdTerm *ast.Term) (*ast.Term, error) {

	var resource string
	if err := ast.As(resourceIdTerm.Value, &resource); err != nil {
		return nil, err
	}

	var permission string
	if err := ast.As(permissionTerm.Value, &permission); err != nil {
		return nil, err
	}

	var subject string
	if err := ast.As(subjectTerm.Value, &subject); err != nil {
		return nil, err
	}

	// Check if it is already cached, assume they never become invalid.
	var cacheKey = checkPermissionCacheKeyType(fmt.Sprintf("%s#%s@%s", subject, permission, resource))
	cached, ok := bctx.Cache.Get(cacheKey)
	if ok {
		return ast.NewTerm(cached.(ast.Value)), nil
	}

	body := openfga.CheckRequest{
		TupleKey: &openfga.TupleKey{
			User:     openfga.PtrString(subject),
			Relation: openfga.PtrString(permission),
			Object:   openfga.PtrString(resource),
		},
	}

	client := opaopenfga.GetOpenFGAClient()
	if client == nil {
		return nil, errors.New("openfga client not configured")
	}

	data /*response*/, _, err := client.OpenFgaApi.Check(bctx.Context).Body(body).Execute()

	if err != nil {
		return nil, err
	}

	result := ast.Boolean(*data.Allowed)
	bctx.Cache.Put(cacheKey, result)

	return ast.NewTerm(result), nil
}
