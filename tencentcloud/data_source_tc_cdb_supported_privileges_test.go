package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbSupportedPrivilegesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbSupportedPrivilegesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cdb_supported_privileges.supported_privileges")),
			},
		},
	})
}

const testAccCdbSupportedPrivilegesDataSource = `

data "tencentcloud_cdb_supported_privileges" "supported_privileges" {
  instance_id = ""
        }

`
