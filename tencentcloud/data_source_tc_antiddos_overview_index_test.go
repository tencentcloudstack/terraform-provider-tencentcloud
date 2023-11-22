package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosOverviewIndexDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosOverviewIndexDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_antiddos_overview_index.overview_index"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_overview_index.overview_index", "all_ip_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_overview_index.overview_index", "antiddos_ip_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_overview_index.overview_index", "attack_ip_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_overview_index.overview_index", "block_ip_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_overview_index.overview_index", "antiddos_domain_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_overview_index.overview_index", "attack_domain_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_overview_index.overview_index", "max_attack_flow"),
				),
			},
		},
	})
}

const testAccAntiddosOverviewIndexDataSource = `
data "tencentcloud_antiddos_overview_index" "overview_index" {
	start_time = "2023-11-20 12:32:12"
	end_time = "2023-11-21 12:32:12"
}
`
