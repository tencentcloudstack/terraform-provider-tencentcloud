package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTeoDnsSec_basic -v
func TestAccTencentCloudTeoDnsSec_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PRIVATE) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoDnsSec,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDnsSecExists("tencentcloud_teo_dns_sec.basic"),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_sec.basic", "zone_id", defaultZoneId),
					resource.TestCheckResourceAttr("tencentcloud_teo_dns_sec.basic", "status", "disabled"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_dns_sec.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDnsSecExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := TeoService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		agents, err := service.DescribeTeoDnsSec(ctx, rs.Primary.ID)
		if agents == nil {
			return fmt.Errorf("zone DnsSec %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTeoDnsSecVar = `
variable "zone_id" {
  default = "` + defaultZoneId + `"
}
`

const testAccTeoDnsSec = testAccTeoDnsSecVar + `

resource "tencentcloud_teo_dns_sec" "basic" {
  zone_id = var.zone_id
  status  = "disabled"
}
`
