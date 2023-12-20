package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixMysqlRoStopReplicationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlRoStopReplication,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_stop_replication.ro_stop_replication", "id")),
			},
		},
	})
}

const testAccMysqlRoStopReplication = `

resource "tencentcloud_mysql_ro_stop_replication" "ro_stop_replication" {
  instance_id = "cdb-fitq5t9h"
}

`
