package enums

type ConditionalAccessPolicyState string

const (
	StateEnabled   ConditionalAccessPolicyState = "Enabled"
	StateReporting ConditionalAccessPolicyState = "Reporting"
	StateDisabled  ConditionalAccessPolicyState = "Disabled"
)

type DevicePlatform string

const (
	DevicePlatformAll          DevicePlatform = "All"
	DevicePlatformAndroid      DevicePlatform = "Android"
	DevicePlatformIOS          DevicePlatform = "iOS"
	DevicePlatformMacOS        DevicePlatform = "macOS"
	DevicePlatformLinux        DevicePlatform = "Linux"
	DevicePlatformWindows      DevicePlatform = "Windows"
	DevicePlatformWindowsPhone DevicePlatform = "Windows_Phone"
)

func GetAllDevicePlatforms() []DevicePlatform {
	return []DevicePlatform{DevicePlatformAndroid, DevicePlatformIOS, DevicePlatformMacOS, DevicePlatformLinux, DevicePlatformWindows, DevicePlatformWindowsPhone}
}

type ClientType string

const (
	ClientTypeEasSupported   ClientType = "EasSupported"
	ClientTypeEasUnsupported ClientType = "EasUnsupported"
	ClientTypeOtherLegacy    ClientType = "OtherLegacy"
	ClientTypeLegacySmtp     ClientType = "LegacySmtp"
	ClientTypeLegacyPop      ClientType = "LegacyPop"
	ClientTypeLegacyImap     ClientType = "LegacyImap"
	ClientTypeLegacyMapi     ClientType = "LegacyMapi"
	ClientTypeLegacyOffice   ClientType = "LegacyOffice"
	ClientTypeNative         ClientType = "Native"
	ClientTypeBrowser        ClientType = "Browser"
)

func GetAllClientTypes() []ClientType {
	return []ClientType{ClientTypeEasSupported, ClientTypeEasUnsupported, ClientTypeOtherLegacy, ClientTypeLegacySmtp, ClientTypeLegacyPop, ClientTypeLegacyImap, ClientTypeLegacyMapi, ClientTypeLegacyOffice, ClientTypeNative, ClientTypeBrowser}
}

type GuestOrExternalUserType string

const (
	GuestOrExternalUserTypeInternalGuest          GuestOrExternalUserType = "InternalGuest"
	GuestOrExternalUserTypeB2bCollaborationGuest  GuestOrExternalUserType = "B2bCollaborationGuest"
	GuestOrExternalUserTypeB2bCollaborationMember GuestOrExternalUserType = "B2bCollaborationMember"
	GuestOrExternalUserTypeB2bDirectConnectUser   GuestOrExternalUserType = "B2bDirectConnectUser"
	GuestOrExternalUserTypeOtherExternalUser      GuestOrExternalUserType = "OtherExternalUser"
	GuestOrExternalUserTypeServiceProvider        GuestOrExternalUserType = "ServiceProvider"
)

type ApplicationACR string

const (
	UserActionDeviceRegisterOrJoin ApplicationACR = "urn:user:registerdevice"
	AuthenticationContextPIM       ApplicationACR = "c1"
)

type ConditionalAccessControl string

const (
	ConditionalAccessBlock                     ConditionalAccessControl = "Block"
	ConditionalAccessMfa                       ConditionalAccessControl = "Mfa"
	ConditionalAccessRequireCompliantDevice    ConditionalAccessControl = "RequireCompliantDevice"
	ConditionalAccessRequireDomainJoinedDevice ConditionalAccessControl = "RequireDomainJoinedDevice"
	ConditionalAccessRequireApprovedApp        ConditionalAccessControl = "RequireApprovedApp"
)

type ApplicationMode string

const (
	ApplicationModeCloudApps             ApplicationMode = "CloudApps"
	ApplicationModeUserActions           ApplicationMode = "UserActions"
	ApplicationModeAuthenticationContext ApplicationMode = "PIM Authentication Context"
)

const DynamicGroupType string = "DynamicMembership"
