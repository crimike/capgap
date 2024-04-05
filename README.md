# CAPGAP

Conditional Access Policies within Azure(Entra ID) are used as if-then statements when users access resources, enforcing different authentication options(MFA, trusted device, etc) or outright Blocking access.
This tool is meant to help with overly complicated Conditional Access Policies within Azure. It does this by trying (almost) all combinations and finding common ways to bypass multiple policies at the same time.

## Usage
```
CapGap is meant to discover Azure Conditional Access Policy bypasses for certain combinations.
  -aad
        Whether to use AAD Graph or MS Graph - current default is AAD Graph (default true)
  -accessToken string
        JWT access token for the specified scope
  -appId string
        Application ID for which to check gaps
  -load
        If present, conditional access policies, users, apps and locations will be loaded from the file given(JSON format)
  -log string
        Specify log filename to log to instead of STDOUT
  -msgraph
        Whether to use AAD Graph or MS Graph - current default is AAD Graph
  -save
        If enabled, saves the conditional access policies, users, apps and locations to file(JSON format) - useful during testing
  -tenant string
        Specify tenant ID 
  -userId string
        User ObjectId for which to check gaps
  -v    Verbose logging
```


## Current caveats

Currently the following properties of a policy are parsed:
* User/group/role
* Application
* DevicePlatform
* Location
* Client Type

* Parsing the device filter has not been implemented yet - first idea is to search for certain keywords(that might be user-editable)
* User or signin risk not implemented yet
* Every Conditional Access Policy is treated as blocking(meaning that if one applies, the assumption is that the controls are good enough)
* Uses AADGraph
* Only access token login is currently implemented

## Next items - TODO

In no particular order:
* Implement Microsoft Graph support - either via the API or the Graph Go library
* Implement different login options(device code, user/password, refresh token)
* Expand custom authentication strengths
* Parse controls to see any gaps there
* Find gaps per Conditional Access Policy
* Parse DeviceFilter
* Take into account user/signin risk
* PrettyPrint bypasses
* recursive roles and groups
* Report for dynamic groups included/excluded - to check if they contain user-editable properties in the filters

## Acknowledgements/inspirations

* Retrieval of CAPs and locations was copied from Dirk-jan's [ROADTools](https://github.com/dirkjanm/ROADtools)

