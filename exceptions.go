package main

import (
	"errors"
	"fmt"
)

var ECallFromUser = errors.New("Environmental call from user mode")
var ECallFromSupervisor = errors.New("Environmental call from supervisor mode")
var ECallFromReserved = errors.New("Can't make a call from mode")
var ECallFromMachine = errors.New("Environmental call from machine mode")

func IllegalInst(inst uint32) {
	panic(fmt.Errorf("Illegal instruction: %#x.", inst))
}
