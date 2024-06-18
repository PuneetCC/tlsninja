package tlsninja

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type LambdaRequest struct {
	Function string
	Region   string
	Client   *lambda.Lambda
}

type ILambdaResponse struct {
	StatusCode int               `json:"statusCode"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
}

func NewLambdaRequest(function, region string) LambdaRequest {
	sess := session.Must(session.NewSession())
	client := lambda.New(sess, aws.NewConfig().WithRegion(region))

	return LambdaRequest{Function: function, Region: region, Client: client}
}

func (p *LambdaRequest) Do(config IRequestConfig) (*IRequestResponse, error) {
	payload, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	fmt.Println(config)

	fmt.Println(string(payload))
	result, err := p.Client.Invoke(&lambda.InvokeInput{
		FunctionName: &p.Function,
		Payload:      payload,
	})
	if err != nil {
		return nil, err
	}

	var lambdaResponse ILambdaResponse
	err = json.Unmarshal(result.Payload, &lambdaResponse)
	if err != nil {
		return nil, err
	}

	return &IRequestResponse{
		StatusCode: lambdaResponse.StatusCode,
		Body:       []byte(lambdaResponse.Body),
		Headers:    lambdaResponse.Headers,
	}, nil
}
