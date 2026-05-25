package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/auth"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/config"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/handlers"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/services"
)

// RouterDeps holds all infrastructure dependencies needed to wire handlers.
type RouterDeps struct {
	Config *config.Config
	Pool   *pgxpool.Pool
	Redis  *redis.Client
}

// newRouter builds and returns a configured *gin.Engine.
func newRouter(deps RouterDeps) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     deps.Config.CORSAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/health", handlers.Health)

	// Build auth dependencies.
	jwtSvc, err := auth.NewJWTService(deps.Config.JWTPrivateKey, deps.Config.JWTPublicKey)
	if err != nil {
		panic("router: failed to initialise JWT service: " + err.Error())
	}
	otpSvc := auth.NewOTPService(deps.Redis)
	fbClient := auth.NewFacebookClient()

	userRepo := repositories.NewUserRepository(deps.Pool)
	tokenRepo := repositories.NewRefreshTokenRepository(deps.Pool)

	authSvc := services.NewAuthService(userRepo, tokenRepo, jwtSvc, otpSvc, fbClient)
	authHandler := handlers.NewAuthHandler(authSvc)

	v1 := r.Group("/api/v1")

	// Public auth routes — no middleware.
	authGroup := v1.Group("/auth")
	{
		reg := authGroup.Group("/register")
		reg.POST("/init", authHandler.RegisterInit)
		reg.POST("/verify-otp", authHandler.RegisterVerifyOTP)
		reg.POST("/complete", authHandler.RegisterComplete)

		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/login/facebook", authHandler.FacebookLogin)
		authGroup.POST("/refresh", authHandler.RefreshToken)

		pwReset := authGroup.Group("/password-reset")
		pwReset.POST("/request", authHandler.PasswordResetRequest)
		pwReset.POST("/verify-otp", authHandler.PasswordResetVerifyOTP)
		pwReset.POST("/complete", authHandler.PasswordResetComplete)
	}

	// Protected auth route — requires valid access token.
	protected := v1.Group("", auth.Middleware(jwtSvc))
	{
		protected.POST("/auth/logout", authHandler.Logout)
	}

	return r
}
