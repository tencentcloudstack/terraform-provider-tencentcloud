package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudWafInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_waf_instance.instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_waf_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafInstance = `

resource "tencentcloud_waf_instance" "instance" {
  goods {
		goods_num = 
		goods_detail {
			time_span = 
			time_unit = ""
			sub_product_code = ""
			pid = 
			instance_name = ""
			auto_renew_flag = 
			real_region = 
			label_types = 
			label_counts = 
			cur_deadline = ""
			instance_id = ""
		}
		goods_category_id = 
		region_id = 

  }
}

`
