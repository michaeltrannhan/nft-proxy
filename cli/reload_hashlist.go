package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	nft_proxy "github.com/alphabatem/nft-proxy"
	services "github.com/alphabatem/nft-proxy/service"
	"github.com/babilu-online/common/context"
	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	HashlistPath string
	APIEndpoint  string
	WorkerCount  int
}

// Hashlist represents a collection of NFT hashes
type Hashlist []string

func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	return &Config{
		HashlistPath: "./hashlist.json",
		APIEndpoint:  "https://api.degencdn.com/v1/nfts/%s/image.jpg",
		WorkerCount:  5,
	}, nil
}

func initializeContext() (*context.Context, error) {
	mainContext, err := context.NewCtx(
		&services.SqliteService{},
		&services.SolanaImageService{},
		&services.ImageService{},
		&services.ResizeService{},
		&services.SolanaService{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize context: %w", err)
	}

	if err := mainContext.Run(); err != nil {
		return nil, fmt.Errorf("failed to run context: %w", err)
	}

	return mainContext, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, err := initializeContext()
	if err != nil {
		log.Fatalf("Failed to initialize context: %v", err)
	}

	if err := run(ctx, cfg); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run(ctx *context.Context, cfg *Config) error {
	db := ctx.Service(services.SQLITE_SVC).(*services.SqliteService)

	hashes, err := loadHashlist(cfg.HashlistPath)
	if err != nil {
		return fmt.Errorf("failed to load hashlist: %w", err)
	}

	log.Printf("Processing %d mints", len(hashes))

	// Get initial count
	var count int64
	if err := db.Db().Model(&nft_proxy.SolanaMedia{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to get initial count: %w", err)
	}
	log.Printf("Initial record count: %d", count)

	// Delete existing records
	if err := deleteExistingRecords(db, hashes); err != nil {
		return fmt.Errorf("failed to delete existing records: %w", err)
	}

	// Get final count
	if err := db.Db().Model(&nft_proxy.SolanaMedia{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to get final count: %w", err)
	}
	log.Printf("Final record count: %d", count)

	return nil
}

func deleteExistingRecords(db *services.SqliteService, hashes Hashlist) error {
	if len(hashes) == 0 {
		return nil
	}

	query := `mint IN ("` + strings.Join(hashes, `","`) + `")`
	return db.Db().Where(query).Delete(&nft_proxy.SolanaMedia{}).Error
}

func reloadRemote(hashes Hashlist, cfg *Config) error {
	client := &http.Client{Timeout: 5 * time.Second}
	var wg sync.WaitGroup
	errors := make(chan error, len(hashes))
	semaphore := make(chan struct{}, cfg.WorkerCount)

	for _, hash := range hashes {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore

		go func(h string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			url := fmt.Sprintf(cfg.APIEndpoint, h)
			resp, err := client.Get(url)
			if err != nil {
				errors <- fmt.Errorf("failed to load hash %s: %w", h, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				errors <- fmt.Errorf("failed to load hash %s: status %d", h, resp.StatusCode)
			}
		}(hash)
	}

	wg.Wait()
	close(errors)

	// Collect errors
	var errs []error
	for err := range errors {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors during reload", len(errs))
	}

	return nil
}

func reloadLocally(img *services.SolanaImageService, hashes Hashlist, cfg *Config) error {
	var wg sync.WaitGroup
	errors := make(chan error, len(hashes))
	semaphore := make(chan struct{}, cfg.WorkerCount)

	for _, hash := range hashes {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore

		go func(h string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			if _, err := img.Media(h, true); err != nil {
				errors <- fmt.Errorf("failed to load media for hash %s: %w", h, err)
			}
		}(hash)
	}

	wg.Wait()
	close(errors)

	// Collect errors
	var errs []error
	for err := range errors {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered %d errors during reload", len(errs))
	}

	return nil
}

func loadHashlist(location string) (Hashlist, error) {
	data, err := os.ReadFile(location)
	if err != nil {
		return nil, fmt.Errorf("failed to read hashlist file: %w", err)
	}

	var hashlist Hashlist
	if err := json.Unmarshal(data, &hashlist); err != nil {
		return nil, fmt.Errorf("failed to parse hashlist JSON: %w", err)
	}

	return hashlist, nil
}
