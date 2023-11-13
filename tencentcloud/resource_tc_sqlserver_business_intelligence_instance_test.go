package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverBusinessIntelligenceInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverBusinessIntelligenceInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_business_intelligence_instance.business_intelligence_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_business_intelligence_instance.business_intelligence_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverBusinessIntelligenceInstance = `

resource "tencentcloud_sqlserver_business_intelligence_instance" "business_intelligence_instance" {
  zone = "ap-guangzhou-1"
  memory = 10
  storage = 100
  cpu = 2
  machine_type = "CLOUD_SSD"
  project_id = 0
  goods_num = 1
  subnet_id = "subnet-bdoe83fa"
  vpc_id = "vpc-dsp338hz"
  d_b_version = ""
  security_group_list = 
  weekly = 
  start_time = ""
  span = 
  resource_tags {
		tag_key = ""
		tag_value = ""

  }
}

`
