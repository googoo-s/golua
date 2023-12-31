package binchunk

import (
	"encoding/binary"
	"math"
)

type reader struct {
	data []byte
}

// readByte 从字节流里读取一个字节
func (self *reader) readByte() byte {
	b := self.data[0]
	self.data = self.data[1:]
	return b
}

// readUint32 使用小端方式从字节流里读取一个cint存储类型（占4个字节，映射为Go语言uint32类型）的整数
func (self *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(self.data)
	self.data = self.data[4:]
	return i
}

// readUint64 使用小端方式从字节流里读取一个cint存储类型（占4个字节，映射为Go语言uint32类型）的整数
func (self *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(self.data)
	self.data = self.data[8:]
	return i
}

// readLuaInteger 从字节流里读取一个Lua整数（占8个字节，映射为Go语言int64类型
func (self *reader) readLuaInteger() int64 {
	return int64(self.readUint64())
}

// readLuaNumber 从字节流里读取一个Lua浮点数（占8个字节，映射为Go语言float64类型）
func (self *reader) readLuaNumber() float64 {
	return math.Float64frombits(self.readUint64())
}

// readString()方法从字节流里读取字符串（映射为Go语言string类型）
func (self *reader) readString() string {
	size := uint(self.readByte())

	if size == 0 {
		return ""
	}

	if size == 0xff {
		size = uint(self.readUint64())
	}
	bytes := self.readBytes(size - 1)
	return string(bytes)
}

// readBytes()方法从字节流里读取n个字节
func (self *reader) readBytes(n uint) []byte {
	bytes := self.data[:n]
	self.data = self.data[n:]
	return bytes
}

// checkHeader()方法从字节流里读取并检查二进制chunk头部的各个字段，
// 如果发现某个字段和期望不符，则调用panic函数终止加载
func (self *reader) checkHeader() {
	if string(self.readBytes(4)) != LUA_SIGNATURE {
		panic(" not a precompiled chunk")
	} else if self.readByte() != LUAC_VERSION {
		panic("version mismatch")
	} else if self.readByte() != LUAC_FORMAT {
		panic("format mismatch")
	} else if string(self.readBytes(6)) != LUAC_DATA {
		panic("corrupted!")
	} else if self.readByte() != CINT_SIZE {
		panic("int size mismatch")
	} else if self.readByte() != CSIZET_SIZE {
		panic("size_t size mismatch")
	} else if self.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch")
	} else if self.readByte() != LUA_INTEGER_SIZE {
		panic("lua_Integer size mismatch")
	} else if self.readByte() != LUA_NUMBER_SIZE {
		panic("Lua_Number size mismatch")
	} else if self.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch")
	} else if self.readLuaNumber() != LUAC_NUM {
		panic("float format nismatch")
	}
}

// readProto 读取原型
func (self *reader) readProto(parentSource string) *Prototype {
	source := self.readString()
	if source == "" {
		source = parentSource
	}

	return &Prototype{
		Source:          source,
		LineDefined:     self.readUint32(),
		LastLineDefined: self.readUint32(),
		NumParams:       self.readByte(),
		IsVararg:        self.readByte(),
		MaxStackSize:    self.readByte(),
		Code:            self.readCode(),
		Constants:       self.readConstants(),
		Upvalues:        self.readUpvalues(),
		Protos:          self.readProtos(source),
		LineInfo:        self.readLineInfo(),
		LocVars:         self.readLocVars(),
		UpvalueNames:    self.readUpvalueName(),
	}
}

// readCode()方法从字节流里读取指令表
func (self *reader) readCode() []uint32 {
	code := make([]uint32, self.readUint32())
	for i := range code {
		code[i] = self.readUint32()
	}
	return code
}

// readConstants 从字节流中读取常量表
func (self *reader) readConstants() []interface{} {
	constants := make([]interface{}, self.readUint32())
	for i := range constants {
		constants[i] = self.readConstant()
	}
	return constants
}

// readConstant()方法从字节流里读取一个常量
func (self *reader) readConstant() interface{} {
	switch self.readByte() {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return self.readByte() != 0
	case TAG_INTEGER:
		return self.readLuaInteger()
	case TAG_NUMBER:
		return self.readLuaNumber()
	case TAG_SHORT_STR, TAG_LONG_STR:
		return self.readString()
	default:
		panic("corrupted!")
	}
}

// readConstant()方法从字节流里读取一个常量
func (self *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, self.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: self.readByte(),
			Idx:     self.readByte(),
		}
	}
	return upvalues
}

// readProtos()方法从字节流里读取子函数原型表
func (self *reader) readProtos(parentSource string) []*Prototype {
	protos := make([]*Prototype, self.readUint32())
	for i := range protos {
		protos[i] = self.readProto(parentSource)
	}
	return protos
}

// readLineInfo 从字节流中读取行号
func (self *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, self.readUint32())
	for i := range lineInfo {
		lineInfo[i] = self.readUint32()
	}
	return lineInfo
}

// readLocVars()方法从字节流里读取局部变量表
func (self *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, self.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: self.readString(),
			StartPC: self.readUint32(),
			EndPC:   self.readUint32(),
		}
	}
	return locVars
}

// readUpvalueNames()方法从字节流里读取Upvalue名列表
func (self *reader) readUpvalueName() []string {
	upvalueNames := make([]string, self.readUint32())
	for i := range upvalueNames {
		upvalueNames[i] = self.readString()
	}
	return upvalueNames
}
