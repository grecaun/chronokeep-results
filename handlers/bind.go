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
	group.POST("/event-year/event", h.GetEventYears)
	group.POST("/event-year/add", h.AddEventYear)
	group.PUT("/event-year/update", h.UpdateEventYear)
	group.DELETE("/event-year/delete", h.DeleteEventYear)
	// Event handlers
	group.POST("/event", h.GetEvent)
	group.GET("/event/all", h.GetEvents)
	group.GET("/event/my", h.GetMyEvents)
	group.POST("/event/add", h.AddEvent)
	group.PUT("/event/update", h.UpdateEvent)
	group.DELETE("/event/delete", h.DeleteEvent)
	// Result handlers
	group.POST("/results", h.GetResults)
	group.POST("/results/all", h.GetAllResults)
	group.POST("/results/finish", h.GetFinishResults)
	group.POST("/results/bib", h.GetBibResults)
	group.POST("/results/add", h.AddResults)
	group.DELETE("/results/delete", h.DeleteResults)
	// Participants handlers
	group.POST("/participants", h.GetParticipants)
	group.POST("/participants/add", h.AddParticipants)
	group.DELETE("/participants/delete", h.DeleteParticipants)
	// BibChip handlers
	group.POST("/bibchips", h.GetBibChips)
	group.POST("/bibchips/add", h.AddBibChips)
	group.DELETE("/bibchips/delete", h.DeleteBibChips)
	// Account Login
	group.POST("/account/login", h.Login)
	group.POST("/account/refresh", h.Refresh)
	// Blocked/banned emails/phone numbers
	group.POST("/blocked/phones/add", h.AddBannedPhone)
	group.GET("/blocked/phones/get", h.GetBannedPhones)
	group.POST("/blocked/emails/add", h.AddBannedEmail)
	group.GET("/blocked/emails/get", h.GetBannedEmails)
	group.POST("/blocked/emails/unblock", h.RemoveBannedEmail)
	// Certificate image
	group.GET("/certificate/:name/:event/:time/:date", h.GetCertificate)
	// SMS Subscriptions
	group.POST("/sms", h.GetSmsSubscriptions)
	group.POST("/sms/add", h.AddSmsSubscription)
	group.POST("/sms/remove", h.RemoveSmsSubscription)
	// Segments
	group.POST("/segments", h.GetSegments)
	group.POST("/segments/add", h.AddSegments)
	group.DELETE("/segments/delete", h.DeleteSegments)
}

func (h Handler) BindRestricted(group *echo.Group) {
	// Account handlers
	group.POST("/account", h.GetAccount)
	group.GET("/account/all", h.GetAccounts)
	group.POST("/account/logout", h.Logout)
	group.POST("/account/add", h.AddAccount)
	group.PUT("/account/update", h.UpdateAccount)
	group.PUT("/account/password", h.ChangePassword)
	group.PUT("/account/email", h.ChangeEmail)
	group.POST("/account/unlock", h.Unlock)
	group.DELETE("/account/delete", h.DeleteAccount)
	group.PUT("/account/link", h.LinkAccounts)
	group.PUT("/account/unlink", h.UnlinkAccounts)
	// Key handlers
	group.POST("/key", h.GetKeys)
	group.POST("/key/add", h.AddKey)
	group.PUT("/key/update", h.UpdateKey)
	group.DELETE("/key/delete", h.DeleteKey)
	// Event handlers (restricted for website use)
	group.POST("/r/event", h.RGetEvents)
	group.POST("/r/event/add", h.RAddEvent)
	group.PUT("/r/event/update", h.RUpdateEvent)
	group.DELETE("/r/event/delete", h.RDeleteEvent)
	// Event Year handlers (restricted for website use)
	group.POST("/r/event-year", h.RGetEventYears)
	group.POST("/r/event-year/add", h.RAddEventYear)
	group.PUT("/r/event-year/update", h.RUpdateEventYear)
	group.DELETE("/r/event-year/delete", h.RDeleteEventYear)
	// Participants handlers
	group.POST("/r/participants", h.RGetParticipants)
	group.POST("/r/participants/add", h.RAddParticipant)
	group.POST("/r/participants/add-many", h.RAddManyParticipants)
	group.DELETE("/r/participants/delete", h.RDeleteParticipants)
	group.POST("r/participants/update", h.RUpdateParticipant)
	group.POST("r/participants/update-many", h.RUpdateManyParticipants)
	// Unblock phone should be restricted to admins only
	group.POST("/blocked/phones/unblock", h.RemoveBannedPhone)
}
