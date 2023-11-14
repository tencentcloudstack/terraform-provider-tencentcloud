package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcrReplicationInstanceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrReplicationInstanceDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tcr_replication_instance.replication_instance")),
			},
		},
	})
}

const testAccTcrReplicationInstanceDataSource = `

data "tencentcloud_tcr_replication_instance" "replication_instance" {
  registry_id = "tcr-xx"
  replication_registry_id = "tcr-xx-1"
  replication_region_id = 1
  show_replication_log = false
        tags = {
    "createdBy" = "terraform"
  }
}

`
