package handlers

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/alexellis/faas/gateway/requests"
    "github.com/aws/aws-sdk-go/service/ecs"
)

func getServiceList(ecsClient *ecs.EcsClient) ([]requests.Function, error) {
    var functions []requests.functions

    return functions, nil    
}
