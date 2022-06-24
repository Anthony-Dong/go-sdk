# PB Codec

## Feature

1. 可以实现 PB2/PB3 报文解析( without idl), 是一个debug工具, 方便定位和解决线上问题, 目前测试兼容情况还是非常OK!
- 例如用户抓取了数据包，但是想知道内容，此时在线上，内部抓包工具无法解析(解析失败)，此时就需要这种工具

2. 关于PB的编码可以看我的文章: [PB介绍](https://anthony-dong.github.io/2022/01/16/cc45d69abc6417303d451f43acf099d9/)
   
3. 具体类型支持情况
- string: 只支持文本编码, 不等价于兼容utf8编码
- int: 全部转换为无符号类型
- bool: 全部转换为int类型, 1表示true, 0表示false
- enum: 全部转换为int类型
- double: 支持逻辑比较trick (因为double在编码的时候会转换成uint64, 所以当数据大于4e18的时候自动转换为double类型), [参考: IEEE 标准 754 浮点数](https://en.wikipedia.org/wiki/IEEE_754)
- list(repeated): 支持 packed & unpacked 编码
  - unpacked: 当元素个数只有一个, 无法区分是否为list, 所以大于1个的时候会自动merge
  - packed: 支持
- map: 支持
  - 特殊case，用户定义了`repeated message{key:1, value: 2}`这种也会被误认为是map
- bytes: 不支持
  - 原因: 因为bytes、message、string等类型的wire type 都是 bytes, 识别逻辑属于auto, 解析优先级为message、string、packed、bytes
- group: 支持, 其他问题同list类型问题

4. 遗留问题: 是否需要支持packed编码(因为PB3默认)? 还是全部将bytes解析为 packed(这种根据业务其实很少用bytes类型)?

## Example

1. define protobuf 2
```protobuf
syntax = "proto2";
import "google/protobuf/any.proto";

message TestPB2Data {
  enum EnumType {
    UnknownType = 0; // 必须以0开始！
    Test1Type = 1;
    Test2Type = 2;
  }
  message TestPB2Obj {
    optional int64 t_int64 = 1;
  }
  optional string t_string = 1;
  optional int64 t_int64 = 2;
  required bool t_bool = 3;
  required fixed64 t_fix64 = 4;
  repeated int64 t_list_i64 = 5[packed = true];
  map<int64, string> t_map = 6;
  optional EnumType t_enum = 7;
  optional TestPB2Obj t_obj = 8 ;
  repeated TestPB2Obj t_list_obj = 9 ;
  map<string, TestPB2Data> t_map_obj = 10;
  repeated string  t_list_string = 11;
  optional google.protobuf.Any any = 12;
  repeated group Result = 13 {
    required string url = 2;
    optional string title = 3;
    repeated string snippets = 4;
  }
  optional double t_double = 14;
  optional bytes t_bytes=15;
}
```
2. new message and encode pb

```go
func NewPB2Data(t *testing.T) []byte {
	any, err := anypb.New(&test.TestPB2Data_TestPB2Obj{TInt64: commons.Int64Ptr(1)})
	if err != nil {
		t.Fatal(err)
	}
	result := test.TestPB2Data{
		TString:  commons.StringPtr("1234"),
		TInt64:   commons.Int64Ptr(1),
		TBool:    commons.BoolPtr(true),
		TFix64:   commons.Uint64Ptr(1),
		TListI64: []int64{1, 2, 3},
		TMap:     map[int64]string{1: "1", 2: "2"},
		TEnum:    test.TestPB2Data_Test1Type.Enum(),
		TObj: &test.TestPB2Data_TestPB2Obj{
			TInt64: commons.Int64Ptr(1),
		},
		TListObj: []*test.TestPB2Data_TestPB2Obj{
			{
				TInt64: commons.Int64Ptr(1),
			},
			{
				TInt64: commons.Int64Ptr(2),
			},
		},
		TMapObj: map[string]*test.TestPB2Data{
			"k1": &test.TestPB2Data{
				TString: commons.StringPtr("1234"),
				TInt64:  commons.Int64Ptr(1),
				TBool:   commons.BoolPtr(false),
				TFix64:  commons.Uint64Ptr(1),
			},
			"k2": &test.TestPB2Data{
				TString: commons.StringPtr("k2"),
				TInt64:  commons.Int64Ptr(2),
				TBool:   commons.BoolPtr(true),
				TFix64:  commons.Uint64Ptr(1),
			},
		},
		TListString: []string{"1", "2"},
		Any:         any,
		Result: []*test.TestPB2Data_Result{
			{
				Url:      commons.StringPtr("1"),
				Title:    commons.StringPtr("Title"),
				Snippets: []string{"1", "2"},
			},
			{
				Url:      commons.StringPtr("2"),
				Title:    commons.StringPtr("Title2"),
				Snippets: []string{"12", "22"},
			},
		},
		TDouble: commons.Float64Ptr(1.0111),
		TBytes:  []byte{0x91, 0x02, 0x03, 0x04},
	}
	marshal, err := proto.Marshal(&result)
	if err != nil {
		t.Fatal(err)
	}
	return marshal
}


func TestDecodeMessage(t *testing.T) {
    message, err := DecodeMessage(context.Background(), codec.NewBuffer(NewPB2Data(t)))
    if err != nil {
        t.Fatal(err)
    }
    t.Log(commons.ToPrettyJsonString(message))
}
```


3. decode pb binary payload to json (without idl)
```json
{
  "1": "1234",
  "2": 1,
  "3": 1,
  "4": 1,
  "5": [
    1,
    2,
    3
  ],
  "6": {
    "1": "1",
    "2": "2"
  },
  "7": 1,
  "8": {
    "1": 1
  },
  "9": [
    {
      "1": 1
    },
    {
      "1": 2
    }
  ],
  "10": {
    "k1": {
      "1": "1234",
      "2": 1,
      "3": 0,
      "4": 1
    },
    "k2": {
      "1": "k2",
      "2": 2,
      "3": 1,
      "4": 1
    }
  },
  "11": [
    "1",
    "2"
  ],
  "12": {
    "1": "type.googleapis.com/TestPB2Data.TestPB2Obj",
    "2": {
      "1": 1
    }
  },
  "13": [
    {
      "2": "1",
      "3": "Title",
      "4": [
        "1",
        "2"
      ]
    },
    {
      "2": "2",
      "3": "Title2",
      "4": [
        "12",
        "22"
      ]
    }
  ],
  "14": 1.0111,
  "15": [
    273,
    3,
    4
  ]
}
```

## BUG

待补充 ....
