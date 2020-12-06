package apidefinition_test

import (
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/foryouandyourcustomers/azapim/internal/apidefinition"
)

func TestGetOpenAPISpecFromFile(t *testing.T) {
	a := apidefinition.Definition{}
	a.OpenAPISpecPath = "file://./testdata/swagger.json"

	// load the openapi spec file
	err := a.GetOpenAPISpec()
	if err != nil {
		t.Error("Unable to load testdata")
	}

	// if the contains wasnt loaded we wont find the sample api string
	if strings.Contains(a.OpenAPISpec, "Sample API") != true {
		t.Error("Unable to parse api spec")
	}

	// check if the api format is correct
	if a.OpenAPIFormat != apimanagement.Openapijson {
		t.Error("Invalid openapi format")
	}
}

func TestGetOpenAPISpecFromUrl(t *testing.T) {
	a := apidefinition.Definition{}
	a.OpenAPISpecPath = "https://myurl/swagger.json"

	// load the openapi spec file
	err := a.GetOpenAPISpec()
	if err != nil {
		t.Error("Unable to load testdata")
	}

	// check if spec is set to url
	if a.OpenAPISpecPath != a.OpenAPISpec {
		t.Error("Invalid openapi spec value")
	}

	// check if the api format is correct
	if a.OpenAPIFormat != apimanagement.OpenapijsonLink {
		t.Error("Invalid openapi format")
	}
}

func TestGetXMLPolicyFromFile(t *testing.T) {
	a := apidefinition.Definition{}
	a.XMLPolicyPath = "file://./testdata/policy.xml"

	// load the openapi spec file
	err := a.GetXMLPolicy()
	if err != nil {
		t.Error("Unable to load xml policy")
	}

	// if the contains wasnt loaded we wont find the sample api string
	if strings.Contains(a.XMLPolicy, "load from file") != true {
		t.Error("Unable to parse xml policy")
	}

	// check if the api format is correct
	if a.XMLPolicyFormat != apimanagement.XML {
		t.Error("Invalid policy format")
	}
}

func TestGetXMLPolicyFromUrl(t *testing.T) {
	a := apidefinition.Definition{}
	a.XMLPolicyPath = "https://myurl/policy.xml"

	// load the openapi spec file
	err := a.GetXMLPolicy()
	if err != nil {
		t.Error("Unable to load testdata")
	}

	// check if spec is set to url
	if a.XMLPolicyPath != a.XMLPolicy {
		t.Error("Invalid xml policy value")
	}

	// check if the api format is correct
	if a.XMLPolicyFormat != apimanagement.XMLLink {
		t.Error("Invalid policy format")
	}
}
