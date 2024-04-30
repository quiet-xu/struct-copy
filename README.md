# struct-copy

结构体复制，代替繁琐的手动复制的方式

**有什么功能?**

1.支持 指针 转 非指针
例如 `` A *string  => A string  ``

2.支持 非指针 转 指针
例如 `` A string => A *string ``

3.支持 非本地包 转 本地包
例如 `` A xxx.Object => A Object ``

4.支持 本地包 转 非本地包
例如 `` A Object => A xxx.Object ``

5.支持 多层copy
   ```bash
    例如 X 是以下结构
    type X struct {
      Obj Object //第一层
    }
    type Object struct{
      Obj2 Object2 //第二层
    }
    type Object2 struct{
      Str string //第三层
    }
    根据 ``depth``深度来copy
    
   ```
6. 支持自定义基础类型与基础类型互转


**使用说明**

1. 如何下载？
   ```bash
    go get github.com/quiet-xu/struct-copy@latest
   ```
   2. 如何使用？
      ``` golang
      Copy(&被赋值体, 赋值体)

      type Number uint16
      type Str string
      type Fl float32
   
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
      BB Number
      CC string
      DD Fl
      }
      func TestName(t *testing.T) {
            var one A1 //要复制的
            one.C = []string{"1", "2"}
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
            one.L = []int32{1, 2, 3, 4}
            one.M = []int64{1, 2, 3, 4}
            one.N = []*int32{
            PInt32(1),
            PInt32(2),
            }
            one.BB = 31
            one.CC = "奥斯卡"
            one.DD = 33.11
            var two B1 //被赋值的
            two.AA = "111"
            err := Copy(&two, one)
            t.Log(err)
      }
      ```