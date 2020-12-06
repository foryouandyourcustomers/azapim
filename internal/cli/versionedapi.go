package cli

import (
	"github.com/foryouandyourcustomers/azapim/internal/apidefinition"
	ucli "github.com/urfave/cli/v2"
)

var (
	apiDef apidefinition.Definition

	// VersionedAPICli contains the versionedapi cli definition
	VersionedAPICli = []*ucli.Command{
		{
			Name:     "versionedapi",
			Category: "Apis",
			Usage:    "Manage versioned apis",
			Flags: []ucli.Flag{
				&ucli.StringFlag{
					Name:        "apiid",
					Usage:       "name (api id) of the api to deploy",
					Required:    true,
					EnvVars:     []string{"APIID"},
					Destination: &apiDef.APIID,
				},
			},
			Subcommands: []*ucli.Command{
				{
					Name:  "create",
					Usage: "Create or Update a versioned api",
					Action: func(c *ucli.Context) error {
						apiDef.SetDefaults()
						err := apiDef.GetOpenAPISpec()
						if err != nil {
							return ucli.Exit(err, 1)
						}
						err = apiDef.GetXMLPolicy()
						if err != nil {
							return ucli.Exit(err, 1)
						}
						err = apimClient.CreateOrUpdate(&apiDef)
						if err != nil {
							return ucli.Exit(err, 1)
						}
						return nil
					},
					Flags: []ucli.Flag{
						&ucli.StringFlag{
							Name:        "openapispec",
							Usage:       "Url or path to openapi spec definition (file:// or https://)",
							Required:    true,
							EnvVars:     []string{"OPENAPISPEC"},
							Destination: &apiDef.OpenAPISpecPath,
						},
						&ucli.StringFlag{
							Name:        "xmlpolicy",
							Usage:       "Url or path to xml policy (file:// or https://)",
							Required:    false,
							EnvVars:     []string{"XMLPOLICY"},
							Destination: &apiDef.XMLPolicyPath,
						},
						&ucli.StringFlag{
							Name:        "apipath",
							Usage:       "the api path relative to the apim service url",
							Required:    true,
							EnvVars:     []string{"APIPATH"},
							Destination: &apiDef.APIPath,
						},
						&ucli.StringFlag{
							Name:        "apiversion",
							Usage:       "version number for the versioned api deplopyment",
							Required:    true,
							EnvVars:     []string{"APIVERSION"},
							Destination: &apiDef.APIVersion,
						},
						&ucli.StringFlag{
							Name:        "apiserviceurl",
							Usage:       "Absolute URL of the backend service implementing this API",
							Required:    true,
							EnvVars:     []string{"APISERVICEURL"},
							Destination: &apiDef.APIServiceURL,
						},
						&ucli.StringFlag{
							Name:        "apiproducts",
							Usage:       "Comma separated list of products to assign the API to, Attention: tool isnt removing API from ANY products at the moment",
							Required:    false,
							EnvVars:     []string{"APIPRODUCTS"},
							Destination: &apiDef.APIProductsRaw,
						},
						&ucli.StringFlag{
							Name:        "apidisplayname",
							Usage:       "Display name in the API management service ",
							Required:    false,
							EnvVars:     []string{"APIDISPLAYNAME"},
							Destination: &apiDef.APIDisplayName,
						},
					},
				},
			},
		},
	}
)
