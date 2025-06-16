package main

// This is intergration test code
// We will launch a service, and with each testcase, we will prepare some row in db, then
// send a request to the service and check the response if it matchs the expected result. And we
// will clean the data after each testcase
import (
	"fmt"

	"github.com/duytacong24895/go-crud-generator/tests/pkg/database"
	"github.com/duytacong24895/go-crud-generator/tests/testcases"
)

func main() {
	fmt.Println("Start testing")
	db := database.InitDB()

	statistics := &Statistics{}

	statistics.On(testcases.NewTestCaseCreateUser(db, "http://localhost:8080/crud").RunTest())
	statistics.On(testcases.NewTestCaseUpdateUser(db, "http://localhost:8080/crud").RunTest())
	// statistics.On(testcases.NewTestCaseDetailUser(db, "http://localhost:8080/crud").RunTest())
	// statistics.On(testcases.NewTestCaseDeleteUser(db, "http://localhost:8080/crud").RunTest())
	// statistics.On(testcases.NewTestCaseGetListUserWithFilter(db, "http://localhost:8080/crud").RunTest())
	// statistics.On(testcases.NewTestCaseGetListUserWithNestedFilter(db, "http://localhost:8080/crud").RunTest())
	fmt.Println("Finish testing")
	// Print the statistics
	statistics.Print()
	statistics.Reset()
}

type Statistics struct {
	NumPassed int
	NumErrors int
	NumTests  int
}

func (s *Statistics) On(passed bool, err error) {
	s.NumTests++
	if passed {
		s.NumPassed++
	}

	if err != nil {
		s.NumErrors++
	}
}
func (s *Statistics) Print() {
	fmt.Printf("Total tests: %d, Passed: %d, Failed: %d, Errors: %d\n", s.NumTests, s.NumPassed, s.NumTests-s.NumPassed, s.NumErrors)
}
func (s *Statistics) Reset() {
	s.NumPassed = 0
	s.NumErrors = 0
	s.NumTests = 0
}
