package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbInstanceParamResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbInstanceParam,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_instance_param.instance_param", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_instance_param.instance_param",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbInstanceParam = `

resource "tencentcloud_cdb_instance_param" "instance_param" {
  instance_ids = &lt;nil&gt;
  param_list {
		name = &lt;nil&gt;
		current_value = &lt;nil&gt;

  }
  template_id = &lt;nil&gt;
  wait_switch = 0
  not_sync_ro = false
  not_sync_dr = false
}

`
