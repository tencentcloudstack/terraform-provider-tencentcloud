package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverStartXeventResource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_start_xevent.start_xevent", "id"),
				),
			},
		},
	})
}

const testAccSqlserverStartXevent = `
resource "tencentcloud_sqlserver_start_xevent" "start_xevent" {
  instance_id = "mssql-gyg9xycl"
  event_config {
    event_type = "blocked"
    threshold  = 0
  }
}
`
