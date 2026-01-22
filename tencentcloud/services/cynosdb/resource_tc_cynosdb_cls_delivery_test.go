package cynosdb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCynosdbClsDeliveryResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClsDelivery,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cls_delivery.example", "id"),
				),
			},
			{
				Config: testAccCynosdbClsDeliveryUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cls_delivery.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_cls_delivery.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbClsDelivery = `
resource "tencentcloud_cynosdb_cls_delivery" "example" {
  instance_id = "cynosdbmysql-ins-m2903cxq"
  cls_info_list {
    region   = "ap-guangzhou"
    topic_id = "a9d582f8-8c14-462c-94b8-bbc579a04f02"
    group_id = "67fca013-379b-4bc6-8e72-390227d869c4"
  }

  running_status = false
}
`

const testAccCynosdbClsDeliveryUpdate = `
resource "tencentcloud_cynosdb_cls_delivery" "example" {
  instance_id = "cynosdbmysql-ins-m2903cxq"
  cls_info_list {
    region   = "ap-guangzhou"
    topic_id = "a9d582f8-8c14-462c-94b8-bbc579a04f02"
    group_id = "67fca013-379b-4bc6-8e72-390227d869c4"
  }

  running_status = true
}
`
