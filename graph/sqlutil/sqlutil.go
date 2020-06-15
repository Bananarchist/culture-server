package sqlutil

import (
	"fmt"
	"reflect"
	"strconv"
)

func fieldNotNil(input interface{}) bool {
	if input == nil {
		return false
	}
	switch reflect.TypeOf(input).Kind() {
	case reflect.Chan, reflect.Ptr, reflect.Map, reflect.Array, reflect.Slice:
		return !reflect.ValueOf(input).IsNil()
	}
	return true
}

func appendToSliceIfNotNil(testCase, val interface{}, list []interface{}) []interface{} {
	if fieldNotNil(testCase) {
		return append(list, val)
	}
	return list
}

// GetArgsArr returns a list without the nil *string args
func GetArgsArr(values ...interface{}) (list []interface{}) {
	for _, col := range values {
		list = appendToSliceIfNotNil(col, col, list)
	}
	return
}

// GetColumnsArr returns a list without the nil *string args
func GetColumnsArr(val2test map[interface{}]interface{}) (list []interface{}) {
	for v, t := range val2test {
		list = appendToSliceIfNotNil(t, v, list)
	}
	return
}

// SQLUpdateArgsString take a list of columns and output them as a "x=$1, y=$2" set string
func SQLUpdateArgsString(columns []interface{}) (r string) {
	for i, col := range columns {
		if i > 0 {
			r += ", "
		}
		r += col.(string) + "=$" + strconv.Itoa(i+1)
	}
	return
}

// SQLInsertArgsString take a number and produce a list of sql args
// eg: f(3) = "$1, $2, $3"
func SQLInsertArgsString(argc int) (r string) {
	accString := func(acc string, val int) string {
		return acc + "$" + strconv.Itoa(val)
	}
	if argc < 1 {
		return
	}
	l := make([]int, argc-1)
	for i := range l {
		r = accString(r, i+1)
		if argc > 1 {
			r += ", "
		}
	}
	r = accString(r, argc)
	return
}

//GetRowValuesStringForColumnsList returns a row eg "$1, $2"
func GetRowValuesStringForColumnsList(offset int, columns []string) (rowStr string) {
	for i := range columns {
		if i > 0 {
			rowStr += ", "
		}
		rowStr += fmt.Sprintf("$%d", offset+i+1)
	}
	return
}

//GetRowsStringForColumnsList returns a rows list eg "($1, $2), ($3, $4)"
func GetRowsStringForColumnsList(totalRows, currentRow int, columns []string) string {
	columnsLength := len(columns)
	if columnsLength == 0 {
		return ""
	}
	currentRowMultiplied := (currentRow - 1) * columnsLength
	rowStr := fmt.Sprintf("(%s)", GetRowValuesStringForColumnsList(currentRowMultiplied, columns))
	if currentRow == totalRows {
		return rowStr
	}
	return fmt.Sprintf("%s, %s", rowStr, GetRowsStringForColumnsList(totalRows, currentRow+1, columns))
}
