package tencentcloud

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

/*
If you want to run through the test case,
the following must be changed to your resource id
*/
const appid string = "1259649581"

// 172.16.0.0/16
const DefaultVpcId = "vpc-h70b6b49"

const DefaultSubnetId = "subnet-1uwh63so"

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
