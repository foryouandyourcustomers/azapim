# azapim

Utility to create or update a VERSIONED api with a given openapi spec and xml policy.
This cli tool is very simple and can only create or update versioned APIs , it can't do revisions, it can't do product assignments etc.

The idea is to execute it inside a pipeline to register updates of microservices after deployments.

The utility first tries to use the login from the azure cli.
If this fails it will try to retrieve credentials from the [runtime environment](https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-environment-based-authentication).

## Installation

```bash
go get github.com/foryouandyourcustomers/azapim/cmd/azapim
```

or download the latest release.

## Usage

```bash
$ ./azapim -h
Usage of ./azapim.linux:
  -apidisplayname string
        the display name of the api  (env var: APIDISPLAYNAME)
  -apidisplayname string
        name (api id) of the api to deploy (env var: APIID)
  -apipath string
        the api path relative to the apim service (env var: APIPATH)
  -apiversion string
        version number for the versioned api deplopyment (env var: APIVERSION)
  -openapispec string
        path to the openapi spec, either file://  or https:// (env var: OPENAPISPEC)
  -resourcegroup string
        Name of the resource group the APIM is in (env var: RESOURCEGROUP)
  -servicename string
        Name of the api management service (env var: APIMGMT)
  -subscription string
        Subscription of the API management service (env var: SUBSCRIPTION)
  -xmlpolicy string
        path to the openapi spec , either file://  or https:// (env var: XMLPOLICY) - OPTIONAL
```

If you specify https endpoints for the openapispec or the xml policy the data is downloaded from the APIM service directly!

### Examples

```bash
# create or update the api "httpbin" with v1 and openapi spec from a https endpoint and the default xml policy
./azapim \
  -apidisplayname "httpbin api"
  -apiid "httpbin"
  -apipath "/httpbin"
  -apiversion "v1"
  -openapispec https://my.url/openapispec.json
```
