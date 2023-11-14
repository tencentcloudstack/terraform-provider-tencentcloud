package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverCompleteExpansionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverCompleteExpansion,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_complete_expansion.complete_expansion", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_complete_expansion.complete_expansion",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverCompleteExpansion = `

resource "tencentcloud_sqlserver_complete_expansion" "complete_expansion" {
  instance_id = "mssql-i1z41iwd"
}

`
