package internal

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	AppName    string   `json:"app_name"`
	AuthDomain string   `json:"auth_domain"`
	Languages  []string `json:"languages"`
	Domains    []string `json:"domains"`
	Server     struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
	Recaptcha struct {
		Enabled        bool    `json:"enabled"`
		SecretKey      string  `json:"secret_key"`
		VerifyURL      string  `json:"verify_url"`
		ExpectedAction string  `json:"expected_action"`
		MinimumScore   float64 `json:"minimum_score"`
		TimeoutSeconds int     `json:"timeout_seconds"`
	} `json:"recaptcha"`
	Database struct {
		Type     string `json:"type"`
		Host     string `json:"host"`
		User     string `json:"user"`
		Password string `json:"password"`
		DbName   string `json:"dbname"`
		Port     int    `json:"port"`
		SSLMode  string `json:"sslmode"`
		DSN      string `json:"dsn"`
	} `json:"database"`
}

func LoadConfiguration(file string) Config {
	var config Config

	configFile, err := os.Open(file)
	if err != nil {
		logrus.WithError(err).Fatal("problem loading config file")
	}
	defer configFile.Close()

	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		logrus.WithError(err).Fatal("problem parsing config file")
	}

	return config
}

func (c Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func NewDatabaseConnection(cfg Config) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	switch cfg.DatabaseType() {
	case "postgres":
		db, err = gorm.Open(postgres.Open(cfg.PostgresDSN()), &gorm.Config{TranslateError: true})
	default:
		db, err = gorm.Open(mysql.Open(cfg.MySQLDSN()), &gorm.Config{TranslateError: true})
	}
	if err != nil {
		logrus.WithError(err).Fatal("unable to connect to database")
	}

	return db
}

func (c Config) DatabaseType() string {
	switch strings.ToLower(strings.TrimSpace(c.Database.Type)) {
	case "postgres", "postgresql":
		return "postgres"
	default:
		return "mysql"
	}
}

func (c Config) MySQLDSN() string {
	if dsn := strings.TrimSpace(c.Database.DSN); dsn != "" {
		return dsn
	}

	loc := url.QueryEscape("UTC")
	tz := url.QueryEscape("'+00:00'")

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=%s&time_zone=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DbName,
		loc,
		tz,
	)
}

func (c Config) PostgresDSN() string {
	if dsn := strings.TrimSpace(c.Database.DSN); dsn != "" {
		return dsn
	}

	sslMode := strings.TrimSpace(c.Database.SSLMode)
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		c.Database.Host,
		c.Database.User,
		c.Database.Password,
		c.Database.DbName,
		c.Database.Port,
		sslMode,
	)
}
