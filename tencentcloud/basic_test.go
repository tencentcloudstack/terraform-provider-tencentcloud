package tencentcloud

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const appid string = "1259649581"

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv(PROVIDER_SECRET_ID); v == "" {
		t.Fatalf("%v must be set for acceptance tests\n", PROVIDER_SECRET_ID)
	}
	if v := os.Getenv(PROVIDER_SECRET_KEY); v == "" {
		t.Fatalf("%v must be set for acceptance tests\n", PROVIDER_SECRET_KEY)
	}
	if v := os.Getenv(PROVIDER_REGION); v == "" {
		log.Println("[INFO] Test: Using ap-guangzhou as test region")
		os.Setenv(PROVIDER_REGION, "ap-guangzhou")
	}
}

func testAccPreSetRegion(region string) {
	os.Setenv(PROVIDER_REGION, region)
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
