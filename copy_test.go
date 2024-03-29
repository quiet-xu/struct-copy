package struct_copy

import (
	"testing"
	"time"
)

type A1 struct {
	A  []string
	B  []C1
	C  []string
	D  string
	E  *[]C1
	F  Datetime
	G  C1
	H  []*string
	I  []string
	K  []*C1
	J  []C1
	L  []int32
	M  []int64
	N  []*int32
	AA string
}
type B1 struct {
	A []string
	B []C1
	C *[]string
	D *string
	E []C1
	//F Datetime
	G  C1
	H  []string
	I  []*string
	K  []C1
	J  []*C1
	L  []int
	M  []int
	N  []int
	AA string
}
type C1 struct {
	C1A string
}

func TestName(t *testing.T) {

	var one A1 //要复制的

	//one.C = []string{"1", "2"}
	//one.D = PStr("1")

	//one.D = "1"
	one.A = []string{"1", "2"}
	one.B = []C1{{C1A: "1"}, {C1A: "2"}}
	one.E = &[]C1{{C1A: "1"}, {C1A: "3"}}
	one.F = Datetime(time.Now())
	one.G = C1{C1A: "123123"}
	one.H = []*string{PStr("4"), PStr("5")}
	one.I = []string{"1", "2"}
	one.K = []*C1{{C1A: "23"}}
	one.J = []C1{{C1A: "2323"}, {C1A: "2323"}}
	one.L = []int32{1, 2, 3, 4}
	one.M = []int64{1, 2, 3, 4}
	one.N = []*int32{
		PInt32(1),
		PInt32(2),
	}
	var two B1 //被赋值的
	two.AA = "111"
	c := NewCopy(1 << 4)

	//err := c.SetOmitBlank().Copy(&two, one)
	//err := StructCopy(&two, one, 1<<4)
	//t.Log(err, two)
	err := c.SliceCopy(two.A, one.A)
	t.Log(err, two)

	//addr.Field(0).Set(reflect.ValueOf("1"))
	//t.Log(addr.Interface())
}
func PStr(a string) *string {
	return &a
}

func PInt(a int) *int {
	return &a
}
func PInt64(a int64) *int64 {
	return &a
}
func PInt32(a int32) *int32 {
	return &a
}
