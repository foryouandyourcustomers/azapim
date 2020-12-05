package apidefinition

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/apimanagement/mgmt/apimanagement"
	"github.com/prometheus/common/log"
)

// Definition allows to set all required values for regsitering and updating an API
// in the api management service
type Definition struct {
	OpenAPISpecPath string
	OpenAPISpec     string
	OpenAPIFormat   apimanagement.ContentFormat

	XMLPolicyPath   string
	XMLPolicy       string
	XMLPolicyFormat apimanagement.PolicyContentFormat

	APIID               string
	APIUniqueID         string
	APIDisplayName      string
	APIVersion          string
	APIVersioningScheme apimanagement.VersioningScheme
	APIPath             string
	APIRevision         string
	APIServiceURL       string
	APIProductsRaw      string
	APIProducts         []string

	APIProtocols         []apimanagement.Protocol
	SubscriptionRequired bool
}

// SetDefaults depending on the given values
func (api *Definition) SetDefaults() {
	api.APIProtocols = append(api.APIProtocols, apimanagement.ProtocolHTTPS)
	api.APIVersioningScheme = apimanagement.VersioningSchemeSegment
	api.SubscriptionRequired = true
	api.APIRevision = "1"
	api.APIUniqueID = fmt.Sprintf("%s-%s", api.APIID, api.APIVersion)
	if len(api.APIProductsRaw) > 0 {
		api.APIProducts = strings.Split(strings.TrimSpace(api.APIProductsRaw), ",")
	}
}

// GetOpenAPISpec retrieves the openapi spec file either from file or url. if unable to load spec throws an exception
func (api *Definition) GetOpenAPISpec() {

	if strings.HasPrefix(api.OpenAPISpecPath, "https://") || strings.HasPrefix(api.OpenAPISpecPath, "http://") {
		log.Infof("OpenApi Spec will be downloaded by APIM during create/update from '%s'", api.OpenAPISpecPath)
		api.OpenAPIFormat = apimanagement.OpenapijsonLink
		api.OpenAPISpec = api.OpenAPISpecPath
	} else {
		log.Infof("Load openapi spec from file: %s", api.OpenAPISpecPath)

		f := api.OpenAPISpecPath
		if strings.HasPrefix(f, "file://") {
			f = f[7:]
		}
		file, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		api.OpenAPISpec = string(file)
		api.OpenAPIFormat = apimanagement.Openapijson
	}
}

// GetXMLPolicy retrives the xml policy either from file or from url. if not specified loads default, empty xml policy
func (api *Definition) GetXMLPolicy() {
	if api.XMLPolicyPath == "" {
		log.Info("No xml policy given, load default policy")
		api.XMLPolicyFormat = apimanagement.XML
		api.XMLPolicyPath = "none (default policy)"
		api.XMLPolicy = `<policies>
<inbound>
<base />
</inbound>
<backend>
<base />
</backend>
<outbound>
<base />
</outbound>
<on-error>
<base />
</on-error>
</policies>`
	} else if strings.HasPrefix(api.XMLPolicyPath, "https://") || strings.HasPrefix(api.XMLPolicyPath, "http://") {
		log.Infof("Xml Policy will be downloaded by APIM during create/update from '%s'", api.XMLPolicyPath)
		api.XMLPolicyFormat = apimanagement.XMLLink
		api.XMLPolicy = api.XMLPolicyPath
	} else {
		log.Infof("Load XML policy from file: %s", api.XMLPolicyPath)

		f := api.XMLPolicyPath
		if strings.HasPrefix(f, "file://") {
			f = f[7:]
		}
		file, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		api.XMLPolicy = string(file)
		api.XMLPolicyFormat = apimanagement.XML
	}
}
