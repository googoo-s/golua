package vm

const (
	MAXARG_Bx  = 1<<18 - 1      // 262143
	MAXARG_sBx = MAXARG_Bx >> 1 // 131071
)

// 指令
//
//	31       22       13       5    0
//	 +-------+^------+-^-----+-^-----
//	 |b=9bits |c=9bits |a=8bits|op=6|
//	 +-------+^------+-^-----+-^-----
//	 |    bx=18bits    |a=8bits|op=6|
//	 +-------+^------+-^-----+-^-----
//	 |   sbx=18bits    |a=8bits|op=6|
//	 +-------+^------+-^-----+-^-----
//	 |    ax=26bits            |op=6|
//	 +-------+^------+-^-----+-^-----
//	31      23      15       7      0
type Instruction uint32

// Opcode()方法从指令中提取操作F
func (self Instruction) Opcode() int {
	return int(self & 0x3F)
}

// ABC()方法从iABC模式指令中提取参数
func (self Instruction) ABC() (a, b, c int) {
	//  8、9、9
	a = int(self >> 6 & 0xFF)
	c = int(self >> 14 & 0x1FF)
	b = int(self >> 23 & 0x1FF)
	return

}

// ABx()方法从iABx模式指令中提取参数，
func (self Instruction) ABx() (a, bx int) {
	// 8 和 18
	a = int(self >> 6 & 0xFF)
	bx = int(self >> 14)
	return
}

// AsBx()方法从iAsBx模式指令中提取参数
func (self Instruction) AsBx() (a, sbx int) {

	// SBx 操作数（共 18 个比特）表示的是**有符号整数**。有很多种方式可以把有符号整数编码成比特序列，
	// 比如 **2 的补码**（Two's Complement）等。Lua 虚拟机这里采用了一种叫作
	// **偏移二进制码（Offset Binary，也叫作 Excess-K）的编码模式。**
	// > 如果把 sBx 解释成无符号整数时它的值是 x，那么解释成有符号整数时它的值就是 x-K。那么 K 是什么呢？
	// K 取 sBx 所能表示的最大无符号整数值的一半。也就是上面代码中的 MAXARG_sBx
	// 8 和 18
	a, bx := self.ABx()
	return a, bx - MAXARG_sBx
}

// Ax()方法从iAx模式指令中提取参数

func (self Instruction) Ax() int {
	//  26 个比特
	return int(self >> 6)
}

// OpName 返回操作码名称
func (self Instruction) OpName() string {
	return opcodes[self.Opcode()].name
}

// OpMode 编码模式
func (self Instruction) OpMode() byte {
	return opcodes[self.Opcode()].opMode
}

// BMode 操作数B的使用模式
func (self Instruction) BMode() byte {
	return opcodes[self.Opcode()].argBMode
}

// CMode 操作数C的使用模式
func (self Instruction) CMode() byte {
	return opcodes[self.Opcode()].argCMode
}
