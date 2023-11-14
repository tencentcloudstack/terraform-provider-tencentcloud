package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainNoPrimaryKeyTablesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainNoPrimaryKeyTablesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_no_primary_key_tables.no_primary_key_tables")),
			},
		},
	})
}

const testAccDbbrainNoPrimaryKeyTablesDataSource = `

data "tencentcloud_dbbrain_no_primary_key_tables" "no_primary_key_tables" {
  instance_id = ""
  date = ""
  product = ""
          }

`
