package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tournament/api" // <== rest api
	"tournament/pkg/db"
	"tournament/site" // <== site pages(in our case only one page)

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

var (
	server   *http.Server
	osSignal chan os.Signal
)

func main() {
	defer fmt.Println("Server shutdown")
	osSignal = make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background()) // main ctx if we need to cancel things around, do not confuse with gin ctx, gin creates ctx for every request(add db timeout on queries for each request above gin request ctx)
	defer cancel()

	_ = godotenv.Load()
	// ignore .env errors, use default values instead
	// in docker specify env variables by ENV command

	databaseUrl := os.Getenv("DATABASE_URL")
	if "" == databaseUrl {
		databaseUrl = "postgresql://postgres:postgres@localhost:5432/tournament" // default value
	}

	conn, err := connectDB(ctx, databaseUrl)
	if err != nil {
		log.Fatalf("Error on postgres database: %v\n", err)
	}
	defer conn.Close()

	app := gin.New()
	app.Use(gin.Recovery()) // recovery middleware
	setupRouter(app)

	serverHost := os.Getenv("SERVER_HOST")
	server = &http.Server{
		Addr:    serverHost + ":8080",
		Handler: app,
	}

	// Start the server concurrently
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Unexpected server error because of: %v\n", err)
		}
	}()

	// Catch the exit signal, graceful shutdown
	<-osSignal

	fmt.Println("Terminating server")
	server.Shutdown(context.Background())
}

func connectDB(ctx context.Context, url string) (*pgxpool.Pool, error) {
	conn, err := db.Connect(ctx, url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func setupRouter(r *gin.Engine) {

	// api
	r.GET("/api", api.RootIndex)
	r.GET("/api/teams", api.GetTeamList)
	r.POST("/api/teams", api.CreateTeam)
	// PUT /api/teams/:id // skip
	// DELETE /api/teams/:id // skip
	r.POST("/api/teams/generate", api.GenerateTeams) // generates 16 ucl(uefa champions league) teams <=== there is no point in task about it, I just added it to myself to add all teams at once
	r.GET("/api/divisions", api.GetDivisions)
	r.POST("/api/divisions/prepare", api.PrepareDivisions)
	r.POST("/api/divisions/start", api.StartDivisions)

	r.GET("/api/playoff", api.GetPlayoffs)

	r.POST("/api/playoff/prepare", api.PreparePlayoff)
	r.POST("/api/playoff/start", api.StartPlayoff)

	r.POST("/api/cleanup", api.Cleanup)

	// site
	r.LoadHTMLGlob("./templates/**/*")
	r.Static("/assets", "./assets")

	// frontend
	frontendRoutes := []string{
		"/", // welcome or redirect to teams
		"/teams",
		"/divisions",
		"/playoff",
		"/cleanup",
	}

	for _, route := range frontendRoutes {
		r.GET(route, site.AppHome)
	}
}
