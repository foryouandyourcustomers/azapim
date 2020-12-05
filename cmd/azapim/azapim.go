package main

import (
	"context"
	"flag"
	"os"

	"github.com/foryouandyourcustomers/azapim/internal/apidefinition"
	"github.com/foryouandyourcustomers/azapim/internal/apimclient"
)

var (
	apimClient apimclient.ApimClient
	apiDef     apidefinition.Definition
)

func init() {
	flag.StringVar(&apimClient.Subscription, "subscription", "", "Subscription of the API management service (env var: SUBSCRIPTION)")
	flag.StringVar(&apimClient.ResourceGroup, "resourcegroup", "", "Name of the resource group the APIM is in (env var: RESOURCEGROUP)")
	flag.StringVar(&apimClient.ServiceName, "servicename", "", "Name of the api management service (env var: APIMGMT)")
	flag.StringVar(&apiDef.OpenAPISpecPath, "openapispec", "", "path to the openapi spec, either file://  or https:// (env var: OPENAPISPEC)")
	flag.StringVar(&apiDef.XMLPolicyPath, "xmlpolicy", "", "path to the openapi spec , either file://  or https:// (env var: XMLPOLICY) - OPTIONAL")
	flag.StringVar(&apiDef.APIID, "apiid", "", "name (api id) of the api to deploy (env var: APIID)")
	flag.StringVar(&apiDef.APIDisplayName, "apidisplayname", "", "the display name of the api  (env var: APIDISPLAYNAME)")
	flag.StringVar(&apiDef.APIPath, "apipath", "", "the api path relative to the apim service (env var: APIPATH)")
	flag.StringVar(&apiDef.APIVersion, "apiversion", "", "version number for the versioned api deplopyment (env var: APIVERSION)")
	flag.StringVar(&apiDef.APIServiceURL, "apiserviceurl", "", "Absolute URL of the backend service implementing this API (env var: APISERVICEURL)")
	flag.StringVar(&apiDef.APIProductsRaw, "apiproducts", "", "Comma separated list of products to assign the API to, Attention: tool isnt removing API from ANY products at the moment (env var: APIPRODUCTS) - OPTIONAL")
	flag.Parse()

	if os.Getenv("SUBSCRIPTION") != "" {
		apimClient.Subscription = os.Getenv("SUBSCRIPTION")
	}
	if os.Getenv("RESOURCEGROUP") != "" {
		apimClient.ResourceGroup = os.Getenv("RESOURCEGROUP")
	}
	if os.Getenv("APIMGMT") != "" {
		apimClient.ServiceName = os.Getenv("APIMGMT")
	}
	if os.Getenv("OPENAPISPEC") != "" {
		apiDef.OpenAPISpecPath = os.Getenv("OPENAPISPEC")
	}
	if os.Getenv("XMLPOLICY") != "" {
		apiDef.XMLPolicyPath = os.Getenv("XMLPOLICY")
	}
	if os.Getenv("APIID") != "" {
		apiDef.APIID = os.Getenv("APIID")
	}
	if os.Getenv("APIDISPLAYNAME") != "" {
		apiDef.APIDisplayName = os.Getenv("APIDISPLAYNAME")
	}
	if os.Getenv("APIPATH") != "" {
		apiDef.APIPath = os.Getenv("APIPATH")
	}
	if os.Getenv("APIVERSION") != "" {
		apiDef.APIVersion = os.Getenv("APIVERSION")
	}
	if os.Getenv("APISERVICEURL") != "" {
		apiDef.APIServiceURL = os.Getenv("APISERVICEURL")
	}
	if os.Getenv("APIPRODUCTS") != "" {
		apiDef.APIProductsRaw = os.Getenv("APIPRODUCTS")
	}

	if (apimClient.Subscription == "") ||
		(apimClient.ResourceGroup == "") ||
		(apimClient.ServiceName == "") ||
		(apiDef.OpenAPISpecPath == "") ||
		(apiDef.APIID == "") ||
		(apiDef.APIDisplayName == "") ||
		(apiDef.APIPath == "") ||
		(apiDef.APIVersion == "") ||
		(apiDef.APIServiceURL == "") {
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func main() {
	// initialize apim client
	apimClient.Ctx = context.Background()
	apimClient.Authenticate()

	// retrieve api definiton and xml policy
	apiDef.SetDefaults()
	apiDef.GetOpenAPISpec()
	apiDef.GetXMLPolicy()

	// create or update the specified api service
	apimClient.CreateOrUpdate(&apiDef)
}
