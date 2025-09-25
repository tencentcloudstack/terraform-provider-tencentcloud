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
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cls_delivery.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cls_delivery.example", "cls_info_list"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cls_delivery.example", "log_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cls_delivery.example", "enable_cls_delivery"),
				),
			},
			{
				Config: testAccCynosdbClsDeliveryUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cls_delivery.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cls_delivery.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cls_delivery.example", "cls_info_list"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cls_delivery.example", "log_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cls_delivery.example", "enable_cls_delivery"),
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
  instance_id = "cynosdbmysql-ins-anknkhvi"
  cls_info_list {
    topic_operation = "reuse"
    group_operation = "reuse"
    region          = "ap-guangzhou"
    topic_id        = "8e38f7c1-17ec-4acb-a4cb-7dc14baaef47"
    group_id        = "7e3bb8b7-81d5-4e6b-b150-f139b90c146e"
  }
  log_type            = "slow"
  enable_cls_delivery = true
}
`

const testAccCynosdbClsDeliveryUpdate = `
resource "tencentcloud_cynosdb_cls_delivery" "example" {
  instance_id = "cynosdbmysql-ins-anknkhvi"
  cls_info_list {
    topic_operation = "reuse"
    group_operation = "reuse"
    region          = "ap-guangzhou"
    topic_id        = "8e38f7c1-17ec-4acb-a4cb-7dc14baaef47"
    group_id        = "7e3bb8b7-81d5-4e6b-b150-f139b90c146e"
  }
  log_type            = "slow"
  enable_cls_delivery = false
}
`
