package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverStartXeventResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverStartXevent,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_start_xevent.start_xevent", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_start_xevent.start_xevent",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverStartXevent = `

resource "tencentcloud_sqlserver_start_xevent" "start_xevent" {
  instance_id = ""
  event_config {
		event_type = ""
		threshold = 

  }
}

`
