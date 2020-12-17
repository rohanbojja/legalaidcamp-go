package repositories

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/rohanbojja/legalaidcamp-go/entity"
	"log"
)

type CourtCaseRepository interface {
	Save(courtCase entity.CourtCase) (entity.CourtCase,error)
	FindByID(courtCaseID string) (entity.CourtCase,error)
	FindAll() ([]entity.CourtCase,error)
}

type courtCaseRepository struct {
	firebaseApp *firebase.App
}

func NewCourtCaseRepository(app *firebase.App) CourtCaseRepository {
	return &courtCaseRepository{firebaseApp: app}
}

func (c *courtCaseRepository) Save(courtCase entity.CourtCase) (entity.CourtCase, error) {
	ctx := context.Background()
	store, err := c.firebaseApp.Firestore(ctx)
	if err != nil {
		return courtCase, err
	}
	if _, err := store.Collection("cases").Doc(courtCase.CourtCaseId).Set(ctx, courtCase); err != nil {
		return entity.CourtCase{}, err
	}
	log.Printf("Court case ID: %s\n", courtCase.CourtCaseId)
	return courtCase, nil
}

func (c *courtCaseRepository) FindByID(courtCaseID string) (entity.CourtCase, error) {
	ctx := context.Background()
	var courtCase entity.CourtCase
	store, err := c.firebaseApp.Firestore(ctx)
	if err != nil {
		return courtCase, err
	}
	courtCaseDoc, err := store.Collection("cases").Doc(courtCaseID).Get(ctx)
	if err != nil {
		return courtCase, err
	}
	if err := courtCaseDoc.DataTo(&courtCase); err != nil {
		return entity.CourtCase{}, err
	}
	return courtCase, nil
}

func (c *courtCaseRepository) FindAll() ([]entity.CourtCase, error) {
	ctx := context.Background()
	var allcases []entity.CourtCase
	store, err := c.firebaseApp.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	collectionRef := store.Collection("cases").DocumentRefs(ctx)
	docs, _ := collectionRef.GetAll()
	for _, e := range docs {
		courtCase := &entity.CourtCase{}
		snap, _ := e.Get(ctx)
		if err := snap.DataTo(&courtCase); err != nil {
			return nil, err
		}
		allcases = append(allcases, *courtCase)
		fmt.Println(courtCase.AreaOfLaw + courtCase.State)
	}
	return allcases, nil
}
