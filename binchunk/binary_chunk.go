package binchunk


const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)


// 常量类型
const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

// binaryChunk 总体结构
type binaryChunk struct {
	// 头部
	header
	// 主函数upvalues 和闭包有关
	sizeUpvalues byte
	// 主函数原型
	mainFunc *Prototype
}

// 头部
type header struct {
	// 签名，相对于Java的模式,主要起快速识别文件格式的作用
	signature [4]byte
	// 版本号，记录二进制chunk文件所对应的Lua版本号
	version byte
	// 格式号，如果和虚拟机本身的格式号不匹配，就拒绝加载该文件
	format byte
	// LUAC_DATA 格式号之后的6个字节在Lua官方实现里叫作LUAC_DATA。
	// 其中前两个字节是0x1993，这是Lua 1.0发布的年份；
	// 后四个字节依次是回车符（0x0D）、换行符（0x0A）、替换符（0x1A）和另一个换行符
	// 6个字节主要起进一步校验的作用。如果Lua虚拟机在加载二进制chunk时发现这6个字节和预期的不一样，就会认为文件已经损坏，拒绝加载
	luacData [6]byte
	// 接下来的5个字节分别记录cint、size_t、Lua虚拟机指令、Lua整数和Lua浮点数这5种数据类型在二进制chunk里占用的字节数
	cintSize        byte
	sizetSize       byte
	instructionSize byte
	luaIntegerSize  byte
	luaNumberSize   byte
	// LUAC_INT， n个字节存放Lua整数值0x5678，存储这个Lua整数的目的是为了检测二进制chunk的大小端方式
	luacInt int64
	// LUAC_NUM，头部的最后n个字节存放Lua浮点数370.5，头部的最后n个字节存放Lua浮点数370.5
	luacNum float64
}


type Prototype struct {
	// 源文件名，只有在主函数原型里，该字段才真正有值，在其他嵌套的函数原型里，该字段存放空字符串
	Source string
	// Prototype,用于记录原型对应的函数在源文件中的起止行号。如果是普通的函数，起止行号都应该大于0；
	//  如果是主函数，则起止行号都是0
	LineDefined     uint32
	LastLineDefined uint32
	// 函数固定参数个数。这里的固定参数，是相对于变长参数（Vararg）而言的
	NumParams byte
	// 用来记录函数是否为Vararg函数，即是否有变长参数
	IsVararg byte
	// 记录的是寄存器数量。Lua编译器会在编译函数时将这个数量计算好，并以字节类型保存在函数原型
	MaxStackSize byte
	// 函数基本信息之后是指令表
	Code []uint32
	// 指令表之后是常量表。常量表用于存放Lua代码里出现的字面量，包括nil、布尔值、整数、浮点数和字符串五种。
	// 每个常量都以1字节tag开头，用来标识后续存储的是哪种类型的常量值。常量tag值,tag类型为TAG_NIL，TAG_BOOLEAN，，
	Constants []interface{}
	// Upvalues 占有两个字节
	Upvalues []Upvalue
	// 子函数原型表
	Protos []*Prototype
	// 行号表 子函数原型表之后是行号表，其中行号按cint类型存储。行号表中的行号和指令表中的指令一一对应
	LineInfo []uint32
	// 号表之后是局部变量表，用于记录局部变量名，
	// 表中每个元素都包含变量名（按字符串类型存储）和起止指令索引（按cint类型存储
	LocVars []LocVar
	// 数原型的最后一部分内容是Upvalue名列表。该列表中的元素（按字符串类型存储）
	// 和前面Upvalue表中的元素一一对应，分别记录每个Upvalue在源代码中的名字
	UpvalueNames []string
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}


// Upvalue 和闭包相关
type Upvalue struct {
	Instack byte
	Idx     byte
}



// 解析二进制chunk
func Undump(data []byte) * Prototype {
	reader := &reader{data}
	reader.checkHeader()
	reader.readByte()
	return reader.readProto("")
}