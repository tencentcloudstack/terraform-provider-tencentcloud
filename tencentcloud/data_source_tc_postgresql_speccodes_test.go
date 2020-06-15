package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudPostgresqlSpeccodes_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudPostgresqlSpeccodesConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_postgresql_speccodes.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_speccodes.foo", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_speccodes.foo", "list.0.id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_speccodes.foo", "list.0.memory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_speccodes.foo", "list.0.storage_min"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_speccodes.foo", "list.0.storage_max"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_speccodes.foo", "list.0.cpu_number"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_speccodes.foo", "list.0.qps"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_speccodes.foo", "list.0.version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_postgresql_speccodes.foo", "list.0.version_name"),
				),
			},
		},
	})
}

const testAccTencentCloudPostgresqlSpeccodesConfigBasic = `
data "tencentcloud_postgresql_speccodes" "foo"{
	availability_zone = "ap-guangzhou-3"
}
`
