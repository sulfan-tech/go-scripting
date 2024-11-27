package entities

import "time"

type User struct {
	Uid                 string    `firestore:"uid,omitempty" dynamodbav:"uid,omitempty" json:"uid"`
	IsCDP               bool      `firestore:"isCDP,omitempty" dynamodbav:"isCDP,omitempty" json:"isCDP"`
	IsCorporate         bool      `firestore:"isCorporate,omitempty" dynamodbav:"isCorporate,omitempty" json:"isCorporate"`
	EmployeeCode        string    `firestore:"employeeCode,omitempty" dynamodbav:"employeeCode,omitempty" json:"employeeCode"`
	CorporateName       string    `firestore:"corporateName,omitempty" dynamodbav:"corporateName,omitempty" json:"corporateName"`
	StartDate           time.Time `firestore:"startDate,omitempty" dynamodbav:"startDate,omitempty" json:"startDate"`
	ExpiredMembership   time.Time `firestore:"expiredMembership,omitempty" dynamodbav:"expiredMembership,omitempty" json:"expiredMembership"`
	UserAppId           string    `firestore:"userAppId,omitempty" dynamodbav:"userAppId,omitempty" json:"userAppId"`
	Name                string    `firestore:"name,omitempty" dynamodbav:"name,omitempty" json:"name"`
	Phone               string    `firestore:"phone,omitempty" dynamodbav:"phone,omitempty" json:"phone"`
	IsCheckin           bool      `firestore:"isCheckin,omitempty" dynamodbav:"isCheckin,omitempty" json:"isCheckin"`
	LastDateCheckin     time.Time `firestore:"lastDateCheckin,omitempty" dynamodbav:"lastDateCheckin,omitempty" json:"lastDateCheckin"`
	LastLocationCheckin string    `firestore:"lastLocationCheckin,omitempty" dynamodbav:"lastLocationCheckin,omitempty" json:"lastLocationCheckin"`
	EndFreeze           time.Time `firestore:"endFreeze,omitempty" dynamodbav:"endFreeze,omitempty" json:"endFreeze"`
	StartFreeze         time.Time `firestore:"startFreeze,omitempty" dynamodbav:"startFreeze,omitempty" json:"startFreeze"`
	Email               string    `firestore:"email,omitempty" dynamodbav:"email,omitempty" json:"email"`
	NameLower           string    `firestore:"nameLower,omitempty" dynamodbav:"nameLower,omitempty" json:"nameLower"`
	IsEmailVerify       bool      `firestore:"isEmailVerify,omitempty" dynamodbav:"isEmailVerify,omitempty" json:"isEmailVerify"`
	ThumbLevel          int32     `firestore:"thumbLevel,omitempty" dynamodbav:"thumbLevel,omitempty" json:"thumbLevel"`
	LastDeviceId        string    `firestore:"lastDeviceId,omitempty" dynamodbav:"lastDeviceId,omitempty" json:"lastDeviceId"`
	CreatedDate         time.Time `firestore:"createdDate,omitempty" dynamodbav:"createdDate,omitempty" json:"createdDate"`
	LastLogin           time.Time `firestore:"lastLogin,omitempty" dynamodbav:"lastLogin,omitempty" json:"lastLogin"`
	PhotoUrl            string    `firestore:"photoUrl,omitempty" dynamodbav:"photoUrl,omitempty" json:"photoUrl"`
	// TODO PhotoIdentify can be list or string
	PhotoIdentify       interface{} `firestore:"photoIdentify,omitempty" dynamodbav:"photoIdentify,omitempty" json:"photoIdentify"`
	IsDeleted           bool        `firestore:"isDeleted,omitempty" dynamodbav:"isDeleted,omitempty" json:"isDeleted"`
	PromoLocation       string      `firestore:"promoLocation,omitempty" dynamodbav:"promoLocation,omitempty" json:"promoLocation"`
	DateOfBirth         time.Time   `firestore:"dateOfBirth,omitempty" dynamodbav:"dateOfBirth,omitempty" json:"dateOfBirth"`
	Remarks             string      `firestore:"remarks,omitempty" dynamodbav:"remarks,omitempty" json:"remarks"`
	EmergencyAddress    string      `firestore:"emergencyAddress,omitempty" dynamodbav:"emergencyAddress,omitempty" json:"emergencyAddress"`
	EmergencyName       string      `firestore:"emergencyName,omitempty" dynamodbav:"emergencyName,omitempty" json:"emergencyName"`
	EmergencyPhone      string      `firestore:"emergencyPhone,omitempty" dynamodbav:"emergencyPhone,omitempty" json:"emergencyPhone"`
	Height              float32     `firestore:"height,omitempty" dynamodbav:"height,omitempty" json:"height"`
	Weight              float32     `firestore:"weight,omitempty" dynamodbav:"weight,omitempty" json:"weight"`
	FCMToken            string      `firestore:"fcmToken,omitempty" dynamodbav:"fcmToken,omitempty" json:"fcmToken"`
	RefCode             string      `firestore:"refCode,omitempty" dynamodbav:"refCode,omitempty" json:"refCode"`
	CRMId               string      `firestore:"crmId,omitempty" dynamodbav:"crmId,omitempty" json:"crmId"`
	RegUserAppId        string      `firestore:"regUserAppId,omitempty" dynamodbav:"regUserAppId,omitempty" json:"regUserAppId"`
	Gender              string      `firestore:"gender,omitempty" dynamodbav:"gender,omitempty" json:"gender"`
	Source              string      `firestore:"source,omitempty" dynamodbav:"source,omitempty" json:"source"`
	TransId             string      `firestore:"transId,omitempty" dynamodbav:"transId,omitempty" json:"transId"`
	TransDate           string      `firestore:"transDate,omitempty" dynamodbav:"transDate,omitempty" json:"transDate"`
	IsCorporateVerified bool        `firestore:"isCorporateVerified,omitempty" dynamodbav:"isCorporateVerified,omitempty" json:"isCorporateVerified"`
	OldNumber           string      `firestore:"oldNumber,omitempty" dynamodbav:"oldNumber,omitempty" json:"oldNumber"`
	Partners            string      `firestore:"partners,omitempty" dynamodbav:"partners,omitempty" json:"partners"`
	TypeUser            string      `firestore:"typeUser,omitempty" dynamodbav:"typeUser,omitempty" json:"typeUser"`
	UserOS              string      `firestore:"userOS,omitempty" dynamodbav:"userOS,omitempty" json:"userOS"`
	LastActive          string      `firestore:"lastActive,omitempty" dynamodbav:"lastActive,omitempty" json:"lastActive"`
	PackageName         string      `firestore:"packageName,omitempty" dynamodbav:"packageName,omitempty" json:"packageName"`
	// TODO: Firestore packagePrice data value exists string and number
	PackagePrice                  interface{}            `firestore:"packagePrice,omitempty" dynamodbav:"packagePrice,omitempty" json:"packagePrice"`
	LastSelectedLocation          string                 `firestore:"lastSelectedLocation,omitempty" dynamodbav:"lastSelectedLocation,omitempty" json:"lastSelectedLocation"`
	LastSelectedLocationOccupancy map[string]interface{} `json:"lastSelectedLocationOccupancy"`
	// TODO: Firestore salesName data value exists string and bool
	SalesName            interface{} `firestore:"salesName,omitempty" dynamodbav:"salesName,omitempty" json:"salesName"`
	Provider             string      `firestore:"provider,omitempty" dynamodbav:"provider,omitempty" json:"provider"`
	BlockInfo            string      `firestore:"blockInfo,omitempty" dynamodbav:"blockInfo,omitempty" json:"blockInfo"`
	BlockClasses         time.Time   `firestore:"blockClasses,omitempty" dynamodbav:"blockClasses,omitempty" json:"blockClasses"`
	ExpiredNotes         string      `firestore:"expiredNotes,omitempty" dynamodbav:"expiredNotes,omitempty" json:"expiredNotes"`
	ExpiredUpdate        time.Time   `firestore:"expiredUpdate,omitempty" dynamodbav:"expiredUpdate,omitempty" json:"expiredUpdate"`
	DataSyncCS           bool        `firestore:"dataSyncCS" dynamodbav:"dataSyncCS,omitempty" json:"dataSyncCS"`
	AgreementVersion     string      `firestore:"agreementVersion,omitempty" dynamodbav:"agreementVersion,omitempty" json:"agreementVersion"`
	RemarksCheckin       string      `firestore:"remarksCheckin,omitempty" json:"remarksCheckin"`
	LockerCheckin        string      `firestore:"lockerCheckin,omitempty" json:"lockerCheckin"`
	MembershipStatus     int         `firestore:"membershipStatus,omitempty" dynamodbav:"membershipStatus,omitempty" json:"membershipStatus"`
	IsAppReviewAvailable bool        `firestore:"isAppReviewAvailable,omitempty" json:"isAppReviewAvailable,omitempty"`
	DeviceType           string      `firestore:"deviceType,omitempty" json:"deviceType,omitempty"`
	FreeTrialID          string      `firestore:"freeTrialId,omitempty" json:"freeTrialId,omitempty"`
	PTFollowOnEligible   bool        `firestore:"PTFollowOnEligible,omitempty" json:"ptFollowOnEligible,omitempty"`
	PTFollowOnPeriod     int         `firestore:"PTFollowOnPeriod,omitempty" json:"ptFollowOnPeriod,omitempty"`
}

type UserLogsMembership struct {
	ExecutionId  string    `firestore:"executionId,omitempty" json:"executionId"`
	DateTime     time.Time `firestore:"dateTime,omitempty" json:"dateTime"`
	NewDate      time.Time `firestore:"newDate,omitempty" json:"newDate"`
	PreviousDate time.Time `firestore:"previousDate,omitempty" json:"previousDate"`
	TypeChange   string    `firestore:"typeChange,omitempty" json:"typeChange"`
}

type UserPTPackage struct {
	ID              string
	Action          string    `firestore:"action" json:"action"`
	CreatedDate     time.Time `firestore:"createdDate" json:"createdDate"`
	ExpiredDate     time.Time `firestore:"expiredDate" json:"expiredDate"`
	Image           string    `firestore:"image" json:"image"`
	IsActive        bool      `firestore:"isActive" json:"isActive"`
	IsFreeze        bool      `firestore:"isFreeze" json:"isFreeze"`
	IsHide          bool      `firestore:"isHide" json:"isHide"`
	IsPromo         bool      `firestore:"isPromo" json:"isPromo"`
	IsLockedUntil   time.Time `firestore:"isLockedUntil" json:"isLockedUntil"`
	Location        string    `firestore:"location" json:"location"`
	Muscle          []string  `firestore:"muscle" json:"muscle"`
	Oid             string    `firestore:"oid" json:"oid"`
	PackageName     string    `firestore:"packageName" json:"packageName"`
	PrevExpiredDate time.Time `firestore:"prevExpiredDate" json:"prevExpiredDate"`
	Session         []string  `firestore:"session" json:"session"`
	SessionLeft     int       `firestore:"sessionLeft" json:"sessionLeft"`
	SessionTotal    int       `firestore:"sessionTotal" json:"sessionTotal"`
	StartDate       time.Time `firestore:"startDate" json:"startDate"`
	Title           string    `firestore:"title" json:"title"`
	Type            string    `firestore:"type" json:"type"`
}
