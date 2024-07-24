package handlers

import (
	"chronokeep/results/types"
	"net/http"
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
)

var reg = regexp.MustCompile(`[-\(\)+]`)

func (h Handler) GetSmsSubscriptions(c echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.GetSmsSubscriptionsRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// Get Key
	mkey, err := database.GetKeyAndAccount(*k)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
	}
	if mkey == nil || mkey.Key == nil || mkey.Account == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key/Account Not Found", nil)
	}
	// Check for expired key
	if mkey.Key.Expired() {
		return getAPIError(c, http.StatusUnauthorized, "Expired Key", nil)
	}
	// And Event for verification of whether or not we can allow access to this key
	year := ""
	if request.Year != nil {
		year = *request.Year
	}
	// Check for host being allowed.
	if !mkey.Key.IsAllowed(c.Request().Referer()) {
		return getAPIError(c, http.StatusUnauthorized, "Host Not Allowed", nil)
	}
	mult, err := database.GetEventAndYear(request.Slug, year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	if mult.Event.AccessRestricted && mkey.Account.Identifier != mult.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Restricted Event", nil)
	}
	subs, err := database.GetSubscribedPhones(mult.EventYear.Identifier)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Subscriptions", err)
	}
	return c.JSON(http.StatusOK, types.GetSmsSubscriptionsResponse{
		Subscriptions: subs,
	})
}

func (h Handler) AddSmsSubscription(c echo.Context) error {
	var request types.AddSmsSubscriptionRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// And Event for verification of whether or not we can allow access to this key
	year := ""
	if request.Year != nil {
		year = *request.Year
	}
	mult, err := database.GetEventAndYear(request.Slug, year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	// Verify we're within the DaysAllowed period for the event.
	curTime := time.Now()
	if mult.EventYear.DateTime.AddDate(0, 0, mult.EventYear.DaysAllowed).Before(curTime) {
		return getAPIError(c, http.StatusForbidden, "Time Limit to Subscribe Exceeded", nil)
	}
	// If the event is restricted check the key, key check isn't necessary otherwise
	if mult.Event.AccessRestricted {
		// Get Key from Authorization Header
		k, err := retrieveKey(c.Request())
		if err != nil {
			return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
		}
		if k == nil {
			return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
		}
		// Get Key
		mkey, err := database.GetKeyAndAccount(*k)
		if err != nil {
			return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Key/Account", err)
		}
		if mkey == nil || mkey.Key == nil || mkey.Account == nil {
			return getAPIError(c, http.StatusUnauthorized, "Key/Account Not Found", nil)
		}
		// Check for expired key
		if mkey.Key.Expired() {
			return getAPIError(c, http.StatusUnauthorized, "Expired Key", nil)
		}
		// Check for host being allowed.
		if !mkey.Key.IsAllowed(c.Request().Referer()) {
			return getAPIError(c, http.StatusUnauthorized, "Host Not Allowed", nil)
		}
		if mkey.Account.Identifier != mult.Event.AccountIdentifier {
			return getAPIError(c, http.StatusUnauthorized, "Restricted Event", nil)
		}
	}
	if (request.Bib == nil) && (request.First == nil || request.Last == nil) {
		return getAPIError(c, http.StatusBadRequest, "No Participant Identified", nil)
	}
	var phone = reg.ReplaceAllString(request.Phone, "")
	if len(phone) != 10 {
		return getAPIError(c, http.StatusBadRequest, "Invalid Phone Number", nil)
	}
	bib := ""
	first := ""
	last := ""
	if request.Bib != nil {
		bib = *request.Bib
	}
	if request.First != nil {
		first = *request.First
	}
	if request.Last != nil {
		last = *request.Last
	}
	if len(bib+first+last) < 1 {
		return getAPIError(c, http.StatusBadRequest, "Bib or First/Last Must Be Set", nil)
	}
	err = database.AddSubscribedPhone(mult.EventYear.Identifier, types.SmsSubscription{
		Bib:   bib,
		First: first,
		Last:  last,
		Phone: request.Phone,
	})
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Subscription", err)
	}
	return c.NoContent(http.StatusOK)
}

func (h Handler) RemoveSmsSubscription(c echo.Context) error {
	var request types.RemoveSmsSubscriptionRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	// And Event for verification of whether or not we can allow access to this key
	year := ""
	if request.Year != nil {
		year = *request.Year
	}
	mult, err := database.GetEventAndYear(request.Slug, year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	var phone = reg.ReplaceAllString(request.Phone, "")
	if len(phone) != 10 {
		return getAPIError(c, http.StatusBadRequest, "Invalid Phone Number", nil)
	}
	err = database.RemoveSubscribedPhone(mult.EventYear.Identifier, phone)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Removing Subscription", err)
	}
	return c.NoContent(http.StatusOK)
}
