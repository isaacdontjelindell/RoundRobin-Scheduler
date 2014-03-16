package main

import "fmt"

const (
    READY = 0
    RUNNING = 1
    WAITING = 2
)

type Process struct {
    Name               string
    Status             int
    Times              []int
    StateTimeRemaining int
}

func changeProcState(proc *Process) int {
    if len(proc.Times) == 0 {
        return 1
    }

    newStateTime := proc.Times[0]
    proc.Times = proc.Times[1:] // remove from head
    proc.StateTimeRemaining = newStateTime
    return 1
}

func testProcess() {
    var p Process
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
