package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresDBInstanceVersionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresDBInstanceVersionsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgres_d_b_instance_versions.d_b_instance_versions")),
			},
		},
	})
}

const testAccPostgresDBInstanceVersionsDataSource = `

data "tencentcloud_postgres_d_b_instance_versions" "d_b_instance_versions" {
  }

`
