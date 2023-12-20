package clb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localclb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/clb"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudLB_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckLBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbBasic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("tencentcloud_lb.classic"),
					resource.TestCheckResourceAttr("tencentcloud_lb.classic", "name", "tf-ci-test"),
					resource.TestCheckResourceAttr("tencentcloud_lb.classic", "type", "OPEN"),
					resource.TestCheckResourceAttrSet("tencentcloud_lb.classic", "forward"),
					resource.TestCheckResourceAttrSet("tencentcloud_lb.classic", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_lb.classic", "project_id"),
				),
			},
			{
				Config: testAccLbBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("tencentcloud_lb.classic"),
					resource.TestCheckResourceAttr("tencentcloud_lb.classic", "name", "tf-ci-test-update"),
					resource.TestCheckResourceAttr("tencentcloud_lb.classic", "type", "OPEN"),
					resource.TestCheckResourceAttrSet("tencentcloud_lb.classic", "forward"),
					resource.TestCheckResourceAttrSet("tencentcloud_lb.classic", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_lb.classic", "project_id"),
				),
			},
		},
	})
}

func testAccCheckLBDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clbService := localclb.NewClbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_lb" {
			continue
		}

		_, err := clbService.DescribeLoadBalancerById(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("clb instance still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

const testAccLbBasic = `
resource "tencentcloud_lb" "classic" {
  type    = "OPEN"
  forward = "APPLICATION"
  name    = "tf-ci-test"
}
`

const testAccLbBasicUpdate = `
resource "tencentcloud_lb" "classic" {
  type    = "OPEN"
  forward = "APPLICATION"
  name    = "tf-ci-test-update"
}
`
