package apimclient

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/apimanagement/mgmt/apimanagement"

	"github.com/Azure/go-autorest/autorest/azure/auth"
	log "github.com/sirupsen/logrus"

	"github.com/foryouandyourcustomers/azapim/internal/apidefinition"
)

// ApimClient represents the azure api management service clients
type ApimClient struct {
	Ctx               context.Context
	APIClient         apimanagement.APIClient
	VersionSetClient  apimanagement.APIVersionSetClient
	PolicyClient      apimanagement.APIPolicyClient
	ProductsAPIClient apimanagement.ProductAPIClient
	ProductClient     apimanagement.ProductClient
	ServiceClient     apimanagement.ServiceClient
	Subscription      string
	ResourceGroup     string
	ServiceName       string
}

// Authenticate against the definied azure subscription
func (apim *ApimClient) Authenticate() {
	apim.APIClient = apimanagement.NewAPIClient(apim.Subscription)
	apim.VersionSetClient = apimanagement.NewAPIVersionSetClient(apim.Subscription)
	apim.PolicyClient = apimanagement.NewAPIPolicyClient(apim.Subscription)
	apim.ProductsAPIClient = apimanagement.NewProductAPIClient(apim.Subscription)
	apim.ProductClient = apimanagement.NewProductClient(apim.Subscription)
	apim.ServiceClient = apimanagement.NewServiceClient(apim.Subscription)

	a, err := auth.NewAuthorizerFromCLI()
	if err != nil {
		// looking at the newauthorizerfromenvrionment funciton it
		// seems that thing never returns an error whatsoever!
		log.Debug("Unable to create authorizer from az cli. Lets load the authorizer from the environment and hope for the best!")
		a, _ = auth.NewAuthorizerFromEnvironment()

	}
	apim.APIClient.Authorizer = a
	apim.VersionSetClient.Authorizer = a
	apim.PolicyClient.Authorizer = a
	apim.ProductsAPIClient.Authorizer = a
	apim.ProductClient.Authorizer = a
	apim.ServiceClient.Authorizer = a

	// increase the polling timeout for the service client to 30 minutes
	apim.ServiceClient.Client.PollingDuration = 30 * time.Minute
}

// CreateOrUpdate - create or update the specified api
func (apim *ApimClient) CreateOrUpdate(a *apidefinition.Definition) error {
	log.Infof("Creating/Updating API versionset: '%s'", a.APIID)
	versionSet, err := apim.CreateOrUpdateVersionSet(a.APIDisplayName, a.APIVersioningScheme, a.APIID)
	if err != nil {
		return err
	}
	log.Infof("Created/Updated API versionset: '%s'", *versionSet.ID)

	log.Infof("Creating/Updating API: '%s' with version '%s' (unique id: %s)", a.APIDisplayName, a.APIVersion, a.APIUniqueID)
	api, err := apim.CreateOrUpdateAPI(
		a.OpenAPIFormat,
		a.APIDisplayName,
		a.OpenAPISpec,
		a.APIProtocols,
		a.APIPath,
		a.SubscriptionRequired,
		a.APIVersion,
		versionSet,
		a.APIRevision,
		a.APIUniqueID,
		a.APIServiceURL,
	)
	if err != nil {
		return err
	}
	log.Info("Created/Updated API '%s'", *api.ID)

	log.Info("Creating/Updating API Policy")
	policy, err := apim.CreateOrUpdatePolicy(a.XMLPolicyFormat, a.XMLPolicy, a.APIUniqueID)
	if err != nil {
		return err
	}
	log.Infof("Created/Updated API versionset: '%s'", *policy.ID)

	for _, v := range a.APIProducts {
		log.Infof("Assign API to product '%s'", v)
		_, err := apim.AssignToProduct(v, a.APIUniqueID)
		if err != nil {
			return err
		}
		log.Info("Assigned API to product")
	}
	return nil
}
