package repositories

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/rohanbojja/legalaidcamp-go/entity"
	"log"
)

type LawyerRepository interface {
	Save(lawyer entity.Lawyer) (entity.Lawyer, error)
	FindByID(lawyerID string) (entity.Lawyer, error)
	FindForCourtCase(courtCase entity.CourtCase) ([]entity.Lawyer, error)
	FindAll() ([]entity.Lawyer, error)
}

type lawyerRepository struct {
	firebaseApp *firebase.App
}

func NewLawyerRepository(app *firebase.App) LawyerRepository {
	return &lawyerRepository{firebaseApp: app}
}

func (l *lawyerRepository) FindAll() ([]entity.Lawyer, error) {
	ctx := context.Background()
	var allLawyers []entity.Lawyer
	store, err := l.firebaseApp.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	collectionRef := store.Collection("lawyers").DocumentRefs(ctx)
	docs, _ := collectionRef.GetAll()
	for _, e := range docs {
		lawyer := &entity.Lawyer{}
		snap, _ := e.Get(ctx)
		if err := snap.DataTo(&lawyer); err != nil {
			return nil, err
		}
		allLawyers = append(allLawyers, *lawyer)
		log.Printf("uid: %v\n", lawyer.UserID)
	}
	return allLawyers, nil
}

func (l *lawyerRepository) Save(lawyer entity.Lawyer) (entity.Lawyer, error) {
	ctx := context.Background()
	//Validate Entity here?
	store, err := l.firebaseApp.Firestore(ctx)
	if err != nil {
		log.Fatalln("Firestore instance couldn't be created.")
		return lawyer, err
	}
	if _, err := store.Collection("lawyers").Doc(lawyer.UserID).Set(ctx, lawyer); err != nil {
		return entity.Lawyer{}, err
	}
	if err != nil {
		return lawyer, err
	}
	return lawyer, nil
}

func (l *lawyerRepository) FindByID(lawyerID string) (entity.Lawyer, error) {
	ctx := context.Background()
	var lawyer entity.Lawyer
	store, err := l.firebaseApp.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
		return lawyer, err
	}
	doc, err := store.Collection("lawyers").Doc(lawyerID).Get(ctx)
	if err != nil {
		return lawyer, err
	}
	if err := doc.DataTo(&lawyer); err != nil {
		return entity.Lawyer{}, err
	}
	return lawyer, nil
}

func (l *lawyerRepository) FindForCourtCase(courtCase entity.CourtCase) ([]entity.Lawyer, error) {
	ctx := context.Background()
	var assignedLawyers []entity.Lawyer
	var assignedLawyersSameCity []entity.Lawyer
	store, err := l.firebaseApp.Firestore(ctx)
	if err != nil {
		log.Fatalln("Firestore instance couldn't be created.")
		return nil, err
	}
	eligibleLawyers, err := store.Collection("lawyers").Where("StateOfPractice", "==", courtCase.State).Where("AreasOfLaw", "array-contains", courtCase.AreaOfLaw).Limit(25).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	for _, s := range eligibleLawyers {
		lawyer := entity.Lawyer{}
		err := s.DataTo(&lawyer)
		if err != nil {
			return nil, err
		}
		for _, lang := range lawyer.Languages {
			if lang == courtCase.Language {
				if lawyer.City == courtCase.City {
					assignedLawyersSameCity = append(assignedLawyersSameCity, lawyer)
				} else {
					assignedLawyers = append(assignedLawyers, lawyer)
				}
				break
			}
		}
	}

	var lawyerPrioritySlice []entity.Lawyer
	lawyerPrioritySlice = append(lawyerPrioritySlice, assignedLawyersSameCity...)
	lawyerPrioritySlice = append(lawyerPrioritySlice, assignedLawyers...)
	return lawyerPrioritySlice, nil
}
