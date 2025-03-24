package uploads

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ////////////////////////////////////////////////////////////////////////////////////////
// Action Mapper
// ////////////////////////////////////////////////////////////////////////////////////////
func DoAction(signal string, action string, payload string, request *events.APIGatewayProxyRequest) string {
	switch action {
	case "GetPreSignedURLTest":
		return GetPreSignedURLTest(request, payload)
	case "GetPreSignedURL":
		return GetPreSignedURL(request, payload)
	default:
		return "{\"signal\":\"error\",\"action\": \"Your action request is invalid.\"}"
	}
}

// ////////////////////////////////////////////////////////////////////////////////////////
// Opens a connection to S3 and retrieves a Signed URL
// ContentType is the content type header the client must provide
// to use the generated signed URL. ContentType string  will look like image/png  mime type
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"uploads","action":"GetPreSignedURLTest","file_ext":"png","content_type":"image/png"}
// curl https://signals.commhubapi.com -d '{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"uploads","action":"GetPreSignedURLTest","file_ext":"png","content_type":"image/png"}'
// ////////////////////////////////////////////////////////////////////////////////////////
func GetPreSignedURLTest(request *events.APIGatewayProxyRequest, payload string) string {

	type ReqUpload struct {
		APIToken string `json:"api_token"`
		USRToken string `json:"user_token"`
		Signal   string `json:"signal"`
		Action   string `json:"action"`
		FileExt  string `json:"file_ext"`
		ContType string `json:"content_type"`
	}
	var rup ReqUpload

	pByte := []byte(payload)
	if err := json.Unmarshal(pByte, &rup); err != nil {
		return "{\"signal\": \"error\",\"action\": \"" + err.Error() + "\"}"
	}

	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create S3 service client
	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String("golang-test-one"),
		Key:    aws.String("DigMQ/PiOvLVHGSBgsLw+LrUt4RBpFYkhdvX+2hm"),
	})
	str, err := req.Presign(15 * time.Minute)

	if err != nil {
		return "{\"signal\": \"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
		return "{\"signal\": \"" + rup.Signal + "\",\"URL\": \"" + str + "\"}"
	}

}

// ////////////////////////////////////////////////////////////////////////////////////////
// Opens a connection to S3 and retrieves a Signed URL
// ContentType is the content type header the client must provide
// to use the generated signed URL. ContentType string  will look like image/png  mime type
// {"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"uploads","action":"GetPreSignedURLTest","file_ext":"png","content_type":"image/png"}
// curl https://signals.commhubapi.com -d '{"api_token":"37b176076fc74698be5aed02f74cbf15","signal":"uploads","action":"GetPreSignedURL","user_token":"37b176076fc74698be5aed02f74cbf15","bucket_name":"testbucket","folder_name":"somefoldername","file_name":"somefile","file_ext":"png","content_type":"image/png"}'
// ////////////////////////////////////////////////////////////////////////////////////////
func GetPreSignedURL(request *events.APIGatewayProxyRequest, payload string) string {

	type ReqUpload struct {
		APIToken string `json:"api_token"`
		Signal   string `json:"signal"`
		Action   string `json:"action"`
		USRToken string `json:"user_token"`
		BktName  string `json:"bucket_name"`
		FldrName string `json:"folder_name"`
		FileName string `json:"file_name"`
		FileExt  string `json:"file_ext"`
		ContType string `json:"content_type"`
	}
	var input ReqUpload

	pByte := []byte(payload)
	if err := json.Unmarshal(pByte, &input); err != nil {
		return "{\"signal\": \"error\",\"action\": \"" + err.Error() + "\"}"
	}

	// Initialize a session in us-east-1 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)

	// Create S3 service client
	svc := s3.New(sess)
	fullpath := input.FldrName + "/" + input.FileName
	bucket := input.BktName

	// Key:    aws.String("DigMQ/PiOvLVHGSBgsLw+LrUt4RBpFYkhdvX+2hm"),

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fullpath),
	})
	str, err := req.Presign(15 * time.Minute)

	if err != nil {
		return "{\"signal\": \"error\",\"action\": \"" + err.Error() + "\"}"
	} else {
		return "{\"signal\": \"" + input.Signal + "\",\"URL\": \"" + str + "\"}"
	}

}
