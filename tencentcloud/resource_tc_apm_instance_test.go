package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
				Config: testAccApmInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_apm_instance.instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_apm_instance.instance", "name", "terraform-for-test"),
				),
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
  name = "terraform-test"
  description = "for terraform test"
  trace_duration = 15
  span_daily_counters = 20
}

`

const testAccApmInstanceUpdate = `

resource "tencentcloud_apm_instance" "instance" {
  name = "terraform-for-test"
  description = "for terraform test"
  trace_duration = 15
  span_daily_counters = 20
}

`
