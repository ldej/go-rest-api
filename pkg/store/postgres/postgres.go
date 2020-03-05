package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq" // import postgres
	config "github.com/spf13/viper"
	"go.uber.org/zap"
)

type Client struct {
	logger *zap.SugaredLogger
	dbName string
	db     *sql.DB
}

func New() (*Client, error) {
	logger := zap.S().With("package", "storage.postgres")

	var dbCreds string
	var dbURL string
	var dbURLOptions string
	var dbCreateCreds string

	// Username
	if username := config.GetString("storage.username"); username != "" {
		dbCreds = username + ":" + config.GetString("storage.password")
	} else {
		return nil, fmt.Errorf("no username specified")
	}

	// Check to see if we have an admin password specified (that we will use to create the database if it does not exist)
	pgPassword := config.GetString("POSTGRES_PASSWORD")
	if pgPassword != "" {
		pgUser := config.GetString("POSTGRES_USER")
		if pgUser == "" {
			pgUser = "postgres"
		}
		dbCreateCreds = fmt.Sprintf("%s:%s", pgUser, pgPassword)
	} else {
		dbCreateCreds = dbCreds
	}

	// Host + Port
	if hostname := config.GetString("storage.host"); hostname != "" {
		dbURL += "@" + hostname
	} else {
		return nil, fmt.Errorf("no hostname specified")
	}
	if port := config.GetString("storage.port"); port != "" {
		dbURL += ":" + port
	}

	// Database Name
	dbName := config.GetString("storage.database")
	if dbName == "" {
		return nil, fmt.Errorf("No database specified")
	}

	// SSL Mode
	if sslMode := config.GetString("storage.sslmode"); sslMode != "" {
		// dbConnection += fmt.Sprintf("sslmode=%s ", sslMode)
		dbURLOptions += fmt.Sprintf("?sslmode=%s", sslMode)
	}

	for retries := config.GetInt("storage.retries"); retries > 0; retries-- {
		createDb, err := sql.Open("postgres", "postgres://"+dbCreateCreds+dbURL+dbURLOptions)
		// Attempt to create the database if it doesn't exist
		if err == nil {
			defer createDb.Close()
			// See if it exists
			var one sql.NullInt64
			err = createDb.QueryRow(`SELECT 1 from pg_database WHERE datname=$1`, dbName).Scan(&one)
			if err == nil {
				break // already exists
			} else if err != sql.ErrNoRows && !strings.Contains(err.Error(), "does not exist") {
				// Some other error besides does not exist
				return nil, fmt.Errorf("could not check for database: %s", err)
			}
		} else if strings.Contains(err.Error(), "permission denied") {
			return nil, fmt.Errorf("could not connect to database: %s", err)
		} else if strings.Contains(err.Error(), "connection refused") {
			logger.Warnw("Connection to database timed out. Sleeping and retry.",
				"storage.host", config.GetString("storage.host"),
				"storage.username", config.GetString("storage.username"),
				"storage.password", "****",
				"storage.port", config.GetInt("storage.port"),
			)
			time.Sleep(config.GetDuration("storage.sleep_between_retries"))
			continue
		} else {
			return nil, err
		}
		logger.Infow("Creating database", "database", dbName)
		_, err = createDb.Exec(`CREATE DATABASE ` + dbName)
		if err != nil {
			return nil, fmt.Errorf("could not create database: %s", err)
		}
		break
	}

	// Build the full DB URL
	fullDbURL := "postgres://" + dbCreds + dbURL + "/" + dbName + dbURLOptions

	// Make the connection using the sqlx connector now that we know the database exists
	db, err := sql.Open("postgres", fullDbURL)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %s", err)
	}

	// Ping the database
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping database %s", err)
	}

	db.SetMaxOpenConns(config.GetInt("storage.max_connections"))

	logger.Debugw("Connected to database server",
		"storage.host", config.GetString("storage.host"),
		"storage.username", config.GetString("storage.username"),
		"storage.port", config.GetInt("storage.port"),
		"storage.database", config.GetString("storage.database"),
	)

	c := &Client{
		logger: logger,
		dbName: dbName,
		db:     db,
	}

	return c, nil
}
