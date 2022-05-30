package struct_copy

import (
	"testing"
	"time"
)

type A1 struct {
	D *string
	C []string
	E *[]C1
	A []string
	B []C1
	G C1
	F Datetime
}
type B1 struct {
	A []string
	B []C1
	E *[]C1
	C *[]string
	D *string
	F Datetime
	G C1
}
type C1 struct {
	A string
}

func TestName(t *testing.T) {
	var one A1 //要复制的

	one.C = []string{"1", "2"}
	one.D = PStr("1")
	one.A = []string{"1", "2"}
	one.B = []C1{{A: "1"}, {A: "2"}}
	one.E = &[]C1{{A: "1"}, {A: "3"}}
	one.F = Datetime(time.Now())
	one.G = C1{A: "123123"}
	var two B1 //被赋值的
	err := StructCopy(&two, one, 1<<4)

	t.Log(err, two)
	//addr.Field(0).Set(reflect.ValueOf("1"))
	//t.Log(addr.Interface())
}
func PStr(a string) *string {
	return &a
}
