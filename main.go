package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/namedotcom/go/namecom"
)

const (
	envsDir = "environments"
)

var (
	logger *slog.Logger
)

func init() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)
	loptions := slog.HandlerOptions{
		Level: lvl,
	}
	if os.Getenv("PRODUCTION") != "" {
		logger = slog.New(slog.NewTextHandler(os.Stderr, &loptions))
	} else {
		lvl.Set(slog.LevelDebug)
		logger = slog.New(slog.NewTextHandler(os.Stdout, &loptions))
	}
	logger.With("app", "echolocation")
}

func main() {
	logger.Info("Starting service")
	name, err := getNameClient()
	if err != nil {
		log.Fatalf("failed to get name.com client: %s", err)
	}

	logger.Info("Getting a list of domains")
	domainList, err := name.ListDomains(&namecom.ListDomainsRequest{})
	if err != nil {
		log.Fatalf("failed to get domains: %s", err)
	}

	for _, d := range domainList.Domains {
		logger.Info("Found domain", "name", d.DomainName, "created", d.CreateDate, "expiry", d.ExpireDate)
	}
}

func getNameClient() (*namecom.NameCom, error) {
	if os.Getenv("PRODUCTION") != "" {
		if err := godotenv.Load(filepath.Join(envsDir, "production.env")); err != nil {
			return nil, fmt.Errorf("unable to load production environment: %w", err)
		}
		logger.Info("Using production client...")
		return namecom.New(os.Getenv("ECHOLOCATION_USER"), os.Getenv("ECHOLOCATION_TOKEN")), nil
	}
	if err := godotenv.Load(filepath.Join(envsDir, "development.env")); err != nil {
		return nil, fmt.Errorf("unable to load development environment: %w", err)
	}
	logger.Info("Using test client...")
	return namecom.Test(os.Getenv("ECHOLOCATION_USER"), os.Getenv("ECHOLOCATION_TOKEN")), nil
}
