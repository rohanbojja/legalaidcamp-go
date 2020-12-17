package entity

type Lawyer struct {
	UserID          string
	PhoneNumber     string
	Verified        bool
	AssignedCases   []string
	ActiveCases     []string
	RejectedCases   []string
	OldCases        []string
	Languages       []int  `json:"languages"`
	AreasOfLaw      []int  `json:"areasOfLaw"`
	StateOfPractice int    `json:"stateOfPractice"`
	BarCouncil      int    `json:"barCouncil"`
	DisplayName     string `json:"displayName"`
	EmailAddress    string `json:"emailAddress"`
	City            string `json:"city"`
	AllowVisits     bool   `json:"allowVisits"`
	AllowCalls      bool   `json:"allowCalls"`
	ProfileStatus   bool   `json:"profileStatus"`
	OfficeAddress   string `json:"officeAddress"`
	OfficePincode   string `json:"officePincode"`
}
