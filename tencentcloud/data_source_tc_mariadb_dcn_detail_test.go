package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbDcnDetailDataSource_basic -v
func TestAccTencentCloudMariadbDcnDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbDcnDetailDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_dcn_detail.dcn_detail"),
				),
			},
		},
	})
}

const testAccMariadbDcnDetailDataSource = `
data "tencentcloud_mariadb_dcn_detail" "dcn_detail" {
  instance_id = "tdsql-9vqvls95"
}
`
