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

func NewTestCaseCreateUser(db *gorm.DB, serverUrl string) pkg.ITestCase {
	return &createUserTestCase{
		db:        db,
		serverUrl: serverUrl,
		Expected:  dummiesdata.Employee_NormalCaseCreateEmployee,
	}
}

type createUserTestCase struct {
	db        *gorm.DB
	serverUrl string
	Expected  any
	Actual    any
}

func (t *createUserTestCase) Name() string {
	return "Normal Case: Create Employee"
}

func (t *createUserTestCase) Preparing() error {
	if err := t.db.Migrator().DropTable(&models.Employee{}); err != nil {
		return err
	}
	if err := t.db.AutoMigrate(&models.Employee{}); err != nil {
		return err
	}
	return nil
}

func (t *createUserTestCase) Cleaning() error {
	return t.db.Migrator().DropTable(&models.Employee{})
}

func (t *createUserTestCase) Do() error {
	fmt.Printf("================[testcase][%s] is running...================\n", t.Name())
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				input := dummiesdata.Employee_NormalCaseCreateEmployee
				isAlive, err := pkg.PostRequest(pkg.NewHTTPClient(t.serverUrl), "/Employee", input, map[string]any{})
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

func (t *createUserTestCase) GetExpected() (any, error) {
	dob, err := time.Parse("2006-01-02", dummiesdata.Employee_NormalCaseCreateEmployee["Dob"].(string))
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"id":     1,
		"Name":   dummiesdata.Employee_NormalCaseCreateEmployee["Name"],
		"Email":  dummiesdata.Employee_NormalCaseCreateEmployee["Email"],
		"Dob":    string(dob.Format("2006-01-02")),
		"Age":    int64(dummiesdata.Employee_NormalCaseCreateEmployee["Age"].(int)),
		"Phone":  dummiesdata.Employee_NormalCaseCreateEmployee["Phone"],
		"Mature": dummiesdata.Employee_NormalCaseCreateEmployee["Mature"],
	}, nil
}

func (t *createUserTestCase) GetActual() (any, error) {
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

func (t *createUserTestCase) CheckResult() (bool, error) {
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

func (t *createUserTestCase) RunTest() (bool, error) {
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
