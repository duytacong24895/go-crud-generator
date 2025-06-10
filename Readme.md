# GO CRUD Generator

CRUD Generator is an open source project that generates RESTful APIs for GORM models and chi
**go version 1.24.1**

## Features

- Automatically generates RESTful APIs for GORM models
- Supports soft delete, create time, and update time
- Supports custom field tags
- Supports custom DTOs that are returned for get detail, get list, and error
- Supports custom middlewares

## Installation

To install CRUD Generator, run the following command:

`
go get -u github.com/duytacong24895/go-crud-generator
`

## How To use

In the following code, we are registering the model User (you can call a chain of RegisterModel to register a list of models), my custom middleware, and a custom struct that handles responses for the list and detail APIs, as well as error requests.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/duytacong24895/template-api-go/config"
	"github.com/duytacong24895/template-api-go/database"
	"github.com/duytacong24895/template-api-go/internal/models"
	crud_generator "github.com/duytacong24895/go-crud-generator"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	cf := config.Load()
	db := database.InitPg(cf.Databases)
  // init
  // RegisterModel: register Model to use
  // RegisterMiddleware : register your own middleware, you can register a chain middleware
  // RegisterDTOForGetDetail : define dto or struct that return out for api get detail 
  // RegisterDTOForGetList : define dto or struct that return out for api get list
  // RegisterDTOForError : define dto or struct that return out when server return an error
	crud_generator.NewCRUDGenerator(r, db).w
		RegisterModel(&models.User{}).
		RegisterMiddleware(
			middleware,
		).
		RegisterDTOForGetDetail(DTOGetDetail).
		RegisterDTOForGetList(DTOGetList).
		RegisterDTOForError(DTOError).
		Run()

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example middleware logic
		fmt.Println("Middleware executed")
		next.ServeHTTP(w, r)
	})
}

func DTOGetDetail(w http.ResponseWriter, r *http.Request, ref any) any {
	return map[string]any{
		"status":  200,
		"message": "Success",
		"data":    ref,
	}
}

func DTOGetList(w http.ResponseWriter, r *http.Request,
	ref any, total, page, pageSize uint) any {
	return map[string]any{
		"status":  200,
		"message": "Success",
		"data":    map[string]any{"list": ref, "total": total, "page": page, "page_size": pageSize},
	}
}

func DTOError(w http.ResponseWriter, r *http.Request,
	err error, msgErr string) any {
	return map[string]any{
		"status":  500,
		"message": msgErr,
		"data":    nil,
	}
}
```

Now you can call the API using a list of URLs and methods.
```
  - Method: Get, URL: .../crud/User/{id} to get detail user

  - Method: Get, URL: .../crud/User to get list user

  - Method: POST URL: .../crud/User to create new user

  - Method: DELETE URL: .../crud/User/{id} to delete one

  - Method: PUT URL.../crud/User/{id} to update one
```

# Get List Api
We also support paging, sorting, and filtering features

**Params**
- page # start from 1
- page_size # page length
- filter # read filter section for more detail
- order_by # same gorm syntax. Example: "id desc" or "id asc" read more at https://gorm.io/docs/query.html#Order

**Filters Query Params**

Each block must be wrapped in [] and include three parts: a column or nested block, an operation, and a value or nested block
```
1. ["age","gt","20"] #valid
  column: age
  operation: gt
  value: 20

2. [["age","gt","20"],"_and",["mature","eq","false"]] # valid
  nested block: ["age","gt","20"]
  operation: _and
  nested block: ["mature","eq","false"]

3. [["age","gt","20"]] #invalid syntax because it's wrapped by [[]] with only one block and no nested block
```

**Supported Operation**
```
	"eq":       "=",
	"gt":       ">",
	"lt":       "<",
	"gte":      ">=",
	"lte":      "<=",
	"ne":       "!=",
	"contain":  "like",
	"ncontain": "not like",
	"bw":       "between",
	"nbw":      "not between",
	"_null":    "is null",
	"_nnull":   "is not null",
	"_and":     "and",
	"_or":      "or",
```

**Example**:

```shell
curl --location --globoff 'localhost:8080/crud/User?page=1&page_size=10&filter=[[%22age%22%2C%22gt%22%2C%2220%22]%2C%22_and%22%2C[%22mature%22%2C%22eq%22%2C%22false%22]]&order_by=id'
```

# Update
Since we are using GORM to interact with the database, we fully adhere to GORM's rules. Let's take a look at the GORM update here
https://gorm.io/docs/update.html#Updates-multiple-columns

**Example**
```shell
curl --location --request PUT 'localhost:8080/crud/User/2' \
--header 'Content-Type: application/json' \
--data '{
	"age" :   51
}'
```

# Soft Delete, Created At, Updated At
We also support soft deletes and automatically manage timing fields in two ways.

1. You can use gorm.Model

``` go
type User struct {
	gorm.Model
	ID              int
	Name            string
	Email           string
	DOB             time.Time
	Age             int
	Phone           string
	Mature          bool
}
```


2. if you want to use your own fields, you can use our tag **crud_generator** to mark them

``` go
type User struct {
	ID              int
	Name            string
	Email           string
	DOB             time.Time
	Age             int
	Phone           string
	Mature          bool
	CustomCreatedAt time.Time `crud_generator:"create_time_field"`
	CustomUpdatedAt time.Time `crud_generator:"update_time_field"`
	CustomDeletedAt time.Time `crud_generator:"soft_delete_field"`
}
```
**Note: The data will be permanently deleted if there is no soft delete marked field.**

# Gen CRUD from DB
If you only have database, and there is no model struct. You can use gorm/gen to gen model struct from db. Then register them to crud_generator. Let follow the documentation here: https://gorm.io/gen/gen_tool.html

## Roadmap
- Write unit tests
- Support for gin, echo ...
- Support gen struct from db (Make out CRUD api from db)
- Support permissions
- Support Upload files

## How to contribute
I warmly welcome everyone to contribute to this project.

What can you do to contribute to the project?
- Request features
- Contribute code and create pull requests

