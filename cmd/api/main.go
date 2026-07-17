package main

import (
	"log"
	"time"

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

func startCleanupWorker(bookingService *services.BookingService) {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			bookingService.SweepStaleBookings()
		}
	}()
}
func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	log.Println("connected to database successfully")

	if err := models.AutoMigrate(db); err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}
	log.Println("all tables migrated successfully")

	redisClient := config.NewRedisClient(cfg)
	razorpayClient := config.NewRazorpayClient(cfg)

	theatreRepo := repositories.NewTheatreRepository(db)
	theatreService := services.NewTheatreService(theatreRepo)
	theatreHandler := handlers.NewTheatreHandler(theatreService)

	screenRepo := repositories.NewScreenRepository(db)
	screenService := services.NewScreenService(screenRepo)
	screenHandler := handlers.NewScreenHandler(screenService)

	seatRepo := repositories.NewSeatRepository(db)
	seatService := services.NewSeatService(seatRepo)
	seatHandler := handlers.NewSeatHandler(seatService)

	movieRepo := repositories.NewMovieRepository(db)
	movieService := services.NewMovieService(movieRepo)
	movieHandler := handlers.NewMovieHandler(movieService)

	showRepo := repositories.NewShowRepository(db)
	showService := services.NewShowService(showRepo)
	showHandler := handlers.NewShowHandler(showService)

	bookingRepo := repositories.NewBookingRepository(db, redisClient, razorpayClient, cfg.RazorpayKeySecret, cfg.RazorpayWebhookSecret)
	bookingService := services.NewBookingService(bookingRepo, cfg.RazorpayKeyID)
	bookingHandler := handlers.NewBookingHandler(bookingService)

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.AdminSignupCode)
	authHandler := handlers.NewAuthHandler(authService)

	router := gin.Default()

	router.Static("/css", "./web/css")
	router.Static("/js", "./web/js")
	router.Static("/images", "./web/images")

	pages := map[string]string{
		"/":                "index.html",
		"/movies":          "movies.html",
		"/shows":           "shows.html",
		"/seatmap":         "seatmap.html",
		"/checkout":        "checkout.html",
		"/booking":         "booking.html",
		"/payment-waiting": "payment-waiting.html",
		"/login":           "login.html",
		"/signup":          "signup.html",
		"/admin":           "admin.html",
	}
	for path, file := range pages {
		filePath := "./web/" + file
		router.GET(path, func(c *gin.Context) { c.File(filePath) })
	}

	routes.RegisterRoutes(router, &routes.Handlers{
		Theatre: theatreHandler,
		Screen:  screenHandler,
		Seat:    seatHandler,
		Movie:   movieHandler,
		Show:    showHandler,
		Booking: bookingHandler,
		Auth:    authHandler,
	}, cfg.JWTSecret)
	startCleanupWorker(bookingService)

	log.Printf("starting server on :%s", cfg.AppPort)
	router.Run(":" + cfg.AppPort)
}
