package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	lb "github.com/zqfan/tencentcloud-sdk-go/services/lb/unversioned"
)

func TestAccTencentCloudLB_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("tencentcloud_lb.classic"),
					resource.TestCheckResourceAttr("tencentcloud_lb.classic", "name", "tf-ci-test"),
					resource.TestCheckResourceAttrSet("tencentcloud_lb.classic", "status"),
				),
			},
		},
	})
}

func testAccCheckLBDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*TencentCloudClient).lbConn
	var lbid string
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_lb" {
			continue
		}
		lbid = rs.Primary.ID
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		req := lb.NewDescribeLoadBalancersRequest()
		req.LoadBalancerIds = []*string{&lbid}
		resp, err := client.DescribeLoadBalancers(req)
		if err != nil {
			return resource.RetryableError(err)
		}
		if *resp.TotalCount != 0 {
			return resource.RetryableError(fmt.Errorf("lb can still be found after deleted"))
		}
		return nil
	})
}

const testAccLBConfig = `
resource "tencentcloud_lb" "classic" {
  type = 2
  forward = 0
  name = "tf-ci-test"
}
`
