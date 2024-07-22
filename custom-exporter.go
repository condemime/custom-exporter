package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "flag"
    "os"
    "os/signal"
    "syscall"
)

func main() {

        SetupCloseHandler()

        var port string
        flag.StringVar(&port, "port", "40080", "HTTP Port")
        var url string
        flag.StringVar(&url, "url", "/metrics", "URL")
        var file string
        flag.StringVar(&file, "file", "custom.prom", "Prom File")
        flag.Parse()
        var hostname string
        hostname, err := os.Hostname()
        if err != nil { panic(err) }
        var listen string
	hostname = "localhost"
        listen = hostname + ":" + port
        var output string
        output = "Running " + listen + url + " (" + file + ")"

        h1 := func(w http.ResponseWriter, _ *http.Request) {
                data, err := ioutil.ReadFile(file)
                if err != nil {
                        fmt.Println("File reading error", err)
                        return
                }
                fmt.Fprint(w, string(data))
        }

        h2 := func(w http.ResponseWriter, _ *http.Request) {
                //io.WriteString(w, os.Getpid())
                fmt.Fprint(w, os.Getpid())
        }

        fmt.Println("", output)

        http.HandleFunc(url, h1)
        http.HandleFunc("/control", h2)
        log.Fatal(http.ListenAndServe(listen, nil))
}

func SetupCloseHandler() {
    c := make(chan os.Signal, 2)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        fmt.Println("\r STOP")
        os.Exit(0)
    }()
}
