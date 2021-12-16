package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudClsLogsetDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsLogSetDataSource_basic,
			},
		},
	})
}

const testAccClsLogSetDataSource_basic = `

resource "tencentcloud_cls_logset" "logset_basic"{
    logset_name = "test"
	tags{
		 key = "aaa"
		 value = "bbb"
      }
}

data "tencentcloud_cls_logsets" "logsets" {
     filters {
                key = "logsetId"
                value = [tencentcloud_cls_logset.logset_basic.id]                                     
             }  
     limit = 1
}
`
