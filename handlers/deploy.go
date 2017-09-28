package handlers

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "regexp"

        "github.com/alexellis/faas/gateway/requests"
        "github.com/stack360/faas-ecs/types"
        "github.com/gorilla/mux"
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/service/ecs"
)

const namespace string = "default"

func ValidateDeployRequest(request *requests.CreateFunctionRequest) error {
        var validDNS = regexp.MustCompile(`^[a-zA-Z\-]+$`)
        matched := validDNS.MatchString(request.Service)
        if matched {
                return nil
        }

        return fmt.Errorf("(%s) must be a valid DNS entry for service name", request.Service)
}

func MakeDeployHandler(ecsClient *ecs.ECS) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request){
                defer r.Body.Close()

                body, _ := ioutil.ReadAll(r.Body)

                request := requests.CreateFunctionRequest{}
                err  json.Unmarshal(body, &request)
                if err != nil {
                        w.WriteHeader(http.StatusBadRequest)
                        return
                }

                if err := ValidateDeployRequest(&request); err != nil {
                        w.WriteHeader(http.StatusBadRequest)
                        w.Write([]byte(err.Error()))
                        return
                }

                taskDefinitionInput := makeTaskDefinitionInput(request)
                registerTaskDefinitionResult, registerTaskDefinitionError := ecsClient.RegisterTaskDefinition(registerTaskDefinitionInput)
                log.Println("Registering Task Definition name: %s, Error: %s", registerTaskDefinitionInput, registerTaskDefinitionError)

                createServiceInput := makeServiceInput(request)
                createServiceResult, createServiceError := ecsClient.CreateService(createServiceInput)
                if createServiceError != nil {
                        w.WriteHeader(http.StatusBadRequest)
                        w.Write([]byte(createServiceError.Error()))
                        return
                }
                log.Println("Successfully created Service: %s", CreateServiceResult)
        }
}

func makeTaskDefinitionInput(request requests.CreateFunctionRequest) *ecs.RegisterTaskDefinitionInput {
        registerTaskDefinitionInput := &ecs.RegisterTaskDefinitionInput{
                ContainerDefinitions: []*ecs.ContainerDefinition{
                        {
                                Name:         aws.String(request.service),
                                Image:        aws.String(request.image),
                                Essential:    aws.Bool(true),
                                PortMappings: []*ecs.PortMapping{
                                        ContainerPort:   aws.Int64(8080),
                                },

                        },
                },
                Family:     aws.String(request.service),
                TaskRoleArn:  aws.String(""),
        }
        return registerTaskDefinitionInput
}

func makeServiceInput(request requests.CreateFunctionRequest) *ecsCreateServiceInput {
        createServiceInput := &ecs.CreateServiceInput{
                DesiredCount:         aws.Int64(1),
                ServiceName:          aws.String(request.service),
                TaskDefinition:       aws.String(request.service),
        }
        return createServiceInput
}
