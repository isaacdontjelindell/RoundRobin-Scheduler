package main

import "fmt"

const (
	READY = 0
	RUN   = 1
	WAIT  = 2
)

type Proc struct {
	Name               string
	Status             int
	Times              []int
	StateTimeRemaining int
}

func changeProcState(proc *Proc) int {
	if len(proc.Times) == 0 {
		return 1
	}

	newStateTime := proc.Times[0]
	proc.Times = proc.Times[1:] // remove from head
	proc.StateTimeRemaining = newStateTime
	return 1
}

func testProc() {
	var p Proc
	p.Name = "testp"
	p.Status = READY
	p.Times = []int{1, 2, 3, 4}
	p.StateTimeRemaining = 1

	fmt.Println(p.Name)
	fmt.Println(p.Status)

	changeProcState(&p)
	fmt.Println(p.Times)
	fmt.Println(p.StateTimeRemaining)

	changeProcState(&p)
	fmt.Println(p.Times)
	fmt.Println(p.StateTimeRemaining)

	changeProcState(&p)
	fmt.Println(p.Times)
	fmt.Println(p.StateTimeRemaining)

	changeProcState(&p)
	fmt.Println(p.Times)
	fmt.Println(p.StateTimeRemaining)

	v := changeProcState(&p)
	if v > 0 {
		fmt.Println("ERROR: no more state times")
	}
	fmt.Println(p.Times)
	fmt.Println(p.StateTimeRemaining)
}
