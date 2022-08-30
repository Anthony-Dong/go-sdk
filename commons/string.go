package commons

import (
	"encoding"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/anthony-dong/go-sdk/commons/bufutils"

	"github.com/anthony-dong/go-sdk/commons/internal/prettyjson"
	"github.com/anthony-dong/go-sdk/commons/internal/unsafe"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

var (
	FormatPrettyJson = prettyjson.Format
)

func Slug(str string) string {
	return slug.Make(str)
}

func GenerateUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func UnsafeBytes(data string) []byte {
	return unsafe.UnsafeBytes(data)
}

func UnsafeString(data []byte) string {
	return unsafe.UnsafeString(data)
}

func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case uint8, uint16, uint32, uint64:
		convertUint64 := func(value interface{}) uint64 {
			switch v := value.(type) {
			case uint8:
				return uint64(v)
			case uint16:
				return uint64(v)
			case uint32:
				return uint64(v)
			case uint64:
				return v
			default:
				panic("ToString uint error")
			}
		}
		return strconv.FormatUint(convertUint64(value), 10)
	case int, int8, int16, int32, int64:
		convertInt64 := func(value interface{}) int64 {
			switch v := value.(type) {
			case int8:
				return int64(v)
			case int16:
				return int64(v)
			case int32:
				return int64(v)
			case int64:
				return v
			case int:
				return int64(v)
			default:
				panic("ToString int error")
			}
		}
		return strconv.FormatInt(convertInt64(value), 10)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		if v == nil {
			return ""
		}
		if str, isOk := value.(fmt.Stringer); isOk {
			return str.String()
		}
		if codec, isOk := value.(encoding.TextMarshaler); isOk {
			if text, err := codec.MarshalText(); err == nil {
				return string(text)
			}
		}
		if codec, isOk := value.(json.Marshaler); isOk {
			if text, err := codec.MarshalJSON(); err == nil {
				return string(text)
			}
		}
		if err, isOK := v.(error); isOK {
			return err.Error()
		}
		if result, err := json.Marshal(v); err == nil {
			return string(result)
		}
		return fmt.Sprintf("%v", value)
	}
}

func NewString(elem byte, len int) string {
	if len == 0 {
		return ""
	}
	buffer := bufutils.NewBuffer()
	defer bufutils.ResetBuffer(buffer)
	for x := 0; x < len; x++ {
		buffer.WriteByte(elem)
	}
	return buffer.String()
}

func FormatFloat(i float64, size int) string {
	return strconv.FormatFloat(i, 'f', -1, size)
}

func ContainsString(str []string, elem string) bool {
	for _, v := range str {
		if v == elem {
			return true
		}
	}
	return false
}

func ToJsonString(v interface{}) string {
	if v == nil {
		return ""
	}
	jsonByte, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(jsonByte)
}

func ToPrettyJsonString(v interface{}) string {
	result, _ := json.MarshalIndent(v, "", "  ")
	return string(result)
}

func LinesToString(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	buffer := bufutils.NewBuffer()
	defer bufutils.ResetBuffer(buffer)
	max := len(lines) - 1
	for index, elem := range lines {
		buffer.WriteString(elem)
		if index != max {
			buffer.WriteByte('\n')
		}
	}
	return buffer.String()
}

func SplitSliceString(slice []string, length int) [][]string {
	if len(slice) == 0 {
		return [][]string{}
	}
	if len(slice) <= length {
		return [][]string{slice}
	}
	cut := 0
	if len(slice)%length == 0 {
		cut = len(slice) / length
	} else {
		cut = len(slice)/length + 1
	}
	result := make([][]string, 0, cut)
	for x := 0; x < cut; x++ {
		end := x*length + length
		if end > len(slice) {
			end = len(slice)
		}
		result = append(result, slice[x*length:end])
	}
	return result
}

// SplitString trim space
func SplitString(str string, sep string) []string {
	split := strings.Split(str, sep)
	result := make([]string, 0, len(split))
	for _, elem := range split {
		elem = strings.TrimSpace(elem)
		if elem == "" {
			continue
		}
		result = append(result, elem)
	}
	return result
}
