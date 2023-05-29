package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlDrInstanceToMaterResource_basic -v
func TestAccTencentCloudNeedFixMysqlDrInstanceToMaterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlDrInstanceToMater,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_dr_instance_to_mater.dr_instance_to_mater", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_dr_instance_to_mater.dr_instance_to_mater",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlDrInstanceToMater = `

resource "tencentcloud_mysql_dr_instance_to_mater" "dr_instance_to_mater" {
  instance_id = ""
}

`
