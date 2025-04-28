package cls_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixClsCloudProductLogTaskV2Resource_basic -v
func TestAccTencentCloudNeedFixClsCloudProductLogTaskV2Resource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsCloudProductLogTaskV2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task_v2.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task_v2.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task_v2.example", "assumer_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task_v2.example", "log_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task_v2.example", "cloud_product_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task_v2.example", "cls_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task_v2.example", "logset_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task_v2.example", "topic_name"),
				),
			},
			{
				Config: testAccClsCloudProductLogTaskV2Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task_v2.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task_v2.example", "extend"),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_cloud_product_log_task_v2.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClsCloudProductLogTaskV2 = `
resource "tencentcloud_cls_cloud_product_log_task_v2" "example" {
  instance_id          = "postgres-mcdstv8l"
  assumer_name         = "PostgreSQL"
  log_type             = "PostgreSQL-SLOW"
  cloud_product_region = "gz"
  cls_region           = "ap-guangzhou"
  logset_name          = "tf-example"
  topic_name           = "tf-example"
}
`

const testAccClsCloudProductLogTaskV2Update = `
resource "tencentcloud_cls_cloud_product_log_task_v2" "example" {
  instance_id          = "postgres-mcdstv8l"
  assumer_name         = "PostgreSQL"
  log_type             = "PostgreSQL-SLOW"
  cloud_product_region = "gz"
  cls_region           = "ap-guangzhou"
  logset_name          = "tf-example"
  topic_name           = "tf-example"
  extend               = "{\"ServiceName\":[\"HDFS\",\"KNOX\",\"YARN\",\"ZOOKEEPER\"],\"Policy\":0}"
}
`
