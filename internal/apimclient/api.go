package apimclient

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/apimanagement/mgmt/apimanagement"
	"github.com/google/uuid"
)

// CreateOrUpdateAPI creates or updates the API definition against the API management service
func (apim *ApimClient) CreateOrUpdateAPI(
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
) (apimanagement.APIContract, error) {
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
	future, err := apim.APIClient.CreateOrUpdate(
		apim.Ctx,
		apim.ResourceGroup,
		apim.ServiceName,
		uid,
		apiProperties,
		uuid.New().String())
	if err != nil {
		return apimanagement.APIContract{}, err
	}
	err = future.WaitForCompletionRef(apim.Ctx, apim.APIClient.Client)
	if err != nil {
		return apimanagement.APIContract{}, err
	}
	contract, err := future.Result(apim.APIClient)
	if err != nil {
		return apimanagement.APIContract{}, err
	}
	return contract, nil
}
