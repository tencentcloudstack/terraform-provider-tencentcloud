package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbIsolateInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbIsolateInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_isolate_instance.isolate_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_isolate_instance.isolate_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbIsolateInstance = `

resource "tencentcloud_cdb_isolate_instance" "isolate_instance" {
  instance_id = ""
}

`
