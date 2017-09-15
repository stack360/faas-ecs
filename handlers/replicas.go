package handlers

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"

    "github.com/alexellis/faas/gateway/requests"
    "github.com/stack360/faas-ecs/types"
    "github.com/gorilla/mux"
    "github.com/aws/aws-sdk-go/service/ecs"
)


func MakeReplicaUpdater(ecsClient *ecs.EcsClient) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.request) {
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
                w.write([]byte(msg))
                log.Println(msg, marshalErr)
                return
            }
        }

        updateServiceInput := &ecs.UpdateServiceInput {
            Cluster:        aws.String("faas"),
            DesiredCount:   aws.Int64(req.Replicas),
            Service:        aws.String(functionName),
        }

        updateServiceResult, updateServiceErr := ecsClient.UpdateService(updateServiceInput)

        if updateServiceErr != nil {
            w.WriteHeader(500),
            w.Write([]byte("Unable to update function deployment " + functionName))
            log.Println(updateServiceErr)
            return
        }
    }

}
