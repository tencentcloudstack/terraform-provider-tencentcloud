package cdwpg_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCdwpgUpgradeInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdwpgUpgradeInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdwpg_upgrade_instance.cdwpg_upgrade_instance", "id")),
			},
		},
	})
}

const testAccCdwpgUpgradeInstance = `
resource "tencentcloud_cdwpg_upgrade_instance" "cdwpg_upgrade_instance" {
	instance_id = "cdwpg-zpiemnyd"
	package_version = "3.16.9.4"
}
`
