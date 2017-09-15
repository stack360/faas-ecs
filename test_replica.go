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

func main() {
    sess, _ := session.NewSession(&aws.Config{
        Region:      aws.String("us-west-1"),
        Credentials: credentials.NewSharedCredentials("", ""),
    })
//    fmt.Println(reflect.TypeOf(int64(2)))
    ecsClient := ecs.New(sess)
    updateServiceInput := &ecs.UpdateServiceInput{
        Cluster:        aws.String("default"),
        DesiredCount:   aws.Int64(1),
        Service:        aws.String("sample-webapp"),
    }
    updateServiceResult, updateServiceErr := ecsClient.UpdateService(updateServiceInput)
    fmt.Printf("update service result: %s, error: %s", updateServiceResult, updateServiceErr)
}
