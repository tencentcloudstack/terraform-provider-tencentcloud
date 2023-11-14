package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcrTagRetentionExecutionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrTagRetentionExecution,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_execution.tag_retention_execution", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcr_tag_retention_execution.tag_retention_execution",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcrTagRetentionExecution = `

resource "tencentcloud_tcr_tag_retention_execution" "tag_retention_execution" {
  registry_id = "tcr-xx"
  retention_id = 1
  dry_run = false
}

`
