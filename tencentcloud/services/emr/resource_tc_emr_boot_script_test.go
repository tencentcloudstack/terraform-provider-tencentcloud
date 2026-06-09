package emr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudEmrBootScriptResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEmrBootScript,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_emr_boot_script.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_emr_boot_script.example", "boot_type", "resourceAfter"),
				),
			},
			{
				ResourceName:            "tencentcloud_emr_boot_script.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pre_executed_file_settings.0.cos_secret_key"},
			},
		},
	})
}

const testAccEmrBootScript = `
resource "tencentcloud_emr_boot_script" "example" {
  instance_id = "emr-qe336v2e"
  boot_type   = "resourceAfter"
  pre_executed_file_settings {
    path          = "test.py"
    bucket        = "test123-1309118522"
    cos_file_name = "test"
    region        = "ap-guangzhou"
  }
}
`
