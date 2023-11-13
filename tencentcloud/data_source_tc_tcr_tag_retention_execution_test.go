package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcrTagRetentionExecutionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrTagRetentionExecutionDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tcr_tag_retention_execution.tag_retention_execution")),
			},
		},
	})
}

const testAccTcrTagRetentionExecutionDataSource = `

data "tencentcloud_tcr_tag_retention_execution" "tag_retention_execution" {
  registry_id = "tcr-xx"
  retention_id = 1
  }

`
