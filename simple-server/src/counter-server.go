package main

import (
    "fmt"
    "net/http"
    "io"
    "sync"
    "time"
    "log"
    "encoding/json"
    "github.com/gorilla/websocket"
    "github.com/gorilla/mux"
)

// Number of calls to local server.
var counter int

// Mutex to protect the numberCalls variable
var lock sync.Mutex

// Websocket
var upgrader = websocket.Upgrader{}

// JSON counter
type Counter struct {
    Counter    int `json:"Counter"`
}

func increaseRequestCounter(){
    // Protect the variable from multi-thread modification with a Mutex
    lock.Lock()
    //Update the simple counter
    counter++
    // Unlock
    lock.Unlock()
}

/**
* -- handlerGeneralWebsite --
* Display a simple text page as an HTTP response
* Increase the counter for each request
*/

func handlerGeneralWebsite(w http.ResponseWriter, r *http.Request) {

    increaseRequestCounter()

    msg := fmt.Sprintf("Hello TheThingsNetwork team! \n\n HTTP Counter: localhost:8080/counter\n JSON Api Counter: localhost:8080/json-api\n Web-Socket Counter: localhost:8080/ws")
    io.WriteString(w, msg)
}

/**
* -- handleHttpCounter --
* Display a simple text page as an HTTP response
*/

func handleHttpCounter(w http.ResponseWriter, req *http.Request) {

    // Protect the variable from multi-thread modification with a Mutex
    lock.Lock()
    msg := fmt.Sprintf("%v", counter)
    // Unlock
    lock.Unlock()

    io.WriteString(w, msg)

}

/**
* -- handleJSONApiCounter --
* Send a JSON object "Counter" as a HTTP response, contening the int counter
*/

func handleJSONApiCounter(w http.ResponseWriter, req *http.Request){

  c := Counter{counter}
  json.NewEncoder(w).Encode(c)

}

/**
* -- handleWebSocketCounter --
* Create a webscoket and broadcast counter value
* Value is update it every 5s
*/
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
/**
* -- main --
* Routes creation and serve
*/
func main() {

    // Create a router
    router := mux.NewRouter()

    //General route, accept paramaters
    router.HandleFunc("/", handlerGeneralWebsite).Methods("GET")

    // Configure HTTP route
    router.HandleFunc("/counter", handleHttpCounter).Methods("GET")

    // Configure a JSON route
    router.HandleFunc("/counter/json-api", handleJSONApiCounter).Methods("GET")

    // Configure websocket route
    router.HandleFunc("/counter/ws", handleWebSocketCounter).Methods("GET")

    //Start to serve (localhost:8080)
    log.Println("http server started on :8000")
    log.Fatal(http.ListenAndServe(":8000", router))
}
