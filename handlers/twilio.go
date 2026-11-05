/*
Chronokeep Results - Race Results API
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
	"unicode"

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
	from := values.Get("From")
	message := values.Get("Body")
	twilio_signature := c.Request().Header.Get("X-Twilio-Signature")
	if from == "" || twilio_signature == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	// Store and compare numbers in one format so a number can be blocked and unblocked
	// no matter which format it was originally added in.
	phone, ok := canonicalPhone(from)
	if !ok {
		phone = from
	}
	params := make(map[string]string)
	for key := range values {
		params[key] = values.Get(key)
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
	words := messageWords(lowerCaseMessage)
	// check if told to unstop
	for _, keyword := range startKeywords {
		if words[keyword] {
			// remove from do not contact list
			for _, variant := range phoneVariants(phone) {
				if err = database.UnblockPhone(variant); err != nil {
					return c.NoContent(http.StatusInternalServerError)
				}
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
		if words[keyword] {
			// add to do not contact list
			err = database.AddBlockedPhone(phone)
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
	for _, blocked := range phones {
		if blocked == phone || samePhone(blocked, phone) {
			return c.NoContent(http.StatusOK)
		}
	}
	if words["help"] {
		return c.String(http.StatusOK, "Reply STOP to stop receiving texts from this number.")
	}
	return c.NoContent(http.StatusOK)
}

// messageWords splits a message into the set of words it contains so keywords are only
// matched as whole words (e.g. "send" must not match the "end" stop keyword).
func messageWords(message string) map[string]bool {
	words := make(map[string]bool)
	for _, word := range strings.FieldsFunc(message, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}) {
		words[word] = true
	}
	return words
}
