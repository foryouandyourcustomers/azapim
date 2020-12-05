package apimclient

import "github.com/Azure/azure-sdk-for-go/profiles/latest/apimanagement/mgmt/apimanagement"

// AssignToProduct assigns an API to a given (existing) product
func (apim *ApimClient) AssignToProduct(p string, id string) (apimanagement.APIContract, error) {
	contract, err := apim.ProductsAPIClient.CreateOrUpdate(apim.Ctx, apim.ResourceGroup, apim.ServiceName, p, id)
	if err != nil {
		return apimanagement.APIContract{}, err
	}
	return contract, nil
}
