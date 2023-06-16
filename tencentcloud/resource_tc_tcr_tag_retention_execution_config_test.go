package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTcrTagRetentionExecutionConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Config: fmt.Sprintf(testAccTcrTagRetentionExecutionConfig, defaultTCRInstanceId),
				Config: testAccTcrTagRetentionRule_manual,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_execution_config.config", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_execution_config.config", "registry_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_execution_config.config", "retention_id"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_tag_retention_execution_config.config", "dry_run", "false"),
					// computed
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_tag_retention_execution_config.config", "execution_id"),
				),
			},
		},
	})
}

const testAccTcrTagRetentionRule_manual = testAccTCRInstance_retention + `

resource "tencentcloud_tcr_namespace" "my_ns" {
  instance_id    = tencentcloud_tcr_instance.mytcr_retention.id
  name           = "tf_test_ns_retention"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "cve-xxxxx"
  }
}

resource "tencentcloud_tcr_tag_retention_rule" "my_rule" {
  registry_id    = tencentcloud_tcr_instance.mytcr_retention.id
  namespace_name = tencentcloud_tcr_namespace.my_ns.name
  retention_rule {
    key   = "nDaysSinceLastPush"
    value = 2
  }
  cron_setting = "weekly"
  disabled     = true
}

resource "tencentcloud_tcr_tag_retention_execution_config" "config" {
  registry_id  = tencentcloud_tcr_tag_retention_rule.my_rule.registry_id
  retention_id = tencentcloud_tcr_tag_retention_rule.my_rule.retention_id
  dry_run      = false
}

`

// const testAccTcrTagRetentionExecutionConfig = `

// resource "tencentcloud_tcr_tag_retention_execution_config" "config" {
//   registry_id = "%s"
//   retention_id = 1
//   dry_run = false
// }

// `
