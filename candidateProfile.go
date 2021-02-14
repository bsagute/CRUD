package model

//For Profile Edit
type EmployeeProfile struct {
	FullName  string `json:"fullName" bson:"fullName"`
	LoginId   string `json:"loginId" bson:"loginId"`
	IsEnabled bool   `json:"isEnabled" bson:"isEnabled"`
	IsDeleted bool   `json:"isDeleted" bson:"isDeleted"`
}
