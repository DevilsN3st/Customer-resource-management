package entrypoint

import (
	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/entrypoint/middleware"
	rest2 "github.com/icrxz/crm-api-core/internal/entrypoint/rest"
)

func LoadRoutes(
	app *gin.Engine,
	pingController rest2.PingController,
	userController rest2.UserController,
	webMessageController rest2.WebMessageController,
	leadController rest2.LeadController,
	customerController rest2.CustomerController,
	tenantController rest2.TenantController,
	authController rest2.AuthController,
	authMiddleware middleware.AuthenticationMiddleware,
	ticketController rest2.TicketController,
	productController rest2.ProductController,
	commentController rest2.CommentController,
	transactionController rest2.TransactionController,
	ticketActionController rest2.TicketActionController,
) {
	authGroup := app.Group("/crm/core/api/v1")
	authGroup.Use(authMiddleware.Authenticate())

	publicGroup := app.Group("/crm/core/api/v1")

	// miscellaneous
	app.GET("/ping", pingController.Pong)

	// user
	publicGroup.POST("/users", userController.CreateUser)
	authGroup.GET("/users", userController.SearchUser)
	authGroup.GET("/users/:userID", userController.GetUser)
	authGroup.PUT("/users/:userID", userController.UpdateUser)
	authGroup.DELETE("/users/:userID", userController.DeleteUser)

	// lead
	authGroup.POST("/leads", leadController.CreateLead)
	authGroup.GET("/leads", leadController.SearchLeads)
	authGroup.GET("/leads/:leadID", leadController.GetLead)
	authGroup.PUT("/leads/:leadID", leadController.UpdateLead)
	authGroup.DELETE("/leads/:leadID", leadController.DeleteLead)
	authGroup.POST("/leads/batch", leadController.CreateBatch)

	// customers
	authGroup.POST("/customers", customerController.CreateCustomer)
	authGroup.GET("/customers", customerController.SearchCustomers)
	authGroup.GET("/customers/:customerID", customerController.GetCustomer)
	authGroup.PUT("/customers/:customerID", customerController.UpdateCustomer)
	authGroup.DELETE("/customers/:customerID", customerController.DeleteCustomer)

	// tenants
	authGroup.POST("/tenants", tenantController.CreateTenant)
	authGroup.GET("/tenants", tenantController.SearchTenants)
	authGroup.GET("/tenants/:tenantID", tenantController.GetTenant)
	authGroup.PUT("/tenants/:tenantID", tenantController.UpdateTenant)
	authGroup.DELETE("/tenants/:tenantID", tenantController.DeleteTenant)

	// auth
	publicGroup.POST("/login", authController.Login)
	authGroup.POST("/logout", authController.Logout)

	// webMessage
	publicGroup.POST("/web/message", webMessageController.ReceiveMessage)

	// tickets
	authGroup.POST("/tickets", ticketController.CreateTicket)
	authGroup.GET("/tickets/:ticketID", ticketController.GetTicket)
	authGroup.PATCH("/tickets/:ticketID", ticketController.UpdateTicket)
	authGroup.GET("/tickets", ticketController.SearchTickets)

	// products
	authGroup.GET("/products/:productID", productController.GetProductByID)

	// comments
	authGroup.GET("/comments/:commentID", commentController.GetByID)
	authGroup.POST("/tickets/:ticketID/comments", commentController.CreateComment)
	authGroup.GET("/tickets/:ticketID/comments", commentController.GetByTicketID)

	// transactions
	authGroup.POST("/tickets/:ticketID/transactions", transactionController.CreateTransaction)
	authGroup.GET("/transactions/:transactionID", transactionController.GetTransaction)
	authGroup.PUT("/transactions/:transactionID", transactionController.UpdateTransaction)
	authGroup.GET("/transactions", transactionController.SearchTransactions)

	// ticket actions
	authGroup.PATCH("/tickets/:ticketID/owner", ticketActionController.ChangeOwner)
	authGroup.PATCH("/tickets/:ticketID/status", ticketActionController.ChangeStatus)
	authGroup.PATCH("/tickets/:ticketID/lead", ticketActionController.ChangeLead)
	authGroup.GET("/tickets/:ticketID/report", ticketActionController.DownloadReport)
}
