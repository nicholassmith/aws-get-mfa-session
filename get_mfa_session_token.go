package main

import (
	"fmt"
	"flag"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
)

func main() {
	token := flag.String("token", "", "MFA Token")
	profile := flag.String("profile", "", "AWS Profile to load")
	flag.Parse()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
     Profile: *profile,
	}))

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


	var mfaToken string

	if *token == "" {
		fmt.Print("Enter MFA token without spaces: ")
		fmt.Scanln(&mfaToken)
	} else {
		mfaToken = *token
	}

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
