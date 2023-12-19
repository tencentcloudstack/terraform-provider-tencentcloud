package acctest

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"

	tcprovider "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud"
	providercommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

var AccProviders map[string]*schema.Provider
var AccProvider *schema.Provider

const (
	ACCOUNT_TYPE_INTERNATIONAL        = "INTERNATIONAL"
	ACCOUNT_TYPE_PREPAY               = "PREPAY"
	ACCOUNT_TYPE_COMMON               = "COMMON"
	ACCOUNT_TYPE_PRIVATE              = "PRIVATE"
	ACCOUNT_TYPE_SUB_ACCOUNT          = "SUB_ACCOUNT"
	ACCOUNT_TYPE_SMS                  = "SMS"
	ACCOUNT_TYPE_SES                  = "SES"
	ACCOUNT_TYPE_TSF                  = "TSF"
	ACCOUNT_TYPE_SSL                  = "SSL"
	ACCOUNT_TYPE_ORGANIZATION         = "ORGANIZATION"
	INTERNATIONAL_PROVIDER_SECRET_ID  = "TENCENTCLOUD_SECRET_ID_INTERNATIONAL"
	INTERNATIONAL_PROVIDER_SECRET_KEY = "TENCENTCLOUD_SECRET_KEY_INTERNATIONAL"
	PREPAY_PROVIDER_SECRET_ID         = "TENCENTCLOUD_SECRET_ID_PREPAY"
	PREPAY_PROVIDER_SECRET_KEY        = "TENCENTCLOUD_SECRET_KEY_PREPAY"
	PRIVATE_PROVIDER_SECRET_ID        = "TENCENTCLOUD_SECRET_ID_PRIVATE"
	PRIVATE_PROVIDER_SECRET_KEY       = "TENCENTCLOUD_SECRET_KEY_PRIVATE"
	COMMON_PROVIDER_SECRET_ID         = "TENCENTCLOUD_SECRET_ID_COMMON"
	COMMON_PROVIDER_SECRET_KEY        = "TENCENTCLOUD_SECRET_KEY_COMMON"
	SUB_ACCOUNT_PROVIDER_SECRET_ID    = "TENCENTCLOUD_SECRET_ID_SUB_ACCOUNT"
	SUB_ACCOUNT_PROVIDER_SECRET_KEY   = "TENCENTCLOUD_SECRET_KEY_SUB_ACCOUNT"
	SMS_PROVIDER_SECRET_ID            = "TENCENTCLOUD_SECRET_ID_SMS"
	SMS_PROVIDER_SECRET_KEY           = "TENCENTCLOUD_SECRET_KEY_SMS"
	TSF_PROVIDER_SECRET_ID            = "TENCENTCLOUD_SECRET_ID_TSF"
	TSF_PROVIDER_SECRET_KEY           = "TENCENTCLOUD_SECRET_KEY_TSF"
	SSL_PROVIDER_SECRET_ID            = "TENCENTCLOUD_SECRET_ID_SSL"
	SSL_PROVIDER_SECRET_KEY           = "TENCENTCLOUD_SECRET_KEY_SSL"
	ORGANIZATION_PROVIDER_SECRET_ID   = "TENCENTCLOUD_SECRET_ID_ORGANIZATION"
	ORGANIZATION_PROVIDER_SECRET_KEY  = "TENCENTCLOUD_SECRET_KEY_ORGANIZATION"
)

func AccPreCheck(t *testing.T) {
	if v := os.Getenv(tcprovider.PROVIDER_SECRET_ID); v == "" {
		t.Fatalf("%v must be set for acceptance tests\n", tcprovider.PROVIDER_SECRET_ID)
	}
	if v := os.Getenv(tcprovider.PROVIDER_SECRET_KEY); v == "" {
		t.Fatalf("%v must be set for acceptance tests\n", tcprovider.PROVIDER_SECRET_KEY)
	}
	if v := os.Getenv(tcprovider.PROVIDER_REGION); v == "" {
		log.Printf("[INFO] Testing: Using %s as test region", DefaultRegion)
		os.Setenv(tcprovider.PROVIDER_REGION, DefaultRegion)
	}
	if v := os.Getenv(COMMON_PROVIDER_SECRET_ID); v != "" {
		secretId := os.Getenv(COMMON_PROVIDER_SECRET_ID)
		os.Setenv(tcprovider.PROVIDER_SECRET_ID, secretId)
	}
	if v := os.Getenv(COMMON_PROVIDER_SECRET_KEY); v != "" {
		secretKey := os.Getenv(COMMON_PROVIDER_SECRET_KEY)
		os.Setenv(tcprovider.PROVIDER_SECRET_KEY, secretKey)
	}
}

func init() {
	AccProvider = tcprovider.Provider()
	AccProviders = map[string]*schema.Provider{
		"tencentcloud": AccProvider,
	}
	envProject := os.Getenv("QCI_JOB_ID")
	envNum := os.Getenv("QCI_BUILD_NUMBER")
	envId := os.Getenv("QCI_BUILD_ID")
	reqCli := fmt.Sprintf("Terraform-%s/%s-%s", envProject, envNum, envId)
	os.Setenv(connectivity.REQUEST_CLIENT, reqCli)
}

func AccStepPreConfigSetTempAKSK(t *testing.T, accountType string) {
	AccPreCheckCommon(t, accountType)
}

func AccStepSetRegion(t *testing.T, region string) {
	os.Setenv(tcprovider.PROVIDER_REGION, region)
}

func AccPreCheckCommon(t *testing.T, accountType string) {
	if v := os.Getenv(tcprovider.PROVIDER_REGION); v == "" {
		log.Printf("[INFO] Testing: Using %s as test region", DefaultRegion)
		os.Setenv(tcprovider.PROVIDER_REGION, DefaultRegion)
	}
	switch {
	case accountType == ACCOUNT_TYPE_INTERNATIONAL:
		secretId := os.Getenv(INTERNATIONAL_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(INTERNATIONAL_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", INTERNATIONAL_PROVIDER_SECRET_ID, INTERNATIONAL_PROVIDER_SECRET_KEY)
		}
		os.Setenv(tcprovider.PROVIDER_SECRET_ID, secretId)
		os.Setenv(tcprovider.PROVIDER_SECRET_KEY, secretKey)
	case accountType == ACCOUNT_TYPE_PREPAY:
		secretId := os.Getenv(PREPAY_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(PREPAY_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", PREPAY_PROVIDER_SECRET_ID, PREPAY_PROVIDER_SECRET_KEY)
		}
		os.Setenv(tcprovider.PROVIDER_SECRET_ID, secretId)
		os.Setenv(tcprovider.PROVIDER_SECRET_KEY, secretKey)
	case accountType == ACCOUNT_TYPE_PRIVATE:
		secretId := os.Getenv(PRIVATE_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(PRIVATE_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", PRIVATE_PROVIDER_SECRET_ID, PRIVATE_PROVIDER_SECRET_KEY)
		}
		os.Setenv(tcprovider.PROVIDER_SECRET_ID, secretId)
		os.Setenv(tcprovider.PROVIDER_SECRET_KEY, secretKey)
	case accountType == ACCOUNT_TYPE_COMMON:
		secretId := os.Getenv(COMMON_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(COMMON_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", COMMON_PROVIDER_SECRET_ID, COMMON_PROVIDER_SECRET_KEY)
		}
		os.Setenv(tcprovider.PROVIDER_SECRET_ID, secretId)
		os.Setenv(tcprovider.PROVIDER_SECRET_KEY, secretKey)
	case accountType == ACCOUNT_TYPE_SUB_ACCOUNT:
		secretId := os.Getenv(SUB_ACCOUNT_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(SUB_ACCOUNT_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", SUB_ACCOUNT_PROVIDER_SECRET_ID, SUB_ACCOUNT_PROVIDER_SECRET_KEY)
		}
		os.Setenv(tcprovider.PROVIDER_SECRET_ID, secretId)
		os.Setenv(tcprovider.PROVIDER_SECRET_KEY, secretKey)
	case accountType == ACCOUNT_TYPE_SMS:
		secretId := os.Getenv(SMS_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(SMS_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", SMS_PROVIDER_SECRET_ID, SMS_PROVIDER_SECRET_KEY)
		}
		os.Setenv(tcprovider.PROVIDER_SECRET_ID, secretId)
		os.Setenv(tcprovider.PROVIDER_SECRET_KEY, secretKey)
	case accountType == ACCOUNT_TYPE_TSF:
		secretId := os.Getenv(TSF_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(TSF_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", TSF_PROVIDER_SECRET_ID, TSF_PROVIDER_SECRET_KEY)
		}
		os.Setenv(tcprovider.PROVIDER_SECRET_ID, secretId)
		os.Setenv(tcprovider.PROVIDER_SECRET_KEY, secretKey)
	case accountType == ACCOUNT_TYPE_SSL:
		secretId := os.Getenv(SSL_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(SSL_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", SSL_PROVIDER_SECRET_ID, SSL_PROVIDER_SECRET_KEY)
		}
		os.Setenv(tcprovider.PROVIDER_SECRET_ID, secretId)
		os.Setenv(tcprovider.PROVIDER_SECRET_KEY, secretKey)
	case accountType == ACCOUNT_TYPE_ORGANIZATION:
		secretId := os.Getenv(ORGANIZATION_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(ORGANIZATION_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", ORGANIZATION_PROVIDER_SECRET_ID, ORGANIZATION_PROVIDER_SECRET_KEY)
		}
		os.Setenv(tcprovider.PROVIDER_SECRET_ID, secretId)
		os.Setenv(tcprovider.PROVIDER_SECRET_KEY, secretKey)
	default:
		if v := os.Getenv(tcprovider.PROVIDER_SECRET_ID); v == "" {
			t.Fatalf("%v must be set for acceptance tests\n", tcprovider.PROVIDER_SECRET_ID)
		}
		if v := os.Getenv(tcprovider.PROVIDER_SECRET_KEY); v == "" {
			t.Fatalf("%v must be set for acceptance tests\n", tcprovider.PROVIDER_SECRET_KEY)
		}
	}
}

func AccCheckTencentCloudDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("data source ID not set")
		}
		return nil
	}
}

func AccPreCheckBusiness(t *testing.T, accountType string) {

	switch accountType {
	case ACCOUNT_TYPE_SES:
		if v := os.Getenv(tcprovider.PROVIDER_SECRET_ID); v == "" {
			t.Fatalf("%v must be set for acceptance tests\n", tcprovider.PROVIDER_SECRET_ID)
		}
		if v := os.Getenv(tcprovider.PROVIDER_SECRET_KEY); v == "" {
			t.Fatalf("%v must be set for acceptance tests\n", tcprovider.PROVIDER_SECRET_KEY)
		}
		os.Setenv(tcprovider.PROVIDER_REGION, DefaultRegionSes)
	default:
		AccPreCheck(t)
	}
}

type TencentCloudClient struct {
	apiV3Conn *connectivity.TencentCloudClient
}

var _ providercommon.ProviderMeta = &TencentCloudClient{}

// GetAPIV3Conn 返回访问云 API 的客户端连接对象
func (meta *TencentCloudClient) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return meta.apiV3Conn
}

func SharedClientForRegion(region string) (interface{}, error) {
	var secretId string
	if secretId = os.Getenv(tcprovider.PROVIDER_SECRET_ID); secretId == "" {
		return nil, fmt.Errorf("%s can not be empty", tcprovider.PROVIDER_SECRET_ID)
	}

	var secretKey string
	if secretKey = os.Getenv(tcprovider.PROVIDER_SECRET_KEY); secretKey == "" {
		return nil, fmt.Errorf("%s can not be empty", tcprovider.PROVIDER_SECRET_KEY)
	}

	securityToken := os.Getenv(tcprovider.PROVIDER_SECURITY_TOKEN)
	protocol := os.Getenv(tcprovider.PROVIDER_PROTOCOL)
	domain := os.Getenv(tcprovider.PROVIDER_DOMAIN)

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
