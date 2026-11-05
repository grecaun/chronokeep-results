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
	"regexp"
	"strings"
)

var phoneFormatting = regexp.MustCompile(`[-\(\).\s]`)

// phoneDigits strips formatting characters from a phone number. It returns the digits,
// whether the number started with a +, and false if any other character was found.
func phoneDigits(phone string) (string, bool, bool) {
	out := phoneFormatting.ReplaceAllString(strings.TrimSpace(phone), "")
	plus := strings.HasPrefix(out, "+")
	out = strings.TrimPrefix(out, "+")
	if out == "" {
		return "", plus, false
	}
	for _, r := range out {
		if r < '0' || r > '9' {
			return "", plus, false
		}
	}
	return out, plus, true
}

// normalizePhone returns the 10 digit US phone number used for SMS subscriptions. A
// leading US country code (1 or +1) is removed. The bool is false if the value is not a
// valid 10 digit phone number.
func normalizePhone(phone string) (string, bool) {
	out, _, ok := phoneDigits(phone)
	if !ok {
		return "", false
	}
	if len(out) == 11 && out[0] == '1' {
		out = out[1:]
	}
	if len(out) != 10 {
		return "", false
	}
	return out, true
}

// canonicalPhone returns a phone number in E.164 format (e.g. +15551234567), which is the
// format Twilio uses for the From field and the desktop application uses when sending.
// 10 digit numbers are assumed to be US numbers. Other international numbers must start
// with a +. The bool is false if the value isn't a valid phone number.
func canonicalPhone(phone string) (string, bool) {
	digits, plus, ok := phoneDigits(phone)
	if !ok {
		return "", false
	}
	switch {
	case len(digits) == 10 && !plus:
		return "+1" + digits, true
	case len(digits) == 11 && digits[0] == '1':
		return "+" + digits, true
	case plus && len(digits) >= 8 && len(digits) <= 15:
		return "+" + digits, true
	}
	return "", false
}

// phoneVariants returns every format a phone number may have been stored in on the
// blocked list, so that a number can be unblocked regardless of how it was added.
func phoneVariants(canonical string) []string {
	digits := strings.TrimPrefix(canonical, "+")
	out := []string{canonical, digits}
	if len(digits) == 11 && digits[0] == '1' {
		out = append(out, digits[1:])
	}
	return out
}

// samePhone reports whether two phone numbers refer to the same number.
func samePhone(one, two string) bool {
	a, ok := canonicalPhone(one)
	if !ok {
		return false
	}
	b, ok := canonicalPhone(two)
	return ok && a == b
}
