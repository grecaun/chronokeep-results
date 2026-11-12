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

package mysql

import (
	"chronokeep/results/database"
	"context"
	"fmt"
	"time"
)

// liveColumn selects whether an event year is live: it counts as live while it has
// received results within database.LiveWindow. It's used in place of the stored live column.
var liveColumn = fmt.Sprintf("COALESCE(results_received_at > NOW() - INTERVAL %d HOUR, FALSE) AS live", int(database.LiveWindow.Hours()))

// SetResultsReceived records that results were just received for an event year.
func (m *MySQL) SetResultsReceived(eventYearID int64) error {
	db, err := m.GetDB()
	if err != nil {
		return err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelfunc()
	_, err = db.ExecContext(
		ctx,
		"UPDATE event_year SET results_received_at=NOW() WHERE event_year_id=?;",
		eventYearID,
	)
	if err != nil {
		return fmt.Errorf("error setting results received: %v", err)
	}
	return nil
}
