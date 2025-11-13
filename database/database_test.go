package database

import (
	"testing"
)

type TestIn interface {
	string
	int
}

type TestStruct struct {
	Name string
	Age  int
}

func TestGetStructNames(t *testing.T) {
	names := getStructNames[TestStruct]()

	want := []string{"Name", "Age"}

	if names[0] != want[0] || names[1] != want[1] {
		t.Errorf(`getStructNames(TestStruct{}) = %v, want %v`, names, want)
	}
}

func TestGetStructValues(t *testing.T) {
	var test TestStruct

	test.Age = 1
	test.Name = "jack"

	f := getStructValues(&test)

	if f.fieldName[0] != "Name" || f.fieldValue[0] != test.Name {
		t.Errorf(`getStructValues(&test) = %s, %s,  want %s, %s`, f.fieldName[0], f.fieldValue[0], "Name", test.Name)
	}
}

func TestGenericInsert(t *testing.T) {

}
