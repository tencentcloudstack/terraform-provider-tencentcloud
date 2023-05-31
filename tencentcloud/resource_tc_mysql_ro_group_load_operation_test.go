package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlRoGroupLoadOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlRoGroupLoadOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_group_load_operation.ro_group_load_operation", "id")),
			},
		},
	})
}

const testAccMysqlRoGroupLoadOperation = `

resource "tencentcloud_mysql_ro_group_load_operation" "ro_group_load_operation" {
	ro_group_id = "cdbrg-bdlvcfpj"
}

`
