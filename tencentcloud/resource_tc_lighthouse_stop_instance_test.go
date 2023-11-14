package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseStopInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseStopInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_stop_instance.stop_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_stop_instance.stop_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseStopInstance = `

resource "tencentcloud_lighthouse_stop_instance" "stop_instance" {
  instance_ids = 
}

`
