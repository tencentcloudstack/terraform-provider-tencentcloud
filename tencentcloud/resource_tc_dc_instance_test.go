package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dc_instance.instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_dc_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcInstance = `

resource "tencentcloud_dc_instance" "instance" {
  direct_connect_name = ""
  access_point_id = ""
  line_operator = ""
  port_type = ""
  circuit_code = ""
  location = ""
  bandwidth = 
  redundant_direct_connect_id = ""
  vlan = 
  tencent_address = ""
  customer_address = ""
  customer_name = ""
  customer_contact_mail = ""
  customer_contact_number = ""
  fault_report_contact_person = ""
  fault_report_contact_number = ""
  sign_law = 
}

`
