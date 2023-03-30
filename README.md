# api-proxy-lambda

An API proxy designed for AWS serverless architecture, works with Lambda and API Gateway.

# Tutorial

## Requirements:

-   An Amazon Web Service account
-   A web browser
-   Latest Go runtime environment (optional, for build on local)

## Deployment steps:

### Setup Lambda function on AWS

1. Log in to the AWS Console.
2. Select a region from the top right corner. Choose a region that balances latency and cost.
3. In the search bar located in the upper left corner, enter "Lambda" and navigate to the Lambda service page.
4. On the left-side panel, select "Functions", then click the "Create function" button on the following page.
5. Select "Author from scratch". Enter a function name, such as "api-proxy" in this tutorial, select "Go 1.x" runtime, and choose "x86_64" architecture. Leave the other settings as default and click "Create function".
6. On the "api-proxy" function page, scroll down to the "Code" tab. Click "Upload from" - ".zip file" on the right, upload the compressed binary file, and wait for it to become ready.
7. Scroll down to "Runtime settings", click "Edit", input "api-proxy-lambda" in the "Handler" textbox, click "Save".
8. Click "Add trigger" on the "Function overview" diagram. On next page, select "API Gateway" and "Create a new API" as source. Choose "HTTP API" for API type and "Open" for security. Click "Add" to finish this page.
9. On "Configuration" tab, select "Environment variables" on the left side panel. Add the following two environment variables:

|Key | Value|
|------------ | ------------|
|DEST_DOMAIN|https://api.somewhere.com/|
|LOG_MODE|true|

The "DEST_DOMAIN" is a domain where you want your traffic be redirct to, it's necessary for this application.
The "LOG_MODE" is optional, if it's set to "true", you will get more logs printout on Cloudwatch.



### Setup API Gateway on AWS

1. Use the search to navigate to the API Gateway service page. You'll see an API was created by Lambda already, which may have a default name "api-proxy-API". Click to enter its detail page.
2. In the "Detail" page, you'll see a default stage was created, with an Invoke URL be like https://random-characters.execute-api.your-region.amazonaws.com/default following.
3. click "Routes" on the left side. You can see a default route leading all incoming traffic from the exact path "/api-proxy", it means all traffic to https://random-characters.execute-api.your-region.amazonaws.com/default/api-proxy will now perform to be a valid route, and only this specific path will work. However, it does not meet our requirements, we will improve it by adding a greedy path variable.
4. Click on the "Create" button, enter "/api-proxy/{proxy+}" in the textbox, and select "ANY" from the drop-down list. Then, click "Save".
5. Click on "Integrations", you may see a green "AWS Lambda" tag is displayed next to the "ANY" string under "/openai-api-proxy", but nothing for the second "ANY" string which under "/{proxy+}".Click on the second "ANY", select your Lambda function (api-proxy) from the drop-down list on the right panel, and click on the "Attach integration" button.

Now that you have completed the implementation of this api-proxy on AWS, it should be functional.
Send an HTTP request to the API Gateway's Invoke URL with the function path, which might be like:
https://random-characters.execute-api.your-region.amazonaws.com/default/api-proxy/and/whatever/you-want-to-add

And this application will redirect your request to:
https://api.somewhere.com/and/whatever/you-want-to-add


### Build the source code and compress binary file

Since the Lambda function supports two different architectures and is based on a Linux environment, the Go build settings should reflect this.

For this tutorial, I use:
```
$env:GOOS = "linux" 
$env:GOARCH = "amd64" 
$env:CGO_ENABLED = "0"`
```

The compressed file should only contain the binary file and should be in zip format before uploading to Lambda.



# Security and Cost Considerations

## Security

Please note that this tutorial does not incorporate any form of authorization on the API Gateway. This means that anyone who obtains your API Gateway's endpoint can access and use your API. We strongly recommend that you keep your endpoint secret until you have implemented adequate security measures to protect your API.

## Cost

Lambda and API Gateway are both paid features, with one or more of them having a limited free quota. The region where you deploy these services also affects the pricing. For more detailed information, please refer to AWS's pricing page.
