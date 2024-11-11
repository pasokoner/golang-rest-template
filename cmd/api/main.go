package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pito-bataan/tourism-be/internal/database"
	"github.com/pito-bataan/tourism-be/internal/env"
	"github.com/pito-bataan/tourism-be/internal/version"

	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	baseURL   string
	httpPort  int
	basicAuth struct {
		username       string
		hashedPassword string
	}
	cookie struct {
		secretKey string
	}
	db struct {
		database string
		password string
		username string
		port     string
		host     string
		schema   string
	}
	jwt struct {
		secretKey string
	}
	notifications struct {
		email string
	}
	// smtp struct {
	// 	host     string
	// 	port     int
	// 	username string
	// 	password string
	// 	from     string
	// }
}

type application struct {
	config config
	db     *database.Queries
	dbPool *pgxpool.Pool
	logger *slog.Logger
	// mailer *smtp.Mailer
	wg sync.WaitGroup
}

func run(logger *slog.Logger) error {
	var cfg config

	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:8080")
	cfg.httpPort = env.GetInt("HTTP_PORT", 4444)
	cfg.basicAuth.username = env.GetString("BASIC_AUTH_USERNAME", "admin")
	cfg.basicAuth.hashedPassword = env.GetString("BASIC_AUTH_HASHED_PASSWORD", "$2a$10$jRb2qniNcoCyQM23T59RfeEQUbgdAXfR6S0scynmKfJa5Gj3arGJa")
	cfg.cookie.secretKey = env.GetString("COOKIE_SECRET_KEY", "hm57szss4dlnbjqj4ewj435f5ky24hed")
	cfg.db.database = env.GetString("DB_DATABASE", "db")
	cfg.db.password = env.GetString("DB_PASSWORD", "pass")
	cfg.db.username = env.GetString("DB_USERNAME", "user")
	cfg.db.port = env.GetString("DB_PORT", "5432")
	cfg.db.host = env.GetString("DB_HOST", "localhost")
	cfg.db.schema = env.GetString("DB_SCHEMA", "public")
	cfg.jwt.secretKey = env.GetString("JWT_SECRET_KEY", "rdquldkj7xctaicp6j5kjxornwlqcs35")
	// cfg.notifications.email = env.GetString("NOTIFICATIONS_EMAIL", "")
	// cfg.smtp.host = env.GetString("SMTP_HOST", "example.smtp.host")
	// cfg.smtp.port = env.GetInt("SMTP_PORT", 25)
	// cfg.smtp.username = env.GetString("SMTP_USERNAME", "example_username")
	// cfg.smtp.password = env.GetString("SMTP_PASSWORD", "pa55word")
	// cfg.smtp.from = env.GetString("SMTP_FROM", "Example Name <no_reply@example.org>")

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", cfg.db.username, cfg.db.password, cfg.db.host, cfg.db.port, cfg.db.database, cfg.db.schema)
	dbPool, err := database.Connect(connStr)

	if err != nil {
		log.Fatal(err)
		return err
	}

	db := database.New(dbPool)

	// mailer, err := smtp.NewMailer(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.from)
	// if err != nil {
	// 	return err
	// }

	app := &application{
		config: cfg,
		db:     db,
		logger: logger,
		dbPool: dbPool,
		// mailer: mailer,
	}

	return app.serveHTTP()
}
