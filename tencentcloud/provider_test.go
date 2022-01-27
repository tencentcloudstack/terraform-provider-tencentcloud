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
}

func testAccPreCheckCommon(t *testing.T, accountType string) {
	if v := os.Getenv(PROVIDER_REGION); v == "" {
		log.Printf("[INFO] Testing: Using %s as test region", defaultRegion)
		os.Setenv(PROVIDER_REGION, defaultRegion)
	}
	switch {
	case accountType == ACCTUNT_TYPE_INTERNATION:
		secretId := os.Getenv(INTERNATION_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(INTERNATION_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", INTERNATION_PROVIDER_SECRET_ID, INTERNATION_PROVIDER_SECRET_KEY)
		}
		os.Setenv(PROVIDER_SECRET_ID, secretId)
		os.Setenv(PROVIDER_SECRET_KEY, secretKey)
	case accountType == ACCTUNT_TYPE_PREPAY:
		secretId := os.Getenv(PREPAY_PROVIDER_SECRET_ID)
		secretKey := os.Getenv(PREPAY_PROVIDER_SECRET_KEY)
		if secretId == "" || secretKey == "" {
			t.Fatalf("%v and %v must be set for acceptance tests\n", PREPAY_PROVIDER_SECRET_ID, PREPAY_PROVIDER_SECRET_KEY)
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
