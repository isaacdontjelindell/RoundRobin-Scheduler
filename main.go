package main

import (
	"fmt"
	"time"
    "os"
    "strconv"
    "strings"
    "io/ioutil"
)

//const CLOCK_SPEED = time.Second / 30
const CLOCK_SPEED= 0
const QUANTUM = 10         // length of quantum before a process will be preempted
const USE_PRIORITY = true  // respect process priority

var quant int = 0          // where we are in the current quantum
var done bool = false      // is the simulation done (are all procs finished)?
var systemTime int = 0     // number of cycles since start of simulation
var idleTime int = 0       // how many idle cycles does the system have?

func main() {
	procList := getInitialProcList()

	runningList := make([]Proc, 0)
	waitingList := make([]Proc, 0)
	readyList := make([]Proc, 0)

	for _, proc := range procList {
        readyList = addToReadyList(readyList, proc)
	}
    fmt.Printf("readyList: %s\n", readyList)

    // prime a proc into the running state
    p := readyList[0] // assumes there's at least one proc - this would be  boring otherwise...
    p.newProcState()
    readyList = readyList[1:]
	runningList = append(runningList, p)

    // start the simulation
	run(readyList, waitingList, runningList)
}

func run(readyList []Proc, waitingList []Proc, runningList []Proc) {
	systemTime++

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
			fmt.Printf("%s enters running state\n", proc.Name)
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
				fmt.Printf("process %s is done\n", runningList[0].Name)
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
            readyList = addToReadyList(readyList, p)
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
        readyList = addToReadyList(readyList, p)
	}

	return waitingList, readyList
}

/* naive priority queuing */
func addToReadyList(readyList []Proc, proc Proc) ([]Proc) {
    if USE_PRIORITY {
        if len(readyList) == 0 { // if no other procs, proc is hightest priority
            readyList = append(readyList, proc)
        } else {
            // find where we should insert this new proc (1 is highest priority)
            insertIndex := -1
            for i, p := range readyList {
                if p.Priority > proc.Priority {
                    insertIndex = i
                    break
                }
            }

            if insertIndex == -1 {
                // if we didn't find an insertion index, it must be lower than everything else
                readyList = append(readyList, proc)
            } else {
                // insert into list at insertIndex
                readyList = append(readyList[:insertIndex],
                    append([]Proc{proc}, readyList[insertIndex:]...)...)
            }
        }

    } else {
        readyList = append(readyList, proc)
    }

    return readyList
}

func printMetrics(doneList []Proc, idleTime int) {
    fmt.Println()
	for _, p := range doneList {
		fmt.Printf("%s waited %d cycles.\n", p.Name, p.WaitTime)
	}
    fmt.Println()
    fmt.Printf("Total system time: %d\n", systemTime)
    fmt.Printf("Idle cycles: %d\n", idleTime)
}

func printReady(readyList []Proc) {
	for _, proc := range readyList {
		fmt.Printf("        %s is waiting\n", proc.Name)
	}
}

func printRunning(runningList []Proc) {
	for _, proc := range runningList {
		fmt.Printf("%s is running\n", proc.Name)
    }
}

func printWaiting(waitingList []Proc) {
	for _, proc := range waitingList {
		fmt.Printf("    %s is doing IO\n", proc.Name)
	}
}

/* proc list can either come from a file with name specified by args[1]
 * or, if a file isn't given, just make a few processes */
func getInitialProcList() []Proc {
    ret := make([]Proc, 0)
    if len(os.Args) > 1 {
        // assume Args[1] is a filename of procs
        filename := os.Args[1]

        content, err := ioutil.ReadFile(filename)
        if err != nil {
            fmt.Printf("Error reading file %s\n", filename)
            os.Exit(1)
        }
        lines := strings.Split(string(content), "\n")

        for i, line := range lines {
            if strings.TrimSpace(line) == "" {  // empty line
                continue
            }
            name := "P" + strconv.Itoa(i+1)
            data := strings.Split(line, " ")
            times := make([]int, 0)
            priority, _ := strconv.Atoi(data[0])
            for _, d := range data[1:] {
                t, _ := strconv.Atoi(d)
                times = append(times, t)
            }

            p := Proc{name, times, 0, false, 0, priority}

            ret = append(ret, p)
        }

    } else {
        p1 := Proc{"P1", []int{7, 2, 9, 6, 10}, 0, false, 0, 0}
        p2 := Proc{"P2", []int{9, 4, 5, 3, 2}, 0, false, 0, 0}
        p3 := Proc{"P3", []int{12, 5, 16, 7, 4}, 0, false, 0, 0}

        ret = append(ret, p1)
        ret = append(ret, p2)
        ret = append(ret, p3)
    }

	return ret
}
