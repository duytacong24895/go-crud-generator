package testcases

import (
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"time"

	crud_generator "github.com/duytacong24895/go-crud-generator"
	"github.com/duytacong24895/go-crud-generator/tests/pkg"
	dummiesdata "github.com/duytacong24895/go-crud-generator/tests/pkg/dummies_data"
	"github.com/duytacong24895/go-crud-generator/tests/pkg/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func NewTestCaseGetListUserWithNestedFilter(db *gorm.DB, serverUrl string) pkg.ITestCase {
	return &getListUserWithNestedFilter{
		db:        db,
		serverUrl: serverUrl,
		Expected:  dummiesdata.Employee_NormalCaseCreateEmployee,
	}
}

type getListUserWithNestedFilter struct {
	db        *gorm.DB
	serverUrl string
	Expected  any
	Actual    any
}

func (t *getListUserWithNestedFilter) Name() string {
	return "Normal Case: get list Employee with nested filter"
}

func (t *getListUserWithNestedFilter) Preparing() error {
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

	// Create another employee with different data
	dobStr2, ok := dummiesdata.Employee_NormalCaseCreateEmployee2["Dob"].(string)
	if !ok {
		return fmt.Errorf("invalid type for Dob, expected string")
	}
	var dob2 time.Time
	var err2 error
	if dob2, err2 = time.Parse("2006-01-02", dobStr2); err2 != nil {
		return fmt.Errorf("invalid date format for Dob: %s, expected format is YYYY-MM-DD", dob2)
	}
	// Create a new employee with the provided data
	var age2 int
	if age2, ok = dummiesdata.Employee_NormalCaseCreateEmployee2["Age"].(int); !ok {
		return fmt.Errorf("invalid type for Age, expected int")
	}
	t.db.Create(&models.Employee{
		Name:   dummiesdata.Employee_NormalCaseCreateEmployee2["Name"].(string),
		Email:  dummiesdata.Employee_NormalCaseCreateEmployee2["Email"].(string),
		Dob:    dob2,
		Age:    int64(age2),
		Phone:  dummiesdata.Employee_NormalCaseCreateEmployee2["Phone"].(string),
		Mature: dummiesdata.Employee_NormalCaseCreateEmployee2["Mature"].(bool),
	})

	// Create another employee with different data
	dobStr3, ok := dummiesdata.Employee_NormalCaseCreateEmployee3["Dob"].(string)
	if !ok {
		return fmt.Errorf("invalid type for Dob, expected string")
	}
	var dob3 time.Time
	var err3 error
	if dob3, err3 = time.Parse("2006-01-02", dobStr3); err3 != nil {
		return fmt.Errorf("invalid date format for Dob: %s, expected format is YYYY-MM-DD", dob2)
	}
	// Create a new employee with the provided data
	var age3 int
	if age3, ok = dummiesdata.Employee_NormalCaseCreateEmployee3["Age"].(int); !ok {
		return fmt.Errorf("invalid type for Age, expected int")
	}
	t.db.Create(&models.Employee{
		Name:   dummiesdata.Employee_NormalCaseCreateEmployee3["Name"].(string),
		Email:  dummiesdata.Employee_NormalCaseCreateEmployee3["Email"].(string),
		Dob:    dob3,
		Age:    int64(age3),
		Phone:  dummiesdata.Employee_NormalCaseCreateEmployee3["Phone"].(string),
		Mature: dummiesdata.Employee_NormalCaseCreateEmployee3["Mature"].(bool),
	})
	return nil
}

func (t *getListUserWithNestedFilter) Cleaning() error {
	return t.db.Migrator().DropTable(&models.Employee{})
}

func (t *getListUserWithNestedFilter) Do() error {
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				var output []map[string]any
				params := map[string]string{
					"filter":    `[["age","lte","18"],"_and",["mature", "eq", "0"]]`,
					"page":      "1",
					"page_size": "10",
					"order_by":  "id desc",
				}
				isAlive, err := pkg.GetWithParams(pkg.NewHTTPClient(t.serverUrl), "/Employee", params, &output)
				if isAlive {
					if err != nil {
						fmt.Println(err)
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

func (t *getListUserWithNestedFilter) GetExpected() (any, error) {
	dob, err := time.Parse("2006-01-02", dummiesdata.Employee_NormalCaseCreateEmployee2["Dob"].(string))
	if err != nil {
		return nil, err
	}
	return []models.Employee{
		{
			ID:     2,
			Name:   dummiesdata.Employee_NormalCaseCreateEmployee2["Name"].(string),
			Email:  dummiesdata.Employee_NormalCaseCreateEmployee2["Email"].(string),
			Dob:    dob,
			Age:    int64(dummiesdata.Employee_NormalCaseCreateEmployee2["Age"].(int)),
			Phone:  dummiesdata.Employee_NormalCaseCreateEmployee2["Phone"].(string),
			Mature: dummiesdata.Employee_NormalCaseCreateEmployee2["Mature"].(bool),
		},
	}, nil
}

func (t *getListUserWithNestedFilter) GetActual() (any, error) {
	var formatedActual []models.Employee
	for _, item := range t.Actual.([]map[string]any) {
		delete(item, "created_at")
		delete(item, "updated_at")
		delete(item, "deleted_at")
		strDob, ok := item["dob"].(string)
		if !ok {
			return nil, fmt.Errorf("dob field is empty in the actual response")
		}
		dob, err := time.Parse("2006-01-02T00:00:00Z", strDob)
		if err != nil {
			return nil, fmt.Errorf("invalid date format for Dob: %s, expected format is YYYY-MM-DD", strDob)
		}
		item["dob"] = string(dob.Format("2006-01-02"))
		formatedActual = append(formatedActual, models.Employee{
			ID:     int64(item["id"].(float64)),
			Name:   item["name"].(string),
			Email:  item["email"].(string),
			Dob:    dob,
			Age:    int64(item["age"].(float64)),
			Phone:  item["phone"].(string),
			Mature: item["mature"].(bool),
		})
	}
	sort.Slice(formatedActual, func(i, j int) bool {
		return formatedActual[i].ID < formatedActual[j].ID
	})
	return formatedActual, nil
}

func (t *getListUserWithNestedFilter) CheckResult() (bool, error) {
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

func (t *getListUserWithNestedFilter) RunTest() (bool, error) {
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
		return true, nil
	}
}
