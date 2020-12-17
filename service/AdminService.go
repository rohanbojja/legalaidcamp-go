package service

import "github.com/rohanbojja/legalaidcamp-go/repositories"

// AdminService interface for all the admin/privileged functions to be implemented
/*
TODO:
- Make repositories and inject with Wire?

 */
type AdminService interface {
	ToggleVerification(lawyerID string) (bool, error)
	// TODO NOT IMPL
	SendText(phoneNumber string, message string) (bool, error)
}

type adminService struct {
	lawyerRepository repositories.LawyerRepository
	courtCaseRepository repositories.CourtCaseRepository
}

// NewAdminService function returns an AdminService with the injected CommonService
func NewAdminService(lawyerRepository repositories.LawyerRepository,courtCaseRepository repositories.CourtCaseRepository) AdminService {
	return &adminService{
		lawyerRepository: lawyerRepository,
		courtCaseRepository: courtCaseRepository,
	}
}

func (service *adminService) ToggleVerification(lawyerID string) (bool, error) {
	lawyer, err := service.lawyerRepository.FindByID(lawyerID)
	if err != nil {
		return false, nil
	}
	lawyer.Verified = !lawyer.Verified
	_, err = service.lawyerRepository.Save(lawyer)
	if err != nil {
		return false, err
	}
	return lawyer.Verified, nil
}

func (service *adminService) SendText(phoneNumber string, message string) (bool, error) {
	panic("implement me")
}
