package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsWorkflowConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsWorkflowConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_workflow_config.workflow_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_workflow_config.workflow_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsWorkflowConfig = `

resource "tencentcloud_mps_workflow_config" "workflow_config" {
  workflow_id = 
}

`
