package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Service struct {
	Name string
	URL  string
}

var services = map[string]Service{
	"users": {
		Name: "user",
		URL:  "http://localhost:8001",
	},
	"products": {
		Name: "product",
		URL:  "http://localhost:8002",
	},
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	api := e.Group("/api")

	for path, service := range services {
		api.Any("/"+path+"/*", proxyHandler(service.URL), middleware.Rewrite(map[string]string{
			"/api/" + path + "/*": "/$1",
		}))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func proxyHandler(targetUrl string) echo.HandlerFunc {
	return func(c echo.Context) error {
		target, err := url.Parse(targetUrl)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Invalid target URL",
			})
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		req := c.Request()
		res := c.Response()

		proxy.ServeHTTP(res, req)
		return nil
	}
}
