// Package amazonssm implements activities for interaction with Amazon Simple Storage Manager (SSM)
package amazonssm

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const (
	ivAction                     = "action"
	ivAwsAccessKeyID             = "awsAccessKeyID"
	ivAwsSecretAccessKey         = "awsSecretAccessKey"
	ivAwsRegion                  = "awsRegion"
	ivParameterName              = "parameterName"
	ivDecryptParameter           = "decryptParameter"
	ivParameterValue             = "parameterValue"
	ivOverwriteExistingParameter = "overwriteExistingParameter"
	ivParameterType              = "parameterType"
	ovResult                     = "result"
)

// log is the default package logger
var log = logger.GetLogger("activity-amazonssm")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// Get the action
	action := context.GetInput(ivAction).(string)
	awsAccessKeyID := context.GetInput(ivAwsAccessKeyID).(string)
	awsSecretAccessKey := context.GetInput(ivAwsSecretAccessKey).(string)
	awsRegion := context.GetInput(ivAwsRegion).(string)

	// Create new credentials using the accessKey and secretKey
	awsCredentials := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")

	// Create a new session to AWS
	awsSession := session.Must(session.NewSession(&aws.Config{
		Credentials: awsCredentials,
		Region:      aws.String(awsRegion),
	}))

	ssmSession := ssm.New(awsSession)

	// See which action needs to be taken
	switch action {
	case "store":
		parameterName := context.GetInput(ivParameterName).(string)
		overwriteExistingParameter := context.GetInput(ivOverwriteExistingParameter).(bool)
		parameterType := context.GetInput(ivParameterType).(string)
		parameterValue := context.GetInput(ivParameterValue).(string)
		dat := make(map[string]interface{})

		val, err := putSSMParameter(ssmSession, parameterName, overwriteExistingParameter, parameterType, parameterValue)
		if err != nil {
			log.Errorf("Error while retrieving parameter from SSM [%s]", err)
			return true, err
		}
		dat[parameterName] = val

		val1, err1 := prepareMapOutput(dat)
		if err1 != nil {
			log.Errorf("Error while retrieving parameter from SSM [%s]", err1)
			return true, err1
		}

		context.SetOutput(ovResult, val1)
		return true, nil
	case "retrieve":
		parameterName := context.GetInput(ivParameterName).(string)
		decryptParameter := context.GetInput(ivDecryptParameter).(bool)
		dat := make(map[string]interface{})

		val, err := getSSMParameter(ssmSession, parameterName, decryptParameter)
		if err != nil {
			log.Errorf("Error while retrieving parameter from SSM [%s]", err)
			return true, err
		}
		dat[parameterName] = val

		val1, err1 := prepareMapOutput(dat)
		if err1 != nil {
			log.Errorf("Error while retrieving parameter from SSM [%s]", err1)
			return true, err1
		}

		context.SetOutput(ovResult, val1)
		return true, nil
	case "retrieveList":
		parameterNames := context.GetInput(ivParameterName).(string)
		decryptParameter := context.GetInput(ivDecryptParameter).(bool)
		parameters := strings.Split(parameterNames, ",")
		dat := make(map[string]interface{})

		for _, parameterName := range parameters {
			val, err := getSSMParameter(ssmSession, parameterName, decryptParameter)
			if err != nil {
				log.Errorf("Error while retrieving parameter from SSM [%s]", err)
				return true, err
			}
			dat[parameterName] = val
		}

		val1, err1 := prepareMapOutput(dat)
		if err1 != nil {
			log.Errorf("Error while retrieving parameter from SSM [%s]", err1)
			return true, err1
		}

		context.SetOutput(ovResult, val1)
		return true, nil
	}

	// Set the output value in the context
	context.SetOutput(ovResult, "NOK")

	return true, nil
}

// getSSMParameter gets a parameter from the AWS Simple Systems Manager service.
func getSSMParameter(ssmSession *ssm.SSM, name string, decrypt bool) (string, error) {
	gpi := &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(decrypt),
	}

	param, err := ssmSession.GetParameter(gpi)
	if err != nil {
		return "", err
	}

	return *param.Parameter.Value, nil
}

// putSSMParameter puts a parameter in the AWS Simple Systems Manager service.
func putSSMParameter(ssmSession *ssm.SSM, name string, overwrite bool, paramtype string, value string) (int64, error) {
	ppi := &ssm.PutParameterInput{
		Name:      aws.String(name),
		Overwrite: aws.Bool(overwrite),
		Type:      aws.String(paramtype),
		Value:     aws.String(value),
	}

	param, err := ssmSession.PutParameter(ppi)
	if err != nil {
		return -1, err
	}

	return *param.Version, nil
}

// Prepare the output format required
func prepareMapOutput(dat map[string]interface{}) (map[string]interface{}, error) {
	jsonString, _ := json.Marshal(dat)
	var resultinterface interface{}
	d := json.NewDecoder(bytes.NewReader(jsonString))
	d.UseNumber()
	err := d.Decode(&resultinterface)
	if err != nil {
		return nil, err
	}
	f := map[string]interface{}{"results": resultinterface}
	return f, nil
}
