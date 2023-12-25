package tencentcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

func init() {
	envProject := os.Getenv("QCI_JOB_ID")
	envNum := os.Getenv("QCI_BUILD_NUMBER")
	envId := os.Getenv("QCI_BUILD_ID")
	reqCli := fmt.Sprintf("Terraform-%s/%s-%s", envProject, envNum, envId)
	_ = os.Setenv(connectivity.REQUEST_CLIENT, reqCli)
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ = Provider()
}
