package main

import (
        "fmt"
//        "reflect"
//        "strings"
        "github.com/aws/aws-sdk-go/service/ecs"
        "github.com/aws/aws-sdk-go/aws"
//        "github.com/aws/aws-sdk-go/aws/awserr"
        "github.com/aws/aws-sdk-go/aws/credentials"
//        "github.com/aws/aws-sdk-go/aws/requests"
        "github.com/aws/aws-sdk-go/aws/session"
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

func main() {
    sess, err := session.NewSession(&aws.Config{
        Region:      aws.String("us-west-1"),
        Credentials: credentials.NewSharedCredentials("", ""),
    })

    ecsClient := ecs.New(sess)
    //  first list services
    res, err := ecsClient.ListServices(&ecs.ListServicesInput{
        Cluster: aws.String("default"),
    })
    //  Declare a new list to store service names
    s := []stringInDynArray{}
    for _, item := range res.ServiceArns {
        //s = append(s, stringInDynArray{strings.Split(*item, "/")[1]})
        s = append(s, stringInDynArray{*item})
    }
    awsStringSlice := make([]*string, 0)
    for _, item := range s {
        awsStringSlice = append(awsStringSlice, aws.String(item.Value()))
    }
    fmt.Printf("aws string slice is %s", awsStringSlice)
    //  s now is a list containing all services from faas cluster
    fmt.Printf("list services result is %s, error is %s", res.ServiceArns, err)

    // describe services
    descServicesInput := &ecs.DescribeServicesInput {
          Services: awsStringSlice,
    }

    descServicesResult, descServicesErr := ecsClient.DescribeServices(descServicesInput)
    fmt.Printf("Desribe services result: %s, error: %s", descServicesResult, descServicesErr)
    desiredCount := descServicesResult.Services[0].DesiredCount
    taskDefinitionName := descServicesResult.Services[1].TaskDefinition
    serviceName := descServicesResult.Services[0].ServiceName
    // get task definition
    taskDefinitionInput := &ecs.DescribeTaskDefinitionInput {
        TaskDefinition: aws.String(*taskDefinitionName),
    }
    result, errrr := ecsClient.DescribeTaskDefinition(taskDefinitionInput)
    image := result.TaskDefinition.ContainerDefinitions[0].Image
    fmt.Printf("\nDescribe Task Definition result: %s, error: %s, image: %s, serviceName: %s, \n Desired count : %s", result, errrr, *image, *serviceName, uint64(*desiredCount))
}
