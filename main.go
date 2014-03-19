package main

import (
    "fmt"
    "time"
)

const CLOCK_SPEED = time.Second
const QUANTUM = 10

func main() {

    procList := getInitialProcList()
    fmt.Println(procList)

    runningList := make([]Proc, 0)
    waitingList := make([]Proc, 0)
    readyList := make([]Proc, 0)

    for _, proc := range(procList) {
        readyList = append(readyList, proc)
    }

    fmt.Println(readyList)

    run(readyList, waitingList, runningList)
}

func run(readyList []Proc, waitingList []Proc, runningList []Proc) {
    //for i := 0; i<40; i++ {
    quant := 0
    for {
        if len(runningList) == 0 {
            // nothing running - try to get something from the ready queue
            if !(len(readyList) == 0) {
                proc := readyList[0]
                readyList = readyList[1:]

                proc.newProcState()
                runningList = append(runningList, proc)
            }
        }

        // if anything is running
        if !(len(runningList) == 0) {
            runningList[0].RemainingStateTime-- // clock tick
            quant++

            if runningList[0].RemainingStateTime < 1 {  // done running
                quant = 0
                v := runningList[0].newProcState()
                if v == 1 { // if process doesn't need any more time
                    fmt.Printf("process %s is done\n", runningList[0])
                    runningList = runningList[1:] // drop it. TODO save finished procs
                } else { // needs IO
                    p := runningList[0]
                    runningList = runningList[1:]
                    waitingList = append(waitingList, p)
                }
            } else if quant == QUANTUM { // quantum is up, preempt the running proc
                p := runningList[0]
                runningList = runningList[1:]
                readyList = append(readyList, p)
                quant = 0
            }
        }

        removeInd := make([]int, 0)
        if !(len(waitingList) == 0) {
            for i, _ := range waitingList {
                waitingList[i].RemainingStateTime--
                if waitingList[i].RemainingStateTime < 1 {
                    removeInd = append(removeInd, i) // save index of proc that needs IO
                    //waitingList = waitingList[1:]
                    //readyList = append(readyList, waitingList[i])
                }
            }
            // There's got to be a better way of doing this
            // These are procs that are done w/ IO, ready to run
            for i:=0; i < len(removeInd); i++ {
                p := waitingList[i]
                waitingList = waitingList[1:]
                readyList = append(readyList, p)
            }
        }

        printRunning(runningList)
        printWaiting(waitingList)
        printReady(readyList)
        time.Sleep(CLOCK_SPEED) // tick tock tick tock

    }
}

func printReady(readyList []Proc) {
    for _, proc := range readyList {
        fmt.Printf("        %s is waiting (remaining: %d)\n", proc, proc.RemainingStateTime)
    }
}

func printRunning(runningList []Proc) {
    for _, proc := range runningList {
        fmt.Printf("%s is running (remaining: %d)\n", proc, proc.RemainingStateTime)
    }
}

func printWaiting(waitingList []Proc) {
    for _, proc := range waitingList {
        fmt.Printf("    %s is doing IO (remaining: %d)\n", proc, proc.RemainingStateTime)
    }
}

func getInitialProcList() []Proc {
    p1 := Proc{"P1", []int{7, 2, 9, 6, 10}, 0}
    p2 := Proc{"P2", []int{9, 4, 5, 3, 2}, 0}
    p3 := Proc{"P3", []int{12, 5, 16, 7, 4}, 0}

    ret := []Proc{p1, p2, p3}

    return ret
}
