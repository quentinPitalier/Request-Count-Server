package main

import (
    "fmt"
    "net/http"
    "io"
    "sync"
    "time"
    "github.com/gorilla/websocket"
    "log"
)

// Number of calls to local server.
var counter int
// Mutex to protect the numberCalls variable
var lock sync.Mutex

// Websocket
var upgrader = websocket.Upgrader{}
// JSON counter
type Counter struct {
        Counter    int `json:"counter"`
}

func importantFunction(name string) {
    lock.Lock()
    defer lock.Unlock()
    fmt.Println(name)
    time.Sleep(1 * time.Second)
}

func handlerGeneralWebsite(w http.ResponseWriter, r *http.Request) {

    // Protect the variable from multi-thread modification with a Mutex
    lock.Lock()
    //Update the simple counter
    counter++
    // Unlock
    lock.Unlock()

    msg := fmt.Sprintf("Hello TheThingsNetwork team! \n\n HTTP Counter: localhost:8080/counter\n JSON Api Counter: localhost:8080/json-api\n Web-Socket Counter: localhost:8080/ws")
    io.WriteString(w, msg)
}

func handleHttpCounter(w http.ResponseWriter, req *http.Request) {

    //Display a simple counter on the page

    // Protect the variable from multi-thread modification with a Mutex
    lock.Lock()
    msg := fmt.Sprintf("%v", counter)
    // Unlock
    lock.Unlock()

    io.WriteString(w, msg)

}

func handleJSONApiCounter(w http.ResponseWriter, req *http.Request){

  //Will handle JSON request

}

func handleWebSocketCounter(w http.ResponseWriter, r *http.Request) {

        // Upgrade initial GET request to a websocket
        ws, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
                log.Fatal(err)
        }

        for {
                // Write counter and send it as JSON object
                err := ws.WriteJSON(counter)
                if err != nil {
                        log.Printf("error: %v", err)
                        break
                }
                // Sleep for 5s before update the counter in the websocket
                time.Sleep(5 * time.Second)
        }

        // Close WebSocket
        defer ws.Close()
}

func main() {

    //General route, accept paramaters
    http.HandleFunc("/", handlerGeneralWebsite)

    // Configure HTTP route
    http.HandleFunc("/counter", handleHttpCounter)

    // Configure a JSON route
    http.HandleFunc("/json-api", handleJSONApiCounter)

    // Configure websocket route
    http.HandleFunc("/ws", handleWebSocketCounter)

    //Start to serve (localhost:8080)
    log.Println("http server started on :8080")
        err := http.ListenAndServe(":8080", nil)
        if err != nil {
                log.Fatal("ListenAndServe: ", err)
        }
}
