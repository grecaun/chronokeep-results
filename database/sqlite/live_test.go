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

package sqlite

import (
	"chronokeep/results/types"
	"context"
	"fmt"
	"testing"
	"time"
)

// setReceivedHoursAgo backdates when an event year last received results.
func setReceivedHoursAgo(t *testing.T, o *SQLite, eventYearID int64, hours int) {
	db, err := o.GetDB()
	if err != nil {
		t.Fatalf("Error getting database: %v", err)
	}
	if _, err := db.ExecContext(context.Background(), fmt.Sprintf("UPDATE event_year SET results_received_at=datetime('now', '-%d hours') WHERE event_year_id=?;", hours), eventYearID); err != nil {
		t.Fatalf("Error backdating results received: %v", err)
	}
}

func TestLiveFromResultsReceived(t *testing.T) {
	db, finalize, err := setupTests(t)
	if err != nil {
		t.Fatalf("Error setting up test. %v", err)
	}
	defer finalize(t)
	setupEventYearTests()
	account, _ := db.AddAccount(accounts[0])
	event, _ := db.AddEvent(types.Event{
		AccountIdentifier: account.Identifier,
		Name:              "Live Event",
		Slug:              "live-event",
		ContactEmail:      "live@test.com",
	})
	// a stored live value no longer counts, only when results were last received
	year, err := db.AddEventYear(types.EventYear{
		EventIdentifier: event.Identifier,
		Year:            "2024",
		DateTime:        time.Date(2024, 6, 1, 8, 0, 0, 0, time.Local),
		Live:            true,
		DaysAllowed:     1,
		RankingType:     "gun",
	})
	if err != nil {
		t.Fatalf("Error adding event year: %v", err)
	}
	checkLive := func(when string, want bool) {
		t.Helper()
		single, err := db.GetEventYear("live-event", "2024")
		if err != nil || single == nil {
			t.Fatalf("%s: error getting event year: %v", when, err)
		}
		if single.Live != want {
			t.Errorf("%s: expected GetEventYear live to be %v, found %v", when, want, single.Live)
		}
		years, err := db.GetEventYears("live-event")
		if err != nil || len(years) != 1 {
			t.Fatalf("%s: error getting event years: %v", when, err)
		}
		if years[0].Live != want {
			t.Errorf("%s: expected GetEventYears live to be %v, found %v", when, want, years[0].Live)
		}
		all, err := db.GetAllEventYears()
		if err != nil || len(all) != 1 {
			t.Fatalf("%s: error getting all event years: %v", when, err)
		}
		if all[0].Live != want {
			t.Errorf("%s: expected GetAllEventYears live to be %v, found %v", when, want, all[0].Live)
		}
		for _, y := range []string{"2024", ""} {
			mult, err := db.GetEventAndYear("live-event", y)
			if err != nil || mult == nil || mult.EventYear == nil {
				t.Fatalf("%s: error getting event and year %q: %v", when, y, err)
			}
			if mult.EventYear.Live != want {
				t.Errorf("%s: expected GetEventAndYear(%q) live to be %v, found %v", when, y, want, mult.EventYear.Live)
			}
		}
	}
	checkLive("before any results", false)
	if err := db.SetResultsReceived(year.Identifier); err != nil {
		t.Fatalf("Error setting results received: %v", err)
	}
	checkLive("just after results", true)
	setReceivedHoursAgo(t, db, year.Identifier, 5)
	checkLive("five hours after results", true)
	setReceivedHoursAgo(t, db, year.Identifier, 7)
	checkLive("seven hours after results", false)
}
