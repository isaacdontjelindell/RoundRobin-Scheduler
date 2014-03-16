package main

import "fmt"


func main() {
    //testProcess()

    procList := getInitialProcList()

    fmt.Println(procList)
}


func getInitialProcList() []Process {
    p1 := Process{"P1", READY, []int{7, 2, 9, 6, 10}, 7}
    p2 := Process{"P2", READY, []int{9, 4, 5, 3, 2}, 9}
    p3 := Process{"P3", READY, []int{12, 5, 16, 7, 4}, 12}

    ret := []Process{p1, p2, p3}

    return ret
}
