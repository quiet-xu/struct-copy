# struct-copy

结构体复制，代替繁琐的手动复制的方式

**有什么功能?**

1.支持 指针 转 非指针
例如 `` A *string  => A string  ``

2.支持 非指针 转 指针
例如 `` A string => A *string ``

3.支持 非本地包 转 本地包
例如 `` A xxx.Object => A Object ``

4.支持 本地宝 转 非本地包
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



**使用说明**

1. 如何下载？
   ```bash
    go get github.com/quiet-xu/struct-copy@latest
   ```
2. 如何使用？
   ```
   StructCopy(&被赋值体, 赋值体, copy深度(int))
   
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
   }
   func PStr(a string) *string {
     return &a
   }

   ```