package struct_copy

import (
	"errors"
	"reflect"
)

var OtherDateType = reflect.TypeOf((*ConvCopy)(nil)).Elem()

// Copy （被赋值的指针参数，需要被复制的值） 默认 1<<4 个深度
func Copy(dst any, src any) (err error) {
	return StructCopy(dst, src, 1<<4)
}

// StructCopy 指针接收 / 非指针发送 / 深度 单位 个
func StructCopy(dst any, src any, depth int) (err error) {
	if src == nil {
		return
	}
	if dst == nil {
		return
	}
	//类型判断并声明  防重复反射
	dstType := reflect.TypeOf(dst)
	srcType := reflect.TypeOf(src)
	if dstType.Kind() != reflect.Ptr && dstType.Elem().Kind() != reflect.Struct {
		err = errors.New("dst value not a struct pointer")
		return
	}
	if srcType.Kind() == reflect.Ptr {
		srcV := reflect.ValueOf(src)
		if srcV.IsZero() {
			return
		}
		src = reflect.ValueOf(src).Elem().Interface()
		srcType = srcType.Elem()

	}
	if srcType.Kind() != reflect.Struct {
		err = errors.New("src value not a struct")
		return
	}
	//使用src构造map 一开始给予空间防止map过大超载加快速度 =>
	v := reflect.ValueOf(dst).Elem()
	err = structChildCopy(v, src, depth, 0)
	return
}
func checkTypeAndConv(newValue reflect.Value, dstField reflect.Value) (fValue any, err error) {
	if newValue.Kind() == reflect.Ptr {
		newValue = newValue.Elem() // 有数据的
	}
	if dstField.Kind() == reflect.Ptr {
		if dstField.Elem() == (reflect.Value{}) {
			dstField = reflect.New(dstField.Type().Elem())

		} else {
			dstField = dstField.Elem()
		}

	}
	switch newValue.Interface().(type) {
	case int32:
		switch dstField.Interface().(type) {
		case int:
			fValue = int(newValue.Interface().(int32))
		case int64:
			fValue = int64(newValue.Interface().(int32))
		default:
			fValue = newValue.Interface()
		}
	case int64:
		switch dstField.Interface().(type) {
		case int:
			fValue = int(newValue.Interface().(int64))
		case int32:
			fValue = int32(newValue.Interface().(int64))
		default:
			fValue = newValue.Interface()
		}
	case int:
		switch dstField.Interface().(type) {
		case int32:
			fValue = int32(newValue.Interface().(int))
		case int64:
			fValue = int64(newValue.Interface().(int))
		default:
			fValue = newValue.Interface()
		}
	default:
		switch newValue.Kind() {
		case reflect.Bool:
			fValue = newValue.Bool()
		case reflect.Int:
			fValue = int(newValue.Int())
		case reflect.Int8:
			fValue = int8(newValue.Int())
		case reflect.Int16:
			fValue = int16(newValue.Int())
		case reflect.Int32:
			fValue = int32(newValue.Int())
		case reflect.Int64:
			fValue = newValue.Int()
		case reflect.Uint:
			fValue = uint(newValue.Uint())
		case reflect.Uint8:
			fValue = uint8(newValue.Uint())
		case reflect.Uint16:
			fValue = uint16(newValue.Uint())
		case reflect.Uint32:
			fValue = uint32(newValue.Uint())
		case reflect.Uint64:
			fValue = newValue.Uint()
		case reflect.Float32:
			fValue = float32(newValue.Float())
		case reflect.Float64:
			fValue = float64(newValue.Float())
		case reflect.String:
			fValue = newValue.String()

		default:
			fValue = newValue.Interface()
		}

	}
	return
}

// structChildCopy 结构体子层copy
func structChildCopy(dst reflect.Value, src any, depth int, current int) (err error) {
	if depth == current {
		return
	}
	convDst := dst //转换后的目标,不用于set
	//如果是个指针则给予新的地址
	if dst.Kind() == reflect.Ptr {
		convDst = reflect.New(dst.Type().Elem())
	}

	valueMap := make(map[string]any)
	getValMap(dst, src, valueMap, depth, 0)
	if convDst.Kind() == reflect.Ptr {
		dst.Set(convDst)
		convDst = convDst.Elem()
	}

	numDstField := reflect.TypeOf(convDst.Interface()).NumField()
	convType := reflect.TypeOf(convDst.Interface())
	for index := 0; index < numDstField; index++ {
		fieldType := convType.Field(index)
		fieldName := fieldType.Name
		var fieldValue reflect.Value

		if dst.Kind() == reflect.Ptr {

			fieldValue = dst.Elem().Field(index)
		} else {
			fieldValue = dst.Field(index)
		}
		//若非暴露的,则不需要set
		if !fieldType.IsExported() {
			continue
		}
		if valueMap[fieldName] == nil {
			if !fieldType.Anonymous {
				continue
			}
		}
		srcVal := valueMap[fieldName]
		//如果是匿名嵌套字段，则要判断src里是否存在同样的内容的结构
		if fieldType.Anonymous && srcVal == nil {
			err = structChildCopy(fieldValue, src, depth, current)
			if err != nil {
				return
			}
			continue
		}
		if isBlank(reflect.ValueOf(srcVal)) {
			continue
		}
		if fieldValue.Type() == reflect.ValueOf(srcVal).Type() {
			fieldValue.Set(reflect.ValueOf(srcVal))
			continue
		}
		valueType := reflect.TypeOf(srcVal).Kind()
		var fValue any
		switch valueType {
		case reflect.Struct:
			//if reflect
			if reflect.TypeOf(srcVal).Implements(OtherDateType) {
				switch fieldValue.Interface().(type) {
				case string:
					fValue = srcVal.(ConvCopy).String()
					setDstVal(fieldValue, reflect.ValueOf(fValue))
				}
			} else {
				structCurrent := current + 1
				err = structChildCopy(fieldValue, srcVal, depth, structCurrent)
				if err != nil {
					return
				}
			}
		case reflect.Slice:
			if reflect.ValueOf(srcVal).Len() == 0 {
				continue
			}
			sliceCurrent := current + 1
			err = sliceCopy(fieldValue, srcVal, depth, sliceCurrent)
			if err != nil {
				return
			}
		default:

			fValue, err = checkTypeAndConv(reflect.ValueOf(srcVal), convDst.Field(index))
			if err != nil {
				return
			}
			if reflect.TypeOf(fieldValue.Interface()).Implements(OtherDateType) {
				switch fValue.(type) {
				case string:
					fValue = fieldValue.Interface().(ConvCopy).CopyStr(fValue.(string))
				}
			}
			if fieldValue.Kind() == reflect.Ptr {
				//如果是个指针，则操作该指针给予新的地址
				item := reflect.New(fieldValue.Type().Elem())
				item.Elem().Set(reflect.ValueOf(fValue))
				setDstVal(fieldValue, item)
				continue
			}
			setDstVal(fieldValue, reflect.ValueOf(fValue))
		}
	}
	return
}

func getValMap(dst reflect.Value, src any, valMap map[string]any, depth int, current int) {
	if reflect.ValueOf(src).Kind() == reflect.Ptr {
		src = reflect.ValueOf(src).Elem().Interface()
	}
	numSrcField := reflect.TypeOf(src).NumField()
	if !reflect.ValueOf(src).IsValid() {
		return
	}
	if valMap == nil {
		valMap = make(map[string]any)
	}
	for index := 0; index < numSrcField; index++ {
		srcTp := reflect.TypeOf(src).Field(index)
		fieldName := srcTp.Name
		item := reflect.ValueOf(src).Field(index)
		if !reflect.TypeOf(src).Field(index).IsExported() {
			continue
		}
		//如果被粘贴的目标是个匿名结构，并且被赋值目标中没有该字段
		if srcTp.Anonymous && dst.FieldByName(fieldName) == (reflect.Value{}) {
			current++
			getValMap(dst, item.Interface(), valMap, depth, current)
		}
		if item.Kind() == reflect.Ptr {
			if item.IsZero() || item.IsNil() {
				//空指针跳过
				continue
			}
			item = item.Elem()
		}
		fieldValue := item.Interface()
		if valMap[fieldName] == nil {
			valMap[fieldName] = fieldValue
		}
	}
}

func sliceCopy(dst reflect.Value, src any, depth int, current int) (err error) {
	if depth == current {
		return
	}
	var basicSlice int
	switch src.(type) {
	case []*string:
		basicSlice++
	case []*int64:
		basicSlice++
	case []*int32:
		basicSlice++
	case []*int:
		basicSlice++
	case []*float64:
		basicSlice++
	case []*float32:
		basicSlice++
	case []*byte:
		basicSlice++
	case []string:
		basicSlice++
	case []int64:
		basicSlice++
	case []int32:
		basicSlice++
	case []int:
		basicSlice++
	case []float64:
		basicSlice++
	case []float32:
		basicSlice++
	case []byte:
		basicSlice++
	}

	//如果是个指针则给予新的地址
	convDst := dst
	if dst.Kind() == reflect.Ptr {
		convDst = reflect.New(dst.Type().Elem())
	}
	if reflect.ValueOf(src).Kind() == reflect.Ptr {
		src = reflect.ValueOf(src).Elem().Interface()
	}
	//基础数组赋值
	if basicSlice > 0 {
		fValue := src
		if convDst.Type() == reflect.ValueOf(fValue).Type() {
			convDst.Set(reflect.ValueOf(fValue))
			dst.Set(convDst)
			return
		}
		fv := reflect.ValueOf(fValue)
		convPrtType := 0
		if convDst.Kind() == reflect.Ptr {
			convPrtType = 1
		} else if convDst.Type().Elem().Kind() == reflect.Ptr {
			convPrtType = 2
		}
		for i := 0; i < fv.Len(); i++ {
			var item reflect.Value
			if convPrtType == 1 {
				item = reflect.New(convDst.Elem().Type().Elem())
			} else if convPrtType == 2 {
				item = reflect.New(convDst.Type().Elem().Elem())
			} else {
				item = reflect.New(convDst.Type().Elem())
			}
			var f any
			f, err = checkTypeAndConv(reflect.ValueOf(fv.Index(i).Interface()), item)
			if err != nil {
				return
			}
			if convPrtType == 1 {
				item.Elem().Set(reflect.ValueOf(f))
				//item.Elem().Set(reflect.ValueOf(fv.Index(i).Interface()))
				convDst.Elem().Set(reflect.Append(convDst.Elem(), item.Elem()))
			} else if convPrtType == 2 {
				item.Elem().Set(reflect.ValueOf(f))
				convDst.Set(reflect.Append(convDst, item))
			} else {
				item.Elem().Set(reflect.ValueOf(f))
				convDst.Set(reflect.Append(convDst, item.Elem()))
			}
		}
		if dst.Type() == convDst.Type() {
			dst.Set(convDst)
		}
		return
	}
	//结构体赋值
	dstItemType := convDst.Type().Elem()
	if dst.Kind() == reflect.Ptr {
		dstItemType = dstItemType.Elem()
	}
	for i := 0; i < reflect.ValueOf(src).Len(); i++ {
		item := reflect.New(dstItemType) //new([]xxx)
		err = structChildCopy(item.Elem(), reflect.ValueOf(src).Index(i).Interface(), depth, 0)
		if err != nil {
			return
		}
		if dst.Kind() == reflect.Ptr {
			convDst.Elem().Set(reflect.Append(convDst.Elem(), item.Elem()))
		} else {
			dst.Set(reflect.Append(dst, item.Elem()))
		}

	}
	if dst.Kind() == reflect.Ptr {
		dst.Set(convDst)
	}
	return
}

func CopyMap(dst interface{}, valueMap map[string]interface{}) (err error) {
	dstType := reflect.TypeOf(dst)
	if dstType.Kind() != reflect.Ptr && dstType.Elem().Kind() != reflect.Struct {
		err = errors.New("dst value not a struct pointer")
		return
	}
	numDstField := dstType.Elem().NumField()
	for index := 0; index < numDstField; index++ {
		fieldName := dstType.Elem().Field(index).Name
		if _, has := valueMap[fieldName]; !has {
			continue
		}
		//需要被copy的类的判断 = >
		valueType := reflect.TypeOf(valueMap[fieldName]).Kind()
		switch valueType {
		case reflect.Struct:
			continue
		case reflect.Slice, reflect.Array:
			continue
		default:
			reflect.ValueOf(dst).Elem().Field(index).Set(reflect.ValueOf(valueMap[fieldName]))
		}
	}
	return
}

func setDstVal(dst reflect.Value, src reflect.Value) {
	if !dst.CanSet() {
		return
	}
	if src.Type() != dst.Type() {
		switch dst.Kind() {
		case reflect.Bool:
			dst.SetBool(src.Bool())
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			dst.SetInt(src.Int())
		case reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			dst.SetUint(src.Uint())
		case reflect.Float32,
			reflect.Float64:
			dst.SetFloat(src.Float())

		case reflect.String:
			dst.SetString(src.String())
		default:
			return
		}
	} else {
		dst.Set(src)
	}

}

// isBlank 非空校验
func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	case reflect.Slice:
		return value.Len() == 0
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

type ConvCopy interface {
	CopyStr(string) any
	String() string
}
