package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoOriginGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoOriginGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_origin_group.origin_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_origin_group.origin_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoOriginGroup = `

resource "tencentcloud_teo_origin_group" "origin_group" {
  zone_id = &lt;nil&gt;
  origin_group_id = &lt;nil&gt;
  origin_group_name = &lt;nil&gt;
  origin_type = &lt;nil&gt;
  configuration_type = &lt;nil&gt;
  origin_records {
		record = &lt;nil&gt;
		port = &lt;nil&gt;
		weight = &lt;nil&gt;
		area = &lt;nil&gt;
		private = &lt;nil&gt;
		private_parameter {
			name = &lt;nil&gt;
			value = &lt;nil&gt;
		}

  }
  }

`
