package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCynosdbInstancesDataSource_full(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbInstancestDataSource_full(),
				Check: resource.ComposeAggregateTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("tencentcloud_cynosdb_readonly_instance.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_instances.instances", "instance_list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_instances.instances", "instance_list.0.instance_name", "tf-cynosdb-readonly-instance"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instances.instances", "instance_list.0.cluster_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instances.instances", "instance_list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instances.instances", "instance_list.0.instance_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instances.instances", "instance_list.0.instance_storage_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instances.instances", "instance_list.0.instance_memory_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instances.instances", "instance_list.0.instance_cpu_core"),
				),
			},
		},
	})
}

func testAccCynosdbInstancestDataSource_full() string {
	return testAccCynosdbReadonlyInstance + `
data "tencentcloud_cynosdb_instances" "instances" {
  instance_id   = tencentcloud_cynosdb_readonly_instance.foo.id
  project_id    = 0
  db_type       = "MYSQL"
  cluster_id    = tencentcloud_cynosdb_readonly_instance.foo.cluster_id
  instance_name = "tf-cynosdb-readonly-instance"
}`
}
