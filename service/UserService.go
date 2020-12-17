package service

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/rohanbojja/legalaidcamp-go/repositories"
	"log"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/rohanbojja/legalaidcamp-go/entity"
)

// UserService interface contains methods dealing with user actions
type UserService interface {
	CreateCase(courtCase entity.CourtCase) (entity.CourtCase, error)
	AssignLawyers(courtCase entity.CourtCase) (int, error) // Return the number of lawyers assigned
	MarkAsCompleted(courtCaseID string) (entity.CourtCase, error)
}

type userService struct {
	courtcaseRepository repositories.CourtCaseRepository
	lawyerRepository    repositories.LawyerRepository
	userRepository      repositories.UserRepository
	firebaseApp         *firebase.App
}

// NewUserService returns UserService with the injected CommonService
func NewUserService(courtcaseRepository repositories.CourtCaseRepository,
	lawyerRepository repositories.LawyerRepository,
	userRepository repositories.UserRepository,
	firebaseApp *firebase.App) UserService {
	return &userService{
		courtcaseRepository: courtcaseRepository,
		lawyerRepository:    lawyerRepository,
		userRepository:      userRepository,
		firebaseApp:         firebaseApp,
	}
}

// Create checks if the user already exists,if not creates a entry for the firebase user

/*
Creates an entry to persist the CourtCase entity passed.
Works, but there's room for improvement!
TODO:
	Limit the number of unsolved cases to 3?

*/

func (service *userService) CreateCase(courtCase entity.CourtCase) (entity.CourtCase, error) {
	ctx := context.Background()

	fmt.Printf("Creating case: %+v\n", courtCase)

	user, err := service.userRepository.FindByID(courtCase.UserID)
	if err != nil {
		user = entity.User{
			UID:         courtCase.UserID,
			DisplayName: courtCase.DisplayName,
			PhoneNumber: courtCase.PhoneNumber,
		}
	}

	store, err := service.firebaseApp.Firestore(ctx)
	if err != nil {
		return entity.CourtCase{}, nil
	}

	courtCaseID := store.Collection("cases").NewDoc().ID
	courtCase.CourtCaseId = courtCaseID
	courtCase.Status = "active"
	persistCourtCase, err := service.courtcaseRepository.Save(courtCase)
	if err != nil {
		return entity.CourtCase{}, err
	}
	lawyersAssigned, err := service.AssignLawyers(courtCase)
	if err != nil {
		return entity.CourtCase{}, err
	}
	user.Cases = append(user.Cases, courtCaseID)
	if _, err := store.Collection("users").Doc(user.UID).Set(ctx, user); err != nil {
		return entity.CourtCase{}, err
	}
	log.Println("Created the case: " + courtCase.CourtCaseId + " with " + strconv.Itoa(lawyersAssigned) + " lawyers")
	return persistCourtCase, nil
}

/*
Assigns Lawyers to the CourtCase entity passed
3/5
TODO:
- Can create custom tree like structure to make queries more efficient, but this'll do for now.
- Limiting query to 25 and hitting up 15 lawyers in the end
*/
func (service *userService) AssignLawyers(courtCase entity.CourtCase) (int, error) {

	if courtCase.Status == "active" {
		// language / state / area of law
		query := fmt.Sprintf("%d_%d_%d", courtCase.Language, courtCase.State, courtCase.AreaOfLaw)
		fmt.Printf("Gathering lawyers: %v\n", query)

	}
	lawyerPrioritySlice, err := service.lawyerRepository.FindForCourtCase(courtCase)
	if err != nil {
		return 0, err
	}
	if len(lawyerPrioritySlice) < 5 {
		go func() {
			ctx := context.Background()
			store, err := service.firebaseApp.Firestore(ctx)
			if err != nil {
				panic("gg")
			}
			_, _ = store.Collection("anomalies").Doc(courtCase.CourtCaseId).Set(ctx, map[string]interface{}{
				"UserID":      courtCase.UserID,
				"CaseID":      courtCase.CourtCaseId,
				"PhoneNumber": courtCase.PhoneNumber,
				"Timestamp":   firestore.ServerTimestamp,
			})
		}()
	}
	for i, e := range lawyerPrioritySlice {
		if i >= 15 {
			break
		}
		e.AssignedCases = append(e.AssignedCases, courtCase.CourtCaseId)
		if _, err := service.lawyerRepository.Save(e); err != nil {
			return 0, err
		}
	}
	return len(lawyerPrioritySlice), nil
}

//TODO
//UNTESTED FUNCTION
func (service *userService) MarkAsCompleted(courtCaseID string) (entity.CourtCase, error) {
	courtCase, err := service.courtcaseRepository.FindByID(courtCaseID)
	if err != nil {
		return entity.CourtCase{}, err
	}
	courtCase.Status = "retired"
	if _, err := service.courtcaseRepository.Save(courtCase); err != nil {
		return entity.CourtCase{}, err
	}
	lawyer, err := service.lawyerRepository.FindByID(courtCase.AssignedLawyerID)
	if err != nil {
		return entity.CourtCase{}, err
	}
	for i, e := range lawyer.ActiveCases {
		if e == courtCaseID {
			_, lawyer.ActiveCases = cut(i, lawyer.ActiveCases)
			lawyer.OldCases = append(lawyer.OldCases, e)
			break
		}
	}
	if _, err := service.lawyerRepository.Save(lawyer); err != nil {
		return entity.CourtCase{}, err
	}
	return courtCase, nil
}
