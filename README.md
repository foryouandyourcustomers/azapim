# azapim

Utility to create or update a VERSIONED api with a given openapi spec

## Installation

```bash
go get github.com/foryouandyourcustomers/azapim/cmd/azapim
```

or download the latest release.

<!-- 
## Usage

```bash
./kvget -h
Usage of ./kvget.linux:
  -s string
        Name of the secret to retrieve (env var: SECRET)
  -v string
        Name of the keyvault (env var: VAULT)
```

Lets retrieve the value of the secret "myawesomesecret" from the keyvault "fyayctestvault"
```bash
# via cli flags
./kvget -s myawesomesecret -v fyayctestvault

# via env vars
SECRET=myawesomesecret VAULT=fyayctestvault ./kvget
```

The utility first tries to use the login from the azure cli.
If this fails it will try to retrieve credentials from the [runtime environment](https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-environment-based-authentication). -->
