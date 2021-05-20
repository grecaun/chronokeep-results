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
	group.POST("/event-year", h.GetEventYear)
	group.PUT("/event-year/add", h.AddEventYear)
	group.PUT("/event-year/update", h.UpdateEventYear)
	group.DELETE("/event-year/delete", h.DeleteEventYear)
	// Event handlers
	group.POST("/event", h.GetEvent)
	group.POST("/event/all", h.GetEvents)
	group.PUT("/event/add", h.AddEvent)
	group.PUT("/event/update", h.UpdateEvent)
	group.DELETE("/event/delete", h.DeleteEvent)
	// Result handlers
	group.POST("/results", h.GetResults)
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
	group.PUT("/account/add", h.AddAccount)
	group.PUT("/account/update", h.UpdateAccount)
	group.DELETE("/account/delete", h.DeleteAccount)
	group.POST("/account/password", h.ChangePassword)
	group.POST("/account/email", h.ChangeEmail)
	group.GET("/account/logout", h.Logout)
	group.POST("/account/unlock", h.Unlock)
	// Key handlers
	group.GET("/key", h.GetKeys)
	group.PUT("/key/add", h.AddKey)
	group.PUT("/key/update", h.UpdateKey)
	group.DELETE("/key/delete", h.DeleteKey)
}
