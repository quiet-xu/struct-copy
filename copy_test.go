package struct_copy

import (
	"testing"
	"time"
)

type A1 struct {
	A []string
	B []C1
	C []string
	D string
	E *[]C1
	F Datetime
	G C1
	H []*string
	I []string
	K []*C1
	J []C1
}
type B1 struct {
	A []string
	B []C1
	C *[]string
	D *string
	E []C1
	//F Datetime
	G C1
	H []string
	I []*string
	K []C1
	J []*C1
}
type C1 struct {
	C1A string
}

func TestName(t *testing.T) {
	var one A1 //要复制的

	//one.C = []string{"1", "2"}
	//one.D = PStr("1")
	one.D = "1"
	one.A = []string{"1", "2"}
	one.B = []C1{{C1A: "1"}, {C1A: "2"}}
	one.E = &[]C1{{C1A: "1"}, {C1A: "3"}}
	one.F = Datetime(time.Now())
	one.G = C1{C1A: "123123"}
	one.H = []*string{PStr("4"), PStr("5")}
	one.I = []string{"1", "2"}
	one.K = []*C1{{C1A: "23"}}
	one.J = []C1{{C1A: "2323"}, {C1A: "2323"}}
	var two B1 //被赋值的
	err := StructCopy(&two, one, 1<<4)
	t.Log(err, two)
	//addr.Field(0).Set(reflect.ValueOf("1"))
	//t.Log(addr.Interface())
}
func PStr(a string) *string {
	return &a
}
