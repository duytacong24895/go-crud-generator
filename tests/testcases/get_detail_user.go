package testcases

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	crud_generator "github.com/duytacong24895/go-crud-generator"
	"github.com/duytacong24895/go-crud-generator/tests/pkg"
	dummiesdata "github.com/duytacong24895/go-crud-generator/tests/pkg/dummies_data"
	"github.com/duytacong24895/go-crud-generator/tests/pkg/models"
	"github.com/go-chi/chi/v5"

	"gorm.io/gorm"
)

func NewTestCaseDetailUser(db *gorm.DB, serverUrl string) pkg.ITestCase {
	return &getDetailUserTestCase{
		db:        db,
		serverUrl: serverUrl,
		Expected:  dummiesdata.Employee_NormalCaseCreateEmployee,
	}
}

type getDetailUserTestCase struct {
	db        *gorm.DB
	serverUrl string
	Expected  any
	Actual    any
}

func (t *getDetailUserTestCase) Name() string {
	return "Normal Case: get detail Employee by id"
}

func (t *getDetailUserTestCase) Preparing() error {
	if err := t.db.Migrator().DropTable(&models.Employee{}); err != nil {
		return err
	}
	if err := t.db.AutoMigrate(&models.Employee{}); err != nil {
		return err
	}

	dobStr, ok := dummiesdata.Employee_NormalCaseCreateEmployee["Dob"].(string)
	if !ok {
		return fmt.Errorf("invalid type for Dob, expected string")
	}
	var dob time.Time
	var err error
	if dob, err = time.Parse("2006-01-02", dobStr); err != nil {
		return fmt.Errorf("invalid date format for Dob: %s, expected format is YYYY-MM-DD", dob)
	}
	// Create a new employee with the provided data
	var age int
	if age, ok = dummiesdata.Employee_NormalCaseCreateEmployee["Age"].(int); !ok {
		return fmt.Errorf("invalid type for Age, expected int")
	}

	t.db.Create(&models.Employee{
		Name:   dummiesdata.Employee_NormalCaseCreateEmployee["Name"].(string),
		Email:  dummiesdata.Employee_NormalCaseCreateEmployee["Email"].(string),
		Dob:    dob,
		Age:    int64(age),
		Phone:  dummiesdata.Employee_NormalCaseCreateEmployee["Phone"].(string),
		Mature: dummiesdata.Employee_NormalCaseCreateEmployee["Mature"].(bool),
	})
	return nil
}

func (t *getDetailUserTestCase) Cleaning() error {
	return t.db.Migrator().DropTable(&models.Employee{})
}

func (t *getDetailUserTestCase) Do() error {
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				var output map[string]any
				isAlive, err := pkg.GetRequest(pkg.NewHTTPClient(t.serverUrl), "/Employee/1", &output)
				if isAlive {
					if err != nil {
						return
					}
					t.Actual = output
					done <- true
				} else {
					fmt.Println("Waiting the server is running...")
					time.Sleep(5 * time.Second)
				}
			}
		}
	}()

	go func() {
		r := chi.NewRouter()
		crud_generator.NewCRUDGenerator(r, t.db).
			RegisterModel(&models.Employee{}).
			Run()

		fmt.Println("Server is running on port 8080")
		http.ListenAndServe(":8080", r)
	}()
	<-done
	return nil
}

func (t *getDetailUserTestCase) GetExpected() (any, error) {
	dob, err := time.Parse("2006-01-02", dummiesdata.Employee_NormalCaseCreateEmployee["Dob"].(string))
	if err != nil {
		return nil, err
	}
	return models.Employee{
		ID:     1,
		Name:   dummiesdata.Employee_NormalCaseCreateEmployee["Name"].(string),
		Email:  dummiesdata.Employee_NormalCaseCreateEmployee["Email"].(string),
		Dob:    dob,
		Age:    int64(dummiesdata.Employee_NormalCaseCreateEmployee["Age"].(int)),
		Phone:  dummiesdata.Employee_NormalCaseCreateEmployee["Phone"].(string),
		Mature: dummiesdata.Employee_NormalCaseCreateEmployee["Mature"].(bool),
	}, nil
}

func (t *getDetailUserTestCase) GetActual() (any, error) {
	delete(t.Actual.(map[string]any), "created_at")
	delete(t.Actual.(map[string]any), "updated_at")
	delete(t.Actual.(map[string]any), "deleted_at")
	strDob, ok := t.Actual.(map[string]any)["dob"].(string)
	if !ok {
		return nil, fmt.Errorf("dob field is empty in the actual response")
	}
	dob, err := time.Parse("2006-01-02T00:00:00Z", strDob)
	if err != nil {
		return nil, fmt.Errorf("invalid date format for Dob: %s, expected format is YYYY-MM-DD", strDob)
	}
	t.Actual.(map[string]any)["dob"] = string(dob.Format("2006-01-02"))
	return models.Employee{
		ID:     1,
		Name:   t.Actual.(map[string]any)["name"].(string),
		Email:  t.Actual.(map[string]any)["email"].(string),
		Dob:    dob,
		Age:    int64(t.Actual.(map[string]any)["age"].(float64)),
		Phone:  t.Actual.(map[string]any)["phone"].(string),
		Mature: t.Actual.(map[string]any)["mature"].(bool),
	}, nil
}

func (t *getDetailUserTestCase) CheckResult() (bool, error) {
	expect, err := t.GetExpected()
	if err != nil {
		return false, err
	}
	actual, err := t.GetActual()

	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(expect, actual), nil
}

func (t *getDetailUserTestCase) RunTest() (bool, error) {
	fmt.Printf("================[testcase][%s] is running...================\n", t.Name())
	defer t.Cleaning()
	time.Sleep(5 * time.Second) // Wait for the server to start

	if err := t.Preparing(); err != nil {
		fmt.Printf("[testcase][%s] got a Error: %v\n", t.Name(), err)
		return false, err
	}
	if err := t.Do(); err != nil {
		fmt.Printf("[testcase][%s] got a Error: %v\n", t.Name(), err)
		return false, err
	}
	result, err := t.CheckResult()
	if err != nil {
		fmt.Printf("[testcase][%s] got a Error: %v\n", t.Name(), err)
		return false, err
	}

	if result {
		fmt.Printf("[testcase][%s] was passed\n", t.Name())
		return true, nil
	} else {
		fmt.Printf("[testcase][%s] was failed\n", t.Name())
		return false, nil
	}
}
