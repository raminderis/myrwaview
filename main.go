package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	"github.com/raminderis/myrwaview/controllers"
	"github.com/raminderis/myrwaview/migrations"
	"github.com/raminderis/myrwaview/models"
	"github.com/raminderis/myrwaview/templates"
	"github.com/raminderis/myrwaview/views"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	//PSQL
	//cfg.PSQL = models.DefaultPostgresConfig()
	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	if cfg.PSQL.Host == "" || cfg.PSQL.Port == "" {
		return cfg, fmt.Errorf("no psql config provided")
	}
	//SMTP
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	//CSRF
	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"
	//Server
	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")
	return cfg, nil
}

func main() {
	//load
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	//DB setup
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	//services setup
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	pwResetService := &models.PasswordResetService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)
	//Setup middleware
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	//CSRF protection
	// csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect([]byte(cfg.CSRF.Key), csrf.Secure(cfg.CSRF.Secure))

	//Setup controller
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: pwResetService,
		EmailService:         emailService,
	}

	//New Router
	r := chi.NewRouter()
	r.Use(csrfMw)
	r.Use(umw.SetUser)
	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))

	usersC.Templates.Signin = views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	r.Get("/signin", usersC.Signin)
	r.Post("/signin", usersC.Authenticate)

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)

	r.Post("/signout", usersC.ProcessSignout)
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "check-your-email.gohtml", "tailwind.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "forgot-pw.gohtml", "tailwind.gohtml"))
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)

	usersC.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, "reset-pw.gohtml", "tailwind.gohtml"))
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)

	r.Get("/currentuser", usersC.CurrentUser)
	r.Route("/currentuser", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found", http.StatusNotFound)
	})
	fmt.Printf("Starting the server on %s...\n", cfg.Server.Address)
	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
