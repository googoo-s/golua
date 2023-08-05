package vm

const (
	// 可以携带 A、B、C 三个操作数，分别占用 8、9、9 个比特
	// 有 39 条使用 iABC 模式
	IABC = iota
	// 可以携带 A 和 Bx 两个操作数，分别占用 8 和 18 个比特
	// 3 条使用 iABx 指令
	IABx
	// 可以携带 A 和 sBx 两个操作数，分别占用 8 和 18 个比特
	// 4 条使用 iAsBx 模式
	IAsBx
	// 只携带一个操作数，占用全部的 26 个比特
	// 1 条使用 iAx 格式（实际上这条指令并不是真正的指令，
	// 只是用来扩展其他指令操作数的）
	IAx
)

const (
	OP_MOVE = iota
	OP_LOADK
	OP_LOADKX
	OP_LOADBOOL
	OP_LOADNIL
	OP_GETUPVAL
	OP_GETTABUP
	OP_GETTABLE
	OP_SETTABUP
	OP_SETUPVAL
	OP_SETTABLE
	OP_NEWTABLE
	OP_SELF
	OP_ADD
	OP_SUB
	OP_MUL
	OP_MOD
	OP_POW
	OP_DIV
	OP_IDIV
	OP_BAND
	OP_BOR
	OP_BXOR
	OP_SHL
	OP_SHR
	OP_UNM
	OP_BNOT
	OP_NOT
	OP_LEN
	OP_CONCAT
	OP_JMP
	OP_EQ
	OP_LT
	OP_LE
	OP_TEST
	OP_TESTSET
	OP_CALL
	OP_TAILCALL
	OP_RETURN
	OP_FORLOOP
	OP_FORPREP
	OP_TFORCALL
	OP_TFORLOOP
	OP_SETLIST
	OP_CLOSURE
	OP_VARARG
	OP_EXTRAARG
)

const (
	// * 不表示任何信息，也就是说不会被使用。比如 MOVE 指令（iABC 模式）只使用 A 和 B 操作数，
	// 不使用 C 操作数（OpArgN 类型
	OpArgN = iota
	// * 操作数也可能表示布尔值、整数值、upvalue 索引、子函数索引等，这些情况都可以归到 OpArgU 类型里
	OpArgU
	// * 在 iABC 模式下表示寄存器索引，
	// 	* 在 iAsBx 模式下表示跳转偏移
	OpArgR
	// * 表示常量表索引或者寄存器索引
	// * 第一种情况是 LOADK 指令（iABx 模式，用于将常量表中的常量加载到寄存器中），
	// 该指令的 Bx 操作数表示**常量表索引**，如果用 Kst (N)表示常量表访问，则 LOADK 指令可以表示为伪代码 R (A) := Kst (Bx)
	// * 第二种情况是部分 iABC 模式指令，**这些指令的 B 或 C 操作数既可以表示常量表索引也可以表示寄存器索引**，
	// 以加法指令 ADD 为例，如果用 RK (N)表示常量表或者寄存器访问，则该指令可以表示为伪代码 R (A):= RK (B)+RK (C)。
	OpArgK
)

type opcode struct {
	// operator is a test (next instruction must be a jump)
	testFlag byte
	// 是否设置寄存器A
	setAFlag byte
	// 操作数B使用类型
	argBMode byte
	//  操作数C使用类型
	argCMode byte
	// 操作模式
	opMode byte
	name   string
}

var opcodes = []opcode{
	//T A    B       C     mode  	   name
	{0, 1, OpArgR, OpArgN, IABC, "MOVE     "},
	{0, 1, OpArgK, OpArgN, IABx, "LOADK    "},
	{0, 1, OpArgN, OpArgN, IABx, "LOADKX  "},
	{0, 1, OpArgU, OpArgU, IABC, "LOADBOOL"},
	{0, 1, OpArgU, OpArgN, IABC, "LOADNIL "},
	{0, 1, OpArgU, OpArgN, IABC, "GETUPVAL"},
	{0, 1, OpArgU, OpArgK, IABC, "GETTABUP"},
	{0, 1, OpArgR, OpArgK, IABC, "GETTABLE"},
	{0, 0, OpArgK, OpArgK, IABC, "SETTABUP"},
	{0, 0, OpArgU, OpArgN, IABC, "SETUPVAL"},
	{0, 0, OpArgK, OpArgK, IABC, "SETTABLE"},
	{0, 1, OpArgU, OpArgU, IABC, "NEWTABLE"},
	{0, 1, OpArgR, OpArgK, IABC, "SELF     "},
	{0, 1, OpArgK, OpArgK, IABC, "ADD      "},
	{0, 1, OpArgK, OpArgK, IABC, "SUB      "},
	{0, 1, OpArgK, OpArgK, IABC, "MUL      "},
	{0, 1, OpArgK, OpArgK, IABC, "MOD      "},
	{0, 1, OpArgK, OpArgK, IABC, "POW      "},
	{0, 1, OpArgK, OpArgK, IABC, "DIV      "},
	{0, 1, OpArgK, OpArgK, IABC, "IDIV     "},
	{0, 1, OpArgK, OpArgK, IABC, "BAND     "},
	{0, 1, OpArgK, OpArgK, IABC, "BOR      "},
	{0, 1, OpArgK, OpArgK, IABC, "BXOR     "},
	{0, 1, OpArgK, OpArgK, IABC, "SHL      "},
	{0, 1, OpArgK, OpArgK, IABC, "SHR      "},
	{0, 1, OpArgR, OpArgN, IABC, "UNM      "},
	{0, 1, OpArgR, OpArgN, IABC, "BNOT     "},
	{0, 1, OpArgR, OpArgN, IABC, "NOT      "},
	{0, 1, OpArgR, OpArgN, IABC, "LEN      "},
	{0, 1, OpArgR, OpArgR, IABC, "CONCAT  "},
	{0, 0, OpArgR, OpArgN, IAsBx, "JMP      "},
	{1, 0, OpArgK, OpArgK, IABC, "EQ       "},
	{1, 0, OpArgK, OpArgK, IABC, "LT       "},
	{1, 0, OpArgK, OpArgK, IABC, "LE       "},
	{1, 0, OpArgN, OpArgU, IABC, "TEST     "},
	{1, 1, OpArgR, OpArgU, IABC, "TESTSET "},
	{0, 1, OpArgU, OpArgU, IABC, "CALL     "},
	{0, 1, OpArgU, OpArgU, IABC, "TAILCALL"},
	{0, 0, OpArgU, OpArgN, IABC, "RETURN  "},
	{0, 1, OpArgR, OpArgN, IAsBx, "FORLOOP "},
	{0, 1, OpArgR, OpArgN, IAsBx, "FORPREP "},
	{0, 0, OpArgN, OpArgU, IABC, "TFORCALL"},
	{0, 1, OpArgR, OpArgN, IAsBx, "TFORLOOP"},
	{0, 0, OpArgU, OpArgU, IABC, "SETLIST "},
	{0, 1, OpArgU, OpArgN, IABx, "CLOSURE "},
	{0, 1, OpArgU, OpArgN, IABC, "VARARG  "},
	{0, 0, OpArgU, OpArgU, IAx, "EXTRAARG"},
}
