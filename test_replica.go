package main

import (
        "fmt"
        "reflect"
//        "strings"
        "github.com/aws/aws-sdk-go/service/ecs"
        "github.com/aws/aws-sdk-go/aws"
//        "github.com/aws/aws-sdk-go/aws/awserr"
        "github.com/aws/aws-sdk-go/aws/credentials"
//        "github.com/aws/aws-sdk-go/aws/requests"
        "github.com/aws/aws-sdk-go/aws/session"
)

func main() {
    sess, _ := session.NewSession(&aws.Config{
        Region:      aws.String("us-west-1"),
        Credentials: credentials.NewSharedCredentials("", ""),
    })
//    fmt.Println(reflect.TypeOf(int64(2)))

    ecsClient := ecs.New(sess)
    /*
    updateServiceInput := &ecs.UpdateServiceInput{
        Cluster:        aws.String("default"),
        DesiredCount:   aws.Int64(1),
        Service:        aws.String("sample-webapp"),
    }
    updateServiceResult, updateServiceErr := ecsClient.UpdateService(updateServiceInput)
    fmt.Printf("update service result: %s, error: %s", updateServiceResult, updateServiceErr)
    */

    registerTaskDefinitionInput := &ecs.RegisterTaskDefinitionInput{
        ContainerDefinitions: []*ecs.ContainerDefinition{
            {
              Command: []*string{
                  aws.String("sleep"),
                  aws.String("360"),
              },
              Cpu:          aws.Int64(10),
              Essential:    aws.Bool(true),
              Image:        aws.String("python-test:latest"),
              Memory:       aws.Int64(10),
              Name:         aws.String("test-deploy"),
              PortMappings:  []*ecs.PortMapping{
                  {
                      ContainerPort:   aws.Int64(8080),
                  },
              },
            },
        },
        Family:       aws.String("test-deploy"),
        TaskRoleArn:  aws.String(""),
    }
    fmt.Println(reflect.TypeOf(registerTaskDefinitionInput))

    registerTaskDefinitionResult, registerTaskDefinitionError := ecsClient.RegisterTaskDefinition(registerTaskDefinitionInput)
    fmt.Printf("result: %s, error: %s", registerTaskDefinitionResult, registerTaskDefinitionError)

    createServiceInput := &ecs.CreateServiceInput{
        DesiredCount:         aws.Int64(1),
        ServiceName:          aws.String("test-create-service"),
        TaskDefinition:       aws.String("test-deploy"),
    }
    createServiceResult, createServiceError := ecsClient.CreateService(createServiceInput)
    fmt.Printf("create service result: %s, error: %s", createServiceResult, createServiceError)
}
