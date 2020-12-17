//+build wireinject

package main

import (
	firebase "firebase.google.com/go/v4"
	"github.com/google/wire"
	"github.com/rohanbojja/legalaidcamp-go/controller"
	"github.com/rohanbojja/legalaidcamp-go/repositories"
	"github.com/rohanbojja/legalaidcamp-go/service"
)

//adminService        service.AdminService
//	authService         service.AuthService
//	courtcaseRepository repositories.CourtCaseRepository
//	lawyerRepository    repositories.LawyerRepository

var conts = wire.NewSet(controller.NewAdminController, controller.NewLawyerController, controller.NewUserController)
var repos = wire.NewSet(repositories.NewCourtCaseRepository, repositories.NewLawyerRepository, repositories.NewUserRepository)
var servs  = wire.NewSet(service.NewAdminService, service.NewAuthService, service.NewLawyerService, service.NewUserService)

func InitializeAdminController(app *firebase.App) controller.AdminController {
	wire.Build(conts, servs ,repos)
	return nil
}

//authService         service.AuthService
//	courtcaseRepository repositories.CourtCaseRepository
//	lawyerRepository    repositories.LawyerRepository
//	lawyerService       service.LawyerService
func InitializeLawyerController(app *firebase.App) controller.LawyerController {
	wire.Build(conts, servs ,repos)
	return nil
}

//userService service.UserService
//	authService service.AuthService
func InitializeUserController(app *firebase.App) controller.UserController {
	wire.Build(conts, servs ,repos)
	return nil
}
