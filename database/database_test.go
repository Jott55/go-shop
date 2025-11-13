package database

import "testing"

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
