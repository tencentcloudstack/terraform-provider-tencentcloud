package cls_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"

	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudClsDashboard_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClsDashboardDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClsDashboard_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsDashboardExists("tencentcloud_cls_dashboard.dashboard"),
					resource.TestCheckResourceAttr("tencentcloud_cls_dashboard.dashboard", "dashboard_name", "tf-dashboard-test"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_dashboard.dashboard", "dashboard_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_dashboard.dashboard", "create_time"),
				),
			},
			{
				Config: testAccClsDashboard_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsDashboardExists("tencentcloud_cls_dashboard.dashboard"),
					resource.TestCheckResourceAttr("tencentcloud_cls_dashboard.dashboard", "dashboard_name", "tf-dashboard-test-updated"),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_dashboard.dashboard",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClsDashboardDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := localcls.NewClsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cls_dashboard" {
			continue
		}

		dashboard, err := service.DescribeClsDashboardById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if dashboard != nil {
			return fmt.Errorf("cls dashboard still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClsDashboardExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cls dashboard %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cls dashboard id is not set")
		}

		service := localcls.NewClsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		dashboard, err := service.DescribeClsDashboardById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if dashboard == nil {
			return fmt.Errorf("cls dashboard not found: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClsDashboard_basic = `
resource "tencentcloud_cls_dashboard" "dashboard" {
  dashboard_name = "tf-dashboard-test"
}
`

const testAccClsDashboard_update = `
resource "tencentcloud_cls_dashboard" "dashboard" {
  dashboard_name = "tf-dashboard-test-updated"
  data           = jsonencode({
    timezone = "browser"
  })
}
`
