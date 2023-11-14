package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudAntiddosBoundipResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosBoundip,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_antiddos_boundip.boundip", "id")),
			},
			{
				ResourceName:      "tencentcloud_antiddos_boundip.boundip",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAntiddosBoundip = `

resource "tencentcloud_antiddos_boundip" "boundip" {
  business = "bgp-multip"
  id = "bgp-000000xe"
  bound_dev_list {
		ip = "1.1.1.1"
		biz_type = "public"
		instance_id = "ins-xxx"
		device_type = "cvm"
		isp_code = 5

  }
  un_bound_dev_list {
		ip = "1.1.1.1"
		biz_type = "public"
		instance_id = "ins-xxx"
		device_type = "cvm"
		isp_code = 5

  }
  copy_policy = ""
}

`
