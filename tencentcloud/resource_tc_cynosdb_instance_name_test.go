package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbInstanceNameResource_basic -v
func TestAccTencentCloudCynosdbInstanceNameResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbInstanceName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_instance_name.instance_name", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_instance_name.instance_name", "instance_id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_instance_name.instance_name", "instance_name", "tf-cynosdb-instance-name"),
				),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_instance_name.instance_name",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbInstanceName = `

resource "tencentcloud_cynosdb_instance_name" "instance_name" {
	instance_id = "cynosdbmysql-ins-afqx1hy0"
	instance_name = "tf-cynosdb-instance-name"
}

`
