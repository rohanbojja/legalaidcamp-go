package entity

//User is Firebase entity
type User struct {
	UID          string `json:"uid"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
	PhoneNumber  string `json:"phoneNumber"`
	PhotoURL     string `json:"photoURL"`
	Cases        []string
	OldCases     []string
}
