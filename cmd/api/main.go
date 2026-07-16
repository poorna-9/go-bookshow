package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/poorna-9/goshow/internal/config"
	"github.com/poorna-9/goshow/internal/handlers"
	"github.com/poorna-9/goshow/internal/models"
	"github.com/poorna-9/goshow/internal/repositories"
	"github.com/poorna-9/goshow/internal/routes"
	"github.com/poorna-9/goshow/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	log.Println("connected to database successfully")

	err = models.AutoMigrate(db)
	if err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}
	log.Println("all tables migrated successfully")

	// Redis
	redisClient := config.NewRedisClient(cfg)

	// Razorpay
	razorpayClient := config.NewRazorpayClient(cfg)

	// Theatre
	theatreRepo := repositories.NewTheatreRepository(db)
	theatreService := services.NewTheatreService(theatreRepo)
	theatreHandler := handlers.NewTheatreHandler(theatreService)

	// Screen
	screenRepo := repositories.NewScreenRepository(db)
	screenService := services.NewScreenService(screenRepo)
	screenHandler := handlers.NewScreenHandler(screenService)

	// Seat
	seatRepo := repositories.NewSeatRepository(db)
	seatService := services.NewSeatService(seatRepo)
	seatHandler := handlers.NewSeatHandler(seatService)

	// Movie
	movieRepo := repositories.NewMovieRepository(db)
	movieService := services.NewMovieService(movieRepo)
	movieHandler := handlers.NewMovieHandler(movieService)

	// Show
	showRepo := repositories.NewShowRepository(db)
	showService := services.NewShowService(showRepo)
	showHandler := handlers.NewShowHandler(showService)

	// Booking
	bookingRepo := repositories.NewBookingRepository(db, redisClient, razorpayClient, cfg.RazorpayKeySecret, cfg.RazorpayWebhookSecret)
	bookingService := services.NewBookingService(bookingRepo, cfg.RazorpayKeyID)
	bookingHandler := handlers.NewBookingHandler(bookingService)

	// Auth
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	router := gin.Default()

	router.Static("/css", "./web/css")
	router.Static("/js", "./web/js")
	router.Static("/images", "./web/images")

	// Serve frontend pages
	router.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	router.GET("/movies", func(c *gin.Context) {
		c.File("./web/movies.html")
	})

	router.GET("/theatres", func(c *gin.Context) {
		c.File("./web/theatres.html")
	})

	router.GET("/checkout", func(c *gin.Context) {
		c.File("./web/checkout.html")
	})

	router.GET("/login", func(c *gin.Context) {
		c.File("./web/login.html")
	})

	router.GET("/signup", func(c *gin.Context) {
		c.File("./web/signup.html")
	})

	routes.RegisterRoutes(router, &routes.Handlers{
		Theatre: theatreHandler,
		Screen:  screenHandler,
		Seat:    seatHandler,
		Movie:   movieHandler,
		Show:    showHandler,
		Booking: bookingHandler,
		Auth:    authHandler,
	}, "")

	log.Printf("starting server on :%s", cfg.AppPort)
	router.Run(":" + cfg.AppPort)
}
