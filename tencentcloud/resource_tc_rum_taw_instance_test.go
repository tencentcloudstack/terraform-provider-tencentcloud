package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRumTawInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumTawInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_rum_taw_instance.taw_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_rum_taw_instance.taw_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccRumTawInstance = `

resource "tencentcloud_rum_taw_instance" "taw_instance" {
  area_id = &lt;nil&gt;
  charge_type = &lt;nil&gt;
  data_retention_days = &lt;nil&gt;
  instance_name = &lt;nil&gt;
  tags {
		key = &lt;nil&gt;
		value = &lt;nil&gt;

  }
  instance_desc = &lt;nil&gt;
  count_num = &lt;nil&gt;
  period_retain = &lt;nil&gt;
  buying_channel = &lt;nil&gt;
            tags = {
    "createdBy" = "terraform"
  }
}

`
