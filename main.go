package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"k8s.io/utils/env"
)

type Engine struct {
	busy      sync.Mutex
	busyUntil time.Time
}

var work Engine

func main() {
	work = Engine{
		busy: sync.Mutex{},
	}

	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	router.GET("/", handler)
	router.GET("/block", blockHandler)

	router.GET("/health", livenessHandler)

	router.Logger.Fatal(router.Start(":8080"))
}

func livenessHandler(c echo.Context) error {
	log.Printf("Health hit by %s check.  Busy until: %s.  Is still busy: %v", c.QueryParam("type"), work.busyUntil.Format(time.RFC3339), work.busyUntil.After(time.Now()))

	if work.busyUntil.IsZero() {
		return c.String(http.StatusOK, "OK")
	}

	work.busy.Lock()

	work.busy.Unlock()

	return c.String(http.StatusOK, "OK")
}

func blockHandler(c echo.Context) error {
	t := c.QueryParam("time")

	if t == "" {
		return c.String(http.StatusOK, "No time provided")
	}

	dur, err := time.ParseDuration(t)
	if err != nil {
		return err
	}

	work.busyUntil = time.Now().Add(dur)

	work.busy.Lock()
	log.Println("Busy until: ", dur.String())

	go func() {
		time.Sleep(dur)

		work.busy.Unlock()
		log.Println("Not busy anymore")
	}()

	return c.String(http.StatusOK, "OK")
}

func handler(c echo.Context) error {
	podIp := env.GetString("POD_IP", "localhost")

	if work.busyUntil.Before(time.Now()) {
		return c.String(http.StatusOK, fmt.Sprintf("Not Busy - %s", podIp))
	}

	return c.String(http.StatusOK, fmt.Sprintf("Busy until: %s - %s", work.busyUntil.Format(time.RFC3339), podIp))
}
