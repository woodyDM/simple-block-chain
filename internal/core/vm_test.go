package core

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

const OpNext = 0x05

var (
	OpNextA = []byte{OpNext}
)

func init() {
	opExecMap[OpNext] = new(OpNextExec)
}

type OpNextExec struct{}

func (o *OpNextExec) exe(v *Vm) error {
	v.op++
	return nil
}

func TestStack(t *testing.T) {
	s := newStack()
	if s.size() != 0 {
		t.Fatalf("size not 0")
	}
	s.push([]byte{1, 2, 3})
	s.push([]byte{2, 3, 4})

	if s.size() != 2 {
		t.Fatalf("size not 2")
	}
	pop, err := s.pop()
	if err != nil {
		t.Fatalf("err not nil")
	}
	bytes.Equal(pop, []byte{2, 3, 4})
	pop2, err2 := s.pop()
	if err2 != nil {
		t.Fatalf("err not nil")
	}
	bytes.Equal(pop2, []byte{1, 2, 3})

	_, err = s.pop()
	if err != EmptyStackErr {
		t.Fatalf("should be empty err")
	}
}

func TestVmExec_OK(t *testing.T) {
	s := Script{{5}, {5}, {5}}
	vm := NewVm(s)
	err := vm.Exec()
	if err != nil {
		t.Fatal("err should be nil")
	}
}

func TestVmExec_OpSizeInvalid(t *testing.T) {
	s := Script{{5}, {5, 2}, {5}}
	vm := NewVm(s)
	err := vm.Exec()
	if err == nil {
		t.Fatal("should err")
	}
	if !strings.Contains(err.Error(), "Invalid opCode size") {
		t.Fatal("error not ok")
	}
	if vm.op != 1 {
		t.Fatal("vm op pos invalid")
	}
}

func TestVmExec_OpInvalid(t *testing.T) {
	s := Script{{5}, {5}, {52}}
	vm := NewVm(s)
	err := vm.Exec()
	if err == nil {
		t.Fatal("should err")
	}
	if !strings.Contains(err.Error(), "Invalid opCode :52") {
		t.Fatal("error not ok")
	}
	if vm.op != 2 {
		t.Fatal("vm op pos invalid")
	}
}

func TestVmExec_OpPushData(t *testing.T) {
	s := Script{OpPushDataA, {2, 3, 4}}
	vm := NewVm(s)
	err := vm.Exec()
	if err != VmExecErr {
		t.Fatal("should err not nil")
	}
	p, _ := vm.stack.peek()
	if !bytes.Equal([]byte{2, 3, 4}, p) {
		t.Fatal("vm stack invalid")
	}
}

func TestVmExec_OpPushData_Error(t *testing.T) {
	s := Script{OpNextA, OpNextA, OpPushDataA}
	vm := NewVm(s)
	err := vm.Exec()
	if err == VmExecErr {
		t.Fatal("should   nil")
	}
	if !strings.Contains(err.Error(), "Next script element is nil") {
		t.Fatal("error msg not ok")
	}
}

func TestVmExec_OpDuplicate(t *testing.T) {
	s := Script{{OpPushData}, {5, 0, 2}, {OpDuplicate}}
	vm := NewVm(s)
	err := vm.Exec()
	if err != VmExecErr {
		t.Fatal("should err nil")
	}
	p, _ := vm.stack.pop()
	if !bytes.Equal([]byte{5, 0, 2}, p) {
		t.Fatal("vm stack data invalid")
	}
	p, _ = vm.stack.pop()
	if !bytes.Equal([]byte{5, 0, 2}, p) {
		t.Fatal("vm stack data invalid")
	}
	if vm.stack.size() != 0 {
		t.Fatal("vm stack size invalid")
	}
}

func TestVmExec_OpSha160(t *testing.T) {
	s := Script{{OpPushData}, {5, 0, 2}, {OpSha160}}
	vm := NewVm(s)
	err := vm.Exec()
	if err != VmExecErr {
		t.Fatal("should err nil", err)
	}
	p, _ := vm.stack.pop()
	hexS := hex.EncodeToString(p)
	if hexS != "d72c354f2dc38f12a84917349c9f6492f0db3d91" {
		t.Fatal("vm stack data invalid")
	}
}

func TestVmExec_OpEqVerify(t *testing.T) {
	s := Script{{OpPushData}, {5, 0, 2}, {OpPushData}, {5, 0, 2}, {OpEqVerify}}
	vm := NewVm(s)
	err := vm.Exec()
	if err != nil {
		t.Fatal("should err be nil")
	}
	if vm.stack.size() != 0 {
		t.Fatal("should size be 0")
	}
}

func TestVmExec_OpEqVerify_error(t *testing.T) {
	s := Script{{OpPushData}, {5, 0, 2}, {OpPushData}, {5, 0, 1}, {OpEqVerify}}
	vm := NewVm(s)
	err := vm.Exec()
	if vm.stack.size() != 0 {
		t.Fatal("should size be 0")
	}
	if err == nil {
		t.Fatal("err should not be nil")
	}
	if !strings.Contains(err.Error(), "Bytes not eq") {
		t.Fatal("err mismatch")
	}
}

func TestVmExec_OpEqVerify_error_length(t *testing.T) {
	s := Script{{OpPushData}, {5, 0, 2}, {OpEqVerify}}
	vm := NewVm(s)
	err := vm.Exec()
	if err == nil {
		t.Fatal("err should not be nil")
	}
	if !strings.Contains(err.Error(), "Invalid opEqVerify") {
		t.Fatal("err mismatch")
	}
}

func TestVmExec_OpDuplicateError(t *testing.T) {
	s := Script{{OpDuplicate}}
	vm := NewVm(s)
	err := vm.Exec()
	if err == nil {
		t.Fatal("  err should not nil")
	}
	if !strings.Contains(err.Error(), "Empty Stack") {
		t.Fatal("err mismatch")
	}
}

func TestVmExec_OpSign(t *testing.T) {
	s := Script{{OpPushData}, {5, 0, 2}, {OpPushData}, {5, 0, 2}, {OpCheckSign}}
	vm := NewVm(s)
	vm.SetEnv(VMEnvHash, []byte("any"))
	vm.CustomExec(OpCheckSign, &OpCheckSignExec{checkFn: func(i []byte, i2 []byte, i3 []byte) bool {
		return true
	}})
	st := opExecMap
	fmt.Sprintln(st)
	err := vm.Exec()
	if err != nil {
		t.Fatal("err should be nil")
	}
	trueCode, _ := vm.stack.pop()
	if !bytes.Equal(trueCode, CodeTrue) {
		t.Fatal("not true")
	}

}
