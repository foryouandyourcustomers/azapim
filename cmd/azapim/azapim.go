package main

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"

	"github.com/Azure/go-autorest/autorest/azure/auth"
)

var (
	apimClient azureApimClient
	apiDef     apiDefinition
)

type azureApimClient struct {
	ctx              context.Context
	apiClient        apimanagement.APIClient
	versionSetClient apimanagement.APIVersionSetClient
	apiPolicyClient  apimanagement.APIPolicyClient
	subscription     string
	resourceGroup    string
	serviceName      string
}

func (apim *azureApimClient) authenticate() {
	apim.apiClient = apimanagement.NewAPIClient(apim.subscription)
	apim.versionSetClient = apimanagement.NewAPIVersionSetClient(apim.subscription)
	apim.apiPolicyClient = apimanagement.NewAPIPolicyClient(apim.subscription)

	a, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		// looking at the newauthorizerfromenvrionment funciton it
		// seems that thing never returns an error whatsoever!
		log.Debug("Unable to create authorizer from az cli. Lets load the authorizer from the environment and hope for the best!")
		a, _ = auth.NewAuthorizerFromEnvironment()

	}
	apim.apiClient.Authorizer = a
	apim.versionSetClient.Authorizer = a
	apim.apiPolicyClient.Authorizer = a
}

func (apim *azureApimClient) createOrUpdate(apidef *apiDefinition) {

	apiVersionSet := apimanagement.APIVersionSetContract{
		APIVersionSetContractProperties: &apimanagement.APIVersionSetContractProperties{
			DisplayName:      &apidef.apiDisplayName,
			VersioningScheme: apidef.apiVersioningScheme,
		},
	}

	log.Infof("Creating/Updating API versionset '%s'", apidef.apiID)
	apiVersionSetContract, err := apim.versionSetClient.CreateOrUpdate(
		apim.ctx,
		apim.resourceGroup,
		apim.serviceName,
		apidef.apiID,
		apiVersionSet,
		uuid.New().String())
	if err != nil {
		log.Fatalf("cannot create api versionset: %v", err)
	}
	log.Infof("Created/Updated API version set '%s'\n", *apiVersionSetContract.ID)

	apiProperties := apimanagement.APICreateOrUpdateParameter{
		APICreateOrUpdateProperties: &apimanagement.APICreateOrUpdateProperties{
			Format:               apidef.openAPIFormat,
			DisplayName:          &apidef.apiDisplayName,
			Value:                &apidef.openAPISpec,
			Protocols:            &apidef.apiProtocols,
			Path:                 &apidef.apiPath,
			SubscriptionRequired: &apidef.subscriptionRequired,
			APIVersion:           &apidef.apiVersion,
			APIVersionSetID:      apiVersionSetContract.ID,
			APIRevision:          &apidef.apiRevision,
		},
	}

	log.Infof("Creating/Updating API '%s'\n", apidef.apiDisplayName)
	future, err := apim.apiClient.CreateOrUpdate(
		apim.ctx,
		apim.resourceGroup,
		apim.serviceName,
		apidef.apiID,
		apiProperties,
		uuid.New().String())
	if err != nil {
		log.Fatalf("cannot create api endpoint: %v\n", err)
	}
	err = future.WaitForCompletionRef(apim.ctx, apim.apiClient.Client)
	if err != nil {
		log.Fatalf("cannot get the api endpoint future response: %v\n", err)
	}

	log.Infof("Update API policy with policy definition\n")
	apiPolicy := apimanagement.PolicyContract{
		PolicyContractProperties: &apimanagement.PolicyContractProperties{
			Format: apidef.xmlPolicyFormat,
			Value:  &apidef.xmlPolicy,
		},
	}
	apiPolicyContract, err := apim.apiPolicyClient.CreateOrUpdate(
		apim.ctx,
		apim.resourceGroup,
		apim.serviceName,
		apidef.apiID,
		apiPolicy,
		uuid.New().String(),
	)
	if err != nil {
		log.Fatalf("cannot create api policy: %v\n", err)
	}
	log.Infof("Created/Updated API Polic set '%s'\n", *apiPolicyContract.ID)
}

type apiDefinition struct {
	openAPISpecPath string
	openAPISpec     string
	openAPIFormat   apimanagement.ContentFormat

	xmlPolicyPath   string
	xmlPolicy       string
	xmlPolicyFormat apimanagement.PolicyContentFormat

	apiID               string
	apiDisplayName      string
	apiVersion          string
	apiVersioningScheme apimanagement.VersioningScheme
	apiPath             string
	apiRevision         string

	apiProtocols         []apimanagement.Protocol
	subscriptionRequired bool
}

func (api *apiDefinition) getOpenAPISpec() {

	if strings.HasPrefix(api.openAPISpecPath, "https://") || strings.HasPrefix(api.openAPISpecPath, "http://") {
		log.Infof("OpenApi Spec will be downloaded by APIM during create/update from '%s'", api.openAPISpecPath)
		api.openAPIFormat = apimanagement.OpenapijsonLink
		api.openAPISpec = api.openAPISpecPath
	} else {
		log.Infof("Load openapi spec from file: %s", api.openAPISpecPath)

		f := api.openAPISpecPath
		if strings.HasPrefix(f, "file://") {
			f = f[7:]
		}
		file, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		api.openAPISpec = string(file)
		api.openAPIFormat = apimanagement.Openapijson
	}
}

func (api *apiDefinition) getXMLPolicy() {
	if api.xmlPolicyPath == "" {
		log.Info("No xml policy given, load default policy")
		api.xmlPolicyFormat = apimanagement.XML
		api.xmlPolicyPath = "none (default policy)"
		api.xmlPolicy = `<policies>
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
	} else if strings.HasPrefix(api.xmlPolicyPath, "https://") || strings.HasPrefix(api.xmlPolicyPath, "http://") {
		log.Infof("Xml Policy will be downloaded by APIM during create/update from '%s'", api.xmlPolicyPath)
		api.xmlPolicyFormat = apimanagement.XMLLink
		api.xmlPolicy = api.xmlPolicyPath
	} else {
		log.Infof("Load XML policy from file: %s", api.xmlPolicyPath)

		f := api.xmlPolicyPath
		if strings.HasPrefix(f, "file://") {
			f = f[7:]
		}
		file, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		api.xmlPolicy = string(file)
		api.xmlPolicyFormat = apimanagement.XML
	}
}

func init() {
	flag.StringVar(&apimClient.subscription, "subscription", "", "Subscription of the API management service (env var: SUBSCRIPTION)")
	flag.StringVar(&apimClient.resourceGroup, "resourcegroup", "", "Name of the resource group the APIM is in (env var: RESOURCEGROUP)")
	flag.StringVar(&apimClient.serviceName, "servicename", "", "Name of the api management service (env var: APIMGMT)")
	flag.StringVar(&apiDef.openAPISpecPath, "openapispec", "", "path to the openapi spec, either file://  or https:// (env var: OPENAPISPEC)")
	flag.StringVar(&apiDef.xmlPolicyPath, "xmlpolicy", "", "path to the openapi spec , either file://  or https:// (env var: XMLPOLICY) - OPTIONAL")
	flag.StringVar(&apiDef.apiID, "apiid", "", "name (api id) of the api to deploy (env var: APIID)")
	flag.StringVar(&apiDef.apiDisplayName, "apidisplayname", "", "the display name of the api  (env var: APIDISPLAYNAME)")
	flag.StringVar(&apiDef.apiPath, "apipath", "", "the api path relative to the apim service (env var: APIPATH)")
	flag.StringVar(&apiDef.apiVersion, "apiversion", "", "version number for the versioned api deplopyment (env var: APIVERSION)")
	flag.Parse()

	if os.Getenv("SUBSCRIPTION") != "" {
		apimClient.subscription = os.Getenv("SUBSCRIPTION")
	}
	if os.Getenv("RESOURCEGROUP") != "" {
		apimClient.resourceGroup = os.Getenv("RESOURCEGROUP")
	}
	if os.Getenv("APIMGMT") != "" {
		apimClient.serviceName = os.Getenv("APIMGMT")
	}
	if os.Getenv("OPENAPISPEC") != "" {
		apiDef.openAPISpecPath = os.Getenv("OPENAPISPEC")
	}
	if os.Getenv("XMLPOLICY") != "" {
		apiDef.xmlPolicyPath = os.Getenv("XMLPOLICY")
	}
	if os.Getenv("APIID") != "" {
		apiDef.apiID = os.Getenv("APIID")
	}
	if os.Getenv("APIDISPLAYNAME") != "" {
		apiDef.apiDisplayName = os.Getenv("APIDISPLAYNAME")
	}
	if os.Getenv("APIPATH") != "" {
		apiDef.apiPath = os.Getenv("APIPATH")
	}
	if os.Getenv("APIVERSION") != "" {
		apiDef.apiVersion = os.Getenv("APIVERSION")
	}

	if (apimClient.subscription == "") ||
		(apimClient.resourceGroup == "") ||
		(apimClient.serviceName == "") ||
		(apiDef.openAPISpecPath == "") ||
		(apiDef.apiID == "") ||
		(apiDef.apiDisplayName == "") ||
		(apiDef.apiPath == "") ||
		(apiDef.apiVersion == "") {
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func main() {
	// initialize apim client
	apimClient.ctx = context.Background()
	apimClient.authenticate()

	// retrieve api definiton and xml policy
	apiDef.apiProtocols = append(apiDef.apiProtocols, apimanagement.ProtocolHTTPS)
	apiDef.apiVersioningScheme = apimanagement.VersioningSchemeSegment
	apiDef.subscriptionRequired = true
	apiDef.apiRevision = "1"
	apiDef.getOpenAPISpec()
	apiDef.getXMLPolicy()

	// create or update the specified api service
	apimClient.createOrUpdate(&apiDef)
}
