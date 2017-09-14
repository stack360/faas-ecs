package handlers

import (
    "strings"
    "encoding/json"
    "log"
    "net/http"

    "github.com/alexellis/faas/gateway/requests"
    "github.com/aws/aws-sdk-go/service/ecs"
    "github.com/aws/aws-sdk-go/aws"
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

func getServiceList(ecsClient *ecs.ECS) ([]requests.Function, error) {
    var functions []requests.Function

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
        name := strings.Split(*item, "/")[1]
        descServicesInput := &ecs.DescribeServicesInput {
            Services: []*string{
                  aws.String(*item),
            },
        }
        descServicesResult, descServicesErr := ecsClient.DescribeServices(descServicesInput)
        if descServicesErr != nil {
            return nil, descServicesErr
        }
        taskDefinitionName := descServicesResult.Services[0].TaskDefinition
        serviceName := descServicesResult.Services[0].ServiceName
        desiredCount := descServicesResult.Services[0].DesiredCount
        taskDefinitionInput := &ecs.DescribeTaskDefinitionInput {
            TaskDefinition: aws.String(*taskDefinitionName),
        }
        taskDefinitionResult, taskDefinitionErr := ecsClient.DescribeTaskDefinition(taskDefinitionInput)
        if taskDefinitionErr != nil {
            return nil, taskDefinitionErr
        }
        image := taskDefinitionResult.TaskDefinition.ContainerDefinitions[0].Image
        function := requests.Function{
            Name:            *serviceName,
            Replicas:        uint64(*desiredCount),
            Image:           *image,
            InvocationCount: 0,
        }
        functions = append(functions, function)
    }
    return functions, nil
}

func MakeFunctionReader(ecsClient *ecs.ECS) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        functions, err := getServiceList(ecsClient)
        if err != nil {
            log.Println(err)
            w.WriteHeader(500)
            w.Write([]byte(err.Error()))
            return
        }

        functionBytes, _ := json.Marshal(functions)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(200)
        w.Write(functionBytes)
    }
}
