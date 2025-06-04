package core

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

const SepOfBetween = "::"

func convertToSqlOperator(op string) (string, error) {
	sqlOp, ok := SuportedOperators[op]
	if !ok {
		return "", fmt.Errorf("unsupported operator: %s", op)
	}
	return string(sqlOp), nil
}

type IFilter interface {
	BuildQuery(db *gorm.DB) (*gorm.DB, error)
	Load(filters string) error
	IsEmpty() bool
}

type filter struct {
	Conditions *Condition
	isEmpty    bool
}

type Condition struct {
	ColumnName string
	Value      string
	Operator   string
	Left       *Condition
	Right      *Condition
	tx         *gorm.DB
}

func (f *filter) BuildQuery(db *gorm.DB) (*gorm.DB, error) {
	if err := f.Conditions.BuildDiveQuery(db); err != nil {
		return nil, fmt.Errorf("failed to build dive query: %w", err)
	}

	var err error
	f.Conditions.tx, err = f.Combine(f.Conditions)
	if err != nil {
		return nil, fmt.Errorf("failed to combine conditions: %w", err)
	}
	return f.Conditions.tx, nil
}

func (f *filter) Combine(con *Condition) (*gorm.DB, error) {
	if con.ColumnName != "" && con.Value != "" {
		return con.tx, nil
	}

	if con.Operator != string(AndOperator) && con.Operator != string(OrOperator) {
		return nil, fmt.Errorf("unsupported operator: %s", con.Operator)
	}

	txLeft, err := f.Combine(con.Left)
	if err != nil {
		return nil, err
	}

	txRight, err := f.Combine(con.Right)
	if err != nil {
		return nil, err
	}

	if con.Operator == string(AndOperator) {
		return txLeft.Where(txRight), nil
	}
	return txLeft.Or(txRight), nil
}

func (c *Condition) BuildDiveQuery(db *gorm.DB) error {
	if c.Left != nil {
		c.Left.BuildDiveQuery(db)
	}
	if c.Right != nil {
		c.Right.BuildDiveQuery(db)
	}
	operator, err := convertToSqlOperator(c.Operator)
	if err != nil {
		return err
	}

	if c.ColumnName != "" && c.Value != "" {
		switch operator {
		case "=":
			c.tx = db.Where(fmt.Sprintf("%s = ?", c.ColumnName), c.Value)
		case ">":
			c.tx = db.Where(fmt.Sprintf("%s > ?", c.ColumnName), c.Value)
		case "<":
			c.tx = db.Where(fmt.Sprintf("%s < ?", c.ColumnName), c.Value)
		case ">=":
			c.tx = db.Where(fmt.Sprintf("%s >= ?", c.ColumnName), c.Value)
		case "<=":
			c.tx = db.Where(fmt.Sprintf("%s <= ?", c.ColumnName), c.Value)
		case "!=":
			c.tx = db.Where(fmt.Sprintf("%s != ?", c.ColumnName), c.Value)
		case "like":
			c.tx = db.Where(fmt.Sprintf("%s like ?", c.ColumnName), c.Value)
		case "not like":
			c.tx = db.Where(fmt.Sprintf("%s not like ?", c.ColumnName), c.Value)
		case "between":
			// Assuming Value is a string with two values separated by a comma
			values := strings.Split(c.Value, SepOfBetween)
			if len(values) != 2 {
				return fmt.Errorf("invalid value for between operator: %s", c.Value)
			}
			c.tx = db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", c.ColumnName), values[0], values[1])
		case "not between":
			// Assuming Value is a string with two values separated by a comma
			values := strings.Split(c.Value, ",")
			if len(values) != 2 {
				return fmt.Errorf("invalid value for not between operator: %s", c.Value)
			}
			c.tx = db.Where(fmt.Sprintf("%s NOT BETWEEN ? AND ?", c.ColumnName), values[0], values[1])
		case "is null":
			c.tx = db.Where(fmt.Sprintf("%s IS NULL", c.ColumnName))
		case "is not null":
			c.tx = db.Where(fmt.Sprintf("%s IS NOT NULL", c.ColumnName))
		default:
			return fmt.Errorf(`%s is unsupported operator or invalid input. Please noted that you can't use OR or AND operator without nested conditions. Example: ["name", "or", "age"] is invalid`, c.Operator)
		}
	}
	return nil
}

func (f *filter) loadCondition(node *Condition, inputArr []interface{}) error {
	isleave, err := isLeave(inputArr)
	if err != nil {
		return err
	}

	if isleave {
		node.ColumnName = inputArr[0].(string)
		node.Operator = inputArr[1].(string)
		node.Value = inputArr[2].(string)
		return nil
	}

	node.Left = &Condition{}
	node.Right = &Condition{}
	node.Operator = inputArr[1].(string)
	f.loadCondition(node.Left, inputArr[0].([]interface{}))
	f.loadCondition(node.Right, inputArr[2].([]interface{}))
	return nil
}

func isLeave(inputArr []interface{}) (bool, error) {
	if len(inputArr) != 3 {
		return false, fmt.Errorf("invalid filters length: %d", len(inputArr))
	}
	return !(reflect.TypeOf(inputArr[0]).Kind() == reflect.Slice || reflect.TypeOf(inputArr[2]).Kind() == reflect.Slice), nil
}

func (f *filter) Load(filters string) error {
	if filters == "" {
		f.isEmpty = true
		return nil
	}

	var inputArr []interface{}
	if err := json.Unmarshal([]byte(filters), &inputArr); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	if len(inputArr) == 0 {
		return fmt.Errorf("filters cannot be empty")
	}

	if err := f.loadCondition(f.Conditions, inputArr); err != nil {
		return fmt.Errorf("error loading conditions: %w", err)
	}
	return nil
}

func (f *filter) IsEmpty() bool {
	return f.isEmpty
}

func NewFilter() IFilter {
	return &filter{
		Conditions: &Condition{},
	}
}
