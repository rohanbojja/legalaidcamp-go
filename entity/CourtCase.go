package entity

type CourtCase struct {
	CourtCaseId      string
	UserID           string
	PhoneNumber      string
	Status           string
	AssignedLawyerID string
	State            int    `json:"state"`
	AreaOfLaw        int    `json:"areaOfLaw"`
	Language         int    `json:"language"`
	DisplayName      string `json:"displayName"`
	City             string `json:"city"`
	Description      string `json:"description"`
	Gender           string `json:"gender"`
}
