package main

import (
	"fmt"
	"os"
	"time"

	"github.com/helloferdie/golib/libhttp"
	"github.com/helloferdie/golib/liblogger"
	"github.com/joho/godotenv"
	"github.com/pquerna/otp/totp"
)

func init() {
	godotenv.Load()
}

func main() {
	args := os.Args
	if len(args) <= 1 {
		liblogger.Log(nil, false).Errorf("%s", "Args is required")
		os.Exit(1)
	}

	repository := os.Args[1]
	if repository == "" {
		liblogger.Log(nil, false).Errorf("%s", "[Args 1] Repository is required")
		os.Exit(1)
	}

	// Generate TOTP
	otpSecret := os.Getenv("otp_secret")
	key, err := totp.GenerateCode(otpSecret, time.Now().UTC())
	if err != nil {
		liblogger.Log(nil, false).Errorf("%s", "Failed to generate OTP")
		liblogger.Log(nil, false).Errorf("%v", err.Error())
		os.Exit(1)
	}

	payload := map[string]interface{}{
		"repository": repository,
		"code":       key,
	}

	url := os.Getenv("webhook_url")
	headerSecretField := os.Getenv("header_secret_field")
	headerSecretValue := os.Getenv("header_secret_value")
	result, responseCode, err := libhttp.Request(url, "POST", payload, map[string]string{
		headerSecretField: headerSecretValue,
	})
	if err != nil {
		txt := fmt.Sprintf("Deploy request received with error code %v", responseCode)
		liblogger.Log(nil, false).Errorln(txt)
		liblogger.Log(nil, false).Errorf("%v", err)
		os.Exit(1)
	}

	txt := fmt.Sprintf("Deploy request received with code %v", responseCode)
	liblogger.Log(nil, false).Infoln(txt)
	liblogger.Log(nil, false).Infof("%v", result)
	os.Exit(0)
}
