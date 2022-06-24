package pb_codec

import (
	"bytes"
	"encoding/json"

	"github.com/anthony-dong/go-sdk/commons"
)

type FieldOrderMap struct {
	list []Field
	data map[Field]interface{}
}

func (f *FieldOrderMap) GetFieldId(v int32) (interface{}, bool) {
	for key, elem := range f.data {
		if key.FieldId == v {
			return elem, true
		}
	}
	return nil, false
}

func (f *FieldOrderMap) ContainsField(v int32) (isExist bool) {
	_, isExist = f.GetFieldId(v)
	return
}

func (f *FieldOrderMap) Size() int {
	return len(f.list)
}

func NewFieldOrderMap(size int) *FieldOrderMap {
	return &FieldOrderMap{
		list: make([]Field, 0, size),
		data: make(map[Field]interface{}, size),
	}
}

func (t FieldOrderMap) MarshalJSON() ([]byte, error) {
	result := bytes.Buffer{}
	result.WriteString("{")
	for index, v := range t.list {
		result.WriteByte('"')
		result.WriteString(v.String())
		result.WriteByte('"')
		result.WriteByte(':')
		marshal, err := json.Marshal(t.data[v])
		if err != nil {
			return nil, err
		}
		result.Write(marshal)
		if index == len(t.list)-1 {
			continue
		}
		result.WriteByte(',')
	}
	result.WriteByte('}')
	return result.Bytes(), nil
}

func (f *FieldOrderMap) Set(field Field, v interface{}) {
	if _, isExist := f.data[field]; isExist {
		// not handler order!! swap...
		f.data[field] = v
		return
	}
	f.list = append(f.list, field)
	f.data[field] = v
}

type Field struct {
	FieldId   int32
	FieldType int8
}

func NewField(fieldId int32, fieldType int8) Field {
	return Field{FieldId: fieldId, FieldType: fieldType}
}

func (t Field) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

func (f Field) MarshalText() (text []byte, err error) {
	return []byte(f.String()), nil
}

func (t Field) String() string {
	return commons.ToString(t.FieldId)
}
