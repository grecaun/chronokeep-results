package types

import (
	"github.com/go-playground/validator/v10"
)

// Result is a structure holding information about a specific time
// result for a specific event.
type Result struct {
	PersonId         string `json:"person_id" validate:"required"`
	Bib              string `json:"bib" validate:"required"`
	First            string `json:"first"`
	Last             string `json:"last"`
	Age              int    `json:"age"`
	Gender           string `json:"gender"`
	AgeGroup         string `json:"age_group"`
	Distance         string `json:"distance"`
	Seconds          int    `json:"seconds" validate:"gte=0"`
	Milliseconds     int    `json:"milliseconds"`
	ChipSeconds      int    `json:"chip_seconds" validate:"gte=0"`
	ChipMilliseconds int    `json:"chip_milliseconds"`
	Segment          string `json:"segment"`
	Location         string `json:"location" validate:"required"`
	Occurence        int    `json:"occurence" validate:"gte=0"`
	Ranking          int    `json:"ranking"`
	AgeRanking       int    `json:"age_ranking"`
	GenderRanking    int    `json:"gender_ranking"`
	Finish           bool   `json:"finish"`
	Type             int    `json:"type"`
	Anonymous        bool   `json:"anonymous"`
	LocalTime        string `json:"local_time"`
	Division         string `json:"division"`
	DivisionRanking  int    `json:"division_ranking"`
}

type ResultVers1 struct {
	Bib              string `json:"bib"`
	First            string `json:"first"`
	Last             string `json:"last"`
	Age              int    `json:"age"`
	Gender           string `json:"gender"`
	AgeGroup         string `json:"age_group"`
	Distance         string `json:"distance"`
	Seconds          int    `json:"seconds"`
	Milliseconds     int    `json:"milliseconds"`
	ChipSeconds      int    `json:"chip_seconds"`
	ChipMilliseconds int    `json:"chip_milliseconds"`
	Segment          string `json:"segment"`
	Location         string `json:"location"`
	Occurence        int    `json:"occurence"`
	Ranking          int    `json:"ranking"`
	AgeRanking       int    `json:"age_ranking"`
	GenderRanking    int    `json:"gender_ranking"`
	Finish           bool   `json:"finish"`
	Type             int    `json:"type"`
	Anonymous        bool   `json:"anonymous"`
}

func (r *Result) ConvertToVers1() ResultVers1 {
	return ResultVers1{
		Bib:              r.Bib,
		First:            r.First,
		Last:             r.Last,
		Age:              r.Age,
		Gender:           r.Gender,
		AgeGroup:         r.AgeGroup,
		Distance:         r.Distance,
		Seconds:          r.Seconds,
		Milliseconds:     r.Milliseconds,
		ChipSeconds:      r.ChipSeconds,
		ChipMilliseconds: r.ChipMilliseconds,
		Segment:          r.Segment,
		Location:         r.Location,
		Occurence:        r.Occurence,
		Ranking:          r.Ranking,
		AgeRanking:       r.AgeRanking,
		GenderRanking:    r.GenderRanking,
		Finish:           r.Finish,
		Type:             r.Type,
		Anonymous:        r.Anonymous,
	}
}

// Validate Ensures valid data in the struct.
func (r *Result) Validate(validate *validator.Validate) error {
	return validate.Struct(r)
}

func (one *Result) Equals(two *Result) bool {
	return one.Bib == two.Bib &&
		one.First == two.First &&
		one.Last == two.Last &&
		one.Age == two.Age &&
		one.Gender == two.Gender &&
		one.AgeGroup == two.AgeGroup &&
		one.Distance == two.Distance &&
		one.Seconds == two.Seconds &&
		one.Milliseconds == two.Milliseconds &&
		one.ChipSeconds == two.ChipSeconds &&
		one.ChipMilliseconds == two.ChipMilliseconds &&
		one.Segment == two.Segment &&
		one.Location == two.Location &&
		one.Occurence == two.Occurence &&
		one.Ranking == two.Ranking &&
		one.AgeRanking == two.AgeRanking &&
		one.GenderRanking == two.GenderRanking &&
		one.Finish == two.Finish &&
		one.Type == two.Type &&
		one.Anonymous == two.Anonymous &&
		one.PersonId == two.PersonId &&
		one.Division == two.Division &&
		one.DivisionRanking == two.DivisionRanking
}

func (one *Result) SamePerson(two *Result) bool {
	return one.Bib == two.Bib &&
		one.First == two.First &&
		one.Last == two.Last &&
		one.Age == two.Age &&
		one.Gender == two.Gender &&
		one.AgeGroup == two.AgeGroup &&
		one.Distance == two.Distance &&
		one.Anonymous == two.Anonymous &&
		one.Division == two.Division
}

func (r *Result) AnonyInt() int {
	if r.Anonymous {
		return 1
	}
	return 0
}
