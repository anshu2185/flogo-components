# flogo-components
Collection custom built flogo components

## Components

### Activities
* [randomnumber](activity/randomnumber): Generate a random unique number between the min and max value
* [randomstring](activity/randomstring): Generate a random string consisting with the length you specify
* [addtodate](activity/addtodate): Adds a specified number of days to either the current date or a chosen date
* [dynamodbinsert](activity/dynamodbinsert): Insert a record in an Amazon DynamoDB
* [dynamodbquery](activity/dynamodbquery): Execute a query against an Amazon DynamoDB
* [gzip](activity/gzip): Activities for reading and writing of gzip format compressed files
* [amazons3](activity/amazons3): Activities for interacting with Amazon Simple Storage Service (S3)

### Apps
* [invoiceservice](apps/invoiceservice): A simple service listening to requests on a port exposed as environment variable, sending back random data and leveraging the [paymentservice](apps/paymentservice)
* [paymentservice](apps/paymentservice): A simple service listening to requests on a port exposed as environment variable, sending back a random date
* [serverless-demo](apps/serverless-demo): A collection of three Flogo apps deployed on AWS Lambda, where one app queries a DynamoDB, one app queries a MySQL instance and one app collects that data and presents it using an API Gateway
* [tci-combinator-app](apps/tci-combinator-app): An API spec and a Flogo app that work in TIBCO Cloud Integation that communicates with the apps from the [serverless-demo](apps/serverless-demo)