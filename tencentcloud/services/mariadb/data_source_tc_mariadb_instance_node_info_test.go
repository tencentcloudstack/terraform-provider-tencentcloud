package mariadb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbInstanceNodeInfoDataSource_basic -v
func TestAccTencentCloudMariadbInstanceNodeInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbInstanceNodeInfoDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_instance_node_info.instance_node_info"),
				),
			},
		},
	})
}

const testAccMariadbInstanceNodeInfoDataSource = `
data "tencentcloud_mariadb_instance_node_info" "instance_node_info" {
  instance_id = "tdsql-9vqvls95"
}
`
