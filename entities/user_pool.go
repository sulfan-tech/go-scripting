package entities

type UserPool struct {
	UserID             string
	AliasAttributes    string
	MFASetting         string
	MFAMethods         string
	AccountStatus      string
	ConfirmationStatus string
	CreatedTime        string
	LastUpdatedTime    string
	UserAttributes     map[string]string
}
