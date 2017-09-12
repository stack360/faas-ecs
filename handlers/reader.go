package handlers

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/alexellis/faas/gateway/requests"
    "github.com/aws/aws-sdk-go/service/ecs"
)

type stringInDynArray struct {
    value string
}

func (s *stringInDynArray) SetValue(value string) {
    s.value = value
}
func (s stringInDynArray) Value() string {
    return s.value
}

func getServiceList(ecsClient *ecs.EcsClient) ([]requests.Function, error) {
    var functions []requests.functions

    // first list all raw services from faas cluster
    res, err := ecsClient.ListServices(&ecs.ListServicesInput{
        Cluster: aws.String("faas"),
    })

    if err != nil {
        return nil, err
    }
    // Declare a new array to store service names
    s := []stringInDynArray{}
    for _, item := range res.ServiceArns {
        s = append(s, stringInDynArray{strings.Split(*item, "/")[1]})
    }
    // s now is a list containing all services from faas cluster
    return functions, nil
}
