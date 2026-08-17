/*
Chronokeep Desktop - Race Scoring Software
Copyright (C) 2026 James Sentinella

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package handlers

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v5"
	log "github.com/sirupsen/logrus"
)

var (
	stopKeywords = []string{
		"stop",
		"stopall",
		"unsubscribe",
		"cancel",
		"end",
		"quit",
	}
	startKeywords = []string{
		"start",
		"unstop",
	}
)

func (h Handler) Twilio(c *echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	from := values["From"][0]
	message := values["Body"][0]
	twilio_signature := c.Request().Header["X-Twilio-Signature"][0]
	params := make(map[string]string)
	for key, pArray := range values {
		params[key] = pArray[0]
	}
	valid_request := twilioRequestValidator.Validate(config.TwilioResponseWebhookURL, params, twilio_signature)
	log.WithFields(log.Fields{
		"params":           params,
		"from":             from,
		"body":             message,
		"twilio signature": twilio_signature,
		"valid_request":    valid_request,
	}).Info("Values found.")
	if !valid_request {
		return c.NoContent(http.StatusUnauthorized)
	}
	lowerCaseMessage := strings.ToLower(strings.TrimSpace(message))
	// check if told to unstop
	// do this before the stop check because stop is a substring of unstop
	for _, keyword := range startKeywords {
		if strings.Contains(lowerCaseMessage, keyword) {
			// remove from do not contact list
			err = database.UnblockPhone(from)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
			// check if it's already being handled by twilio
			// and don't send message if so
			if lowerCaseMessage == keyword {
				return c.NoContent(http.StatusOK)
			}
			// send success message
			return c.String(http.StatusOK, "You have successfully been re-subscribed to messages from this number. Reply STOP to unsubscribe. Msg&Data Rates May Apply.")
		}
	}
	// check if told to stop
	for _, keyword := range stopKeywords {
		if strings.Contains(lowerCaseMessage, keyword) {
			// add to do not contact list
			err = database.AddBlockedPhone(from)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}
			// check if they're already being handled by twilio
			// and don't send message if so
			if lowerCaseMessage == keyword {
				return c.NoContent(http.StatusOK)
			}
			// send success message
			return c.String(http.StatusOK, "You have been successfully unsubscribed. You will not receive any more messages from this number. Reply START to resubscribe.")
		}
	}
	// check if the phone number is on the do not call list
	phones, err := database.GetBlockedPhones()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	for _, phone := range phones {
		if strings.Contains(from, phone) {
			return c.NoContent(http.StatusOK)
		}
	}
	if strings.Contains(lowerCaseMessage, "help") {
		return c.String(http.StatusOK, "Reply STOP to stop receiving texts from this number.")
	}
	return c.NoContent(http.StatusOK)
}

