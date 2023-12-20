package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlIsolateInstanceResource_basic -v
func TestAccTencentCloudMysqlIsolateInstanceResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlIsolateInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_isolate_instance.isolate_instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_isolate_instance.isolate_instance", "status", "5"),
				),
			},
			{
				Config: testAccMysqlIsolateInstanceUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_isolate_instance.isolate_instance", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_isolate_instance.isolate_instance", "status", "1"),
				),
			},
		},
	})
}

const testAccMysqlIsolateInstance = `

resource "tencentcloud_mysql_isolate_instance" "isolate_instance" {
	instance_id = "cdb-fitq5t9h"
	operate     = "isolate"
}

`

const testAccMysqlIsolateInstanceUp = `

resource "tencentcloud_mysql_isolate_instance" "isolate_instance" {
	instance_id = "cdb-fitq5t9h"
	operate     = "recover"
}

`
