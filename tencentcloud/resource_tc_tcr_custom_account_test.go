package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTcrCustomAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrCustomAccount,
				PreConfig: func() {
					testAccStepSetRegion(t, "ap-shanghai")
					testAccPreCheckCommon(t, ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_custom_account.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_custom_account.example", "registry_id"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_custom_account.example", "name", "tf_example_account"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_custom_account.example", "permissions.#"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_custom_account.example", "permissions.0.resource", "tf_test_tcr_namespace"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_custom_account.example", "permissions.0.actions.#"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_tcr_custom_account.example", "permissions.0.actions.*", "tcr:PushRepository"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_tcr_custom_account.example", "permissions.0.actions.*", "tcr:PullRepository"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_custom_account.example", "description", "tf example for tcr custom account"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_custom_account.example", "duration", "10"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_custom_account.example", "disable", "false"),
				),
			},
			{
				ResourceName:      "tencentcloud_tcr_custom_account.custom_account",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcrCustomAccount = `

resource "tencentcloud_tcr_instance" "example" {
  name          = "tf-example-tcr-instance"
  instance_type = "basic"
  delete_bucket = true
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tcr_namespace" "example" {
  instance_id    = tencentcloud_tcr_instance.example.id
  name           = "tf_test_tcr_namespace"
  is_public      = true
  is_auto_scan   = true
  is_prevent_vul = true
  severity       = "medium"
  cve_whitelist_items {
    cve_id = "tf_example_cve_id"
  }
}

resource "tencentcloud_tcr_custom_account" "example" {
  registry_id = tencentcloud_tcr_instance.example.id
  name        = "tf_example_account"
  permissions {
    resource = tencentcloud_tcr_namespace.example.name
    actions  = ["tcr:PushRepository", "tcr:PullRepository"]
  }
  description = "tf example for tcr custom account"
  duration    = 10
  disable     = false
  tags = {
    "createdBy" = "terraform"
  }
}

`
