package apimclient

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/apimanagement/mgmt/apimanagement"
	"github.com/google/uuid"
)

// CreateOrUpdatePolicy updates the xml policy of the given api
func (apim *ApimClient) CreateOrUpdatePolicy(
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
	apiPolicyContract, err := apim.PolicyClient.CreateOrUpdate(
		apim.Ctx,
		apim.ResourceGroup,
		apim.ServiceName,
		uid,
		apiPolicy,
		uuid.New().String(),
	)
	if err != nil {
		return apiPolicyContract, err
	}
	return apiPolicyContract, nil
}
