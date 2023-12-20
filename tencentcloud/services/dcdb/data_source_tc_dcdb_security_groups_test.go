package dcdb_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDCDBSecurityGroupsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceDcdbSecurityGroups_basic, tcacctest.DefaultDcdbSGId, tcacctest.DefaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_security_groups.security_groups"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_security_groups.security_groups", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_security_groups.security_groups", "list.0.create_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_security_groups.security_groups", "list.0.security_group_id", tcacctest.DefaultDcdbSGId),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_security_groups.security_groups", "list.0.security_group_name", tcacctest.DefaultDcdbSGName),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_security_groups.security_groups", "list.0.inbound.#"),
				),
			},
		},
	})
}

const testAccDataSourceDcdbSecurityGroups_basic = `

resource "tencentcloud_dcdb_security_group_attachment" "default" {
  security_group_id = "%s"
  instance_id = "%s"
}

data "tencentcloud_dcdb_security_groups" "security_groups" {
  instance_id = tencentcloud_dcdb_security_group_attachment.default.instance_id
}

`
