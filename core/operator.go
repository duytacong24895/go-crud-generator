package core

type Operator string

const (
	AndOperator Operator = "_and"
	OrOperator  Operator = "_or"
)

var SuportedOperators = map[string]Operator{
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
}
