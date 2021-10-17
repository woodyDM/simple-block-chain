package core

import (
	"bytes"
	"errors"
)

type OpCode byte

const (
	//脚本下一个数据入栈
	OpPushData = 0x00
	//复制栈顶元素
	OpDuplicate = 0x01
	//将栈顶元素转化为公钥地址 （先SHA256,再SHA160)
	OpSha160 = 0x02
	//出栈 1，2 比较栈顶1，2元素，如果为false,则报错
	OpEqVerify = 0x03
	//栈 （A，B,C,D,E)
	//元素A是公钥，元素B是交易签名，通过公钥校验交易签名是否正确。如果校验失败，抛出异常；如果校验成功，然后将true放入栈。
	//栈 （true,C,D,E)
	OpCheckSign = 0x04

	VMEnvHash = "VM_TX_HASH"
)

var (
	EmptyStackErr = errors.New("Empty Stack ")
	VmExecErr     = errors.New("Vm Error with invalid stack ")
	opExecMap     map[OpCode]OpExec
	CodeTrue      = []byte{0x01}
	CodeFalse     = []byte{0x00}
	OpPushDataA   = []byte{OpPushData}
	OpDuplicateA  = []byte{OpDuplicate}
	OpSha160A     = []byte{OpSha160}
	OpEqVerifyA   = []byte{OpEqVerify}
	OpCheckSignA  = []byte{OpCheckSign}
)

func init() {
	opExecMap = make(map[OpCode]OpExec)
	opExecMap[OpPushData] = &OpPushDataExec{}
	opExecMap[OpDuplicate] = &OpDuplicateExec{}
	opExecMap[OpSha160] = &OpSha160Exec{}
	opExecMap[OpEqVerify] = &OpEqVerifyExec{}
	opExecMap[OpCheckSign] = &OpCheckSignExec{
		checkFn: Verify,
	}
}

type OpPushDataExec struct{}
type OpDuplicateExec struct{}
type OpSha160Exec struct{}
type OpEqVerifyExec struct{}
type OpCheckSignExec struct {
	checkFn func(msgHash, sign, pubKey []byte) bool
}

func (o *OpPushDataExec) exe(v *Vm) error {
	if v.op+1 >= len(v.script) {
		return ErrWrapf("Invalid op code!Next script element is nil\n")
	}
	v.stack.push(CopyBytes(v.script[v.op+1]))
	v.op += 2
	return nil
}

func (o *OpDuplicateExec) exe(v *Vm) error {
	p, err := v.stack.peek()
	if err != nil {
		return ErrWrap("Invalid op code!", err)
	}
	v.stack.push(CopyBytes(p))
	v.op += 1
	return nil
}

func (o *OpSha160Exec) exe(v *Vm) error {
	bs, err := v.stack.pop()
	if err != nil {
		return ErrWrap("Invalid opSha160", err)
	}
	bs = Sha256(bs)
	bs = Sha160(bs)
	v.stack.push(bs)
	v.op += 1
	return nil
}

//OpEqVerify Handler
func (o *OpEqVerifyExec) exe(v *Vm) error {
	if v.stack.size() < 2 {
		return ErrWrapf("Invalid opEqVerify \n")
	}
	b1, _ := v.stack.pop()
	b2, _ := v.stack.pop()
	if !bytes.Equal(b1, b2) {
		return ErrWrapf("Bytes not eq\n")
	}
	v.op += 1
	return nil
}

//OpCheckSignExec Handler
func (o *OpCheckSignExec) exe(v *Vm) error {
	if v.stack.size() < 2 {
		return ErrWrapf("Invalid opCheckSign \n")
	}
	publicKey, _ := v.stack.pop()
	sign, _ := v.stack.pop()
	hash, ok := v.GetEnv(VMEnvHash)
	if !ok {
		return ErrWrapf("No hash found!\n")
	}
	h := hash.([]byte)
	if !o.checkFn(h, sign, publicKey) {
		return ErrWrapf("Sign check failed %v %v\n", publicKey, sign)
	}
	v.op += 1
	v.stack.push(CodeTrue)
	return nil
}

type OpExec interface {
	exe(v *Vm) error
}

type Vm struct {
	stack   *Stack
	op      int
	script  Script
	execMap map[OpCode]OpExec
	env     map[string]interface{}
}

func NewVm(script Script) *Vm {
	return &Vm{
		stack:   newStack(),
		script:  script,
		execMap: opExecMap,
	}
}

func (v *Vm) SetEnv(key string, value interface{}) {
	if v.env == nil {
		v.env = make(map[string]interface{})
	}
	v.env[key] = value
}

func (v *Vm) GetEnv(key string) (interface{}, bool) {
	if v.env == nil {
		return nil, false
	}
	value, ok := v.env[key]
	return value, ok
}

//for test
func (v *Vm) CustomExec(op OpCode, exe OpExec) {
	dst := make(map[OpCode]OpExec)
	for k, v := range v.execMap {
		if op == k {
			dst[op] = exe
		} else {
			dst[k] = v
		}
	}
	v.execMap = dst
}

func (v *Vm) Exec() error {
	for v.op < len(v.script) {
		opSlice := v.script[v.op]
		if len(opSlice) != 1 {
			return ErrWrapf("Invalid opCode size:%d  op is :%v ", len(opSlice), opSlice)
		}
		code := OpCode(opSlice[0])
		exec, exist := v.execMap[code]
		if !exist {
			return ErrWrapf("Invalid opCode :%v ", opSlice[0])
		}
		Log.Debug("Run op code: ", code, " at pos [", v.op, "]")
		err := exec.exe(v)
		if err != nil {
			err := ErrWrap("Vm Error", err)
			Log.Debug(err)
			return err
		}
	}
	stackSize := v.stack.size()
	if stackSize == 0 {
		return nil
	}
	topEle, _ := v.stack.peek()
	if stackSize == 1 && bytes.Equal(CodeTrue, topEle) {
		return nil
	}
	return VmExecErr
}

type Stack struct {
	data [][]byte
}

func newStack() *Stack {
	return &Stack{data: make([][]byte, 0)}
}

func (s *Stack) push(d []byte) {
	s.data = append(s.data, d)
}

func (s *Stack) size() int {
	return len(s.data)
}

func (s *Stack) peek() ([]byte, error) {
	if s.size() == 0 {
		return nil, EmptyStackErr
	}
	return s.data[s.size()-1], nil
}

func (s *Stack) pop() ([]byte, error) {
	if s.size() == 0 {
		return nil, EmptyStackErr
	}
	d := s.data[s.size()-1]
	s.data = s.data[:s.size()-1]
	return d, nil
}
