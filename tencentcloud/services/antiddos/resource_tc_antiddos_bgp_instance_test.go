package antiddos_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosBgpInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosBgpInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_antiddos_bgp_instance.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_antiddos_bgp_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAntiddosBgpInstance = `
resource "tencentcloud_antiddos_bgp_instance" "example" {
  instance_charge_type = "POSTPAID_BY_MONTH"
  package_type         = "Standard"
  standard_package_config {
    region                 = "ap-guangzhou"
    protect_ip_count       = 1
    bandwidth              = 100
    elastic_bandwidth_flag = true
  }

  tag_info_list {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
`
