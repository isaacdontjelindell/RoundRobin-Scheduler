package main

import (
    "fmt"
    "time"
)

const CLOCK_SPEED = time.Second/20
const QUANTUM = 10

func main() {

    procList := getInitialProcList()
    fmt.Println(procList)

    runningList := make([]Proc, 0)
    waitingList := make([]Proc, 0)
    readyList := make([]Proc, 0)

    procList[0].newProcState() // yikes - unchecked operation
    runningList = append(runningList, procList[0])
    for _, proc := range procList[1:] {
        readyList = append(readyList, proc)
    }

    fmt.Println(readyList)

    run(readyList, waitingList, runningList)
}

func run(readyList []Proc, waitingList []Proc, runningList []Proc) {
    systemTime := 1
    quant := 0

    for {
        fmt.Printf("%d:\n", systemTime )
        printRunning(runningList)
        printReady(readyList)
        printWaiting(waitingList)

        systemTime++

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
                waitingList = append(waitingList[:i], waitingList[i+1:]...)
                //waitingList = waitingList[1:]
                readyList = append(readyList, p)
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
                p.Preempted = true
                runningList = runningList[1:]
                readyList = append(readyList, p)
                quant = 0
            }
        }

        if len(runningList) == 0 {
            // nothing running - try to get something from the ready queue
            if !(len(readyList) == 0) {
                proc := readyList[0]
                readyList = readyList[1:]

                proc.newProcState()
                runningList = append(runningList, proc)
            }
        }

        time.Sleep(CLOCK_SPEED) // tick tock tick tock
    }
}

func printReady(readyList []Proc) {
    for _, proc := range readyList {
        fmt.Printf("            %s is waiting\n", proc)
    }
}

func printRunning(runningList []Proc) {
    for _, proc := range runningList {
        fmt.Printf("    %s is running\n", proc)
    }
}

func printWaiting(waitingList []Proc) {
    for _, proc := range waitingList {
        fmt.Printf("        %s is doing IO\n", proc)
    }
}

func getInitialProcList() []Proc {
    p1 := Proc{"P1", []int{7, 2, 9, 6, 10}, 0, false}
    p2 := Proc{"P2", []int{9, 4, 5, 3, 2}, 0, false}
    p3 := Proc{"P3", []int{12, 5, 16, 7, 4}, 0, false}

    ret := []Proc{p1, p2, p3}

    return ret
}
