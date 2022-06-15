package struct_copy

import (
	"errors"
	"reflect"
)

//StructCopy 指针接收 / 非指针发送 / 深度 单位 个
func StructCopy(dst interface{}, src interface{}, depth int) (err error) {
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
	}
	if srcType.Kind() != reflect.Struct {
		err = errors.New("src value not a struct")
		return
	}
	//使用src构造map 一开始给予空间防止map过大超载加快速度 =>
	numSrcField := srcType.NumField()
	valueMap := make(map[string]interface{}, numSrcField)
	for index := 0; index < numSrcField; index++ {
		item := reflect.ValueOf(src).Field(index)
		if !srcType.Field(index).IsExported() {
			continue
		}
		if item.Kind() == reflect.Ptr {
			if item.IsZero() {
				//空指针跳过
				continue
			}
			item = item.Elem()
		}
		fieldName := srcType.Field(index).Name
		fieldValue := item.Interface()
		valueMap[fieldName] = fieldValue
	}

	//构造dst = >
	numDstField := dstType.Elem().NumField()
	for index := 0; index < numDstField; index++ {
		fieldName := dstType.Elem().Field(index).Name
		if valueMap[fieldName] == nil || isBlank(reflect.ValueOf(valueMap[fieldName])) {
			continue
		}
		fieldValue := reflect.ValueOf(dst).Elem().Field(index)
		valueType := reflect.TypeOf(valueMap[fieldName]).Kind()
		if fieldValue.Type() == reflect.ValueOf(valueMap[fieldName]).Type() {
			fieldValue.Set(reflect.ValueOf(valueMap[fieldName]))
			continue
		}
		switch valueType {
		case reflect.Struct:

			err = structChildCopy(fieldValue, valueMap[fieldName], depth, 0)
			if err != nil {
				return
			}
		case reflect.Slice, reflect.Array:
			if reflect.ValueOf(valueMap[fieldName]).Len() == 0 {
				continue
			}
			err = sliceCopy(fieldValue, valueMap[fieldName], depth, 0)
			if err != nil {
				return
			}
		default:
			var fValue interface{}
			fValue, err = checkTypeAndConv(reflect.ValueOf(valueMap[fieldName]), reflect.ValueOf(dst))
			if err != nil {
				return
			}
			if fieldValue.Kind() == reflect.Ptr {
				//如果是个指针，则操作该指针给予新的地址
				item := reflect.New(fieldValue.Type().Elem())
				item.Elem().Set(reflect.ValueOf(fValue))
				fieldValue.Set(item)
				continue
			}
			reflect.ValueOf(dst).Elem().Field(index).Set(reflect.ValueOf(fValue))
		}

		//判别类型处理

	}

	return
}
func checkTypeAndConv(newValue reflect.Value, dstField reflect.Value) (fValue interface{}, err error) {
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
		fValue = newValue.Interface()
	}
	return
}

//structChildCopy 结构体子层copy
func structChildCopy(dst reflect.Value, src interface{}, depth int, current int) (err error) {
	if depth == current {
		return
	}
	convDst := dst //转换后的目标,不用于set
	//如果是个指针则给予新的地址
	if dst.Kind() == reflect.Ptr {
		convDst = reflect.New(dst.Type().Elem())
	}
	if reflect.ValueOf(src).Kind() == reflect.Ptr {
		src = reflect.ValueOf(src).Elem().Interface()
	}

	numSrcField := reflect.TypeOf(src).NumField()
	valueMap := make(map[string]interface{}, numSrcField)

	for index := 0; index < numSrcField; index++ {
		fieldName := reflect.TypeOf(src).Field(index).Name
		var fieldValue interface{}
		if !reflect.TypeOf(src).Field(index).IsExported() {
			continue
		} else {
			fieldValue = reflect.ValueOf(src).Field(index).Interface()
		}
		valueMap[fieldName] = fieldValue
	}
	numDstField := reflect.TypeOf(convDst.Interface()).NumField()
	for index := 0; index < numDstField; index++ {
		fieldName := reflect.TypeOf(convDst.Interface()).Field(index).Name
		if valueMap[fieldName] == nil || isBlank(reflect.ValueOf(valueMap[fieldName])) {
			continue
		} else {
			fieldValue := dst.Field(index)
			if fieldValue.Type() == reflect.ValueOf(valueMap[fieldName]).Type() {
				fieldValue.Set(reflect.ValueOf(valueMap[fieldName]))
				continue
			}
			if !reflect.TypeOf(src).Field(index).IsExported() {
				continue
			}
			valueType := reflect.TypeOf(valueMap[fieldName]).Kind()
			switch valueType {
			case reflect.Struct:
				structCurrent := current + 1
				err = structChildCopy(fieldValue, valueMap[fieldName], depth, structCurrent)
				if err != nil {
					return
				}
			case reflect.Slice:
				sliceCurrent := current + 1
				err = sliceCopy(fieldValue, valueMap[fieldName], depth, sliceCurrent)
				if err != nil {
					return
				}
			default:
				var fValue interface{}
				fValue, err = checkTypeAndConv(reflect.ValueOf(valueMap[fieldName]), convDst.Field(index))
				if err != nil {
					return
				}
				convDst.Field(index).Set(reflect.ValueOf(fValue))
				dst = convDst
			}
		}
	}
	return
}

func sliceCopy(dst reflect.Value, src interface{}, depth int, current int) (err error) {
	if depth == current {
		return
	}
	var basicSlice int
	switch src.(type) {
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
		var fValue interface{}
		fValue, err = checkTypeAndConv(reflect.ValueOf(src), reflect.ValueOf(convDst))
		if err != nil {
			return
		}
		if dst.Kind() == reflect.Ptr {
			convDst.Elem().Set(reflect.ValueOf(fValue))
		} else {
			convDst.Set(reflect.ValueOf(fValue))
		}
		dst.Set(convDst)
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
