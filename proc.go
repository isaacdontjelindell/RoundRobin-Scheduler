package main

import "fmt"

type Proc struct {
	Name               string
	Times              []int
    RemainingStateTime int
    Preempted          bool
}

func (p Proc) String() string {
    s := fmt.Sprintf("%s", p.Name)
    return s
}


func (p *Proc) newProcState() int {
    if p.Preempted {
        p.Preempted = false
        return 0
    }
    if len(p.Times) == 0 {
        return 1
    }
    newStateTime := p.Times[0]
    p.Times = p.Times[1:]
    p.RemainingStateTime = newStateTime
    return 0
}

