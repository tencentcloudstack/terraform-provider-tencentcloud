package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbDbInstancesDataSource -v
func TestAccTencentCloudMariadbDbInstancesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMariadbDbInstances,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_db_instances.db_instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_db_instances.db_instances", "instances.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_db_instances.db_instances", "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_db_instances.db_instances", "instances.0.instance_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_db_instances.db_instances", "instances.0.project_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_db_instances.db_instances", "instances.0.region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_db_instances.db_instances", "instances.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_db_instances.db_instances", "instances.0.memory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_db_instances.db_instances", "instances.0.storage"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_db_instances.db_instances", "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_db_instances.db_instances", "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_db_instances.db_instances", "instances.0.db_version_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_db_instances.db_instances", "instances.0.vip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mariadb_db_instances.db_instances", "instances.0.vport"),
				),
			},
		},
	})
}

const testAccDataSourceMariadbDbInstances = testAccMariadbInstance + `

data "tencentcloud_mariadb_db_instances" "db_instances" {
	depends_on = [ tencentcloud_mariadb_instance.instance ]
}

`
