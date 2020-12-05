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
NAME:
   azapim - A new cli application

USAGE:
   azapim [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command
   Apis:
     versionedapi  Manage versioned apis

GLOBAL OPTIONS:
   --subscription ID     Azure Subscription ID of the API management service [$SUBSCRIPTION]
   --resourcegroup Name  Name of the resource group containing the API management service [$RESOURCEGROUP]
   --servicename Name    Name of the API management service [$APIMGMT]
   --help, -h            show help (default: false)
```

If you specify https endpoints for the openapispec or the xml policy the data is downloaded from the APIM service directly!

### Examples

```bash
# create or update the api "httpbin" with v1 and openapi spec from a https endpoint and the default xml policy
./azapim versionedapi create \
  -apidisplayname "httpbin api" \
  -apiid "httpbin" \
  -apipath "/httpbin" \
  -apiserviceurl "https://my.backend.service/httpbin" \
  -apiversion "v1" \
  -openapispec https://my.backend.service/httpbin/openapispec.json \

# create or update v2 with a custom xml policy retrieved from a local file
./azapim versionedapi create \
  -apidisplayname "httpbin api" \
  -apiid "httpbin" \
  -apipath "/httpbin" \
  -apiserviceurl "https://my.backend.service/httpbin-v2" \
  -apiversion "v2" \
  -openapispec https://my.backend.service/httpbin-v2/openapispec.json \
  -xmlpolicy "file://./policy.xml"

# create or update v2 with a custom xml policy retrieved from a local file
# and assign it to the starter and unlimited products
./azapim versionedapi create \
  -apidisplayname "httpbin api" \
  -apiid "httpbin" \
  -apipath "/httpbin" \
  -apiproducts "starter,unlimited" \
  -apiserviceurl "https://my.backend.service/httpbin-v2" \
  -apiversion "v2" \
  -openapispec https://my.backend.service/httpbin-v2/openapispec.json \
  -xmlpolicy "file://./policy.xml
```
