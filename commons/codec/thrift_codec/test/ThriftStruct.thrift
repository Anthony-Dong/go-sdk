namespace go thriftstruct

struct MapData {
    1: optional map<i64,string> M_i64,
    2: optional map<i32,string> M_i32,
    3: optional map<i16,string> M_i16,
    4: optional map<byte,string> M_byte,
    5: optional map<double,string> M_double,
    6: optional map<string,string> M_string,
    7: optional map<bool,string> M_bool,
    9: optional map<Numberz,string> M_enum,
}

struct ListData {
    1: optional list<i64> L_i64,
    2: optional list<i32> L_i32,
    3: optional list<i16> L_i16,
    4: optional list<byte> L_byte,
    5: optional list<double> L_double,
    6: optional list<string> L_string,
    7: optional list<bool> L_bool,
    8: optional list<NormalStruct> L_struct,
    9: optional list<Numberz> L_enum,
    10: optional list<TypedefStruct> L_Ref,
}

struct SetData {
    1: optional set<i64> S_i64,
    2: optional set<i32> S_i32,
    3: optional set<i16> S_i16,
    4: optional set<byte> S_byte,
    5: optional set<double> S_double,
    6: optional set<string> S_string,
    7: optional set<bool> S_bool,
    8: optional set<NormalStruct> S_struct,
    9: optional set<Numberz> S_enum,
    10: optional set<TypedefStruct> S_Ref,
}


struct NormalData {
    1: optional i64 F_i64,
    2: optional i32 F_i32,
    3: optional i16 F_i16,
    4: optional byte F_byte,
    5: optional double F_double,
    6: optional string F_string,
    7: optional bool F_bool,
    8: optional NormalStruct F_struct,
    9: optional Numberz F_enum,
    10: optional binary F_binary,
    11: optional list<string> F_list_string,
    12: optional set<string> F_set_string,
    13: optional map<string,i64> F_map_string,
    14: optional MapData F_MapData,
    15: optional ListData F_ListData,
    16: optional SetData F_SetData,
    17: optional TypedefData F_TypedefData,
}

struct TypedefData {
//    1: optional TypedefStruct F_TypedefStruct,
    2: optional TypedefString F_TypedefString,
    3: optional TypedefEnum F_TypedefEnum,
    4: optional TypedefMap F_TypedefMap,
}

typedef NormalStruct TypedefStruct
typedef string TypedefString
typedef map<string,i64> TypedefMap
typedef Numberz TypedefEnum

struct NormalStruct {
    1: optional string F_1
}

enum Numberz
{
  ONE = 1,
  TWO,
  THREE,
  FIVE = 5,
  SIX,
  EIGHT = 8
}

enum Numberz1{
  PostiveOne=-1,
  Zero = 0,
  ONE = 1,
  TWO =2,
}

struct TestRequest {
    1: optional NormalData Data
}

struct TestResponse {
    1: optional NormalData Data
}

service TestService {
    TestResponse Test (1: TestRequest req),
    oneway void TestOneway (1: TestRequest req),
}