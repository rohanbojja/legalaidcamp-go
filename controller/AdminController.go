package controller

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/rohanbojja/legalaidcamp-go/repositories"
	"github.com/rohanbojja/legalaidcamp-go/service"
)

// TODO:
// - CustomRoles
// - Authentication
// - Auditing

type AdminController interface {
	GetAllCases(ctx *fiber.Ctx) error //Todo: Paginate!
	GetAllLawyers(ctx *fiber.Ctx) error
	GetLawyer(ctx *fiber.Ctx) error
	GetCourtCase(ctx *fiber.Ctx) error
	VerifyLawyer(ctx *fiber.Ctx) error
}

type adminController struct {
	adminService        service.AdminService
	authService         service.AuthService
	courtcaseRepository repositories.CourtCaseRepository
	lawyerRepository    repositories.LawyerRepository
}

func NewAdminController(adminService service.AdminService,
	authService service.AuthService,
	courtcaseRepository repositories.CourtCaseRepository,
	lawyerRepository repositories.LawyerRepository, ) AdminController {
	return &adminController{
		adminService:        adminService,
		authService:         authService,
		courtcaseRepository: courtcaseRepository,
		lawyerRepository:    lawyerRepository,
	}
}

func (c *adminController) GetAllCases(ctx *fiber.Ctx) error {
	idToken := ctx.Get("idtoken")
	if _, err := c.authService.HasRole(idToken, "admin"); err != nil {
		return err
	} else {
		allcases, err := c.courtcaseRepository.FindAll()
		if err != nil {
			return err
		}
		return ctx.Status(200).JSON(allcases)
	}
}

func (c *adminController) GetAllLawyers(ctx *fiber.Ctx) error {

	idToken := ctx.Get("idtoken")
	if _, err := c.authService.HasRole(idToken, "admin"); err != nil {
		return err
	} else {
		allLawyers, err := c.lawyerRepository.FindAll()
		if err != nil {
			return err
		}
		return ctx.Status(200).JSON(allLawyers)
	}
}
func (c *adminController) GetLawyer(ctx *fiber.Ctx) error {
	idToken := ctx.Get("idtoken")
	lawyerID := ctx.Get("lawyerID")
	if _, err := c.authService.HasRole(idToken, "admin"); err != nil {
		return err
	} else {
		lawyer, err := c.lawyerRepository.FindByID(lawyerID)
		if err != nil {
			return err
		}
		return ctx.Status(200).JSON(lawyer)
	}
}
func (c *adminController) GetCourtCase(ctx *fiber.Ctx) error {
	idToken := ctx.Get("idtoken")
	courtCaseID := ctx.Get("courtCaseID")
	if _, err := c.authService.HasRole(idToken, "admin"); err != nil {
		return err
	} else {
		courtCase, err := c.courtcaseRepository.FindByID(courtCaseID)
		if err != nil {
			return err
		}
		return ctx.Status(200).JSON(courtCase)
	}
}

//Toggle verification status
func (c *adminController) VerifyLawyer(ctx *fiber.Ctx) error {
	idToken := ctx.Get("idtoken")
	lawyerID := ctx.Get("lawyerID")
	if _, err := c.authService.HasRole(idToken, "admin"); err != nil {
		return err
	} else {
		status, err := c.adminService.ToggleVerification(lawyerID)
		if err != nil {
			return err
		}
		return ctx.Status(200).JSON(status)
	}
}
