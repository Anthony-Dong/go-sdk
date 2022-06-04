package commons

func Int64Ptr(v int64) *int64 {
	return &v
}

func PtrInt64(p *int64, v ...int64) int64 {
	if p == nil {
		if len(v) > 0 {
			return v[0]
		}
		return 0
	}
	return *p
}

func Int32Ptr(v int32) *int32 {
	return &v
}

func PtrInt32(p *int32, v ...int32) int32 {
	if p == nil {
		if len(v) > 0 {
			return v[0]
		}
		return 0
	}
	return *p
}

func Int16Ptr(v int16) *int16 {
	return &v
}

func PtrInt16(p *int16, v ...int16) int16 {
	if p == nil {
		if len(v) > 0 {
			return v[0]
		}
		return 0
	}
	return *p
}

func Int8Ptr(v int8) *int8 {
	return &v
}

func PtrInt8(p *int8, v ...int8) int8 {
	if p == nil {
		if len(v) > 0 {
			return v[0]
		}
		return 0
	}
	return *p
}

func IntPtr(v int) *int {
	return &v
}

func PtrInt(p *int, v ...int) int {
	if p == nil {
		if len(v) > 0 {
			return v[0]
		}
		return 0
	}
	return *p
}

func Uint64Ptr(v uint64) *uint64 {
	return &v
}

func PtrUint64(p *uint64, v ...uint64) uint64 {
	if p == nil {
		if len(v) > 0 {
			return v[0]
		}
		return 0
	}
	return *p
}

func Uint32Ptr(v uint32) *uint32 {
	return &v
}

func PtrUint32(p *uint32, v ...uint32) uint32 {
	if p == nil {
		if len(v) > 0 {
			return v[0]
		}
		return 0
	}
	return *p
}

func Uint16Ptr(v uint16) *uint16 {
	return &v
}

func PtrUint16(p *uint16, v ...uint16) uint16 {
	if p == nil {
		if len(v) > 0 {
			return v[0]
		}
		return 0
	}
	return *p
}

func Uint8Ptr(v uint8) *uint8 {
	return &v
}

func PtrUint8(p *uint8, v ...uint8) uint8 {
	if p == nil {
		if len(v) > 0 {
			return v[0]
		}
		return 0
	}
	return *p
}

func BytePtr(v byte) *byte {
	return &v
}

func PtrByte(p *byte, v ...byte) byte {
	if p == nil {
		if len(v) > 0 {
			return v[0]
		}
		return 0
	}
	return *p
}

func Float64Ptr(v float64) *float64 {
	return &v
}

func PtrFloat64(p *float64, v ...float64) float64 {
	if p == nil {
		if len(v) > 0 {
			return v[0]
		}
		return 0
	}
	return *p
}

func Float32Ptr(v float32) *float32 {
	return &v
}

func PtrFloat32(p *float32, v ...float32) float32 {
	if p == nil {
		if len(v) > 0 {
			return v[0]
		}
		return 0
	}
	return *p
}

func StringPtr(v string) *string {
	return &v
}

func PtrString(p *string, v ...string) string {
	if p == nil {
		if len(v) > 0 {
			return v[0]
		}
		return ""
	}
	return *p
}

func BoolPtr(v bool) *bool {
	return &v
}

func PtrBool(p *bool, v ...bool) bool {
	if p == nil {
		if len(v) > 0 {
			return v[0]
		}
		return false
	}
	return *p
}
