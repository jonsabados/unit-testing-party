package unit

import (
	"strconv"
	"time"
)

type Generation string

const (
	// see https://www.cnn.com/2013/11/06/us/baby-boomer-generation-fast-facts/index.html
	Greatest   Generation = "Greatest"     // 1924 and earlier
	Silent     Generation = "Silent"       // 1925 - 1945
	BabyBoomer Generation = "Baby Boomer"  // 1946 - 1964
	GenX       Generation = "Generation X" // 1965 - 1980
	Millennial Generation = "Millennial"   // 1981 - 1996
	GenZ       Generation = "Generation Z" // 1997 +
)

type RemoteEmployee struct {
	Status string `json:"status"`
	Data   struct {
		ID             string `json:"id"`
		EmployeeName   string `json:"employee_name"`
		EmployeeSalary string `json:"employee_salary"`
		EmployeeAge    string `json:"employee_age"`
		ProfileImage   string `json:"profile_image"`
	} `json:"data"`
}

type Employee struct {
	ID         string     `json:"id"`
	Name       string     `json:"employee_name"`
	Age        int        `json:"age"`
	Generation Generation `json:"generation"`
}

// MapBirthYear maps a birth year to a generation. It doesn't really need to be its own thing currently and could live
// inside a function to map employees easy enough, but testing gets much easier if we can just test mapping birth year.
// And, this is functionality that really could be used outside of mapping an employee so what is there to loose by
// pulling it out?
func MapBirthYear(birthYear int) Generation {
	if birthYear <= 1924 {
		return Greatest
	} else if birthYear <= 1945 {
		return Silent
	} else if birthYear <= 1964 {
		return BabyBoomer
	} else if birthYear <= 1980 {
		return GenX
	} else if birthYear <= 1996 {
		return Millennial
	} else {
		return GenZ
	}
}

type EmployeeConverter func (employee *RemoteEmployee) (*Employee, error)

// Using a higher order function to produce a type that is just a function might be overkill in this case, could just
// as easily have a ConvertEmployee function that calls MapBirthYear. But, this higher order function technique is
// super useful if you need to call external dependencies from a function who's method signature doesn't allow it. An
// alternate approach would be to define a variable that defaults to the production dependency (if possible), but that
// creates a global singleton thing that can be changed at runtime, so EWWWWWW. Doing it this way tests can just inline
// a function to take the place of MapBirthYear to do whatever behavior desired. This is also a good way to deal with
// things where a struct might be used for dependencies but there is no state and only a single function (nix the struct,
// just pass around a function created by another function who has the dependency in scope).
func NewEmployeeFactory(mapBirthYear func(birthYear int) Generation) EmployeeConverter {
	return func(employee *RemoteEmployee)  (*Employee, error) {
		age, err := strconv.Atoi(employee.Data.EmployeeAge)
		if err != nil {
			return nil, err
		}
		birthYear := time.Now().Year() - age

		return &Employee{
			ID:         employee.Data.ID,
			Name:       employee.Data.EmployeeName,
			Age:        age,
			Generation: mapBirthYear(birthYear),
		}, nil
	}
}
