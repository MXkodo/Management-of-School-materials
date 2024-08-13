package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MXkodo/Management-of-School-materials/config"
	"github.com/MXkodo/Management-of-School-materials/internal/handler"
	"github.com/MXkodo/Management-of-School-materials/internal/repo"
	"github.com/MXkodo/Management-of-School-materials/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func RunApp(cfg config.Config) error {
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseURL())
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	materialRepo := repo.NewPgMaterialRepository(conn)
	materialService := service.NewMaterialService(materialRepo)

	r := gin.Default()
	handler.RegisterRoutes(r, materialService)

	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	log.Println("Server stopped gracefully")
	return nil
}
