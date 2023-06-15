package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
	"fmt"
)

func TestAccTencentCloudTcrTagRetentionExecutionConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTcrTagRetentionExecutionConfig, defaultTCRInstanceId),
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_execution_config.tag_retention_execution_config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_execution_config.tag_retention_execution_config", "registry_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_execution_config.tag_retention_execution_config", "retention_id"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_tag_retention_execution_config.tag_retention_execution_config", "dry_run", "false"),
					// computed
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_execution_config.tag_retention_execution_config", "execution_id"),
				),
			},
		},
	})
}

const testAccTcrTagRetentionExecutionConfig = `

resource "tencentcloud_tcr_tag_retention_execution_config" "tag_retention_execution_config" {
  registry_id = "%s"
  retention_id = 1
  dry_run = false
}

`
