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
		fmt.Println("Error", err)
		return
	}

	fmt.Printf("user %s created %v\n", *result.User.UserName, result.User.CreateDate)

	mfa, err := svc.ListMFADevices(&iam.ListMFADevicesInput{
		UserName: result.User.UserName,
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	for i, device := range mfa.MFADevices {
		if device == nil {
			continue
		}
		fmt.Printf("device info: %s for user: %s, device number: %d\n", *device.SerialNumber, *device.UserName, i)
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
		fmt.Println("Error", err)
		return
	}

	fmt.Printf("export AWS_ACCESS_KEY_ID=%s;\nexport AWS_SECRET_ACCESS_KEY=%s;\nexport AWS_SESSION_TOKEN=%s;\n",
		*sessionToken.Credentials.AccessKeyId,
		*sessionToken.Credentials.SecretAccessKey,
		*sessionToken.Credentials.SessionToken)
}
