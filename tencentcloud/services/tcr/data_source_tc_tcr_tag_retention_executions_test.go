package tcr_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTcrTagRetentionExecutionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTcrTagRetentionExecutionsDataSource, tcacctest.DefaultTCRInstanceId),
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tcr_tag_retention_executions.tag_retention_executions"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tcr_tag_retention_executions.tag_retention_executions", "id"),
					resource.TestCheckResourceAttr("data.tencentcloud_tcr_tag_retention_executions.tag_retention_executions", "registry_id", tcacctest.DefaultTCRInstanceId),
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

const testAccTcrTagRetentionExecutionsDataSource = `

data "tencentcloud_tcr_tag_retention_executions" "tag_retention_executions" {
  registry_id = "%s"
  retention_id = 1
}

`
