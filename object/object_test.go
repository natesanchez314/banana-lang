package object

import "testing"

func TestStringDictKey(t *testing.T) {
	hello1 := &String{Val: "Hello World"}
	hello2 := &String{Val: "Hello World"}
	diff1 := &String{Val: "Goodbye World"}

	if hello1.DictKey() != hello2.DictKey() {
		t.Errorf("Strings with same content have different dict keys.")
	}
	if hello1.DictKey() == diff1.DictKey() {
		t.Errorf("Strings with different content have same dict keys.")
	}
}