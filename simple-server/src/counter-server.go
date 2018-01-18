package main

import (
    "fmt"
    "net/http"
    "io"
)

// Number of calls to local server.
var numberCalls int

func HelloServer(w http.ResponseWriter, req *http.Request) {

    //Update the simple counter
    //Note: no handelling of potential simultaneous calls
    numberCalls = numberCalls + 1

    //Display a simple text & counter
    msg := fmt.Sprintf("Hello TheThingsNetwork team! \n\nNumber of calls to this server: %v", numberCalls)
    io.WriteString(w, msg)
}

func main() {
    // Serve on localhost:8080

    http.HandleFunc("/", HelloServer)
    http.ListenAndServe(":8080", nil)
}
