package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosPendingRiskInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosPendingRiskInfoDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_antiddos_pending_risk_info.pending_risk_info"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_pending_risk_info.pending_risk_info", "attacking_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_pending_risk_info.pending_risk_info", "blocking_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_pending_risk_info.pending_risk_info", "expired_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_pending_risk_info.pending_risk_info", "total"),
				),
			},
		},
	})
}

const testAccAntiddosPendingRiskInfoDataSource = `
data "tencentcloud_antiddos_pending_risk_info" "pending_risk_info" {
}
`
