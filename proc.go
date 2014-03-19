package main

import "fmt"

const (
	READY = 0
	RUN   = 1
	WAIT  = 2
    DONE  = 3
)

type Proc struct {
	Name               string
	Times              []int
    RemainingStateTime int
}

func (p Proc) String() string {
    s := fmt.Sprintf("%s", p.Name)
    return s
}


func (p *Proc) newProcState() int {
    if len(p.Times) == 0 {
        return 1
    }
    newStateTime := p.Times[0]
    p.Times = p.Times[1:]
    p.RemainingStateTime = newStateTime
    return 0
}

