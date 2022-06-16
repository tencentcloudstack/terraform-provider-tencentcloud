package tencentcloud

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

const (
	ACCOUNT_TYPE_INTERNATION        = "INTERNATION"
	ACCOUNT_TYPE_PREPAY             = "PREPAY"
	ACCOUNT_TYPE_COMMON             = "COMMON"
	ACCOUNT_TYPE_PRIVATE            = "PRIVATE"
	ACCOUNT_TYPE_SUB_ACCOUNT        = "SUB_ACCOUNT"
	INTERNATION_PROVIDER_SECRET_ID  = "TENCENTCLOUD_SECRET_ID_INTERNATION"
	INTERNATION_PROVIDER_SECRET_KEY = "TENCENTCLOUD_SECRET_KEY_INTERNATION"
	PREPAY_PROVIDER_SECRET_ID       = "TENCENTCLOUD_SECRET_ID_PREPAY"
	PREPAY_PROVIDER_SECRET_KEY      = "TENCENTCLOUD_SECRET_KEY_PREPAY"
	PRIVATE_PROVIDER_SECRET_ID      = "TENCENTCLOUD_SECRET_ID_PRIVATE"
	PRIVATE_PROVIDER_SECRET_KEY     = "TENCENTCLOUD_SECRET_KEY_PRIVATE"
	COMMON_PROVIDER_SECRET_ID       = "TENCENTCLOUD_SECRET_ID_COMMON"
	COMMON_PROVIDER_SECRET_KEY      = "TENCENTCLOUD_SECRET_KEY_COMMON"
	SUB_ACCOUNT_PROVIDER_SECRET_ID  = "TENCENTCLOUD_SECRET_ID_SUB_ACCOUNT"
	SUB_ACCOUNT_PROVIDER_SECRET_KEY = "TENCENTCLOUD_SECRET_KEY_SUB_ACCOUNT"
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
	case accountType == ACCOUNT_TYPE_INTERNATION:
		secretId := os.Getenv(INTERNATION_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(INTERNATION_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", INTERNATION_PROVIDER_SECRET_ID, INTERNATION_PROVIDER_SECRET_KEY)
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
