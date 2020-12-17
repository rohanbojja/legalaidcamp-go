package repositories

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/rohanbojja/legalaidcamp-go/entity"
	"log"
)

type UserRepository interface {
	Save(user entity.User) (entity.User,error)
	FindByID(UID string) (entity.User,error)
}

type userRepository struct{
	firebaseApp *firebase.App
}
func NewUserRepository(app *firebase.App) UserRepository {
	return &userRepository{firebaseApp: app}
}

func (u userRepository) Save(user entity.User) (entity.User, error) {
	ctx := context.Background()
	//Validate Entity here?
	store, err := u.firebaseApp.Firestore(ctx)
	if err != nil {
		panic("Firestore instance couldn't be created.")
		return entity.User{}, err
	}
	if _, err := store.Collection("users").Doc(user.UID).Set(ctx, user); err != nil {
		return entity.User{}, err
	}
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (u userRepository) FindByID(UID string) (entity.User, error) {
	ctx := context.Background()
	var user entity.User
	store, err := u.firebaseApp.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
		return entity.User{}, err
	}
	doc, err := store.Collection("users").Doc(UID).Get(ctx)
	if err != nil {
		return entity.User{}, err
	}
	if err := doc.DataTo(&user); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

