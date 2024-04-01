package models

import (
	"capgap/enums"
)

type Bypass struct {
	UserId         string
	ApplicationId  string
	LocationId     string
	DevicePlatform enums.DevicePlatform
	ClientType     enums.ClientType
}

// func (b Bypass) PrettyPrint() {

// }
