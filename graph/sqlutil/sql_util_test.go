package sqlutil

import (
	"testing"
)

func Test_FieldNotNil_WithNil(t *testing.T) {
	expected := false
	var input *string
	got := fieldNotNil(input)
	if got != expected {
		t.Errorf("fieldNotNil(nil) != false, instead == %t", got)
	}
}

func Test_FieldNotNil_WithString(t *testing.T) {
	expected := true
	var fakeo interface{}
	fakeo = "hello"
	got := fieldNotNil(&fakeo)
	if got != expected {
		t.Errorf("fieldNotNil(%s) != true, instead == %t", fakeo, got)
	}
}

func Test_AppendToSliceIfNotNil_WithNilReturnsEmpty(t *testing.T) {
	expected := make([]interface{}, 0)
	fakeo := make([]interface{}, 0)
	got := appendToSliceIfNotNil(nil, nil, fakeo)
	if len(got) != len(expected) {
		t.Errorf("appendToSliceIfNotNil(nil, []) != [], instead == %s", got)
	}
}

func Test_AppendToSliceIfNotNil_WithStringReturnsString(t *testing.T) {
	expected := make([]interface{}, 1)
	expected[0] = "rat"
	fakeo := make([]interface{}, 0)
	var faker interface{}
	faker = "rat"
	got := appendToSliceIfNotNil(&faker, faker, fakeo)
	if got[0] != expected[0] {
		t.Errorf("appendToSliceIfNotNil(nil, []) != [%s], instead == %s", expected[0], got)
	}
}

func Test_AppendToSliceIfNotNil_WithStringAndArrReturns2xString(t *testing.T) {
	expected := make([]interface{}, 2)
	expected[0] = "rat"
	expected[1] = "cat"
	fakeo := make([]interface{}, 1)
	fakeo[0] = "rat"
	var faker interface{}
	faker = "cat"
	got := appendToSliceIfNotNil(&faker, faker, fakeo)
	if got[0] != expected[0] || got[1] != expected[1] {
		t.Errorf("appendToSliceIfNotNil(%s, [%s]) != [%s, %s], instead == %s", faker, fakeo[0], expected[0], expected[1], got)
	}
}

func Test_GetArgsArr_ReturnsEmptyArrIfAllNil(t *testing.T) {
	expected := 0
	got := len(GetArgsArr(nil, nil))
	if got != expected {
		t.Errorf("GetArgsArr(nil, nil) length != 0, instead == %d", got)
	}
}

func Test_GetArgsArr_ReturnsEmptyArrIfNoValues(t *testing.T) {
	expected := 0
	got := len(GetArgsArr())
	if got != expected {
		t.Errorf("GetArgsArr() length != 0, instead == %d", got)
	}
}

func Test_GetArgsArr_ReturnsAListWithoutNils(t *testing.T) {
	expected := 2
	var faker, fakeo interface{}
	fakeo = "cat"
	faker = "rat"
	got := len(GetArgsArr(&fakeo, nil, &faker))
	if got != expected {
		t.Errorf("GetArgsArr('cat', nil, 'rat') length != 2, instead == %d", got)
	}
}

func Test_GetColumnsArr_ReturnsAListWithoutNils(t *testing.T) {
	expected := 2
	col1 := "cat"
	col2 := "rat"
	col3 := "bat"
	got := len(GetColumnsArr(map[interface{}]interface{}{col1: &col1, col2: nil, col3: &col3}))
	if got != expected {
		t.Errorf("GetColumnsArr({'%s':'%s', '%s':nil, '%s':'%s'}) length != 2, instead == %d", col1, col1, col2, col3, col3, got)
	}
}

func Test_SQLUpdateArgsString_ReturnsEmptyStringIfNil(t *testing.T) {
	expected := 0
	got := len(SQLUpdateArgsString([]interface{}{}))
	if got != expected {
		t.Errorf("SQLUpdateArgsString([]string{}) did not return len=0, returned %d", got)
	}
}

func Test_SQLUpdateArgsString_ReturnsTwoSetArgs(t *testing.T) {
	expected := "cat=$1, rat=$2"
	got := SQLUpdateArgsString([]interface{}{"cat", "rat"})
	if got != expected {
		t.Errorf("SQLUpdateArgsString([]string{\"cat\", \"rat\"}) did not return %s, returned %s", expected, got)
	}
}

func Test_SQLUpdateArgsString_ReturnsOneSetArgNoComma(t *testing.T) {
	expected := "cat=$1"
	got := SQLUpdateArgsString([]interface{}{"cat"})
	if got != expected {
		t.Errorf("SQLUpdateArgsString([]string{\"cat\"}) did not return %s, returned %s", expected, got)
	}
}

func Test_SQLInsertArgsString_ReturnsEmptyStringIf0(t *testing.T) {
	expected := 0
	got := len(SQLInsertArgsString(0))
	if got != expected {
		t.Errorf("SQLInsertArgsString(0) did not return len=%d, returned %d", expected, got)
	}
}

func Test_SQLInsertArgsString_ReturnsTwoSetArgs(t *testing.T) {
	expected := "$1, $2, $3"
	got := SQLInsertArgsString(3)
	if got != expected {
		t.Errorf("SQLInsertArgsString([]string{3}) did not return %s, returned %s", expected, got)
	}
}

func Test_SQLInsertArgsString_ReturnsOneSetArgNoComma(t *testing.T) {
	expected := "$1"
	got := SQLInsertArgsString(1)
	if got != expected {
		t.Errorf("SQLInsertArgsString([]string{1}) did not return %s, returned %s", expected, got)
	}
}

func Test_GetRowValuesStringForColumnsList_ReturnsEmptyStringWithEmpty(t *testing.T) {
	expected := ""
	got := GetRowValuesStringForColumnsList(0, []string{})
	if got != expected {
		t.Errorf("Returned %s expected %s", got, expected)

	}
}

func Test_GetRowValuesStringForColumnsList_Cash1Cash2(t *testing.T) {
	expected := "$1, $2"
	got := GetRowValuesStringForColumnsList(0, []string{"cat", "rat"})
	if got != expected {
		t.Errorf("Returned %s expected %s", got, expected)

	}
}

func Test_GetRowValuesStringForColumnsList_Cash5Cash6Cash7(t *testing.T) {
	expected := "$5, $6, $7"
	got := GetRowValuesStringForColumnsList(4, []string{"cat", "rat", "bat"})
	if got != expected {
		t.Errorf("Returned %s expected %s", got, expected)
	}
}

func Test_GetRowsStringForColumnsList_NilStringForNilArgs(t *testing.T) {
	expected := ""
	got := GetRowsStringForColumnsList(1, 1, []string{})
	if got != expected {
		t.Errorf("Returned %s expected %s", got, expected)
	}
}

func Test_GetRowsStringForColumnsList_C1C2C3_C4C5C6(t *testing.T) {
	expected := "($1, $2, $3), ($4, $5, $6)"
	got := GetRowsStringForColumnsList(2, 1, []string{"cat", "rat", "bat"})
	if got != expected {
		t.Errorf("Returned %s expected %s", got, expected)
	}
}

func Test_GetRowsStringForColumnsList_C1C2C3(t *testing.T) {
	expected := "($1, $2, $3)"
	got := GetRowsStringForColumnsList(1, 1, []string{"cat", "rat", "bat"})
	if got != expected {
		t.Errorf("Returned %s expected %s", got, expected)
	}
}
