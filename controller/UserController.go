package controller

import (
	"fmt"
	"log"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/rohanbojja/legalaidcamp-go/entity"
	"github.com/rohanbojja/legalaidcamp-go/service"
)

//UserController is..
type UserController interface {
	CreateCase(ctx *fiber.Ctx) error
}

type userController struct {
	userService service.UserService
	authService service.AuthService
}

//New is the constructor for this controller
func NewUserController(userService service.UserService,
	authService service.AuthService, ) UserController {
	return &userController{
		userService: userService,
		authService: authService,
	}
}

func (c *userController) CreateCase(ctx *fiber.Ctx) error {
	header := ctx.Get("idtoken")
	log.Printf("Header" + header)
	courtCase := entity.CourtCase{}

	uid, err := c.authService.GetUID(header)
	if err != nil {
		return ctx.Status(400).JSON(err)
	}
	phoneNumber, err := c.authService.GetPhoneNumber(header)
	if err != nil {
		return ctx.Status(400).JSON(err)
	}

	if err := ctx.BodyParser(&courtCase); err != nil {
		return ctx.Status(400).JSON(err)
	}
	fmt.Printf("Create case: %+v\n", courtCase)
	courtCase.UserID = uid
	courtCase.PhoneNumber = phoneNumber
	createCase, err := c.userService.CreateCase(courtCase)
	if err != nil {
		log.Panicf("Some error while creating the case. %v", err)
		return ctx.Status(500).JSON(err.Error())
	}
	return ctx.Status(201).JSON(createCase)
}
