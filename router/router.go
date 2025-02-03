package router

import (
	"net/http"
	"time"

	. "gitlab.com/CoiaPrant/Sleepy/common/server"
	"gitlab.com/CoiaPrant/Sleepy/router/api"
	"gitlab.com/CoiaPrant/Sleepy/router/handlers"
	"gitlab.com/CoiaPrant/Sleepy/router/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer() *gin.Engine {
	router := gin.New()

	router.HandleMethodNotAllowed = true
	router.SetTrustedProxies([]string{"0.0.0.0/0", "::/0"})

	if gin.Mode() == gin.DebugMode {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())
	router.Use(SameSite)

	if gin.Mode() == gin.DebugMode || AllowCORS {
		router.Use(
			cors.New(cors.Config{
				AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},

				AllowHeaders:  []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
				ExposeHeaders: []string{"Content-Length", "Content-Type", "Content-Disposition", "Set-Authorization"},

				AllowWebSockets:        true,
				AllowBrowserExtensions: true,
				AllowFiles:             true,
				AllowPrivateNetwork:    true,
				AllowCredentials:       true,
				AllowOriginFunc: func(origin string) bool {
					return true
				},

				MaxAge: 24 * time.Hour,
			}),
		)
	}

	{
		// Custom JS/CSS
		router.Static("/custom", "custom")
	}

	{
		// UI
		router := router.Group("")

		router.StaticFile("/", "resources/ui/index.html")

		router.Static("/assets", "resources/ui/assets")
		router.StaticFile("/favicon.ico", "resources/ui/favicon.ico")

		// NoRoute
		handlers.NoRoute.Add("/", NoCache, Redirect)

		// NoMethod
		handlers.NoRoute.Add("/", NoCache, Redirect)
	}

	{
		// WebSocket
		router := router.Group("/ws", NoCache, websocket.Check)

		router.GET("/status", websocket.Status)

		// NoRoute
		handlers.NoRoute.Add("/ws", NoCache, websocket.NoRoute)

		// NoMethod
		handlers.NoMethod.Add("/ws", NoCache, websocket.NoMethod)
	}

	{
		// API
		router := router.Group("/api", NoCache)

		router.GET("/status", api.Status)

		// NoRoute
		handlers.NoRoute.Add("/api", NoCache, api.NoRoute)

		// NoMethod
		handlers.NoMethod.Add("/api", NoCache, api.NoMethod)
	}

	router.NoRoute(handlers.NoRoute.Handlers()...)
	router.NoMethod(handlers.NoMethod.Handlers()...)

	return router
}

func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, must-revalidate")
}

func SameSite(c *gin.Context) {
	if gin.Mode() == gin.DebugMode || AllowCORS {
		c.SetSameSite(http.SameSiteNoneMode)
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
}

func Redirect(c *gin.Context) {
	c.Redirect(http.StatusFound, "/#"+c.Request.RequestURI)
}
