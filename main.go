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
        switch {
            case proc.Status == READY:
                readyQueue <- proc
            case proc.Status == RUN:
                runQueue <- proc
            case proc.Status == WAIT:
                waitQueue <- proc
        }
    }

    for {
        time.Sleep(time.Second)
        tick(readyQueue, runQueue, waitQueue)
    }
}

func tick(readyQueue chan Proc, runQueue chan Proc, waitQueue chan Proc) {
    select {
    case p := <-readyQueue:
        fmt.Println(p)
    default:
        fmt.Println("[tick] nothing in ready queue")
    }

    select {
    case p := <-waitQueue:
        fmt.Println(p)
    default:
        fmt.Println("[tick] nothing in wait queue")
    }

    select {
    case p := <-runQueue:
        fmt.Println(p)
    default:
        fmt.Println("[tick] nothing in run queue")
    }
}


func getInitialProcList() []Proc {
    p1 := Proc{"P1", READY, []int{7, 2, 9, 6, 10}, 7}
    p2 := Proc{"P2", READY, []int{9, 4, 5, 3, 2}, 9}
    p3 := Proc{"P3", READY, []int{12, 5, 16, 7, 4}, 12}

    ret := []Proc{p1, p2, p3}

    return ret
}
