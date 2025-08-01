package handlers

import (
	"chronokeep/results/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetResults(c echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.GetResultsRequest
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
	if mult.Event.AccessRestricted && mkey.Account.Identifier != mult.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Restricted Event", nil)
	}
	years, err := database.GetEventYears(request.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event Years", err)
	}
	distance := ""
	limit := 0
	page := 0
	if request.Distance != nil {
		distance = *request.Distance
	}
	if request.Limit != nil {
		limit = *request.Limit
	}
	if request.Page != nil {
		page = *request.Page
		if page > 0 {
			page--
		}
	}
	results, err := database.GetDistanceResults(mult.EventYear.Identifier, distance, limit, page)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Results", err)
	}
	parts, err := database.GetParticipants(mult.EventYear.Identifier, 0, 0, nil)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Participants", err)
	}
	if request.Version != nil && *request.Version == 1 {
		outRes := make(map[string][]types.ResultVers1)
		for _, result := range results {
			if _, ok := outRes[result.Distance]; !ok {
				outRes[result.Distance] = make([]types.ResultVers1, 0, 1)
			}
			outRes[result.Distance] = append(outRes[result.Distance], result.ConvertToVers1())
		}
		outYears := make([]types.EventYearVers1, 0, 1)
		for _, year := range years {
			outYears = append(outYears, year.ConvertToVers1())
		}
		return c.JSON(http.StatusOK, types.GetResultsResponseVers1{
			Event:   mult.Event.ConvertToVers1(),
			Years:   outYears,
			Results: outRes,
			Count:   len(results),
		})
	}
	outRes := make(map[string][]types.Result)
	for _, result := range results {
		if _, ok := outRes[result.Distance]; !ok {
			outRes[result.Distance] = make([]types.Result, 0, 1)
		}
		outRes[result.Distance] = append(outRes[result.Distance], result)
	}
	outParts := make([]types.ResultParticipant, 0)
	for _, part := range parts {
		if !part.Anonymous {
			outParts = append(outParts, part.ToResultParticipant())
		}
	}
	distances, err := database.GetDistances(mult.EventYear.Identifier)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Distances", err)
	}
	if distances != nil && len(distances) == 0 {
		distances = nil
	}
	return c.JSON(http.StatusOK, types.GetResultsResponse{
		Event:        *mult.Event,
		EventYear:    *mult.EventYear,
		Years:        years,
		Results:      outRes,
		Count:        len(results),
		Participants: outParts,
		Distances:    distances,
	})
}

func (h Handler) GetFinishResults(c echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.GetResultsRequest
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
	if mult.Event.AccessRestricted && mkey.Account.Identifier != mult.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Restricted Event", nil)
	}
	years, err := database.GetEventYears(request.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event Years", err)
	}
	distance := ""
	limit := 0
	page := 0
	if request.Distance != nil {
		distance = *request.Distance
	}
	if request.Limit != nil {
		limit = *request.Limit
	}
	if request.Page != nil {
		page = *request.Page
		if page > 0 {
			page--
		}
	}
	results, err := database.GetFinishResults(mult.EventYear.Identifier, distance, limit, page)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Results", err)
	}
	if request.Version != nil && *request.Version == 1 {
		outRes := make(map[string][]types.ResultVers1)
		for _, result := range results {
			if _, ok := outRes[result.Distance]; !ok {
				outRes[result.Distance] = make([]types.ResultVers1, 0, 1)
			}
			outRes[result.Distance] = append(outRes[result.Distance], result.ConvertToVers1())
		}
		outYears := make([]types.EventYearVers1, 0, 1)
		for _, year := range years {
			outYears = append(outYears, year.ConvertToVers1())
		}
		return c.JSON(http.StatusOK, types.GetResultsResponseVers1{
			Event:   mult.Event.ConvertToVers1(),
			Years:   outYears,
			Results: outRes,
			Count:   len(results),
		})
	}
	outRes := make(map[string][]types.Result)
	for _, result := range results {
		if _, ok := outRes[result.Distance]; !ok {
			outRes[result.Distance] = make([]types.Result, 0, 1)
		}
		outRes[result.Distance] = append(outRes[result.Distance], result)
	}
	return c.JSON(http.StatusOK, types.GetResultsResponse{
		Event:     *mult.Event,
		EventYear: *mult.EventYear,
		Years:     years,
		Results:   outRes,
		Count:     len(results),
	})
}

func (h Handler) GetAllResults(c echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.GetResultsRequest
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
	if mult.Event.AccessRestricted && mkey.Account.Identifier != mult.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Restricted Event", nil)
	}
	years, err := database.GetEventYears(request.Slug)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event Years", err)
	}
	distance := ""
	limit := 0
	page := 0
	if request.Distance != nil {
		distance = *request.Distance
	}
	if request.Limit != nil {
		limit = *request.Limit
	}
	if request.Page != nil {
		page = *request.Page
		if page > 0 {
			page--
		}
	}
	results, err := database.GetAllDistanceResults(mult.EventYear.Identifier, distance, limit, page)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Results", err)
	}
	if request.Version != nil && *request.Version == 1 {
		outRes := make(map[string][]types.ResultVers1)
		for _, result := range results {
			if _, ok := outRes[result.Distance]; !ok {
				outRes[result.Distance] = make([]types.ResultVers1, 0, 1)
			}
			outRes[result.Distance] = append(outRes[result.Distance], result.ConvertToVers1())
		}
		outYears := make([]types.EventYearVers1, 0, 1)
		for _, year := range years {
			outYears = append(outYears, year.ConvertToVers1())
		}
		return c.JSON(http.StatusOK, types.GetResultsResponseVers1{
			Event:   mult.Event.ConvertToVers1(),
			Years:   outYears,
			Results: outRes,
			Count:   len(results),
		})
	}
	outRes := make(map[string][]types.Result)
	for _, result := range results {
		if _, ok := outRes[result.Distance]; !ok {
			outRes[result.Distance] = make([]types.Result, 0, 1)
		}
		outRes[result.Distance] = append(outRes[result.Distance], result)
	}
	return c.JSON(http.StatusOK, types.GetResultsResponse{
		Event:     *mult.Event,
		EventYear: *mult.EventYear,
		Years:     years,
		Results:   outRes,
		Count:     len(results),
	})
}

func (h Handler) GetBibResults(c echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.GetBibResultsRequest
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
	// And Event for verification of whether or not we can allow access to this key
	mult, err := database.GetEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	if mult.Event.AccessRestricted && mkey.Account.Identifier != mult.Event.AccountIdentifier {
		return getAPIError(c, http.StatusUnauthorized, "Restricted Event", nil)
	}
	results, err := database.GetBibResults(mult.EventYear.Identifier, request.Bib)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Results", err)
	}
	person, err := database.GetPerson(request.Slug, request.Year, request.Bib)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Person", err)
	}
	if person == nil {
		return getAPIError(c, http.StatusNotFound, "Person Not Found", nil)
	}
	segments, err := database.GetDistanceSegments(mult.EventYear.Identifier, person.Distance)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Fetching Segments", nil)
	}
	distance, err := database.GetDistance(mult.EventYear.Identifier, person.Distance)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Fetching Distance", nil)
	}
	return c.JSON(http.StatusOK, types.GetBibResultsResponse{
		Event:          *mult.Event,
		EventYear:      *mult.EventYear,
		Results:        results,
		Person:         person,
		SingleDistance: *mult.DistanceCount == 1,
		Segments:       segments,
		Distance:       distance,
	})
}

func (h Handler) AddResults(c echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.AddResultsRequest
	if err := c.Bind(&request); err != nil {
		return getAPIError(c, http.StatusBadRequest, "Invalid Request Body", err)
	}
	var resToAdd []types.Result
	for _, res := range request.Results {
		// Validate all results, only add the results that pass validation.
		if err := res.Validate(h.validate); err == nil {
			// we want seconds to be high if the type is DNF
			if res.Type == 3 || res.Type == 30 {
				res.Seconds = 1000000
			}
			resToAdd = append(resToAdd, res)
		}
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
	if mkey.Key.Type == "read" {
		return getAPIError(c, http.StatusUnauthorized, "Key is ReadOnly", nil)
	}
	// And Event for verification of whether or not we can allow access to this key
	mult, err := database.GetEventAndYear(request.Slug, request.Year)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Retrieving Event/Year", err)
	}
	if mult == nil || mult.Event == nil || mult.EventYear == nil {
		return getAPIError(c, http.StatusNotFound, "Event/Year Not Found", nil)
	}
	// Check if they own this event.
	if mult.Event.AccountIdentifier != mkey.Account.Identifier {
		return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
	}
	results, err := database.AddResults(mult.EventYear.Identifier, resToAdd)
	if err != nil {
		return getAPIError(c, http.StatusInternalServerError, "Error Adding Results", err)
	}
	return c.JSON(http.StatusOK, types.AddResultsResponse{
		Count: len(results),
	})
}

func (h Handler) DeleteResults(c echo.Context) error {
	// Get Key from Authorization Header
	k, err := retrieveKey(c.Request())
	if err != nil {
		return getAPIError(c, http.StatusUnauthorized, "Error Getting Key From Authorization Header", err)
	}
	if k == nil {
		return getAPIError(c, http.StatusUnauthorized, "Key Not Provided in Authorization Header", nil)
	}
	var request types.GetResultsRequest
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
	// For results, let a write key delete.
	if mkey.Key.Type == "read" {
		return getAPIError(c, http.StatusUnauthorized, "Key is ReadOnly/Write", nil)
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
	// Check if the account is an admin or if they own this event.
	if mult.Event.AccountIdentifier != mkey.Account.Identifier {
		return getAPIError(c, http.StatusUnauthorized, "Ownership Error", nil)
	}
	var count int64
	if request.Distance != nil && len(*request.Distance) > 0 {
		count, err = database.DeleteDistanceResults(mult.EventYear.Identifier, *request.Distance)
		if err != nil {
			return getAPIError(c, http.StatusInternalServerError, "Error Deleting Results", err)
		}
	} else {
		count, err = database.DeleteEventResults(mult.EventYear.Identifier)
		if err != nil {
			return getAPIError(c, http.StatusInternalServerError, "Error Deleting Results", err)
		}
	}
	return c.JSON(http.StatusOK, types.AddResultsResponse{
		Count: int(count),
	})
}
