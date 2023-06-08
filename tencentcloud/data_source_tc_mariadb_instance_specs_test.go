package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbInstanceSpecsDataSource_basic -v
func TestAccTencentCloudMariadbInstanceSpecsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbInstanceSpecsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_instance_specs.instance_specs"),
				),
			},
		},
	})
}

const testAccMariadbInstanceSpecsDataSource = `
data "tencentcloud_mariadb_instance_specs" "instance_specs" {
}
`
