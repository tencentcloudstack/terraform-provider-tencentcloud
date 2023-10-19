package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudOceanusResourceConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusResourceConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_oceanus_resource_config.resource_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_oceanus_resource_config.resource_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOceanusResourceConfig = `

resource "tencentcloud_oceanus_resource_config" "resource_config" {
  resource_id = "resource-xxx"
  resource_loc {
		storage_type = 1
		param {
			bucket = "scs-online-1257058945"
			path = "251008563/100000006047/flink-cos-fs-hadoop-oceanus-1-20210304112"
			region = "ap-chengdu"
		}

  }
  remark = "xxx"
  auto_delete = 1
  work_space_id = "space-xxx"
}

`
