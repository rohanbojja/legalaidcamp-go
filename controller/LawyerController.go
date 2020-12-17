package controller

import (
	"fmt"
	"github.com/rohanbojja/legalaidcamp-go/repositories"
	"log"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/rohanbojja/legalaidcamp-go/entity"
	"github.com/rohanbojja/legalaidcamp-go/service"
)

type LawyerController interface {
	CreateLawyer(ctx *fiber.Ctx) error
	ToggleStatus(ctx *fiber.Ctx) error
	AcceptCase(ctx *fiber.Ctx) error
	RejectCase(ctx *fiber.Ctx) error
	RemoveCase(ctx *fiber.Ctx) error
	GetAssignedCases(ctx *fiber.Ctx) error
}

type lawyerController struct {
	authService         service.AuthService
	courtcaseRepository repositories.CourtCaseRepository
	lawyerRepository    repositories.LawyerRepository
	lawyerService       service.LawyerService
}

func NewLawyerController(authService service.AuthService,
	courtcaseRepository repositories.CourtCaseRepository,
	lawyerRepository repositories.LawyerRepository,
	lawyerService service.LawyerService, ) LawyerController {

	return &lawyerController{
		authService:         authService,
		courtcaseRepository: courtcaseRepository,
		lawyerRepository:    lawyerRepository,
		lawyerService:       lawyerService,
	}
}

func (c *lawyerController) CreateLawyer(ctx *fiber.Ctx) error {
	header := ctx.Get("idtoken")
	lawyer := entity.Lawyer{}
	uid, err := c.authService.GetUID(header)
	if err != nil {
		return err
	}
	//Retrieve uid from token
	phoneNumber, err := c.authService.GetPhoneNumber(uid)
	if err != nil {
		return err
	}
	if ctx.BodyParser(&lawyer) != nil {
		return err
	}
	lawyer.UserID = uid
	lawyer.PhoneNumber = phoneNumber
	persistLawyer, err := c.lawyerService.CreateLawyer(lawyer)
	if err != nil {
		return err
	}
	return ctx.Status(201).JSON(persistLawyer)
}

func (c *lawyerController) ToggleStatus(ctx *fiber.Ctx) error {
	fmt.Println("Lawyer controller.")
	header := ctx.Get("idtoken")
	uid, err := c.authService.GetUID(header)
	if err != nil {
		return err
	}
	lawyer, err := c.lawyerRepository.FindByID(uid)
	if err != nil {
		return err
	}
	log.Printf("uid: %s, lawyer: %v", uid, lawyer)
	status, err := c.lawyerService.ToggleStatus(uid, !lawyer.ProfileStatus)
	if err != nil {
		return err
	}
	return ctx.Status(200).JSON(status)
}
func (c *lawyerController) AcceptCase(ctx *fiber.Ctx) error {
	header := ctx.Get("idtoken")
	cid := ctx.Get("cid")
	uid, err := c.authService.GetUID(header)
	if err != nil {
		return err
	}
	caseStatus, err := c.lawyerService.HandleCase(uid, cid, true)
	if err != nil {
		return err
	}
	return ctx.Status(200).JSON(caseStatus)
}

func (c *lawyerController) RejectCase(ctx *fiber.Ctx) error {
	header := ctx.Get("idtoken")
	cid := ctx.Get("cid")
	uid, err := c.authService.GetUID(header)
	if err != nil {
		return err
	}
	casestatus, err := c.lawyerService.HandleCase(uid, cid, false)
	if err != nil {
		return err
	}
	return ctx.Status(200).JSON(casestatus)
}

func (c *lawyerController) RemoveCase(ctx *fiber.Ctx) error {
	header := ctx.Get("idtoken")
	cid := ctx.Get("cid")
	uid, err := c.authService.GetUID(header)
	if err != nil {
		return err
	}
	status, err := c.lawyerService.RemoveCase(uid, cid)
	if err != nil {
		return err
	}
	return ctx.Status(200).JSON(status)
}

func (c *lawyerController) GetAssignedCases(ctx *fiber.Ctx) error {
	var assignedCases []entity.CourtCase
	header := ctx.Get("idtoken")
	uid, err := c.authService.GetUID(header)
	if err != nil {
		return err
	}
	lawyer, _ := c.lawyerRepository.FindByID(uid)
	for _, e := range lawyer.AssignedCases {
		caseDetails, _ := c.courtcaseRepository.FindByID(e)
		assignedCases = append(assignedCases, caseDetails)
	}
	return ctx.Status(200).JSON(assignedCases)
}
