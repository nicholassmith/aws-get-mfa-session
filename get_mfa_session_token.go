package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	// Create a IAM service client.
	svc := iam.New(sess)

	result, err := svc.GetUser(&iam.GetUserInput{})

	if err != nil {
		fmt.Println("Error getting user: ", err)
		return
	}

	mfa, err := svc.ListMFADevices(&iam.ListMFADevicesInput{
		UserName: result.User.UserName,
	})

	if err != nil {
		fmt.Println("Error getting MFA Devices: ", err)
		return
	}

	if len(mfa.MFADevices) == 0 {
		fmt.Println("No devices found")
		return
	}

	serialNumber := *mfa.MFADevices[0].SerialNumber

	stsSvc := sts.New(sess)

	fmt.Print("Enter MFA token without spaces: ")
	var mfaToken string
	fmt.Scanln(&mfaToken)

	sessionToken, err := stsSvc.GetSessionToken(&sts.GetSessionTokenInput{
		SerialNumber: &serialNumber,
		TokenCode:    &mfaToken,
	})

	if err != nil {
		fmt.Println("Error getting session token for input: ", err)
		return
	}

	fmt.Printf("export AWS_ACCESS_KEY_ID=%s;\nexport AWS_SECRET_ACCESS_KEY=%s;\nexport AWS_SESSION_TOKEN=%s;\n",
		*sessionToken.Credentials.AccessKeyId,
		*sessionToken.Credentials.SecretAccessKey,
		*sessionToken.Credentials.SessionToken)
}
