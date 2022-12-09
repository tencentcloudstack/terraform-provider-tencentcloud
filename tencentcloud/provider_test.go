package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

type MappedProviders map[string]terraform.ResourceProvider

var testAccProviders MappedProviders
var testAccProvider *schema.Provider

var providerCache = map[string]MappedProviders{}

var testAccContext = context.TODO()

const (
	ACCOUNT_TYPE_INTERNATIONAL        = "INTERNATIONAL"
	ACCOUNT_TYPE_PREPAY               = "PREPAY"
	ACCOUNT_TYPE_COMMON               = "COMMON"
	ACCOUNT_TYPE_PRIVATE              = "PRIVATE"
	ACCOUNT_TYPE_SUB_ACCOUNT          = "SUB_ACCOUNT"
	ACCOUNT_TYPE_SMS                  = "SMS"
	ACCOUNT_TYPE_SES                  = "SES"
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
)

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"tencentcloud": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ = Provider()
}

func useProvider(name string, override *ConfigureOverride) MappedProviders {
	if override.context != nil {
		if v, ok := override.context.Value(connectivity.LogTitleCtxKey).(string); ok {
			name += v
		}
	}
	providers, ok := providerCache[name]
	if !ok || providers == nil {
		provider := Provider().(*schema.Provider)
		if override != nil {
			provider.ConfigureFunc = CreateConfigureFunc(override)
		}
		providers = MappedProviders{
			"tencentcloud": provider,
		}
		providerCache[name] = providers
	}
	return providers
}

func useAccProvider(t *testing.T, accountTypeOverride ...string) MappedProviders {
	accountType := ""
	if len(accountTypeOverride) > 0 {
		accountType = accountTypeOverride[0]
	}
	override := &ConfigureOverride{}
	testNameRE := regexp.MustCompile(`^TestAcc(\w*)TencentCloud`)
	logTitle := testNameRE.ReplaceAllString(t.Name(), "") + "@"
	override.context = context.WithValue(testAccContext, connectivity.LogTitleCtxKey, logTitle)

	if v := os.Getenv(PROVIDER_REGION); v == "" {
		log.Printf("[INFO] Testing: Using %s as test region", defaultRegion)
		override.region = defaultRegion
	} else {
		override.region = v
	}
	secretId := ""
	secretKey := ""
	switch {
	case accountType == ACCOUNT_TYPE_INTERNATIONAL:
		secretId = os.Getenv(INTERNATIONAL_PROVIDER_SECRET_ID)
		secretKey = os.Getenv(INTERNATIONAL_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", INTERNATIONAL_PROVIDER_SECRET_ID, INTERNATIONAL_PROVIDER_SECRET_KEY)
		}
	case accountType == ACCOUNT_TYPE_PREPAY:
		secretId = os.Getenv(PREPAY_PROVIDER_SECRET_ID)
		secretKey = os.Getenv(PREPAY_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", PREPAY_PROVIDER_SECRET_ID, PREPAY_PROVIDER_SECRET_KEY)
		}
	case accountType == ACCOUNT_TYPE_PRIVATE:
		secretId = os.Getenv(PRIVATE_PROVIDER_SECRET_ID)
		secretKey = os.Getenv(PRIVATE_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", PRIVATE_PROVIDER_SECRET_ID, PRIVATE_PROVIDER_SECRET_KEY)
		}
	case accountType == ACCOUNT_TYPE_COMMON:
		secretId = os.Getenv(COMMON_PROVIDER_SECRET_ID)
		secretKey = os.Getenv(COMMON_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", COMMON_PROVIDER_SECRET_ID, COMMON_PROVIDER_SECRET_KEY)
		}
	case accountType == ACCOUNT_TYPE_SUB_ACCOUNT:
		secretId = os.Getenv(SUB_ACCOUNT_PROVIDER_SECRET_ID)
		secretKey = os.Getenv(SUB_ACCOUNT_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", SUB_ACCOUNT_PROVIDER_SECRET_ID, SUB_ACCOUNT_PROVIDER_SECRET_KEY)
		}
	case accountType == ACCOUNT_TYPE_SMS:
		secretId = os.Getenv(SMS_PROVIDER_SECRET_ID)
		secretKey = os.Getenv(SMS_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", SMS_PROVIDER_SECRET_ID, SMS_PROVIDER_SECRET_KEY)
		}
	default:
		secretId = os.Getenv(PROVIDER_SECRET_ID)
		secretKey = os.Getenv(PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", SUB_ACCOUNT_PROVIDER_SECRET_ID, SUB_ACCOUNT_PROVIDER_SECRET_KEY)
		}
	}

	override.secretId = secretId
	override.secretKey = secretKey
	return useProvider(accountType, override)
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv(PROVIDER_SECRET_ID); v == "" {
		t.Fatalf("%v must be set for acceptance tests\n", PROVIDER_SECRET_ID)
	}
	if v := os.Getenv(PROVIDER_SECRET_KEY); v == "" {
		t.Fatalf("%v must be set for acceptance tests\n", PROVIDER_SECRET_KEY)
	}
	if v := os.Getenv(PROVIDER_REGION); v == "" {
		log.Printf("[INFO] Testing: Using %s as test region", defaultRegion)
		os.Setenv(PROVIDER_REGION, defaultRegion)
	}
	if v := os.Getenv(COMMON_PROVIDER_SECRET_ID); v != "" {
		secretId := os.Getenv(COMMON_PROVIDER_SECRET_ID)
		os.Setenv(PROVIDER_SECRET_ID, secretId)
	}
	if v := os.Getenv(COMMON_PROVIDER_SECRET_KEY); v != "" {
		secretKey := os.Getenv(COMMON_PROVIDER_SECRET_KEY)
		os.Setenv(PROVIDER_SECRET_KEY, secretKey)
	}
}

func testAccStepPreConfigSetTempAKSK(t *testing.T, accountType string) {
	testAccPreCheckCommon(t, accountType)
}

func testAccPreCheckCommon(t *testing.T, accountType string) {
	if v := os.Getenv(PROVIDER_REGION); v == "" {
		log.Printf("[INFO] Testing: Using %s as test region", defaultRegion)
		os.Setenv(PROVIDER_REGION, defaultRegion)
	}
	switch {
	case accountType == ACCOUNT_TYPE_INTERNATIONAL:
		secretId := os.Getenv(INTERNATIONAL_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(INTERNATIONAL_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", INTERNATIONAL_PROVIDER_SECRET_ID, INTERNATIONAL_PROVIDER_SECRET_KEY)
		}
		os.Setenv(PROVIDER_SECRET_ID, secretId)
		os.Setenv(PROVIDER_SECRET_KEY, secretKey)
	case accountType == ACCOUNT_TYPE_PREPAY:
		secretId := os.Getenv(PREPAY_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(PREPAY_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", PREPAY_PROVIDER_SECRET_ID, PREPAY_PROVIDER_SECRET_KEY)
		}
		os.Setenv(PROVIDER_SECRET_ID, secretId)
		os.Setenv(PROVIDER_SECRET_KEY, secretKey)
	case accountType == ACCOUNT_TYPE_PRIVATE:
		secretId := os.Getenv(PRIVATE_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(PRIVATE_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", PRIVATE_PROVIDER_SECRET_ID, PRIVATE_PROVIDER_SECRET_KEY)
		}
		os.Setenv(PROVIDER_SECRET_ID, secretId)
		os.Setenv(PROVIDER_SECRET_KEY, secretKey)
	case accountType == ACCOUNT_TYPE_COMMON:
		secretId := os.Getenv(COMMON_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(COMMON_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", COMMON_PROVIDER_SECRET_ID, COMMON_PROVIDER_SECRET_KEY)
		}
		os.Setenv(PROVIDER_SECRET_ID, secretId)
		os.Setenv(PROVIDER_SECRET_KEY, secretKey)
	case accountType == ACCOUNT_TYPE_SUB_ACCOUNT:
		secretId := os.Getenv(SUB_ACCOUNT_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(SUB_ACCOUNT_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", SUB_ACCOUNT_PROVIDER_SECRET_ID, SUB_ACCOUNT_PROVIDER_SECRET_KEY)
		}
		os.Setenv(PROVIDER_SECRET_ID, secretId)
		os.Setenv(PROVIDER_SECRET_KEY, secretKey)
	case accountType == ACCOUNT_TYPE_SMS:
		secretId := os.Getenv(SMS_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(SMS_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", SMS_PROVIDER_SECRET_ID, SMS_PROVIDER_SECRET_KEY)
		}
		os.Setenv(PROVIDER_SECRET_ID, secretId)
		os.Setenv(PROVIDER_SECRET_KEY, secretKey)
	default:
		if v := os.Getenv(PROVIDER_SECRET_ID); v == "" {
			t.Fatalf("%v must be set for acceptance tests\n", PROVIDER_SECRET_ID)
		}
		if v := os.Getenv(PROVIDER_SECRET_KEY); v == "" {
			t.Fatalf("%v must be set for acceptance tests\n", PROVIDER_SECRET_KEY)
		}
	}
}

func testAccCheckTencentCloudDataSourceID(n string) resource.TestCheckFunc {
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

func testAccPreCheckBusiness(t *testing.T, accountType string) {

	switch accountType {
	case ACCOUNT_TYPE_SES:
		if v := os.Getenv(PROVIDER_SECRET_ID); v == "" {
			t.Fatalf("%v must be set for acceptance tests\n", PROVIDER_SECRET_ID)
		}
		if v := os.Getenv(PROVIDER_SECRET_KEY); v == "" {
			t.Fatalf("%v must be set for acceptance tests\n", PROVIDER_SECRET_KEY)
		}
		os.Setenv(PROVIDER_REGION, defaultRegionSes)
	default:
		testAccPreCheck(t)
	}
}
