# commhub-engine
An API service and data model to support property and restaurant management.

The golang code and RDBMS SQL for generic property and restaurant management work order tracking API. It's based on stateless patterns for communication between client and server. 

It's old and uses a traditional relational DB model, but it can be updated to use a modern NoSQL DB so you won't need the stored procedures.
It was written for AWS Lambda functions to serve as a serverless endpoint for the API calls through the gateway.  But it can be ported to run on any
serverless public cloud.
