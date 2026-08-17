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
	"chronokeep/results/types"
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h Handler) GetMultiResults(c *echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.GetMultiResultsRequest
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
	// Check for host being allowed.
	if !mkey.Key.IsAllowed(c.Request().Referer()) {
		return getAPIError(c, http.StatusUnauthorized, "Host Not Allowed", nil)
	}
	// Check that at least one year was specified.
	if len(request.Years) < 1 {
		return getAPIError(c, http.StatusBadRequest, "No Years Given", nil)
	}
	// Verify key is allowed to access the event.
	event, err := database.GetEvent(request.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event", err)
	}
	if event == nil {
		return getAPIError(c, http.StatusNotFound, "Event Not Found", nil)
	}
	if event.AccessRestricted && mkey.Account.Identifier != event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Restricted Event", nil)
	}
	outRes := make(map[string]map[string][]types.Result)
	for _, year := range request.Years {
		eYear, err := database.GetEventYear(event.Slug, year)
		if err != nil {
			return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Year", err)
		}
		if _, ok := outRes[year]; !ok {
			outRes[year] = make(map[string][]types.Result)
		}
		if eYear != nil {
			results, err := database.GetResults(eYear.Identifier, 0, 0)
			if err != nil {
				return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Results", err)
			}
			for _, res := range results {
				if _, ok := outRes[year][res.Distance]; !ok {
					outRes[year][res.Distance] = make([]types.Result, 0, 1)
				}
				outRes[year][res.Distance] = append(outRes[year][res.Distance], res)
			}
		}
	}
	return c.JSON(http.StatusOK, types.GetMultiResultsResponse{
		Event:   *event,
		Results: outRes,
	})
}

