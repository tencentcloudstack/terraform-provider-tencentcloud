package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverBusinessIntelligenceFileResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverBusinessIntelligenceFile,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_business_intelligence_file.business_intelligence_file", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_business_intelligence_file.business_intelligence_file",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverBusinessIntelligenceFile = `

resource "tencentcloud_sqlserver_business_intelligence_file" "business_intelligence_file" {
  instance_id = "mssql-zjaha891"
  file_u_r_l = ""
  file_type = ""
  remark = ""
}

`
