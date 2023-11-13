package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudApmInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApmInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_apm_instance.instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_apm_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApmInstance = `

resource "tencentcloud_apm_instance" "instance" {
  name = ""
  description = ""
  trace_duration = 
  span_daily_counters = 
  tags = {
    "createdBy" = "terraform"
  }
}

`
