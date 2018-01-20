package main

import (
    "fmt"
    "net/http"
    "io"
    "sync"
    "time"
)

// Number of calls to local server.
var counter int
// Mutex to protect the numberCalls variable
var lock sync.Mutex


func importantFunction(name string) {
    lock.Lock()
    defer lock.Unlock()
    fmt.Println(name)
    time.Sleep(1 * time.Second)
}

func HelloServer(w http.ResponseWriter, req *http.Request) {

    // Protect the variable from multi-thread modification with a Mutex
    lock.Lock()
    //Update the simple counter
    counter++
    // Unlock
    lock.Unlock()

    //Display a simple text & counter
    msg := fmt.Sprintf("Hello TheThingsNetwork team! \n\nNumber of calls to this server: %v", counter)
    io.WriteString(w, msg)

}

func main() {
    // Serve on localhost:8080

    http.HandleFunc("/", HelloServer)
    http.ListenAndServe(":8080", nil)
}
