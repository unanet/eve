package common

import (
	"bytes"
	"fmt"
)

type WhereArg func(*WhereClause)

type WhereCondition string

const (
	ORWhereCondition  WhereCondition = "OR"
	ANDWhereCondition WhereCondition = "AND"
)

func CheckWhereArgs(sql string, whereArgs []WhereArg) (string, []interface{}) {
	if len(whereArgs) == 0 {
		return sql, []interface{}{}
	}
	where := Where()
	for _, condition := range whereArgs {
		condition(where)
	}
	return where.AppendToSql(sql)
}

func getPositionalCounter(start int, count int) []interface{} {
	returnValues := []interface{}{}
	for index := start; index < count+start; index++ {
		returnValues = append(returnValues, fmt.Sprintf("$%d", index))
	}

	return returnValues
}

type Clause struct {
	condition WhereCondition
	value     string
	params    []interface{}
}

type WhereClause struct {
	clauses []Clause
}

func Where() *WhereClause {
	return &WhereClause{}
}

func (w *WhereClause) AddClause(value string, condition WhereCondition, params ...interface{}) *WhereClause {
	w.clauses = append(w.clauses, Clause{condition: condition, value: value, params: params})
	return w
}

func (w *WhereClause) Value() string {
	if len(w.clauses) == 0 {
		return ""
	}

	buffer := bytes.Buffer{}
	positionCounter := 1

	for index, element := range w.clauses {
		positionCounters := getPositionalCounter(positionCounter, len(element.params))
		tempValue := fmt.Sprintf(element.value, positionCounters...)
		positionCounter += len(element.params)
		if index == 0 {
			buffer.WriteString(fmt.Sprintf("WHERE %s", tempValue))
		} else {
			buffer.WriteString(fmt.Sprintf(" %s %s", element.condition, tempValue))
		}
	}

	return buffer.String()
}

func (w *WhereClause) AppendToSql(sql string) (string, []interface{}) {
	return fmt.Sprintf("%s %s", sql, w.Value()), w.Params()
}

func (w WhereClause) Params() []interface{} {
	returnParams := []interface{}{}
	for _, element := range w.clauses {
		returnParams = append(returnParams, element.params...)
	}
	return returnParams
}
