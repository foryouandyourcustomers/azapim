package main

import (
	"context"
	"flag"
	"fmt"
	"io"
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
	policyClient     apimanagement.PolicyClient

	subscription  string
	resourceGroup string
	serviceName   string
}

func (apim *azureApimClient) authenticate() {
	apim.apiClient = apimanagement.NewAPIClient(apim.subscription)
	apim.versionSetClient = apimanagement.NewAPIVersionSetClient(apim.subscription)
	apim.policyClient = apimanagement.NewPolicyClient(apim.subscription)

	a, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		// looking at the newauthorizerfromenvrionment funciton it
		// seems that thing never returns an error whatsoever!
		log.Debug("Unable to create authorizer from az cli. Lets load the authorizer from the environment and hope for the best!")
		a, _ = auth.NewAuthorizerFromEnvironment()

	}
	apim.apiClient.Authorizer = a
	apim.versionSetClient.Authorizer = a
	apim.policyClient.Authorizer = a
}

func (apim *azureApimClient) createOrUpdate(apidef *apiDefinition) {
	apiVersionSetID := fmt.Sprintf("%s%s", apidef.apiName, apidef.apiVersion)

	apiVersionSet := apimanagement.APIVersionSetContract{
		APIVersionSetContractProperties: &apimanagement.APIVersionSetContractProperties{
			DisplayName:      &apidef.apiDisplayName,
			VersioningScheme: apidef.apiVersioningScheme,
		},
	}

	log.Infof("Creating/Updating API versionset '%s'", apiVersionSetID)
	apiVersionSetContract, err := apim.versionSetClient.CreateOrUpdate(
		apim.ctx,
		apim.resourceGroup,
		apim.serviceName,
		apiVersionSetID,
		apiVersionSet,
		uuid.New().String())
	if err != nil {
		log.Fatalf("cannot create api versionset endpoint: %v", err)
	}
	log.Info("Created/Updated API version set '%s'", *apiVersionSetContract.ID)

	apiProperties := apimanagement.APICreateOrUpdateParameter{
		APICreateOrUpdateProperties: &apimanagement.APICreateOrUpdateProperties{
			Format:               apidef.openAPIFormat,
			DisplayName:          &apidef.apiDisplayName,
			Value:                &apidef.openAPISpecPath,
			Protocols:            &apidef.apiProtocols,
			Path:                 &apidef.apiPath,
			SubscriptionRequired: &apidef.subscriptionRequired,
			APIVersion:           &apidef.apiVersion,
			APIVersionSetID:      apiVersionSetContract.ID,
		},
	}

	log.Infof("Creating/Updating API '%s'", apidef.apiDisplayName)
	future, err := apim.apiClient.CreateOrUpdate(
		apim.ctx,
		apim.resourceGroup,
		apim.serviceName,
		apidef.apiName,
		apiProperties,
		uuid.New().String())
	if err != nil {
		log.Fatalf("cannot create api endpoint: %v", err)
	}
	err = future.WaitForCompletionRef(apim.ctx, apim.apiClient.Client)
	if err != nil {
		log.Fatalf("cannot get the api endpoint future response: %v", err)
	}

	//apim.apiClient.Cre-ateOrUpdate(apim.ctx, apim.resourceGroupName, apim.serviceName, api.apiName, parameters apimanagement.APICreateOrUpdateParameter, ifMatch string)
}

type apiDefinition struct {
	openAPISpecPath string
	openAPISpec     io.ReadCloser
	openAPIFormat   apimanagement.ContentFormat

	xmlPolicyPath string
	xmlPolicy     io.ReadCloser

	apiName             string
	apiDisplayName      string
	apiVersion          string
	apiVersioningScheme apimanagement.VersioningScheme
	apiPath             string

	apiProtocols         []apimanagement.Protocol
	subscriptionRequired bool
}

func (api *apiDefinition) getOpenAPISpec() {

	if strings.HasPrefix(api.openAPISpecPath, "https://") || strings.HasPrefix(api.openAPISpecPath, "http://") {
		// log.Infof("Download openapi spec from: %s", api.openApiSpecPath)
		// resp, err := http.Get(api.openAPISpecPath)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// if resp.StatusCode != http.StatusOK {
		// 	log.Fatalf("Unable to download file - Status Code: %d", resp.StatusCode)
		// }
		// api.openAPISpec = resp.Body

		log.Infof("OpenApi Spec will be downloaded by APIM during create/update from '%s'", api.openAPISpecPath)
		api.openAPIFormat = apimanagement.OpenapijsonLink
	} else {
		log.Infof("Load openapi spec from file: %s", api.openAPISpecPath)

		f := api.openAPISpecPath
		if strings.HasPrefix(f, "file://") {
			f = f[7:]
		}
		file, err := os.Open(f) // For read access.
		if err != nil {
			log.Fatal(err)
		}
		api.openAPISpec = file
		api.openAPIFormat = apimanagement.Openapijson
		log.Fatal("Loading from file is currently not implemented !")
	}
}

func (api *apiDefinition) getXMLPolicy() {
	if api.xmlPolicyPath == "" {
		log.Info("No xml policy given, load default policy")
		defXMLPolicy := `<policies>
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
		api.xmlPolicy = ioutil.NopCloser(strings.NewReader(defXMLPolicy))
	}
}

func init() {
	flag.StringVar(&apimClient.subscription, "s", "", "Subscription of the API management service (env var: SUBSCRIPTION)")
	flag.StringVar(&apimClient.resourceGroup, "r", "", "Name of the resource group the APIM is in (env var: RESOURCEGROUP)")
	flag.StringVar(&apimClient.serviceName, "a", "", "Name of the api management service (env var: APIMGMT)")
	flag.StringVar(&apiDef.openAPISpecPath, "o", "", "path to the openapi spec, either file://  or http:// (env var: OPENAPISPEC)")
	flag.StringVar(&apiDef.xmlPolicyPath, "x", "", "path to the openapi spec , either file://  or http:// (env var: XMLPOLICY)")
	flag.StringVar(&apiDef.apiName, "n", "", "name (api id) of the api to deploy (env var: APINAME)")
	flag.StringVar(&apiDef.apiDisplayName, "d", "", "the display name of the api  (env var: APIDISPLAYNAME)")
	flag.StringVar(&apiDef.apiPath, "p", "", "the api path relative to the apim service (env var: APIPATH)")
	flag.StringVar(&apiDef.apiVersion, "v", "", "version number for the versioned api deplopyment (env var: APIVERSION)")
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
	if os.Getenv("APINAME") != "" {
		apiDef.apiName = os.Getenv("APINAME")
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
		(apiDef.apiName == "") ||
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
	apiDef.getOpenAPISpec()
	apiDef.getXMLPolicy()

	// create or update the specified api service
	apimClient.createOrUpdate(&apiDef)

	//apimClient.client.CreateOrUpdate(apimClient.ctx, resourceGroupName string, serviceName string, apiid string, parameters apimanagement.APICreateOrUpdateParameter, ifMatch string)

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(apiDef.xmlPolicy)
	// s := buf.String()
	// fmt.Println(s)

	// r, err := apimClient.client.Get(apimClient.ctx, "dev-mbp-infrastructure", "dev-mbp-apim", "mbp-pdf-html")

	// if err != nil {
	// 	log.Fatal(err)
	// }

}
