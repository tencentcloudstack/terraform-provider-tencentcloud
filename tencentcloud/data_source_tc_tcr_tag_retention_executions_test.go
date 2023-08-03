package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTcrTagRetentionExecutionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrTagRetentionExecutionsDataSource, 
				PreConfig: func() {
					// testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tcr_tag_retention_executions.tag_retention_executions"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_tag_retention_executions.tag_retention_executions", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_tag_retention_executions.tag_retention_executions", "registry_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_tcr_tag_retention_executions.tag_retention_executions", "retention_id", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_tag_retention_executions.tag_retention_executions", "retention_execution_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_tag_retention_executions.tag_retention_executions", "retention_execution_list.0.execution_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_tag_retention_executions.tag_retention_executions", "retention_execution_list.0.retention_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_tag_retention_executions.tag_retention_executions", "retention_execution_list.0.start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_tag_retention_executions.tag_retention_executions", "retention_execution_list.0.end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_tag_retention_executions.tag_retention_executions", "retention_execution_list.0.status"),
				),
			},
		},
	})
}

const testAccTcrTagRetentionExecutionsDataSource = TCRDataSource+`

data "tencentcloud_tcr_tag_retention_executions" "tag_retention_executions" {
  registry_id = local.tcr_id
  retention_id = 1
}

`
