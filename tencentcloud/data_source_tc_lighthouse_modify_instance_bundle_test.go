package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseModifyInstanceBundleDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseModifyInstanceBundleDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_modify_instance_bundle.modify_instance_bundle")),
			},
		},
	})
}

const testAccLighthouseModifyInstanceBundleDataSource = DefaultLighthoustVariables + `

data "tencentcloud_lighthouse_modify_instance_bundle" "modify_instance_bundle" {
  instance_id = var.lighthouse_id
  filters {
	name = "bundle-id"
	values = ["bundle_gen_mc_med2_02"]

  }
}

`
