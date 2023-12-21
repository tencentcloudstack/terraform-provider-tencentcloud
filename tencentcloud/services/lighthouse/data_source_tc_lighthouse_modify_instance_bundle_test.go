package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseModifyInstanceBundleDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseModifyInstanceBundleDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_modify_instance_bundle.modify_instance_bundle")),
			},
		},
	})
}

const testAccLighthouseModifyInstanceBundleDataSource = tcacctest.DefaultLighthoustVariables + `

data "tencentcloud_lighthouse_modify_instance_bundle" "modify_instance_bundle" {
  instance_id = var.lighthouse_id
  filters {
	name = "bundle-id"
	values = ["bundle_gen_mc_med2_02"]

  }
}

`
