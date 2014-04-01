package main

import (
	"fmt"
	"time"
)

const CLOCK_SPEED = time.Second / 30
const QUANTUM = 10 // length of quantum before a process will be preempted
var quant int = 0 // where we are in the current quantum
var done bool = false // is the simulation done (are all procs finished?)
var idleTime int = 0 // how many idle cycles does the system have?

func main() {

	procList := getInitialProcList()
	fmt.Println(procList)

	runningList := make([]Proc, 0)
	waitingList := make([]Proc, 0)
	readyList := make([]Proc, 0)

	// prime a proc into the running state
	procList[0].newProcState() // yikes - unchecked operation. What if no procs??
	runningList = append(runningList, procList[0])
	for _, proc := range procList[1:] {
		readyList = append(readyList, proc)
	}

	run(readyList, waitingList, runningList)
}

func run(readyList []Proc, waitingList []Proc, runningList []Proc) {
    systemTime := 1

	doneList := make([]Proc, 0)

	for !done {
		//fmt.Printf("%d:\n", systemTime ) // for debugging
		printRunning(runningList)
		printReady(readyList)
		printWaiting(waitingList)

		systemTime++

		/* keep track of non-I/O wait time for each proc */
		for i, _ := range readyList {
			readyList[i].WaitTime++
		}

		/* simulate an IO clock cycle - anything that's in an IO wait is one cycle closer to running */
		waitingList, readyList =
            ioTick(waitingList, readyList)

        /* If anything is running, make it tick. Also preempt any processes that have
         * overrun their quantum */
		runningList, doneList, waitingList, readyList =
            runningTick(runningList, doneList, waitingList, readyList)

		/* Check if the system is done */
		if len(readyList) == 0 && len(waitingList) == 0 && len(runningList) == 0 {
			done = true
		}

		/* Nothing running - try to get something from the ready queue */
        runningList, readyList =
            getReadyProc(runningList, readyList)


		time.Sleep(CLOCK_SPEED) // tick tock tick tock

	}

	printMetrics(doneList, idleTime)
}


/* Nothing running - try to get something from the ready queue */
func getReadyProc(runningList []Proc, readyList []Proc) ([]Proc, []Proc) {
    if len(runningList) == 0 {
        if !(len(readyList) == 0) {
            proc := readyList[0]
            readyList = readyList[1:]

            proc.newProcState()
            fmt.Printf("%s enters running state\n", proc)
            runningList = append(runningList, proc)
        } else if !done {
            idleTime++ // idle cycle
        }
    }

    return runningList, readyList
}


/* If anything is running, make it tick. Also preempt any processes that have
 * overrun their quantum */
func runningTick(runningList []Proc,
                    doneList []Proc,
                    waitingList []Proc,
                    readyList []Proc) ([]Proc, []Proc, []Proc, []Proc) {
	if !(len(runningList) == 0) {
		runningList[0].RemainingStateTime-- // clock tick
		quant++

		if runningList[0].RemainingStateTime < 1 { // done running
			quant = 0
			v := runningList[0].newProcState()
			if v == 1 { // if process doesn't need any more time
				fmt.Printf("process %s is done\n", runningList[0])
				doneList = append(doneList, runningList[0])
				runningList = runningList[1:] // drop it.
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

    return runningList, doneList, waitingList, readyList
}


/* simulate an IO clock cycle - anything that's in an IO wait is one cycle closer to running */
func ioTick(waitingList []Proc, readyList []Proc) ([]Proc, []Proc) {
	removeInd := make([]int, 0)
	for i, _ := range waitingList {
		waitingList[i].RemainingStateTime--
		if waitingList[i].RemainingStateTime < 1 {
			removeInd = append(removeInd, i) // save index of proc that needs IO
		}
	}
	// There's got to be a better way of doing this
	// These are procs that are done w/ IO, ready to run
	for i := 0; i < len(removeInd); i++ {
		p := waitingList[i]
		// cut out the i-th element from waitingList
		waitingList = append(waitingList[:i], waitingList[i+1:]...)
		readyList = append(readyList, p)
	}

	return waitingList, readyList
}


func printMetrics(doneList []Proc, idleTime int) {
	for _, p := range doneList {
		fmt.Printf("%s waited %d cycles.\n", p, p.WaitTime)
	}
	fmt.Printf("System had %d idle cycles.\n", idleTime)
}


func printReady(readyList []Proc) {
	for _, proc := range readyList {
		fmt.Printf("        %s is waiting\n", proc)
	}
}


func printRunning(runningList []Proc) {
	for _, proc := range runningList {
		fmt.Printf("%s is running\n", proc)
	}
}


func printWaiting(waitingList []Proc) {
	for _, proc := range waitingList {
		fmt.Printf("    %s is doing IO\n", proc)
	}
}


func getInitialProcList() []Proc {
	p1 := Proc{"P1", []int{7, 2, 9, 6, 10}, 0, false, 0}
	p2 := Proc{"P2", []int{9, 4, 5, 3, 2}, 0, false, 0}
	p3 := Proc{"P3", []int{12, 5, 16, 7, 4}, 0, false, 0}

	ret := []Proc{p1, p2, p3}

	return ret
}
