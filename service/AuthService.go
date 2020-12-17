package service

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"log"
	"math/rand"
	"os"
	"strconv"
)

// AuthService interface contains common operations, this is injected into CommonService.
type AuthService interface {
	GetUID(idToken string) (string, error)
	HasRole(idToken, role string) (bool, error)
	GetPhoneNumber(idToken string) (string, error)
}

type authService struct {
	firebaseService *firebase.App
}

// NewAuthService return an AuthService with the global instance of Firebase injected
func NewAuthService(firebaseApp *firebase.App) AuthService {
	return &authService{
		firebaseService: firebaseApp,
	}
}

func (service *authService) GetUID(idToken string) (string, error) {
	if os.Getenv("SRV_DEV") == "true" {
		return idToken[0:4], nil
	}
	ctx := context.Background()
	auth, err := service.firebaseService.Auth(ctx)
	if err != nil {
		return "", err
	}
	token, err := auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", err
	}
	log.Println("Retrieved UID from ID token.")
	return token.UID, nil
}

/*
todo: Not impl
*/
func (service *authService) HasRole(idToken, role string) (bool, error) {
	if os.Getenv("SRV_DEV") == "true" {
		return true, nil
	}
	ctx := context.Background()
	//initialize auth
	auth, err := service.firebaseService.Auth(ctx)
	if err != nil {
		return false, err
	}
	//Verify token
	token, err := auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return false, err
	}

	claims := token.Claims
	if admin, ok := claims[role]; ok {
		if admin.(bool) {
			return true, nil
		}
	}
	return false, nil
}

func (service *authService) GetPhoneNumber(uid string) (string, error) {
	if os.Getenv("SRV_DEV") == "true" {
		return strconv.Itoa(rand.Int()), nil
	}
	ctx := context.Background()
	auth, err := service.firebaseService.Auth(ctx)
	if err != nil {
		return "", err
	}
	user, err := auth.GetUser(ctx, uid)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
		return "", err
	}
	log.Println("Retrieved phone number from UID.")
	return user.PhoneNumber, nil
}
