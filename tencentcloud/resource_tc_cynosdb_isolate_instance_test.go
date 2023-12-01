package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbIsolateInstanceResource_basic -v
func TestAccTencentCloudCynosdbIsolateInstanceResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbIsolateInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_isolate_instance.isolate_instance", "id"),
				),
			},
			{
				Config: testAccCynosdbIsolateInstanceUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_isolate_instance.isolate_instance", "id"),
				),
			},
		},
	})
}

const testAccCynosdbIsolateInstance = CommonCynosdb + `

resource "tencentcloud_cynosdb_isolate_instance" "isolate_instance" {
	cluster_id = var.cynosdb_cluster_id
	instance_id = "cynosdbmysql-ins-rikr6z4o"
	operate = "isolate"
}

`

const testAccCynosdbIsolateInstanceUp = CommonCynosdb + `

resource "tencentcloud_cynosdb_isolate_instance" "isolate_instance" {
	cluster_id = var.cynosdb_cluster_id
	instance_id = "cynosdbmysql-ins-rikr6z4o"
	operate = "activate"
}

`
