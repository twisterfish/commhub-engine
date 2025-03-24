package main

import (
	"encoding/json"
	//"log"
	"context"
	"datastores"
	"strings"
	"switchboard"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// ////////////////////////////////////////////////////////////////////////////////////////
// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
// ////////////////////////////////////////////////////////////////////////////////////////
type Response events.APIGatewayProxyResponse

// ////////////////////////////////////////////////////////////////////////////////////////
// Main
// ////////////////////////////////////////////////////////////////////////////////////////
func main() {
	lambda.Start(Handler)
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Handler is our lambda handler invoked by the `lambda.Start` function call
// ////////////////////////////////////////////////////////////////////////////////////////
func Handler(ctx context.Context, request *events.APIGatewayProxyRequest) (Response, error) {
	//log.Print("enter")
	//log.Print(request)
	//log.Print(ctx)
	apiresponse := validatePayload(strings.TrimSpace(request.Body), request)
	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            apiresponse, //buf.String(),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
			"X-commhub-Signals-Replay":    "signals-handler",
		},
	}
	return resp, nil
}

// ////////////////////////////////////////////////////////////////////////////////////////
// This validates for well formed JSON then sends it to be processed
// ////////////////////////////////////////////////////////////////////////////////////////
func validatePayload(payload string, request *events.APIGatewayProxyRequest) string {
	/*
		Struct to convert payload JSON string for pertinent values used in validation:
		Token: is the api auth token
		Signal: Tells microservice where to process request
		Action: Tells microservice how to process request
	*/
	type EntryPointData *struct {
		APIToken string `json:"api_token"`
		Signal   string `json:"signal"`
		Action   string `json:"action"`
	}
	var input EntryPointData
	pByte := []byte(payload)
	if err := json.Unmarshal(pByte, &input); err != nil {
		return "{\"signal\":\"error\",\"action\":\"" + err.Error() + "\"}" // Malformed JSON - kick it back
	} else {
		return validateAPIToken(input.APIToken, input.Signal, input.Action, payload, request) // all good
		//return guids.GetGUID() // all good
	}
}

// ////////////////////////////////////////////////////////////////////////////////////////
// This validates for current API token from user
// ////////////////////////////////////////////////////////////////////////////////////////
func validateAPIToken(api_token string, signal string, action string, payload string, request *events.APIGatewayProxyRequest) string {
	// This is a request for an api token - send to switchboard directly
	if signal == "authorize" {
		return switchboard.RouteSignal(signal, action, payload, request)
	}
	// process the validation and data
	type outputData struct {
		Status            string `json:"out_status"`
		Debug_Mode        int    `json:"out_dbg_mode"`
		Sample_Mode       int    `json:"out_smpl_mode"`
		Sample_Rate       int    `json:"out_smpl_rate"`
		Sample_Duration   int    `json:"out_smpl_duration"`
		Sample_Start_Time int    `json:"out_smpl_time_start"`
	}
	var output outputData
	db, err := datastores.OpenRDS()
	//defer db.Close()
	if err != nil {
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	}
	query := "CALL commhub_junction.validate_api_token(\"" + strings.TrimSpace(api_token) + "\")"
	results, err := db.Query(query)
	if err != nil {
		//results.Close()
		return "{\"signal\":\"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
		for results.Next() {
			// for each record, scan the result into our  struct
			err = results.Scan(&output.Status, &output.Debug_Mode, &output.Sample_Mode, &output.Sample_Rate, &output.Sample_Duration, &output.Sample_Start_Time)
			switch strings.TrimSpace(output.Status) {
			case "1":
				results.Close()
				return switchboard.RouteSignal(signal, action, payload, request)
			case "2":
				results.Close()
				return switchboard.RouteSignal(signal, action, payload, request)
			case "3":
				results.Close()
				return switchboard.RouteSignal(signal, action, payload, request)
			case "4":
				results.Close()
				return switchboard.RouteSignal(signal, action, payload, request)
			default:
				return "{\"signal\":\"error\",\"action\":\"invalid-api-entry: " + strings.TrimSpace(output.Status) + "\"}"
			}
		}
	}
	//results.Close()
	return switchboard.RouteSignal(signal, action, payload, request)
}
