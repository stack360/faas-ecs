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
        taskDefinitionName = descServicesResult.services[0].TaskDefinition
        taskDefinitionInput := &ecs.DescribeTaskDefinitionInput {
            TaskDefinition: aws.String(*taskDefinitionName),
        }
        taskDefinitionResult, taskDefinitionErr := ecsClient.DescribeTaskDefinition(taskDefinitionInput)
        if taskDefinitionErr != nil {
            return nil, taskDefinitionErr
        }
        image := taskDefinitionResult.TaskDefinition.ContainerDefinitions[0].Image
        serviceName := descServicesResult.services[0].ServiceName
        function := requests.Function{
            Name:            *serviceName,
            Replicas:        nil,
            Image:           *image,
            InvocationCount: 0,
        }
        functions = append(functions, function)
    }
    return functions, nil
}
