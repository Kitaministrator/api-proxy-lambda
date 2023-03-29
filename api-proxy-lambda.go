package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Body           string `json:"body"`
	Resource       string `json:"resource"`
	RequestPath    string `json:"path"`
	HttpMethod     string `json:"httpMethod"`
	Headers        map[string]string
	QueryString    string `json:"queryStringParameters"`
	PathParameters map[string]string
}

type Response events.APIGatewayProxyResponse

// Set envs, all of these envs can only be set in AWS Lambda Function-Configuration-Environment variables, not in the code
// destDomain is the destination you want to redirect to, which is necessary,
// with only the domain name should be set, no need to add path
var destDomain = os.Getenv("DEST_DOMAIN") // https://api.somedestination.com/
var logMode = os.Getenv("LOG_MODE")       // string, true or false or empty

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (Response, error) {
	// Get the path parameters from the request
	distPath := req.PathParameters["proxy"]
	newUrl := destDomain + distPath

	// Use new destination to create a new request, modify the request header if necessary
	newReq, err := http.NewRequest(req.HTTPMethod, newUrl, strings.NewReader(req.Body))
	if err != nil {
		return Response{}, err
	}

	// Set the headers of the new request, replace the host header with the new destination
	for k, v := range req.Headers {
		newReq.Header.Set(k, v)
	}
	newReq.Header.Set("Host", destDomain)

	// Send the new request to the new destination
	client := http.Client{}
	response, err := client.Do(newReq)
	if err != nil {
		return Response{}, err
	}

	defer response.Body.Close()

	// Create the response to be returned to the sender
	forwardResponse := Response{
		StatusCode: response.StatusCode,
		Headers:    make(map[string]string),
	}

	// Copy the headers from the response to the forward response
	for k, v := range response.Header {
		forwardResponse.Headers[k] = strings.Join(v, ",")
	}

	// Read the response body and copy it to the forward response
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return Response{}, err
	}
	forwardResponse.Body = string(bodyBytes)

	// Debug output
	if logMode == "true" {
		// check basic properties
		log.Println("New Domain: " + destDomain)                                       // https://api.somedestination.com/
		log.Println("New URL: " + newUrl)                                              // https://api.somedestination.com/proxy_path
		log.Println("Method: " + req.HTTPMethod)                                       // GET or POST or ...
		log.Println("req.Path: " + req.Path)                                           // /stage/route/proxy_path
		log.Println("RequestContext.Protocol: " + req.RequestContext.Protocol)         // HTTP/1.1
		log.Println("RequestContext.DomainName: " + req.RequestContext.DomainName)     // mydomain.execute-api.someregion.amazonaws.com
		log.Println("RequestContext.Prefix: " + req.RequestContext.DomainPrefix)       // mydomain
		log.Println("RequestContext.ResourceID: " + req.RequestContext.ResourceID)     // ANY /route/{proxy+}
		log.Println("RequestContext.ResourcePath: " + req.RequestContext.ResourcePath) // /route/{proxy+}
		for k, v := range req.PathParameters {
			log.Print("k: " + k + ", v: " + v) // k: proxy, v: proxy_path
		}

		// print the request header use httputil
		dumpReqHeader, err := httputil.DumpRequest(newReq, true)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("========  Start to print Request Header  ========")
		log.Println(string(dumpReqHeader))
		log.Println("========  End of printing Request Header  ========")

		// print the request body use httputil
		dumpReqBody, err := httputil.DumpRequest(newReq, true)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("========  Start to print Request Body  ========")
		log.Println(string(dumpReqBody))
		log.Println("========  End of printing Request Body  ========")

		// print the final response send back to api gateway
		// it's an apigatewayproxyresponse type, cannot use the fancy way to print
		log.Println("========  Start to print Forward Response  ========")
		log.Printf("Forward response: %d %s %s", forwardResponse.StatusCode, forwardResponse.Headers, forwardResponse.Body)
		log.Println("========  End of printing Forward Response  ========")
	}

	return forwardResponse, nil
}

func main() {
	lambda.Start(HandleRequest)
}
