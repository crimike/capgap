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

An AAD Graph token can be retrieved as follows:
```
PS> (Get-AzAccessToken -ResourceUrl https://graph.windows.net).Token
```

### Examples 
The following checks one specific entry, namely the user represented by the ObjectId connecting to the app represented by the ApplicationID. The 3 combinations listed are 3 possibilities to access the app as that user without triggering any Conditional Access Policies to be applied to the SignIn
```
term> ./capgap -load -userId 12345678-1234-1234-1234-123456789ab  -appId 9876543-1234-1234-1234-123456789ab
INFO: 16:19:15 log.go:47: Starting CAPGAP at 16:19:15.077541
Bypasses for user John going to InternalAppA
ClientType: Any; DevicePlatform: Any; Location: MainOffice (1.1.1.1/24, )
ClientType: Any; DevicePlatform: macOS; Location: SecondaryOffice (8.8.8.8/28, 9.9.9.9/28, )
ClientType: Any; DevicePlatform: Windows_Phone; Location: MainOffice (1.1.1.1/24, )
INFO: 16:19:15 log.go:52: Parsing finished at 16:19:15.277670
```

If not loading from file, accessToken and tenant need to be supplied:
```
term> ./capgap -save -accessToken eyJ0.... -tenant 99999999-1234-1234-1234-123456789ab -userId 12345678-1234-1234-1234-123456789ab  -appId 9876543-1234-1234-1234-123456789ab
```

An entire report can be generated as such.
```
term> ./capgap -load -report report.txt
```


## Current caveats

* Note that the tool does not establish if the users actually have access to the application, but it checks for paths where Conditional Access Policies might not apply
* Parsing the device filter has not been implemented yet - first idea is to search for certain keywords(that might be user-editable)
* User or signin risk not implemented yet
* Every Conditional Access Policy is treated as blocking(meaning that if one applies, the assumption is that the controls are good enough)
* Uses AADGraph
* Only access token login is currently implemented


## Internals

* For users, the ObjectID is the PK. For applications, it is the ApplicationID. For Locations, it is the PolicyId
* When getting the full report(no appId/userId) depending on your tenant, it might be time/memory consuming. Some optimizations have been added, including in the report generation, but overall, the parsing basically tries all combinations.
Currently the following properties of a policy are parsed:
* User/group/role
* Application
* DevicePlatform
* Location
* Client Type



## Next items - TODO

In no particular order:
* Implement Microsoft Graph support - either via the API or the Graph Go library
* Implement different login options(device code, user/password, refresh token)
* Expand custom authentication strengths
* Parse controls to see any gaps there
* Find gaps per Conditional Access Policy
* Parse DeviceFilter
* Take into account user/signin risk
* recursive roles and groups
* Report for dynamic groups included/excluded - to check if they contain user-editable properties in the filters
* Checks for well-known apps
* Report for PIM

## Acknowledgements/inspirations

* Retrieval of CAPs and locations was copied from Dirk-jan's [ROADTools](https://github.com/dirkjanm/ROADtools)

