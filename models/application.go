package models

type Application struct {
	//commenting out objectId to avoid confusion, since the applicationID is mentioned in conditional access policies
	// ObjectId      string
	ApplicationId string
	DisplayName   string
}
