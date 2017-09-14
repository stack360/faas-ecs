package main

import (
        "fmt"
        "log"
        "net/http"
        "time"

        "github.com/stack360/faas-ecs/handlers"
        "github.com/gorilla/mux"

        "github.com/aws/aws-sdk-go/service/ecs"
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/session"
)

func main() {
    sess, err := session.NewSession(&aws.Config{
        Region:      aws.String("us-west-1"),
        Credentials: credentials.NewSharedCredentials("", ""),
    })

    r := mux.NewRouter()
    ecsClient := ecs.New(sess)

    r.HandleFunc("/system/functions", handlers.MakeFunctionReader(ecsClient)).Methods("GET")
    r.HandleFunc("/system/functions", handlers.MakeDeployHandler(ecsClient)).Methods("POST")
    r.HandleFunc("/system/functions", handlers.MakeDeleteHandler(ecsClient)).Methods("DELETE")

    r.HandleFunc("/system/function/{name:[-a-zA-Z_0-9]+}", handlers.MakeReplicaReader(ecsClient)).Methods("GET")
    R.HandleFunc("/system/scale-function/{name:[-a-zA-Z_0-9]+}", handlers.MakeReplicaUpdater(ecsClient)).Methods("POST")

    functionProxy := handlers.MakeProxy()
    r.HandleFunc("/function/{name:[-a-zA-Z_0-9]+}", functionProxy)
    r.HandleFunc("/function/{name:[-a-zA-Z_0-9]+}/", functionProxy)

    readTimeout := 8 * time.Second
    writeTimeout := 8 * time.Second
    tcpPort := 8082

    s := &http.Server{
        Addr:           fmt.Sprintf(":%d", tcpPort),
        ReadTimeout:    readTimeout,
        WriteTimeout:   writeTimeout,
        MaxHeaderBytes: http.DefaultMaxHeaderBytes,
        Handler:        r,
    }

    log.Fatal(s.ListenandServe())
}
