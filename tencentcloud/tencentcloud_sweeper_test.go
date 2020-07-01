package tencentcloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedClientForRegion(region string) (interface{}, error) {
	var secretId string
	if secretId = os.Getenv(PROVIDER_SECRET_ID); secretId == "" {
		return nil, fmt.Errorf("%s can not be empty", PROVIDER_SECRET_ID)
	}

	var secretKey string
	if secretKey = os.Getenv(PROVIDER_SECRET_KEY); secretKey == "" {
		return nil, fmt.Errorf("%s can not be empty", PROVIDER_SECRET_KEY)
	}

	securityToken := os.Getenv(PROVIDER_SECURITY_TOKEN)
	protocol := os.Getenv(PROVIDER_PROTOCOL)
	domain := os.Getenv(PROVIDER_DOMAIN)

	client := &connectivity.TencentCloudClient{
		Credential: common.NewTokenCredential(
			secretId,
			secretKey,
			securityToken,
		),
		Region:   region,
		Protocol: protocol,
		Domain:   domain,
	}

	var tcClient TencentCloudClient
	tcClient.apiV3Conn = client

	return &tcClient, nil
}
