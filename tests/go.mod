module github.com/duytacong24895/go-crud-generator/tests

go 1.24.1

require (
	github.com/duytacong24895/go-crud-generator v0.0.0-20241022120000-abcdef123456
	github.com/go-chi/chi/v5 v5.2.1
	github.com/go-resty/resty/v2 v2.16.5
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.30.0
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)

replace github.com/duytacong24895/go-crud-generator => ./..
