package types

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRequestYearGetDateTime(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  time.Time
	}{
		{"rfc3339", "2024-05-06T07:08:09Z", time.Date(2024, 5, 6, 7, 8, 9, 0, time.UTC)},
		{"slash with offset", "2024/05/06 07:08:09 -07:00", time.Date(2024, 5, 6, 7, 8, 9, 0, time.FixedZone("", -7*3600))},
		// Previously this format returned time.Now() because the error check was inverted.
		{"slash without offset", "2024/05/06 07:08:09", time.Date(2024, 5, 6, 7, 8, 9, 0, time.UTC)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := RequestYear{Year: "2024", DateTime: c.input}
			got := req.GetDateTime()
			assert.True(t, c.want.Equal(got), "expected %v, got %v", c.want, got)
			year := req.ToYear()
			assert.True(t, c.want.Equal(year.DateTime), "expected %v, got %v", c.want, year.DateTime)
		})
	}
}

func TestRequestYearGetDateTimeInvalidFallsBackToNow(t *testing.T) {
	// Previously an unparseable value returned the zero time.
	for _, input := range []string{"", "not a date", "2024-13-45"} {
		before := time.Now()
		got := RequestYear{DateTime: input}.GetDateTime()
		after := time.Now()
		assert.False(t, got.IsZero(), "input %q", input)
		assert.False(t, got.Before(before) || got.After(after), "input %q gave %v", input, got)
	}
}

func TestRequestYearToYearCopiesFields(t *testing.T) {
	req := RequestYear{
		Year:        "2024-half",
		DateTime:    "2024/05/06 07:08:09",
		Live:        true,
		DaysAllowed: 3,
		RankingType: "chip",
	}
	year := req.ToYear()
	assert.Equal(t, "2024-half", year.Year)
	assert.True(t, year.Live)
	assert.Equal(t, 3, year.DaysAllowed)
	assert.Equal(t, "chip", year.RankingType)
}

func TestYearValidation(t *testing.T) {
	v := validator.New()
	for _, year := range []string{"2024", "2024-2", "2024-half", "A1"} {
		assert.NoError(t, (&RequestYear{Year: year}).Validate(v), year)
		assert.NoError(t, (&EventYear{Year: year, DateTime: time.Now()}).Validate(v), year)
	}
	for _, year := range []string{"", "2024 2", "2024/2", "2024_2", "../etc"} {
		assert.Error(t, (&RequestYear{Year: year}).Validate(v), year)
		assert.Error(t, (&EventYear{Year: year}).Validate(v), year)
	}
}

func TestEventYearValidateDoesNotPanic(t *testing.T) {
	// The DateTime field previously carried a string-only `datetime` validator tag
	// which caused validator to panic on a time.Time field.
	v := validator.New()
	assert.NotPanics(t, func() {
		_ = (&EventYear{Year: "2024", DateTime: time.Now()}).Validate(v)
	})
}

func TestRequestKeyValidate(t *testing.T) {
	v := validator.New()
	for _, kType := range []string{"read", "write", "delete"} {
		k := RequestKey{Type: kType}
		assert.NoError(t, k.Validate(v))
		assert.Equal(t, kType, k.Type)
	}
	// "default" must be converted to "read" on the request itself, otherwise a key of
	// type "default" is stored, which write checks (Type == "read") don't treat as read only.
	k := RequestKey{Type: "default"}
	assert.NoError(t, k.Validate(v))
	assert.Equal(t, "read", k.Type)
	assert.Equal(t, "read", k.ToKey().Type)
	for _, kType := range []string{"", "admin", "READ", "default "} {
		k := RequestKey{Type: kType}
		assert.Error(t, k.Validate(v), kType)
	}
}

func TestKeyValidate(t *testing.T) {
	v := validator.New()
	k := Key{Type: "default"}
	assert.NoError(t, k.Validate(v))
	assert.Equal(t, "read", k.Type)
	k = Key{Type: "invalid"}
	assert.Error(t, k.Validate(v))
}

func TestKeyIsAllowed(t *testing.T) {
	cases := []struct {
		allowed string
		referer string
		want    bool
	}{
		{"", "", true},
		{"", "https://anything.com/page", true},
		{"chronokeep.com", "https://chronokeep.com/results", true},
		{"chronokeep.com", "http://chronokeep.com", true},
		{"chronokeep.com", "https://chronokeep.com:8080/results", true},
		{"chronokeep.com", "https://evil.com/chronokeep.com", false},
		{"chronokeep.com", "https://chronokeep.com.evil.com/", false},
		{"chronokeep.com", "", false},
		{"chronokeep.com", "chronokeep.com", false},
		{"chronokeep.com", "ftp://chronokeep.com", false},
		{"a.com,b.com", "https://b.com/", true},
		// Whitespace after commas used to cause the second host to never match.
		{"a.com, b.com", "https://b.com/", true},
		{" a.com ,b.com", "https://a.com/", true},
		// Host names are case insensitive.
		{"Chronokeep.com", "https://chronokeep.COM/", true},
		{"a.com,,b.com", "https://c.com/", false},
	}
	for _, c := range cases {
		k := Key{AllowedHosts: c.allowed}
		assert.Equal(t, c.want, k.IsAllowed(c.referer), "allowed %q referer %q", c.allowed, c.referer)
	}
}

func TestKeyExpired(t *testing.T) {
	assert.False(t, Key{}.Expired())
	future := time.Now().Add(time.Hour)
	assert.False(t, Key{ValidUntil: &future}.Expired())
	past := time.Now().Add(-time.Second)
	assert.True(t, Key{ValidUntil: &past}.Expired())
}

func TestRequestKeyValidUntil(t *testing.T) {
	// RFC3339 is used as is.
	k := RequestKey{ValidUntil: "2030-01-02T03:04:05Z"}
	if assert.NotNil(t, k.GetValidUntil()) {
		assert.True(t, time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC).Equal(*k.GetValidUntil()))
	}
	// Date only values are valid until the end of the day in local time.
	endOfDay := time.Date(2030, 1, 2, 23, 59, 59, 0, time.Local)
	for _, input := range []string{"2030/01/02", "2030/1/2", "01/02/2030", "1/2/2030", "2030-01-02", "2030-1-2", "01-02-2030", "1-2-2030"} {
		k := RequestKey{ValidUntil: input}
		got := k.GetValidUntil()
		if assert.NotNil(t, got, input) {
			assert.True(t, endOfDay.Equal(*got), "input %q gave %v", input, got)
		}
		key := k.ToKey()
		if assert.NotNil(t, key.ValidUntil, input) {
			assert.True(t, endOfDay.Equal(*key.ValidUntil), "input %q gave %v", input, key.ValidUntil)
		}
	}
	for _, input := range []string{"", "never", "2030/13/45"} {
		k := RequestKey{ValidUntil: input}
		assert.Nil(t, k.GetValidUntil(), input)
		assert.Nil(t, k.ToKey().ValidUntil, input)
	}
}

func TestRequestKeyToKeyTrimsHosts(t *testing.T) {
	k := RequestKey{Name: "n", Value: "v", Type: "write", AllowedHosts: "  a.com,b.com  "}.ToKey()
	assert.Equal(t, "a.com,b.com", k.AllowedHosts)
	assert.Equal(t, "n", k.Name)
	assert.Equal(t, "v", k.Value)
	assert.Equal(t, "write", k.Type)
}

func TestKeyEqual(t *testing.T) {
	now := time.Now()
	same := now
	later := now.Add(time.Minute)
	a := Key{Name: "a", Value: "v", Type: "read"}
	b := a
	assert.True(t, a.Equal(&b))
	a.ValidUntil = &now
	assert.False(t, a.Equal(&b))
	b.ValidUntil = &same
	assert.True(t, a.Equal(&b))
	b.ValidUntil = &later
	assert.False(t, a.Equal(&b))
}

func TestEventValidate(t *testing.T) {
	v := validator.New()
	valid := func() Event {
		return Event{
			Name:         "Test Event's Run / Walk & Roll 5",
			Slug:         "Test-Event-5",
			ContactEmail: "test@example.com",
			Type:         "distance",
		}
	}
	e := valid()
	assert.NoError(t, e.Validate(v))
	assert.Equal(t, "test-event-5", e.Slug, "slug should be lower cased")
	for _, eType := range []string{"distance", "time", "backyardultra"} {
		e = valid()
		e.Type = eType
		assert.NoError(t, e.Validate(v), eType)
	}
	// Contact email is optional.
	e = valid()
	e.ContactEmail = ""
	assert.NoError(t, e.Validate(v))
	e = valid()
	e.Website = "https://example.com"
	e.Image = "https://example.com/img.png"
	assert.NoError(t, e.Validate(v))

	invalid := map[string]func(*Event){
		"empty slug":       func(e *Event) { e.Slug = "" },
		"slug with space":  func(e *Event) { e.Slug = "bad slug" },
		"slug with slash":  func(e *Event) { e.Slug = "bad/slug" },
		"slug with dot":    func(e *Event) { e.Slug = "bad.slug" },
		"empty name":       func(e *Event) { e.Name = "" },
		"name with symbol": func(e *Event) { e.Name = "Bad <script>" },
		"empty type":       func(e *Event) { e.Type = "" },
		"unknown type":     func(e *Event) { e.Type = "relay" },
		"bad website":      func(e *Event) { e.Website = "not a url" },
		"bad image":        func(e *Event) { e.Image = "not a url" },
		"bad email":        func(e *Event) { e.ContactEmail = "not-an-email" },
	}
	for name, mod := range invalid {
		e := valid()
		mod(&e)
		assert.Error(t, e.Validate(v), name)
	}
}

func TestParticipantValidate(t *testing.T) {
	v := validator.New()
	for _, bd := range []string{"1/2/2000", "01/02/2000", "2000/1/2", "2000/01/02"} {
		p := Participant{Distance: "5K", Birthdate: bd}
		assert.NoError(t, p.Validate(v), bd)
	}
	future := time.Now().AddDate(1, 0, 0)
	for _, bd := range []string{"", "yesterday", "2000-01-02", "13/45/2000", future.Format("1/2/2006"), future.Format("2006/1/2")} {
		p := Participant{Distance: "5K", Birthdate: bd}
		assert.Error(t, p.Validate(v), bd)
	}
	p := Participant{Birthdate: "1/2/2000"}
	assert.Error(t, p.Validate(v), "distance is required")
}

func TestAccountValidate(t *testing.T) {
	v := validator.New()
	for _, aType := range []string{"admin", "free", "paid", "registration"} {
		a := Account{Name: "n", Email: "a@b.com", Type: aType}
		assert.NoError(t, a.Validate(v), aType)
	}
	assert.Error(t, (&Account{Name: "n", Email: "a@b.com", Type: "root"}).Validate(v))
	assert.Error(t, (&Account{Name: "n", Email: "a@b.com"}).Validate(v))
	assert.Error(t, (&Account{Name: "n", Email: "nope", Type: "free"}).Validate(v))
	assert.Error(t, (&Account{Email: "a@b.com", Type: "free"}).Validate(v))
}

func TestAccountPasswordIsHashed(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	assert.NoError(t, err)
	assert.True(t, (&Account{Password: string(hash)}).PasswordIsHashed())
	for _, pw := range []string{"", "password", "$2a$", "$2a$10$short"} {
		assert.False(t, (&Account{Password: pw}).PasswordIsHashed(), pw)
	}
	// A hash of the empty string is not considered a valid password hash.
	empty, _ := bcrypt.GenerateFromPassword([]byte(""), bcrypt.MinCost)
	assert.False(t, (&Account{Password: string(empty)}).PasswordIsHashed())
}

func TestAccountEquals(t *testing.T) {
	a := Account{Name: "n", Email: "e@e.com", Type: "free", Identifier: 1}
	b := a
	b.Identifier = 2
	assert.True(t, a.Equals(&b))
	assert.False(t, a.Equals(nil))
	b.Type = "paid"
	assert.False(t, a.Equals(&b))
}

func TestResultValidate(t *testing.T) {
	v := validator.New()
	r := Result{PersonId: "1", Bib: "1", Location: "Finish"}
	assert.NoError(t, r.Validate(v))
	for name, mod := range map[string]func(*Result){
		"no person id":       func(r *Result) { r.PersonId = "" },
		"no bib":             func(r *Result) { r.Bib = "" },
		"no location":        func(r *Result) { r.Location = "" },
		"negative seconds":   func(r *Result) { r.Seconds = -1 },
		"negative chip":      func(r *Result) { r.ChipSeconds = -1 },
		"negative occurence": func(r *Result) { r.Occurence = -1 },
	} {
		bad := r
		mod(&bad)
		assert.Error(t, bad.Validate(v), name)
	}
}

func TestResultConvertToVers1(t *testing.T) {
	r := Result{Bib: "10", First: "F", Last: "L", Age: 30, Seconds: 100, Milliseconds: 5, Location: "Finish", Finish: true, Anonymous: true}
	v1 := r.ConvertToVers1()
	assert.Equal(t, r.Bib, v1.Bib)
	assert.Equal(t, r.Seconds, v1.Seconds)
	assert.Equal(t, r.Milliseconds, v1.Milliseconds)
	assert.Equal(t, r.Finish, v1.Finish)
	assert.Equal(t, r.Anonymous, v1.Anonymous)
	assert.Equal(t, 1, r.AnonyInt())
	r.Anonymous = false
	assert.Equal(t, 0, r.AnonyInt())
}

func TestSegmentAndDistanceValidate(t *testing.T) {
	v := validator.New()
	s := Segment{Location: "L", DistanceName: "5K", Name: "Mile 1", DistanceValue: 1, DistanceUnit: "miles"}
	assert.NoError(t, s.Validate(v))
	s.DistanceValue = 0
	assert.Error(t, s.Validate(v))
	d := Distance{Name: "5K", Certification: "CERT"}
	assert.NoError(t, d.Validate(v))
	d.Certification = ""
	assert.Error(t, d.Validate(v))
}
