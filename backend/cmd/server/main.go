package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"aiweekly/backend/internal/config"
	"aiweekly/backend/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := config.ConnectMySQL(cfg)
	if err != nil {
		log.Fatalf("connect mysql failed: %v", err)
	}

	if err := config.AutoMigrate(db); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	corsCfg := cors.Config{
		AllowOrigins:     cfg.CORSAllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}
	r.Use(cors.New(corsCfg))

	router.Register(r, db, cfg)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server run failed: %v", err)
	}
}
