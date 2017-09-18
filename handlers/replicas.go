package handlers

import (
        "encoding/json"
        "io/ioutil"
        "log"
        "net/http"

        "github.com/alexellis/faas/gateway/requests"
        "github.com/stack360/faas-ecs/types"
        "github.com/gorilla/mux"
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/service/ecs"
)


func MakeReplicaUpdater(ecsClient *ecs.ECS) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
        log.Println("Update replicas")

                vars := mux.Vars(r)
                functionName := vars["name"]

                req := types.ScaleServiceRequest{}
                if r.Body != nil {
                        defer r.Body.Close()
                        bytesIn, _ := ioutil.ReadAll(r.Body)
                        marshalErr := json.Unmarshal(bytesIn, &req)
                        if marshalErr != nil {
                                w.WriteHeader(http.StatusBadRequest)
                                msg := "Cannot parse request. Please pass valid JSON"
                                w.Write([]byte(msg))
                                log.Println(msg, marshalErr)
                                return
                        }
                }

                updateServiceInput := &ecs.UpdateServiceInput {
                        Cluster:        aws.String("faas"),
                        DesiredCount:   aws.Int64(int64(req.Replicas)),
                        Service:        aws.String(functionName),
                }
                updateServiceResult, updateServiceErr := ecsClient.UpdateService(updateServiceInput)
                log.Println("update service result is %s", updateServiceResult)

                if updateServiceErr != nil {
                        w.WriteHeader(500)
                        w.Write([]byte("Unable to update function deployment " + functionName))
                        log.Println(updateServiceErr)
                        return
                }
        }

}

func MakeReplicaReader(ecsClient *ecs.ECS) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
        log.Println("Update replicas")

                vars := mux.Vars(r)
                functionName := vars["name"]

                functions, err := getServiceList(ecsClient)
                if err != nil {
                        w.WriteHeader(500)
                        return
                }

                var found *requests.Function
                for _, function:= range functions {
                        if function.Name == functionName {
                              found = &function
                              break
                        }
                }

                if found == nil {
                        w.WriteHeader(404)
                        return
                }

                functionBytes, _ := json.Marshal(found)
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(200)
                w.Write(functionBytes)
        }

}
