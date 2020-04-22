package data

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type WhereArg func(*WhereClause)

type LogicalOperator string

const (
	ORLogicalOperator  LogicalOperator = "OR"
	ANDLogicalOperator LogicalOperator = "AND"
)

func CheckWhereArgs(sql string, whereArgs []WhereArg) (string, []interface{}) {
	if len(whereArgs) == 0 {
		return sql, []interface{}{}
	}
	where := &WhereClause{}
	for _, condition := range whereArgs {
		condition(where)
	}
	return where.AppendToSql(sql)
}

func Where(key string, value interface{}) WhereArg {
	return func(clause *WhereClause) {
		clause.AddClause(fmt.Sprintf("%s=?", key), ANDLogicalOperator, value)
	}
}

func WhereIn(key string, values []interface{}) WhereArg {
	return func(clause *WhereClause) {
		clause.AddClause(fmt.Sprintf("%s in (?)", key), ANDLogicalOperator, values...)
	}
}

func AndWhere(key string, value interface{}) WhereArg {
	return func(clause *WhereClause) {
		clause.AddClause(fmt.Sprintf("%s=?", key), ANDLogicalOperator, value)
	}
}

func OrWhere(key string, value interface{}) WhereArg {
	return func(clause *WhereClause) {
		clause.AddClause(fmt.Sprintf("%s=?", key), ORLogicalOperator, value)
	}
}

type Clause struct {
	operator LogicalOperator
	value    string
	params   []interface{}
}

type WhereClause struct {
	clauses []Clause
}

func (w *WhereClause) AddClause(value string, operator LogicalOperator, params ...interface{}) *WhereClause {
	w.clauses = append(w.clauses, Clause{operator: operator, value: value, params: params})
	return w
}

func (w *WhereClause) Value() string {
	if len(w.clauses) == 0 {
		return ""
	}

	buffer := bytes.Buffer{}

	paramCount := 0
	for index, element := range w.clauses {
		placeHolderRegex := regexp.MustCompile("\\?")
		placeHolderCount := placeHolderRegex.FindAllStringIndex(element.value, -1)
		var tempValue string
		if len(placeHolderCount) == 1 {
			var replacementParams []string
			for range element.params {
				paramCount += 1
				replacementParams = append(replacementParams, fmt.Sprintf("$%d", paramCount))
			}
			tempValue = strings.Replace(element.value, "?", strings.Join(replacementParams, ","), 1)
		} else {
			tempValue = placeHolderRegex.ReplaceAllStringFunc(element.value, func(s string) string {
				paramCount += 1
				return fmt.Sprintf("$%d", paramCount)
			})
		}
		if index == 0 {
			buffer.WriteString(fmt.Sprintf("WHERE %s", tempValue))
		} else {
			buffer.WriteString(fmt.Sprintf(" %s %s", element.operator, tempValue))
		}
	}

	return buffer.String()
}

func (w *WhereClause) AppendToSql(sql string) (string, []interface{}) {
	return fmt.Sprintf("%s %s", sql, w.Value()), w.Params()
}

func (w WhereClause) Params() []interface{} {
	var returnParams []interface{}
	for _, element := range w.clauses {
		returnParams = append(returnParams, element.params...)
	}
	return returnParams
}
