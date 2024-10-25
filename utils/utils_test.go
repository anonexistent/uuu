package utils_test

//	go test -v ./...

import (
	"testing"
	"uuu/utils"
)

func TestFoo10(t *testing.T) {
	t.Skip()
	if utils.Foo() != 10 {
		t.Fail()
	}
}
func TestFoo11(t *testing.T) {
	t.Skip()
	if utils.Foo() != 11 {
		t.Fail()
	}
}

func TestTestStatus(t *testing.T) {
	//	example close db connection
	defer func() {
		t.Log("defer")
	}()

	t.Log(1)
	t.Fail()
	t.Log(2)
	t.FailNow()
	t.Log(3)

}
