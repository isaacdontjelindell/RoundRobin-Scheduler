package main

import (
    "fmt"
    "time"
)


func main() {
    readyQueue := make(chan Proc, 20)
    runQueue := make(chan Proc, 20)
    waitQueue := make(chan Proc, 20)

    procList := getInitialProcList()
    fmt.Println(procList)

    for _, proc := range(procList) {
        dispatch(proc, readyQueue, runQueue, waitQueue)
    }

    tick(readyQueue, runQueue, waitQueue)
}

func dispatch(p Proc, readyQueue chan Proc, runQueue chan Proc, waitQueue chan Proc) {
    switch {
        case p.Status == READY:
            readyQueue <- p
        case p.Status == RUN:
            runQueue <- p
        case p.Status == WAIT:
            waitQueue <- p
    }
}


func tick(readyQueue chan Proc, runQueue chan Proc, waitQueue chan Proc) {
    fmt.Printf("Welcome to tick\n")

    system_time := 0

    for {
        select {
        case p := <-runQueue:
            if p.StateTimeRemaining > 0 {
                fmt.Printf("Process %s is running\n", p.Name)
                p.StateTimeRemaining--
                runQueue <- p
            } else {
                v := changeProcState(&p)
                if v > 0 {
                    fmt.Printf("Process %s is done\n", p.Name)
                } else {
                    dispatch(p, readyQueue, runQueue, waitQueue)
                }
            }
        default:
           fmt.Println("nothing in the runQueue")
           readyProc := <-readyQueue
           runQueue <- readyProc
        }

        select {
        case p := <-waitQueue:
            if p.StateTimeRemaining > 0 {
                fmt.Printf("Process %s is waiting\n", p.Name)
                p.StateTimeRemaining--
                waitQueue <- p
            } else {
                changeProcState(&p) // don't care about return (wait -> done not possible)
                fmt.Printf("Process %s is done waiting\n", p.Name)
                dispatch(p, readyQueue, runQueue, waitQueue)
            }
        default:
            fmt.Println("nothing in the waitQueue")
        }

        system_time++
        time.Sleep(time.Second)
    }

}


func getInitialProcList() []Proc {
    p1 := Proc{"P1", READY, []int{7, 2, 9, 6, 10}, 7}
    p2 := Proc{"P2", READY, []int{9, 4, 5, 3, 2}, 9}
    p3 := Proc{"P3", READY, []int{12, 5, 16, 7, 4}, 12}

    ret := []Proc{p1, p2, p3}

    return ret
}
