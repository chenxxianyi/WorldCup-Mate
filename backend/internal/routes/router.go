package routes

import (
	"time"

	"worldcup-mate/internal/handlers"
	"worldcup-mate/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.Use(middleware.CORS())
	r.Use(middleware.SecurityHeaders())

	api := r.Group("/api")

	// Auth (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", middleware.RateLimit("auth_register", 10, time.Hour), handlers.Register)
		auth.POST("/login", middleware.RateLimit("auth_login", 10, 5*time.Minute), handlers.Login)
		auth.POST("/logout", handlers.Logout)
	}

	// Matches (public)
	matches := api.Group("/matches")
	{
		matches.GET("", handlers.ListMatches)
		matches.GET("/today", handlers.GetTodayMatches)
		matches.GET("/tomorrow", handlers.GetTomorrowMatches)
		matches.GET("/upcoming", handlers.GetUpcomingMatches)
		matches.GET("/live", handlers.GetLiveMatches)
		matches.GET("/recommended", handlers.GetRecommendedMatches)
		matches.GET("/timeline", handlers.GetTimeline)
		matches.GET("/progress", handlers.GetTournamentProgress)
		matches.GET("/:id/lineups", handlers.GetMatchLineups)
		matches.GET("/:id", handlers.GetMatchDetail)
		matches.GET("/by-team/:teamId", handlers.GetMatchesByTeam)
		matches.GET("/by-group/:groupId", handlers.GetMatchesByGroup)
		matches.GET("/by-stage/:stage", handlers.GetMatchesByStage)
	}

	// Teams (public)
	teams := api.Group("/teams")
	{
		teams.GET("", handlers.ListTeams)
		teams.GET("/:id", handlers.GetTeamDetail)
		teams.GET("/:id/matches", handlers.GetTeamMatches)
		teams.GET("/:id/players", handlers.GetTeamPlayers)
	}

	// Groups (public)
	groups := api.Group("/groups")
	{
		groups.GET("", handlers.ListGroups)
		groups.GET("/:id", handlers.GetGroupDetail)
		groups.GET("/:id/standings", handlers.GetGroupStandings)
	}

	// Standings (public)
	standings := api.Group("/standings")
	{
		standings.GET("", handlers.ListStandings)
		standings.GET("/best-third", handlers.GetBestThird)
	}

	// Cities & Stadiums (public)
	api.GET("/cities", handlers.ListCities)
	api.GET("/stadiums", handlers.ListStadiums)
	api.GET("/stadiums/:id", handlers.GetStadiumDetail)
	api.GET("/sync/status", handlers.GetSyncStatus)

	// Post-match summary read is public, but must not trigger AI generation.
	matches.GET("/:id/post-match-summary", handlers.GetPostMatchSummary)

	// Authenticated routes
	authRequired := api.Group("")
	authRequired.Use(middleware.JWTAuth())
	{
		authRequired.GET("/user/profile", handlers.GetProfile)
		authRequired.PUT("/user/profile", handlers.UpdateProfile)
		authRequired.PUT("/user/password", handlers.ChangePassword)
		authRequired.POST("/user/avatar", handlers.UploadAvatar)

		// Favorites
		fav := authRequired.Group("/favorites")
		{
			fav.POST("/teams/:teamId", handlers.AddFavoriteTeam)
			fav.DELETE("/teams/:teamId", handlers.RemoveFavoriteTeam)
			fav.GET("/teams", handlers.ListFavoriteTeams)
			fav.POST("/matches/:matchId", handlers.AddFavoriteMatch)
			fav.DELETE("/matches/:matchId", handlers.RemoveFavoriteMatch)
			fav.GET("/matches", handlers.ListFavoriteMatches)
		}

		// Reminders
		rem := authRequired.Group("/reminders")
		{
			rem.POST("/batch", handlers.CreateReminderBatch)
			rem.POST("", handlers.CreateReminder)
			rem.GET("", handlers.ListReminders)
			rem.PUT("/:id", handlers.UpdateReminder)
			rem.DELETE("/:id", handlers.DeleteReminder)
		}

		// Notifications
		notif := authRequired.Group("/notifications")
		{
			notif.GET("", handlers.ListNotifications)
			notif.GET("/unread-count", handlers.CountUnreadNotifications)
			notif.PUT("/:id/read", handlers.MarkNotificationRead)
			notif.PUT("/read-all", handlers.MarkAllNotificationsRead)
		}

		// AI chat history
		aiAuth := authRequired.Group("/ai")
		aiAuth.Use(middleware.RateLimit("ai", 30, time.Hour))
		{
			aiAuth.POST("/match-insight", handlers.AIMatchInsight)
			aiAuth.POST("/today-recommendations", handlers.AITodayRecommendations)
			aiAuth.POST("/group-analysis", handlers.AIGroupAnalysis)
			aiAuth.POST("/explain", handlers.AIExplain)
			aiAuth.POST("/share-copy", handlers.AIShareCopy)
			aiAuth.POST("/chat", handlers.AIChat)
			aiAuth.POST("/chat/stream", handlers.AIChatStream)
			aiAuth.GET("/conversations", handlers.AIListConversations)
			aiAuth.GET("/conversations/:id", handlers.AIGetConversation)
			aiAuth.DELETE("/conversations/:id", handlers.AIDeleteConversation)
		}

		authRequired.POST("/matches/:id/post-match-summary/generate", middleware.RateLimit("ai_post_match_summary", 10, time.Hour), handlers.GeneratePostMatchSummary)
	}

	// Admin routes
	admin := api.Group("/admin")
	{
		admin.POST("/login", middleware.RateLimit("admin_login", 5, 10*time.Minute), handlers.AdminLogin)

		adminAuth := admin.Group("")
		adminAuth.Use(middleware.JWTAuth(), middleware.AdminAuth())
		{
			adminAuth.GET("/dashboard", handlers.AdminDashboard)

			adminAuth.GET("/teams", handlers.AdminListTeams)
			adminAuth.POST("/teams", handlers.AdminCreateTeam)
			adminAuth.PUT("/teams/:id", handlers.AdminUpdateTeam)
			adminAuth.DELETE("/teams/:id", handlers.AdminDeleteTeam)
			adminAuth.GET("/teams/:id/player-mapping", handlers.AdminGetTeamPlayerMapping)
			adminAuth.PUT("/teams/:id/player-mapping", handlers.AdminUpsertTeamPlayerMapping)
			adminAuth.POST("/teams/:id/sync-players", handlers.AdminSyncTeamPlayers)

			adminAuth.GET("/groups", handlers.AdminListGroups)
			adminAuth.POST("/groups", handlers.AdminCreateGroup)
			adminAuth.PUT("/groups/:id", handlers.AdminUpdateGroup)

			adminAuth.GET("/cities", handlers.AdminListCities)
			adminAuth.POST("/cities", handlers.AdminCreateCity)
			adminAuth.PUT("/cities/:id", handlers.AdminUpdateCity)
			adminAuth.DELETE("/cities/:id", handlers.AdminDeleteCity)

			adminAuth.GET("/stadiums", handlers.AdminListStadiums)
			adminAuth.POST("/stadiums", handlers.AdminCreateStadium)
			adminAuth.PUT("/stadiums/:id", handlers.AdminUpdateStadium)
			adminAuth.DELETE("/stadiums/:id", handlers.AdminDeleteStadium)

			adminAuth.GET("/matches", handlers.AdminListMatches)
			adminAuth.POST("/matches", handlers.AdminCreateMatch)
			adminAuth.PUT("/matches/:id", handlers.AdminUpdateMatch)
			adminAuth.DELETE("/matches/:id", handlers.AdminDeleteMatch)
			adminAuth.PUT("/matches/:id/score", handlers.AdminUpdateMatchScore)
			adminAuth.PUT("/matches/:id/status", handlers.AdminUpdateMatchStatus)
			adminAuth.POST("/matches/:id/sync-lineups", handlers.AdminSyncMatchLineups)
			adminAuth.GET("/matches/:id/external-match-mapping", handlers.AdminGetMatchExternalMapping)
			adminAuth.PUT("/matches/:id/external-match-mapping", handlers.AdminUpsertMatchExternalMapping)
			adminAuth.POST("/matches/import", handlers.AdminImportMatches)
			adminAuth.POST("/sync/matches", handlers.AdminSyncMatches)
			adminAuth.POST("/sync/players", handlers.AdminSyncAllPlayers)
			adminAuth.POST("/sync/lineups/live-window", handlers.AdminSyncLiveWindowLineups)

			adminAuth.GET("/standings", handlers.AdminListStandings)
			adminAuth.POST("/standings/recalculate", handlers.AdminRecalculateStandings)
			adminAuth.PUT("/standings/:id", handlers.AdminUpdateStanding)

			adminAuth.GET("/users", handlers.AdminListUsers)
			adminAuth.PUT("/users/:id/status", handlers.AdminUpdateUserStatus)
		}
	}

	return r
}
