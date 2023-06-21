package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbRestartInstanceResource_basic -v
func TestAccTencentCloudCynosdbRestartInstanceResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbRestartInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_restart_instance.restart_instance", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_restart_instance.restart_instance", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_restart_instance.restart_instance", "status"),
				),
			},
		},
	})
}

const testAccCynosdbRestartInstance = CommonCynosdb + `

resource "tencentcloud_cynosdb_restart_instance" "restart_instance" {
	instance_id = var.cynosdb_cluster_instance_id
}

`
