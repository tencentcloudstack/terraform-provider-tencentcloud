package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbFlushBinlogResource_basic -v
func TestAccTencentCloudMariadbFlushBinlogResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbFlushBinlog,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_flush_binlog.flush_binlog", "id"),
				),
			},
		},
	})
}

const testAccMariadbFlushBinlog = `
resource "tencentcloud_mariadb_flush_binlog" "flush_binlog" {
  instance_id = "tdsql-9vqvls95"
}
`
