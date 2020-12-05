package apimclient

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/apimanagement/mgmt/apimanagement"
	"github.com/google/uuid"
)

// CreateOrUpdateVersionSet create or updates the specified version set of the api
func (apim *ApimClient) CreateOrUpdateVersionSet(
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

	apiVersionSetContract, err := apim.VersionSetClient.CreateOrUpdate(
		apim.Ctx,
		apim.ResourceGroup,
		apim.ServiceName,
		id,
		apiVersionSet,
		uuid.New().String())
	if err != nil {
		return apiVersionSetContract, err
	}
	return apiVersionSetContract, nil
}
