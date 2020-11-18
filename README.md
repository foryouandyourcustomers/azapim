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
  -apiproducts string
        Comma separated list of products to assign the API to, Attention: tool isnt removing API from ANY products at the moment (env var: APIPRODUCTS) - OPTIONAL
  -apiserviceurl string
        Absolute URL of the backend service implementing this API (env var: APISERVICEURL)
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
  -apidisplayname "httpbin api" \
  -apiid "httpbin" \
  -apipath "/httpbin" \
  -apiserviceurl "https://my.backend.service/httpbin" \
  -apiversion "v1" \
  -openapispec https://my.backend.service/httpbin/openapispec.json \

# create or update v2 with a custom xml policy retrieved from a local file
./azapim \
  -apidisplayname "httpbin api" \
  -apiid "httpbin" \
  -apipath "/httpbin" \
  -apiserviceurl "https://my.backend.service/httpbin-v2" \
  -apiversion "v2" \
  -openapispec https://my.backend.service/httpbin-v2/openapispec.json \
  -xmlpolicy "file://./policy.xml"

# create or update v2 with a custom xml policy retrieved from a local file
# and assign it to the starter and unlimited products
./azapim \
  -apidisplayname "httpbin api" \
  -apiid "httpbin" \
  -apipath "/httpbin" \
  -apiproducts "starter,unlimited" \
  -apiserviceurl "https://my.backend.service/httpbin-v2" \
  -apiversion "v2" \
  -openapispec https://my.backend.service/httpbin-v2/openapispec.json \
  -xmlpolicy "file://./policy.xml
```
