package service

import (
	"os"
	"os/exec"

	"github.com/helloferdie/golib/liblogger"
	"github.com/helloferdie/golib/libresponse"
	"github.com/helloferdie/golib/libvalidator"
	"github.com/helloferdie/pusher/pkg/deploy/request"
	"github.com/pquerna/otp/totp"
)

// Deploy -
func Deploy(r *request.Deploy) *libresponse.Response {
	resp, err := libvalidator.Validate(r)
	if err != nil {
		return resp
	}

	secret := os.Getenv("otp_secret")
	valid := totp.Validate(r.Code, secret)
	if !valid {
		resp.Error = "deployer.otp.invalid"
		return resp.ErrorUnauthorized()
	}

	repository := os.Getenv("dir_repository") + "/" + r.Repository
	if _, err := os.Stat(repository); os.IsNotExist(err) {
		resp.Error = "deployer.repository.not_found"
		return resp.ErrorInternal()
	}

	if r.Async {
		liblogger.Log(nil, false).Infoln("Docker pull - Async Request Received")
		go Exec(repository)
	} else {
		err = Exec(repository)
		if err != nil {
			return resp.ErrorUpdate()
		}
	}

	return resp.SuccessDefault()
}

// Exec -
func Exec(repository string) error {
	dockerUser := os.Getenv("docker_user")
	dockerToken := os.Getenv("docker_token")

	cmd := exec.Command("sudo", "./dockerpull", repository, dockerUser, dockerToken)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		liblogger.Log(nil, false).Errorf("%v\n", err.Error())
		liblogger.Log(nil, false).Infoln("Docker pull - Fail")
		return err
	}

	liblogger.Log(nil, false).Infoln(string(stdout))
	liblogger.Log(nil, false).Infoln("Docker pull - Done")
	return nil
}
