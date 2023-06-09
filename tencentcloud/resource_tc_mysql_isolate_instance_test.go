package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlIsolateInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlIsolateInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_isolate_instance.isolate_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_isolate_instance.isolate_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlIsolateInstance = `

resource "tencentcloud_mysql_isolate_instance" "isolate_instance" {
  instance_id = ""
}

`
