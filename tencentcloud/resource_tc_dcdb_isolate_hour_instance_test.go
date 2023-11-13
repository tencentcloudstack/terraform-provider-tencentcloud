package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbIsolateHourInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbIsolateHourInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_isolate_hour_instance.isolate_hour_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_isolate_hour_instance.isolate_hour_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbIsolateHourInstance = `

resource "tencentcloud_dcdb_isolate_hour_instance" "isolate_hour_instance" {
  instance_ids = 
}

`
