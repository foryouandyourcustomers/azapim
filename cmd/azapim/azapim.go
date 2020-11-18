package main

import (
	"context"
	"flag"
	"fmt"
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

func (apim *azureApimClient) createOrUpdateVersionSet(
	dn string,
	vs apimanagement.VersioningScheme,
	id string,
) (apimanagement.APIVersionSetContract, error) {
	apiVersionSet := apimanagement.APIVersionSetContract{
		APIVersionSetContractProperties: &apimanagement.APIVersionSetContractProperties{
			DisplayName:      &dn,
			VersioningScheme: vs,
		},
	}

	apiVersionSetContract, err := apim.versionSetClient.CreateOrUpdate(
		apim.ctx,
		apim.resourceGroup,
		apim.serviceName,
		id,
		apiVersionSet,
		uuid.New().String())
	if err != nil {
		return apiVersionSetContract, err
	}
	return apiVersionSetContract, nil
}

func (apim *azureApimClient) createOrUpdateAPI(
	cf apimanagement.ContentFormat,
	dn string,
	va string,
	pr []apimanagement.Protocol,
	pa string,
	sr bool,
	ve string,
	c apimanagement.APIVersionSetContract,
	re string,
	uid string,
	su string,
) error {
	apiProperties := apimanagement.APICreateOrUpdateParameter{
		APICreateOrUpdateProperties: &apimanagement.APICreateOrUpdateProperties{
			Format:               cf,
			DisplayName:          &dn,
			Value:                &va,
			Protocols:            &pr,
			Path:                 &pa,
			SubscriptionRequired: &sr,
			APIVersion:           &ve,
			APIVersionSetID:      c.ID,
			APIRevision:          &re,
			ServiceURL:           &su,
		},
	}
	future, err := apim.apiClient.CreateOrUpdate(
		apim.ctx,
		apim.resourceGroup,
		apim.serviceName,
		uid,
		apiProperties,
		uuid.New().String())
	if err != nil {
		return err
	}
	err = future.WaitForCompletionRef(apim.ctx, apim.apiClient.Client)
	if err != nil {
		return err
	}
	return nil

}

func (apim *azureApimClient) createOrUpdatePolicy(
	f apimanagement.PolicyContentFormat,
	p string,
	uid string,

) (apimanagement.PolicyContract, error) {
	apiPolicy := apimanagement.PolicyContract{
		PolicyContractProperties: &apimanagement.PolicyContractProperties{
			Format: f,
			Value:  &p,
		},
	}
	apiPolicyContract, err := apim.apiPolicyClient.CreateOrUpdate(
		apim.ctx,
		apim.resourceGroup,
		apim.serviceName,
		uid,
		apiPolicy,
		uuid.New().String(),
	)
	if err != nil {
		return apiPolicyContract, err
	}
	return apiPolicyContract, nil
}

func (apim *azureApimClient) createOrUpdate(a *apiDefinition) {

	log.Infof("Creating/Updating API versionset: '%s'", a.apiID)
	versionSet, err := apim.createOrUpdateVersionSet(a.apiDisplayName, a.apiVersioningScheme, a.apiID)
	if err != nil {
		log.Fatalf("cannot create/update version set: %v\n", err)
	}
	log.Infof("Created/Updated API versionset: '%s'", *versionSet.ID)

	log.Infof("Creating/Updating API versionset: '%s' with version '%s' (unique id: %s)", a.apiDisplayName, a.apiVersion, a.apiUniqueID)
	err = apim.createOrUpdateAPI(
		a.openAPIFormat,
		a.apiDisplayName,
		a.openAPISpec,
		a.apiProtocols,
		a.apiPath,
		a.subscriptionRequired,
		a.apiVersion,
		versionSet,
		a.apiRevision,
		a.apiUniqueID,
		a.apiServiceURL,
	)
	if err != nil {
		log.Fatalf("cannot create/update API: %v\n", err)
	}
	log.Info("Created/Updated API")

	log.Infof("Creating/Updating API Policy for %s", a.apiUniqueID)
	policy, err := apimClient.createOrUpdatePolicy(a.xmlPolicyFormat, a.xmlPolicy, a.apiUniqueID)
	if err != nil {
		log.Fatalf("cannot create/update policy: %v\n", err)
	}
	log.Infof("Created/Updated API versionset: '%s'", *policy.ID)
}

type apiDefinition struct {
	openAPISpecPath string
	openAPISpec     string
	openAPIFormat   apimanagement.ContentFormat

	xmlPolicyPath   string
	xmlPolicy       string
	xmlPolicyFormat apimanagement.PolicyContentFormat

	apiID               string
	apiUniqueID         string
	apiDisplayName      string
	apiVersion          string
	apiVersioningScheme apimanagement.VersioningScheme
	apiPath             string
	apiRevision         string
	apiServiceURL       string

	apiProtocols         []apimanagement.Protocol
	subscriptionRequired bool
}

func (api *apiDefinition) setDefaults() {
	api.apiProtocols = append(api.apiProtocols, apimanagement.ProtocolHTTPS)
	api.apiVersioningScheme = apimanagement.VersioningSchemeSegment
	api.subscriptionRequired = true
	api.apiRevision = "1"
	api.apiUniqueID = fmt.Sprintf("%s-%s", api.apiID, api.apiVersion)
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
	flag.StringVar(&apiDef.apiServiceURL, "apiserviceurl", "", "Absolute URL of the backend service implementing this API (env var: APISERVICEURL)")
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
	if os.Getenv("APISERVICEURL") != "" {
		apiDef.apiServiceURL = os.Getenv("APISERVICEURL")
	}

	if (apimClient.subscription == "") ||
		(apimClient.resourceGroup == "") ||
		(apimClient.serviceName == "") ||
		(apiDef.openAPISpecPath == "") ||
		(apiDef.apiID == "") ||
		(apiDef.apiDisplayName == "") ||
		(apiDef.apiPath == "") ||
		(apiDef.apiVersion == "") ||
		(apiDef.apiServiceURL == "") {
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func main() {
	// initialize apim client
	apimClient.ctx = context.Background()
	apimClient.authenticate()

	// retrieve api definiton and xml policy
	apiDef.setDefaults()
	apiDef.getOpenAPISpec()
	apiDef.getXMLPolicy()

	// create or update the specified api service
	apimClient.createOrUpdate(&apiDef)
}
