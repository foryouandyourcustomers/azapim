package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

var (
	apimClient azureApimClient
	apiDef     apiDefinition

	subscriptionID           string
	resourceGroupName        string
	apiManagementServiceName string

	openAPISpecPath string
	xmlPolicyPath   string
)

type azureApimClient struct {
	ctx    context.Context
	client apimanagement.APIClient

	subscription  string
	resourceGroup string
	serviceName   string
}

func (apim *azureApimClient) authenticate(sub string, rg string, srv string) {
	apim.subscription = sub
	apim.resourceGroup = rg
	apim.serviceName = srv

	apim.client = apimanagement.NewAPIClient(sub)

	a, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		// looking at the newauthorizerfromenvrionment funciton it
		// seems that thing never returns an error whatsoever!
		log.Debug("Unable to create authorizer from az cli. Lets load the authorizer from the environment and hope for the best!")
		a, _ = auth.NewAuthorizerFromEnvironment()

	}
	apim.client.Authorizer = a
}

type apiDefinition struct {
	openAPISpec io.ReadCloser
	xmlPolicy   io.ReadCloser
}

func (api *apiDefinition) getOpenAPISpec(path string) {

	if strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") {
		log.Infof("Download openapi spec from: %s", path)
		resp, err := http.Get(path)
		if err != nil {
			log.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Unable to download file - Status Code: %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		api.openAPISpec = resp.Body
	} else {
		file_path := path
		if strings.HasPrefix(file_path, "file://") {
			file_path = file_path[7:]
		}
		fmt.Print(file_path)
	}

	// if strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://") {
	// 	resp, err := c.httpClient.Get(name)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return resp.Body, nil
	// }
	// return os.Open(name)

	//	fmt.Print(path[7:])

}

func init() {
	flag.StringVar(&subscriptionID, "s", "", "Subscription of the API management service (env var: SUBSCRIPTION)")
	flag.StringVar(&resourceGroupName, "r", "", "Name of the resource group the APIM is in (env var: RESOURCEGROUP)")
	flag.StringVar(&apiManagementServiceName, "a", "", "Name of the api management service (env var: APIMGMT)")
	flag.StringVar(&openAPISpecPath, "o", "", "path to the openapi spec (either file://  or http:// (env var: OPENAPISPEC)")
	flag.StringVar(&xmlPolicyPath, "x", "", "path to the openapi spec (either file://  or http:// (env var: XMLPOLICY)")
	flag.Parse()

	if os.Getenv("SUBSCRIPTION") != "" {
		subscriptionID = os.Getenv("SUBSCRIPTION")
	}
	if os.Getenv("RESOURCEGROUP") != "" {
		resourceGroupName = os.Getenv("RESOURCEGROUP")
	}
	if os.Getenv("APIMGMT") != "" {
		apiManagementServiceName = os.Getenv("APIMGMT")
	}
	if os.Getenv("OPENAPISPEC") != "" {
		openAPISpecPath = os.Getenv("OPENAPISPEC")
	}
	if os.Getenv("XMLPOLICY") != "" {
		xmlPolicyPath = os.Getenv("XMLPOLICY")
	}

	if (subscriptionID == "") || (resourceGroupName == "") || (apiManagementServiceName == "") || (openAPISpecPath == "") {
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func main() {
	// initialize apim client
	apimClient.ctx = context.Background()
	apimClient.authenticate(subscriptionID, resourceGroupName, apiManagementServiceName)

	// retrieve api definiton and xml policy
	apiDef.getOpenAPISpec(openAPISpecPath)

	r, err := apimClient.client.Get(apimClient.ctx, "dev-mbp-infrastructure", "dev-mbp-apim", "mbp-pdf-html")

	if err != nil {
		log.Fatal(err)
	}

	d := *r.DisplayName

	fmt.Print(d)
}
