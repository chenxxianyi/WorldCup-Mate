package routes

import (
	"time"

	"worldcup-mate/internal/handlers"
	"worldcup-mate/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.JSONLogger())
	r.Use(middleware.Recover())
	r.Use(middleware.CORS())

	// Health checks (OBS-04).
	r.GET("/health/live", handlers.HealthLive)
	r.GET("/health/ready", handlers.HealthReady)

	api := r.Group("/api")

	// Auth (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", middleware.RateLimit("rl:register:ip", 5, time.Minute), handlers.Register)
		auth.POST("/login", middleware.RateLimit("rl:login:ip", 20, time.Minute), handlers.Login)
		auth.POST("/refresh", handlers.Refresh)
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
		matches.GET("/progress", handlers.GetTournamentProgress)
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

	// Home aggregation (DATA-10).
	api.GET("/home", handlers.HomeAggregate)

	// Competitions (public)
	api.GET("/competitions", handlers.ListCompetitions)
	api.GET("/competitions/:code/overview", handlers.CompetitionOverview)
	api.GET("/competitions/:code/standings", handlers.GetCompetitionStandings)
	api.GET("/featured", handlers.PublicFeatured)

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
	}

	// Admin routes
	admin := api.Group("/admin")
	{
		admin.POST("/login", handlers.AdminLogin)

		adminAuth := admin.Group("")
		adminAuth.Use(middleware.JWTAuth(), middleware.AdminAuth())
		{
			adminAuth.GET("/dashboard", handlers.AdminDashboard)

			adminAuth.GET("/competitions", handlers.AdminListCompetitions)
			adminAuth.POST("/competitions", handlers.AdminCreateCompetition)
			adminAuth.PUT("/competitions/:id", handlers.AdminUpdateCompetition)

			adminAuth.GET("/featured", handlers.AdminListFeatured)
			adminAuth.GET("/featured/:code/matches", handlers.AdminFeaturedMatches)
			adminAuth.PUT("/featured/:code", handlers.AdminUpsertFeatured)

			adminAuth.GET("/teams", handlers.AdminListTeams)
			adminAuth.POST("/teams", handlers.AdminCreateTeam)
			adminAuth.PUT("/teams/:id", handlers.AdminUpdateTeam)
			adminAuth.DELETE("/teams/:id", handlers.AdminDeleteTeam)

			adminAuth.GET("/groups", handlers.AdminListGroups)
			adminAuth.POST("/groups", handlers.AdminCreateGroup)
			adminAuth.PUT("/groups/:id", handlers.AdminUpdateGroup)
			adminAuth.DELETE("/groups/:id", handlers.AdminDeleteGroup)

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
			adminAuth.POST("/matches/import", handlers.AdminImportMatches)
			adminAuth.POST("/sync/matches", handlers.AdminSyncMatches)

			adminAuth.GET("/standings", handlers.AdminListStandings)
			adminAuth.POST("/standings/recalculate", handlers.AdminRecalculateStandings)
			adminAuth.POST("/standings/league/recalculate", handlers.AdminRecalculateLeagueStanding)
			adminAuth.PUT("/standings/:id", handlers.AdminUpdateStanding)

			adminAuth.GET("/users", handlers.AdminListUsers)
			adminAuth.PUT("/users/:id/status", handlers.AdminUpdateUserStatus)

			adminAuth.GET("/sync/history", handlers.AdminListSyncHistory)
			adminAuth.GET("/reminders", handlers.AdminListReminders)
			adminAuth.POST("/reminders/:id/retry", handlers.AdminRetryReminder)
		}
	}

	return r
}
