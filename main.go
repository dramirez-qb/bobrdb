package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

// ========== middlewares

func noRoute(c *gin.Context) {
	c.String(http.StatusNotFound, "nothing here")
}

// NewRouter just build our router
func NewRouter() *gin.Engine {
	router := gin.Default()

	exp := os.Getenv("APP_EXPIRATION")
	if exp == "" {
		exp = "5"
	}

	purge := os.Getenv("APP_PURGE")
	if purge == "" {
		purge = "10"
	}
	expv, _ := strconv.ParseFloat(exp, 10)
	purgev, _ := strconv.ParseFloat(purge, 10)

	ourCache := cache.New(time.Duration(expv)*time.Minute, time.Duration(purgev)*time.Minute)

	// no route, bad url
	router.NoRoute(noRoute)

	// Ping handler
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Health handler
	router.GET("/healthz", func(c *gin.Context) {
		c.String(200, "healthy")
	})

	// Simple group: v1
	api := router.Group("/api")
	{
		api.GET("/:key", func(c *gin.Context) {
			key := c.Param("key")
			if x, found := ourCache.Get(key); found {
				c.JSON(http.StatusOK, gin.H{"status": "ok", "key": key, "value": x.(string)})
			} else {
				c.JSON(http.StatusNotFound, gin.H{"msg": "wrong", "value": nil})
			}
		})
		api.POST("/:key/:value", func(c *gin.Context) {
			key := c.Param("key")
			value := c.Param("value")
			if _, found := ourCache.Get(key); found {
				ourCache.Replace(key, value, cache.DefaultExpiration)
			} else {
				ourCache.Set(key, value, cache.DefaultExpiration)
			}
			c.JSON(http.StatusOK, gin.H{"status": "created", "key": key, "value": value})
		})
		api.PUT("/:key/:value", func(c *gin.Context) {
			key := c.Param("key")
			value := c.Param("value")
			if _, found := ourCache.Get(key); found {
				ourCache.Replace(key, value, cache.DefaultExpiration)
			} else {
				ourCache.Set(key, value, cache.DefaultExpiration)
			}
			c.JSON(http.StatusOK, gin.H{"status": "updated", "key": key, "value": value})
		})
		api.DELETE("/:key", func(c *gin.Context) {
			key := c.Param("key")
			if x, found := ourCache.Get(key); found {
				ourCache.Delete(key)
				c.JSON(http.StatusOK, gin.H{"status": "deleted", "key": key, "value": x})
			} else {
				c.JSON(http.StatusNotFound, gin.H{"msg": "wrong", "key": key, "value": nil})
			}
		})
		api.PATCH("/:key", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "PATCH"})
		})
		api.HEAD("/:key", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "HEAD"})
		})
		api.OPTIONS("/:key", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "OPTIONS"})
		})
	}
	return router
}

func startServer() {
	// set config

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// get and start router
	router := NewRouter()
	router.Run(":" + port)
}

func main() {
	startServer()
}
