package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"zhanbu/config"
	"zhanbu/internal/handler"
	"zhanbu/internal/middleware"
	"zhanbu/internal/repository"
	"zhanbu/internal/service"
	"zhanbu/pkg/utils"
	"gorm.io/gorm"
)

// SetupRouter initializes all routes and returns a Gin engine.
func SetupRouter(db *gorm.DB, cfg *config.Config, logger zerolog.Logger) *gin.Engine {
	r := gin.New()

	// Set mode
	gin.SetMode(cfg.Server.Mode)

	// Global middleware
	r.Use(middleware.CORSMiddleware(&cfg.CORS))
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(gin.Recovery())

	// Rate limiting
	r.Use(middleware.RateLimitMiddleware(cfg.RateLimit.APIPerMinute, 60))

	// Initialize dependencies
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	userRepo := repository.NewUserRepository(db)
	verRepo := repository.NewVerificationRepository(db)
	emailService := service.NewEmailService(verRepo, &cfg.SMTP, &cfg.Security)
	authService := service.NewAuthService(userRepo, jwtManager, &cfg.JWT, emailService, &cfg.Security, &cfg.SMTP)
	authHandler := handler.NewAuthHandler(authService)
	authMiddleware := middleware.AuthMiddleware(jwtManager)

	// Health check
	healthHandler := handler.NewHealthHandler("1.0.0")
	healthHandler.SetMode(cfg.Server.Mode)
	r.GET("/api/health", healthHandler.Check)

	// Auth routes (public)
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/verify-email", authHandler.VerifyEmail)
		authGroup.POST("/resend-verification", authHandler.ResendVerification)
	}

	// Auth routes (authenticated)
	authProtected := r.Group("/api/auth")
	authProtected.Use(authMiddleware)
	{
		authProtected.POST("/refresh", authHandler.RefreshToken)
		authProtected.GET("/profile", authHandler.GetProfile)
		authProtected.PUT("/profile", authHandler.UpdateProfile)
	}

	// Divination record repository (shared)
	divinationRepo := repository.NewDivinationRepository(db)

	// Tarot routes
	tarotRepo := repository.NewTarotRepository(db)
	tarotService := service.NewTarotService(tarotRepo)
	tarotService.SetRecordSaver(divinationRepo)
	tarotHandler := handler.NewTarotHandler(tarotService)
	tarotGroup := r.Group("/api/tarot")
	{
		tarotGroup.GET("/cards", tarotHandler.GetCards)
		tarotGroup.GET("/cards/:id", tarotHandler.GetCardByID)
		tarotGroup.GET("/spreads", tarotHandler.GetSpreads)
		tarotGroup.POST("/draw", authMiddleware, tarotHandler.Draw)
	}

	// Horoscope routes
	horoscopeTemplates, err := service.NewJSONTemplateStore("data/horoscope")
	if err != nil {
		logger.Error().Err(err).Msg("failed to load horoscope templates")
		// Fallback: create with nil templates (will return empty strings)
	}
	horoscopeService := service.NewHoroscopeService(horoscopeTemplates)
	horoscopeHandler := handler.NewHoroscopeHandler(horoscopeService)
	horoscopeGroup := r.Group("/api/horoscope")
	{
		horoscopeGroup.GET("/:zodiac", horoscopeHandler.GetHoroscope)
	}

	// AI interpretation routes (authenticated + rate limited)
	var aiProvider service.AIProvider
	if cfg.AI.APIKey != "" {
		var providerErr error
		aiProvider, providerErr = service.NewOpenAIProvider(&cfg.AI)
		if providerErr != nil {
			logger.Error().Err(providerErr).Msg("failed to initialize AI provider")
		}
	} else {
		logger.Warn().Msg("AI provider disabled: ZHANBU_AI_API_KEY is not set")
	}

	aiService := service.NewAIService(aiProvider, divinationRepo)
	aiHandler := handler.NewAIHandler(aiService)
	aiGroup := r.Group("/api/ai")
	aiGroup.Use(authMiddleware)
	{
		aiGroup.POST("/interpret", aiHandler.Interpret)
	}

	// LiuYao (六爻) routes
	liuyaoService := service.NewLiuYaoService()
	liuyaoHandler := handler.NewLiuYaoHandler(liuyaoService, divinationRepo)
	liuyaoGroup := r.Group("/api/liuyao")
	{
		liuyaoGroup.GET("/hexagrams", liuyaoHandler.GetHexagrams)
		liuyaoGroup.POST("/throw", authMiddleware, liuyaoHandler.Throw)
	}

	// LiuYao v2 (高岛易断) routes
	liuyaoV2Service, err := service.NewLiuYaoV2Service(cfg.LiuYao.Method)
	if err != nil {
		logger.Error().Err(err).Msg("failed to initialize LiuYao v2 service")
	} else {
		liuyaoV2Handler := handler.NewLiuYaoV2Handler(liuyaoV2Service, divinationRepo)
		liuyaoV2Group := r.Group("/api/liuyao/v2")
		{
			liuyaoV2Group.GET("/hexagrams", liuyaoV2Handler.GetHexagrams)
			liuyaoV2Group.GET("/hexagrams/:id", liuyaoV2Handler.GetHexagramByID)
			liuyaoV2Group.GET("/config", liuyaoV2Handler.GetConfig)
			liuyaoV2Group.POST("/throw", authMiddleware, liuyaoV2Handler.Throw)
		}
	}

	// BaZi (八字) routes
	baziService := service.NewBaZiService()
	baziHandler := handler.NewBaZiHandler(baziService, divinationRepo)
	baziGroup := r.Group("/api/bazi")
	{
		baziGroup.POST("/calculate", authMiddleware, baziHandler.Calculate)
	}

	// History routes (authenticated)
	historyService := service.NewHistoryService(divinationRepo)
	historyHandler := handler.NewHistoryHandler(historyService)
	historyGroup := r.Group("/api/history")
	historyGroup.Use(authMiddleware)
	{
		historyGroup.GET("", historyHandler.List)
		historyGroup.GET(":id", historyHandler.GetDetail)
		historyGroup.DELETE(":id", historyHandler.Delete)
	}

	// Email verification routes (if SMTP enabled)
	logger.Info().Bool("smtp_enabled", cfg.SMTP.Enabled).Msg("checking SMTP config")
	if cfg.SMTP.Enabled {
		emailHandler := handler.NewEmailHandler(emailService)
		emailGroup := r.Group("/api/email")
		{
			emailGroup.POST("/send-code", emailHandler.SendCode)
			emailGroup.POST("/verify-code", emailHandler.VerifyCode)
		}
	}

	return r
}
