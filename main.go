package main

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	/*
		TODO
			Profiles for testing, staging and prod
			Rate limiting
			Auditing
			MSG91 API
			End-points for admin panel, lawyer panel and user panel(primitive).
			Validation on all entities, throw bad request for those cases
			Write tests
	*/
	ctx := context.Background()
	SrvEnv := os.Args[1]

	//Set dev or prod environment
	if SrvEnv == "dev" {
		if err := os.Setenv("SRV_DEV", "true"); err != nil {
			log.Panic(err)
		}
		if err := os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8082"); err != nil {
			log.Panic(err)
		}
		if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/Users/rohanbojja/Creds/legalaidcamp.json"); err != nil {
			log.Panic(err)
		}
	}


	conf := &firebase.Config{
		//Retrieve from viper config
		DatabaseURL:   "https://legalaidcamp-b5e4d.firebaseio.com/",
		StorageBucket: "gs://legalaidcamp-b5e4d.appspot.com",
	}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	userController := InitializeUserController(app)
	adminController := InitializeAdminController(app)
	lawyerController := InitializeLawyerController(app)

	server := fiber.New()
	server.Use(cors.New())
	server.Static("/", "./static")

	//Api group
	api := server.Group("/api")

	//Check environment
	api.Get("/info", func(ctx *fiber.Ctx) error {
		env:= map[string]interface{}{
			"DEV" : os.Getenv("SRV_DEV"),
			"FIRESTORE_EMULATOR_HOST" : os.Getenv("FIRESTORE_EMULATOR_HOST"),
		}
		return ctx.Status(200).JSON(env)
	})

	//End-point definitions

	//User end-points
	UserAPI := api.Group("/user")
	//Create a case
	UserAPI.Post("", userController.CreateCase)

	//Lawyers end-points
	lawyerAPI := api.Group("/lawyers")
	//Create a lawyer
	lawyerAPI.Post("", lawyerController.CreateLawyer)
	//Toggle profile status
	lawyerAPI.Get("/toggle", lawyerController.ToggleStatus)
	//Accept case
	lawyerAPI.Get("/cases/accept", lawyerController.AcceptCase)
	//Reject case
	lawyerAPI.Get("/cases/reject", lawyerController.RejectCase)
	//FindByID all assigned cases
	lawyerAPI.Get("/cases/assigned", lawyerController.GetAssignedCases)
	//Upload photocopy of Bar Council license
	//TODO IMPL

	//Admin end-points
	adminAPI := api.Group("/admin")
	//FindByID all cases
	adminAPI.Get("/cases", adminController.GetAllCases)
	//FindByID all lawyers
	adminAPI.Get("/lawyers", adminController.GetAllLawyers)
	//Verify lawyer
	adminAPI.Get("/verify", adminController.VerifyLawyer)
	//FindByID a lawyer by id
	adminAPI.Get("/lawyer", adminController.GetLawyer)
	//FindByID a court case by id
	adminAPI.Get("/case", adminController.GetCourtCase)

	// Only for testing, disable in production.
	if err := server.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
	log.Println("We good.")
	//Generating a few cases in dev mode,

	if os.Getenv("SRV_DEV") == "true"{
		/*
		TODO
			Inflate Firestore Emulator with content!
		 */
	}
}
