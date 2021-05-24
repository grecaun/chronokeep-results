package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Handler Struct for using methods for handling information.
type Handler struct {
	validate *validator.Validate
}

func (h Handler) Bind(group *echo.Group) {
	// Event Year handlers
	group.GET("/event-year", h.GetEventYear)
	group.GET("/event-year/event", h.GetEventYears)
	group.POST("/event-year/add", h.AddEventYear)
	group.PUT("/event-year/update", h.UpdateEventYear)
	group.DELETE("/event-year/delete", h.DeleteEventYear)
	// Event handlers
	group.GET("/event", h.GetEvent)
	group.GET("/event/all", h.GetEvents)
	group.GET("/event/my", h.GetMyEvents)
	group.POST("/event/add", h.AddEvent)
	group.PUT("/event/update", h.UpdateEvent)
	group.DELETE("/event/delete", h.DeleteEvent)
	// Result handlers
	group.GET("/results", h.GetResults)
	group.PUT("/results/add", h.AddResults)
	group.DELETE("/results/delete", h.DeleteResults)
	// Account Login
	group.POST("/account/login", h.Login)
	group.POST("/account/refresh", h.Refresh)
}

func (h Handler) BindRestricted(group *echo.Group) {
	// Account handlers
	group.GET("/account", h.GetAccount)
	group.GET("/account/all", h.GetAccounts)
	group.GET("/account/logout", h.Logout)
	group.POST("/account/add", h.AddAccount)
	group.PUT("/account/update", h.UpdateAccount)
	group.PUT("/account/password", h.ChangePassword)
	group.PUT("/account/email", h.ChangeEmail)
	group.POST("/account/unlock", h.Unlock)
	group.DELETE("/account/delete", h.DeleteAccount)
	// Key handlers
	group.GET("/key", h.GetKeys)
	group.POST("/key/add", h.AddKey)
	group.PUT("/key/update", h.UpdateKey)
	group.DELETE("/key/delete", h.DeleteKey)
	// Event handlers (restricted for website use)
	group.GET("/r/event", h.RGetEvents)
	group.POST("/r/event/add", h.RAddEvent)
	group.PUT("/r/event/update", h.RUpdateEvent)
	group.DELETE("/r/event/delete", h.RDeleteEvent)
	// Event Year handlers (restricted for website use)
	group.GET("/r/event-year", h.RGetEventYears)
	group.POST("/r/event-year/add", h.RAddEventYear)
	group.PUT("/r/event-year/update", h.RUpdateEventYear)
	group.DELETE("/r/event-year/delete", h.RDeleteEventYear)
}
