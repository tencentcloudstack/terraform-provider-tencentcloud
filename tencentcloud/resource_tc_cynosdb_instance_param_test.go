package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbInstanceParamResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbInstanceParam,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_instance_param.instance_param", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_instance_param.instance_param",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbInstanceParam = `

resource "tencentcloud_cynosdb_instance_param" "instance_param" {
  cluster_id = ""
  instance_ids = 
  cluster_param_list {
		param_name = ""
		current_value = ""
		old_value = ""

  }
  instance_param_list {
		param_name = ""
		current_value = ""
		old_value = ""

  }
  is_in_maintain_period = ""
}

`
