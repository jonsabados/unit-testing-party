package unit

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMapBirthYear(t *testing.T) {
	testCases := []struct {
		desc           string
		input          int
		expectedResult Generation
	}{
		{
			"really old",
			1890,
			Greatest,
		},
		{
			"Greatest Edge",
			1294,
			Greatest,
		},
		{
			"Silent start",
			1925,
			Silent,
		},
		{
			"Silent end",
			1945,
			Silent,
		},
		{
			"Baby Boomer start",
			1956,
			BabyBoomer,
		},
		{
			"Baby Boomer end",
			1964,
			BabyBoomer,
		},
		{
			"GenX start",
			1965,
			GenX,
		},
		{
			"GenX end",
			1980,
			GenX,
		},
		{
			"Millennial start",
			1981,
			Millennial,
		},
		{
			"Millennial end",
			1996,
			Millennial,
		},
		{
			"GenZ start",
			1997,
			GenZ,
		},
		{
			"Young GenZ end",
			2015,
			GenZ,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			asserter := assert.New(t)

			asserter.Equal(tc.expectedResult, MapBirthYear(tc.input))
		})
	}
}

func TestNewEmployeeGargabeAge(t *testing.T) {
	asserter := assert.New(t)

	mapBirthYear := func(birthYear int) Generation {
		asserter.Fail("should not have reached this point")
		return GenZ
	}

	input := RemoteEmployee{
		Status: "blah",
		Data: struct {
			ID             string `json:"id"`
			EmployeeName   string `json:"employee_name"`
			EmployeeSalary string `json:"employee_salary"`
			EmployeeAge    string `json:"employee_age"`
			ProfileImage   string `json:"profile_image"`
		}{
			"1",
			"Bob",
			"123",
			"NotANumber",
			"foo",
		},
	}

	_, err := NewEmployeeFactory(mapBirthYear)(&input)
	asserter.EqualError(err, "strconv.Atoi: parsing \"NotANumber\": invalid syntax")
}

func TestNewEmployeeFactoryHappyPath(t *testing.T) {
	asserter := assert.New(t)

	expectedBirthYear := time.Now().Year() - 20
	expectedGeneration := GenZ
	mapBirthYear := func(birthYear int) Generation {
		asserter.Equal(expectedBirthYear, birthYear)
		return expectedGeneration
	}

	input := RemoteEmployee{
		Status: "blah",
		Data: struct {
			ID             string `json:"id"`
			EmployeeName   string `json:"employee_name"`
			EmployeeSalary string `json:"employee_salary"`
			EmployeeAge    string `json:"employee_age"`
			ProfileImage   string `json:"profile_image"`
		}{
			"1",
			"Bob",
			"123",
			"20",
			"foo",
		},
	}

	res, err := NewEmployeeFactory(mapBirthYear)(&input)
	asserter.NoError(err)
	asserter.Equal(&Employee{
		ID:         "1",
		Name:       "Bob",
		Age:        20,
		Generation: expectedGeneration,
	}, res)
}
