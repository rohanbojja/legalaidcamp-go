package service

import (
	"github.com/rohanbojja/legalaidcamp-go/repositories"
	"log"

	"github.com/rohanbojja/legalaidcamp-go/entity"
)

// LawyerService contains functions that pertain to actions relating to the entity "Lawyer"
type LawyerService interface {
	CreateLawyer(lawyer entity.Lawyer) (entity.Lawyer, error)
	UpdateProfile(lawyer entity.Lawyer) (entity.Lawyer, error)
	HandleCase(lawyerID string, courtCaseID string, action bool) (bool, error)
	RemoveCase(lawyerID string, courtCaseID string) (bool, error)
	ToggleStatus(lawyerID string, status bool) (bool, error)
	//TODO NOT IMPL
	UploadLicense()
}

type lawyerService struct {
	lawyerRepository    repositories.LawyerRepository
	courtCaseRepository repositories.CourtCaseRepository
}

// NewLawyerService returns a LawyerService with the injected CommonService
func NewLawyerService(lawyerRepository repositories.LawyerRepository) LawyerService {
	service := lawyerService{
		lawyerRepository: lawyerRepository,
	}
	return &service
}

//Create and activate the lawyer profile ( Not verify)
/*
TODO
	Upload photocopy of proof of practice?
	Resfresh list and remove deprecated cases

*/
func (service *lawyerService) CreateLawyer(lawyer entity.Lawyer) (entity.Lawyer, error) {
	lawyer.Verified = false
	lawyer.ProfileStatus = true
	log.Printf("Creating lawyer: %+v\n", lawyer)
	if _, err := service.lawyerRepository.Save(lawyer); err != nil {
		return entity.Lawyer{}, err
	}
	return lawyer, nil
}

func (service *lawyerService) UpdateProfile(lawyer entity.Lawyer) (entity.Lawyer, error) {
	lawyer2, err := service.lawyerRepository.FindByID(lawyer.UserID)
	if err != nil {
		return entity.Lawyer{}, err
	}

	//Hack to not have these value change
	lawyer.Verified = lawyer2.Verified

	if _, err := service.lawyerRepository.Save(lawyer); err != nil {
		return entity.Lawyer{}, err
	}
	return lawyer, nil
}

func cut(i int, xs []string) (string, []string) {
	y := xs[i]
	ys := append(xs[:i], xs[i+1:]...)
	return y, ys
}

func (service *lawyerService) RemoveCase(lawyerID string, courtCaseID string) (bool, error) {
	lawyer, err := service.lawyerRepository.FindByID(lawyerID)
	if err != nil {
		return false, err
	}
	success := false
	for i, e := range lawyer.ActiveCases {
		if e == courtCaseID {
			courtCase, err := service.courtCaseRepository.FindByID(e)
			if err != nil {
				log.Printf("Courtcase not found\n")
				return false, err
			}
			_, lawyer.ActiveCases = cut(i, lawyer.ActiveCases)
			courtCase.Status = "active"
			courtCase.AssignedLawyerID = ""
			lawyer.RejectedCases = append(lawyer.RejectedCases, e)
			if _, err := service.courtCaseRepository.Save(courtCase); err != nil {
				return false, err
			}
			if _, err := service.lawyerRepository.Save(lawyer); err != nil {
				return false, err
			}
			success = true
			break
		}
	}
	return success, nil
}

func (service *lawyerService) HandleCase(lawyerID, courtCaseID string, action bool) (bool, error) {
	//Firestore transaction
	lawyer, err := service.lawyerRepository.FindByID(lawyerID)
	if err != nil {
		panic(err)
		return false, err
	}
	success := false
	for i, e := range lawyer.AssignedCases {
		if e == courtCaseID {
			courtCase, err := service.courtCaseRepository.FindByID(e)
			if err != nil {
				log.Printf("Courtcase not found\n")
				return false, nil
			}
			_, lawyer.AssignedCases = cut(i, lawyer.AssignedCases)
			if action {
				if courtCase.Status == "active" {
					courtCase.Status = "assigned"
					courtCase.AssignedLawyerID = lawyerID
					lawyer.ActiveCases = append(lawyer.ActiveCases, e)
					if _, err := service.courtCaseRepository.Save(courtCase); err != nil {
						return false, err
					}
				}

			} else {
				lawyer.RejectedCases = append(lawyer.RejectedCases, e)
			}
			success = true
			if _, err := service.lawyerRepository.Save(lawyer); err != nil {
				return false, err
			}
			break
		}
	}
	if success {
		return success, nil
	}
	return success, nil
}

func (service *lawyerService) ToggleStatus(lawyerID string, status bool) (bool, error) {
	lawyer, err := service.lawyerRepository.FindByID(lawyerID)
	if err != nil {
		return false, err
	}
	log.Printf("current lawyer status:%v and status: %v", lawyer.ProfileStatus, status)
	lawyer.ProfileStatus = status
	if _, err := service.lawyerRepository.Save(lawyer); err != nil {
		return false, err
	}
	return lawyer.ProfileStatus, nil
}

func (service *lawyerService) UploadLicense() {
	panic("implement me")
}
