package main

import (
        "fmt"
//        "reflect"
        "strings"
        "github.com/aws/aws-sdk-go/service/ecs"
        "github.com/aws/aws-sdk-go/aws"
//        "github.com/aws/aws-sdk-go/aws/awserr"
        "github.com/aws/aws-sdk-go/aws/credentials"
//        "github.com/aws/aws-sdk-go/aws/request"
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
        s = append(s, stringInDynArray{strings.Split(*item, "/")[1]})
    }
    //  s now is a list containing all services from faas cluster
    fmt.Printf("String array is %s, %s", s[1].Value())
    fmt.Printf("list services result is %s, error is %s", res.ServiceArns, err)

    // describe services
    descServicesInput := &ecs.DescribeServicesInput {
        Services: []*string{
              aws.String("arn:aws:ecs:us-west-1:204338471371:service/sample-webapp"),
        },
    }

    descServicesResult, descServicesErr := ecsClient.DescribeServices(descServicesInput)
    fmt.Printf("Desribe services result: %s, error: %s", descServicesResult, descServicesErr)

    // get task definition
    taskDefinitionInput := &ecs.DescribeTaskDefinitionInput {
        TaskDefinition: aws.String("console-sample-app-static:1"),
    }

    result, errrr := ecsClient.DescribeTaskDefinition(taskDefinitionInput)
    fmt.Printf("\nDescribe Task Definition result: %s, error: %s", result, errrr)
}
