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

func NewTestCaseUpdateUser(db *gorm.DB, serverUrl string) pkg.ITestCase {
	return &updateUserTestCase{
		db:        db,
		serverUrl: serverUrl,
		Expected:  dummiesdata.Employee_NormalCaseCreateEmployee,
	}
}

type updateUserTestCase struct {
	db        *gorm.DB
	serverUrl string
	Expected  any
	Actual    any
}

func (t *updateUserTestCase) Name() string {
	return "Normal Case: update Employee"
}

func (t *updateUserTestCase) Preparing() error {
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

func (t *updateUserTestCase) Cleaning() error {
	return t.db.Migrator().DropTable(&models.Employee{})
}

func (t *updateUserTestCase) Do() error {
	fmt.Printf("================[testcase][%s] is running...================\n", t.Name())
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				input := dummiesdata.Employee_NormalCaseUpdateEmployee
				isAlive, err := pkg.PutRequest(pkg.NewHTTPClient(t.serverUrl), "/Employee/1", input, map[string]any{})
				if isAlive {
					if err != nil {
						return
					}
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

func (t *updateUserTestCase) GetExpected() (any, error) {
	dob, err := time.Parse("2006-01-02", dummiesdata.Employee_NormalCaseCreateEmployee["Dob"].(string))
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"id":     1,
		"Name":   dummiesdata.Employee_NormalCaseCreateEmployee["Name"],
		"Email":  dummiesdata.Employee_NormalCaseCreateEmployee["Email"],
		"Dob":    string(dob.Format("2006-01-02")),
		"Age":    int64(dummiesdata.Employee_NormalCaseUpdateEmployee["Age"].(int)),
		"Phone":  dummiesdata.Employee_NormalCaseUpdateEmployee["Phone"],
		"Mature": dummiesdata.Employee_NormalCaseUpdateEmployee["Mature"],
	}, nil
}

func (t *updateUserTestCase) GetActual() (any, error) {
	var res = new(models.Employee)
	if err := t.db.Model(&models.Employee{}).First(&res).Error; err != nil {
		return nil, err
	}
	t.Expected = map[string]any{
		"id":     1,
		"Name":   res.Name,
		"Email":  res.Email,
		"Dob":    res.Dob.Format("2006-01-02"),
		"Age":    res.Age,
		"Phone":  res.Phone,
		"Mature": res.Mature,
	}
	return t.Expected, nil
}

func (t *updateUserTestCase) CheckResult() (bool, error) {
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

func (t *updateUserTestCase) RunTest() (bool, error) {
	fmt.Printf("================[testcase][%s] is running...================\n", t.Name())
	defer t.Cleaning()

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
