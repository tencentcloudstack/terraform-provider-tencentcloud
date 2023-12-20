package dayuv2_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataDayuEip = "data.tencentcloud_dayu_eip.test"

func TestAccTencentCloudDataDayuEip(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_INTERNATIONAL) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDayuEipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataDayuEip,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDayuEipExists("tencentcloud_dayu_eip.test_eip"),
					resource.TestCheckResourceAttr(testDataDayuEip, "list.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudDataDayuEip = `
resource "tencentcloud_dayu_eip" "test_eip" {
	resource_id = "bgpip-000004xg"
	eip = "162.62.163.50"
	bind_resource_id = "ins-4m0jvxic"
	bind_resource_region = "hk"
	bind_resource_type = "cvm"
  }

data "tencentcloud_dayu_eip" "test" {
	resource_id = tencentcloud_dayu_eip.test_eip.resource_id
}
`
