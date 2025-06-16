package pkg

type ITestCase interface {
	Preparing() error
	Cleaning() error
	Name() string
	Do() error
	GetExpected() (any, error)
	GetActual() (any, error)
	CheckResult() (bool, error)
	RunTest() (bool, error)
}

// type NewTestCaseDTO struct {
// 	Name        string
// 	Preparing   func()
// 	Cleaning    func()
// 	Do          func() error
// 	CheckResult func() (bool, error)
// }

// func NewTestCase(input *NewTestCaseDTO) ITestCase {
// 	return &testCase{
// 		name:        input.Name,
// 		preparing:   input.Preparing,
// 		cleaning:    input.Cleaning,
// 		do:          input.Do,
// 		checkResult: input.CheckResult,
// 	}
// }

// type testCase struct {
// 	name        string
// 	preparing   func()
// 	cleaning    func()
// 	do          func() error
// 	checkResult func() (bool, error)
// }

// func (t *testCase) Preparing() {
// 	println("Preparing " + t.name)
// 	t.preparing()
// }

// func (t *testCase) Cleaning() {
// 	println("Cleaning " + t.name)
// 	t.cleaning()
// }

// func (t *testCase) Do() error {
// 	return t.do()
// }

// func (t *testCase) CheckResult() (bool, error) {
// 	return t.checkResult()
// }

// func (t *testCase) Name() string {
// 	return t.name
// }
