package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/webdav"
)

func healthSkipper(c echo.Context) bool {
	if c.Path() == "/health" {
		return true
	}
	return false
}

func basicAuthValidator(authString string) func(string, string, echo.Context) (bool, error) {
	return func(username string, password string, c echo.Context) (bool, error) {
		if username+":"+password == authString {
			return true, nil
		}
		return false, nil
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	basicAuth := os.Getenv("BASIC_AUTH")
	if basicAuth == "" {
		log.Fatal("Environment variable \"BASIC_AUTH\" is required")
	}

	directory := os.Getenv("DIRECTORY")
	if directory == "" {
		directory = "./"
	}
	log.Printf("Port: %s, Directory: %s", port, directory)

	e := echo.New()
	e.HideBanner = true
	authConfig := middleware.BasicAuthConfig{
		Skipper:   healthSkipper,
		Validator: basicAuthValidator(basicAuth),
	}

	srv := &webdav.Handler{
		FileSystem: webdav.Dir("./"),
		LockSystem: webdav.NewMemLS(),
	}

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.BasicAuthWithConfig(authConfig))

	e.GET("health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	e.Any("/*", echo.WrapHandler(srv))
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
