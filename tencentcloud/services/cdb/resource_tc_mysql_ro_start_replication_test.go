package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlRoStartReplicationResource_basic -v
func TestAccTencentCloudNeedFixMysqlRoStartReplicationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlRoStartReplication,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_start_replication.ro_start_replication", "id"),
				),
			},
		},
	})
}

const testAccMysqlRoStartReplication = `

resource "tencentcloud_mysql_ro_start_replication" "ro_start_replication" {
  instance_id = "cdb-fitq5t9h"
}

`
