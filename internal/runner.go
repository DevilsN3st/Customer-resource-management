package internal

import (
	"context"
	"github.com/icrxz/crm-api-core/config"
	entrypoint2 "github.com/icrxz/crm-api-core/internal/entrypoint"
	"github.com/icrxz/crm-api-core/internal/entrypoint/middleware"
	rest2 "github.com/icrxz/crm-api-core/internal/entrypoint/rest"
	bucket2 "github.com/icrxz/crm-api-core/internal/repository/bucket"
	database2 "github.com/icrxz/crm-api-core/internal/repository/database"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
)

func RunApp() error {
	_ = context.Background()

	appConfig, err := config.Load()
	if err != nil {
		return err
	}

	// database
	mongoDB, err := database2.NewDBInstance(appConfig.Database)
	if err != nil {
		return err
	}

	// bucket
	s3Client, err := bucket2.NewS3Bucket(context.Background(), appConfig.AttachmentsBucket)
	if err != nil {
		return err
	}

	attachmentBucket := bucket2.NewAttachmentBucket(s3Client, appConfig.AttachmentsBucket.Name)

	// repositories
	userRepository := database2.NewUserRepository(mongoDB)
	leadRepository := database2.NewLeadRepository(mongoDB)
	customerRepository := database2.NewCustomerRepository(mongoDB)
	tenantRepository := database2.NewTenantRepository(mongoDB)
	ticketRepository := database2.NewTicketRepository(mongoDB)
	productRepository := database2.NewProductRepository(mongoDB)
	commentRepository := database2.NewCommentRepository(mongoDB)
	transactionRepository := database2.NewTransactionRepository(mongoDB)
	attachmentRepository := database2.NewAttachmentRepository(mongoDB)

	// services
	userService := application.NewUserService(userRepository)
	leadService := application.NewLeadService(leadRepository)
	customerService := application.NewCustomerService(customerRepository)
	tenantService := application.NewTenantService(tenantRepository)
	authService := application.NewAuthService(userRepository, appConfig.SecretKey())
	productService := application.NewProductService(productRepository)
	ticketService := application.NewTicketService(customerService, ticketRepository, productService, userService)
	commentService := application.NewCommentService(commentRepository, attachmentRepository, attachmentBucket)
	transactionService := application.NewTransactionService(transactionRepository)
	reportService := application.NewReportService(
		appConfig.ReportFolder,
		ticketService,
		productService,
		customerService,
		commentService,
		leadService,
		tenantService,
		attachmentBucket,
	)
	ticketActionService := application.NewTicketActionService(ticketRepository, commentService, reportService)

	// controllers
	pingController := rest2.NewPingController()
	userController := rest2.NewUserController(userService)
	leadController := rest2.NewLeadController(leadService)
	customerController := rest2.NewCustomerController(customerService)
	tenantController := rest2.NewTenantController(tenantService)
	webMessageController := rest2.NewWebMessageController()
	authController := rest2.NewAuthController(authService)
	ticketController := rest2.NewTicketController(ticketService)
	productController := rest2.NewProductController(productService)
	commentController := rest2.NewCommentController(commentService)
	transactionController := rest2.NewTransactionController(transactionService)
	ticketActionController := rest2.NewTicketActionController(ticketActionService)

	// middlewares
	authMiddleware := middleware.NewAuthenticationMiddleware(authService)

	router := gin.Default()
	router.Use(entrypoint2.CustomErrorEncoder())

	entrypoint2.LoadRoutes(
		router,
		pingController,
		userController,
		webMessageController,
		leadController,
		customerController,
		tenantController,
		authController,
		authMiddleware,
		ticketController,
		productController,
		commentController,
		transactionController,
		ticketActionController,
	)

	return router.Run()
}
