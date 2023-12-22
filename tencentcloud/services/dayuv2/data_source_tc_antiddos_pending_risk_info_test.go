package dayuv2_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosPendingRiskInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosPendingRiskInfoDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_antiddos_pending_risk_info.pending_risk_info"),
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
