package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ARTMUC/magic-video/api/middleware"
	"github.com/ARTMUC/magic-video/internal/config"
	"github.com/ARTMUC/magic-video/internal/cronjobs"
	"github.com/ARTMUC/magic-video/internal/domain/composition"
	"github.com/ARTMUC/magic-video/internal/domain/customer"
	"github.com/ARTMUC/magic-video/internal/domain/mail"
	order2 "github.com/ARTMUC/magic-video/internal/domain/order"
	"github.com/ARTMUC/magic-video/internal/logger"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cnfg, err := config.LoadConfig(true)
	if err != nil {
		panic(err)
	}

	logger.Init()

	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(
		[]any{
			customer.Customer{},
			customer.CustomerAccess{},
			composition.Image{},
			mail.MailLog{},
			order2.Order{},
			order2.OrderLine{},
			order2.OrderPayment{},
			order2.OrderTransaction{},
			order2.Product{},
			composition.VideoComposition{},
		}...)
	if err != nil {
		panic(err)
	}

	//repositories
	//transactionProvider := base.NewTransactionProvider(db)
	//customerRepository := customer.NewCustomerRepo(db)
	//customerAccessRepository := customer.NewCustomerAccessRepo(db)
	//mailLogRepository := mail.NewMailLogRepo(db)
	//orderLineRepository := order2.NewOrderLineRepo(db)
	//orderPaymentRepository := order2.NewOrderPaymentRepo(db)
	//orderRepository := order2.NewOrderRepo(db)
	//orderTransactionRepository := order2.NewOrderTransactionRepo(db)
	//productRepository := order2.NewProductRepo(db)
	//videoCompositionRepository := composition.NewVideoCompositionRepository(db)
	//
	//// mailer
	//mailSender := mailer.NewBrevoEmailSender(cnfg.BrevoEmailClientConfig(), mailLogRepository)
	//customerAccessMailSender := mailer.NewCustomerAccessEmailSender(
	//	cnfg.ServerConfig(),
	//	mailSender,
	//	mailLogRepository,
	//)

	//crud
	//customerCrud := crud.NewCustomerCrud(customerRepository)

	// services
	//sessionService := service.NewSessionService(cnfg.SessionConfig(), customerCrud)
	//customerService := customer.NewCustomerService(
	//	customerRepository,
	//	customerAccessRepository,
	//	customerAccessMailSender,
	//	cnfg.ServerConfig(),
	//	cnfg.EncryptionConfig(),
	//)
	//orderService := order2.NewOrderService(
	//	transactionProvider,
	//	videoCompositionRepository,
	//	productRepository,
	//	orderRepository,
	//	orderLineRepository,
	//)
	//paymentService := order2.NewPaymentService(
	//	cnfg.Przelewy24ClientConfig(),
	//	transactionProvider,
	//	orderTransactionRepository,
	//	orderPaymentRepository,
	//	orderRepository,
	//)
	//videoCompositionService := job.NewVideoCompositionService()

	// controllers
	//customerAuthController := customerauth.NewCustomerAuthController(customerService, sessionService)
	//orderController := order.NewOrderController(customerService, sessionService, orderService, paymentService)

	router := chi.NewMux()
	router.Use(middleware.PanicRecovery)
	router.Use(middleware.RateLimiter)
	router.Use(middleware.ExtractRequest)

	//api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))
	//customerauth.RegisterRoutes(api, customerAuthController)
	//order.RegisterRoutes(api, orderController)

	c, err := cronjobs.Start(
		[]cronjobs.Func{
			{
				Cron: "* * * * *",
				F: func() {
					//err := videoCompositionService.Create()
					//if err != nil {
					//	logger.Log.Error("Failed to create video compositions", zap.Error(err))
					//}
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Stop()

	logger.Log.Info(fmt.Sprintf("Server running on http://localhost:%s", cnfg.ServerConfig().Port()))
	log.Fatal(http.ListenAndServe(":"+cnfg.ServerConfig().Port(), router))
}
