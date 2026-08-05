package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/randibudi/airline-voucher/backend/internal/app"
	"github.com/randibudi/airline-voucher/backend/internal/domain"
	"github.com/randibudi/airline-voucher/backend/internal/httpapi"
	"github.com/randibudi/airline-voucher/backend/internal/repository"
	"github.com/randibudi/airline-voucher/backend/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	databasePath := os.Getenv("DATABASE_PATH")
	if databasePath == "" {
		databasePath = "vouchers.db"
	}

	db, err := repository.OpenSQLite(ctx, databasePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	voucherRepository := repository.NewSQLite(db)
	if err := voucherRepository.Initialize(ctx); err != nil {
		return err
	}

	generator := domain.NewDefaultSeatGenerator()
	voucherService := service.NewVoucher(voucherRepository, generator)
	handler := httpapi.NewHandler(voucherService)
	application := app.New(handler)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := application.ShutdownWithContext(shutdownCtx); err != nil {
			log.Printf("shutdown server: %v", err)
		}
	}()

	if err := application.Listen(":8080"); err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	return nil
}
