package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbIsolateInstanceResource_basic -v
func TestAccTencentCloudCynosdbIsolateInstanceResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
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

const testAccCynosdbIsolateInstance = tcacctest.CommonCynosdb + `

resource "tencentcloud_cynosdb_isolate_instance" "isolate_instance" {
	cluster_id = var.cynosdb_cluster_id
	instance_id = "cynosdbmysql-ins-rikr6z4o"
	operate = "isolate"
}

`

const testAccCynosdbIsolateInstanceUp = tcacctest.CommonCynosdb + `

resource "tencentcloud_cynosdb_isolate_instance" "isolate_instance" {
	cluster_id = var.cynosdb_cluster_id
	instance_id = "cynosdbmysql-ins-rikr6z4o"
	operate = "activate"
}

`
