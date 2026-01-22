package cls_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixClsCloudProductLogTaskResource_basic -v
func TestAccTencentCloudNeedFixClsCloudProductLogTaskResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsCloudProductLogTask,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task.example", "assumer_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task.example", "log_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task.example", "cloud_product_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task.example", "cls_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task.example", "logset_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task.example", "topic_name"),
				),
			},
			{
				Config: testAccClsCloudProductLogTaskUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_cloud_product_log_task.example", "extend"),
				),
			},
			{
				ResourceName:            "tencentcloud_cls_cloud_product_log_task.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"logset_name", "topic_name"},
			},
		},
	})
}

const testAccClsCloudProductLogTask = `
resource "tencentcloud_cls_cloud_product_log_task" "example" {
  instance_id          = "postgres-mcdstv8l"
  assumer_name         = "PostgreSQL"
  log_type             = "PostgreSQL-SLOW"
  cloud_product_region = "gz"
  cls_region           = "ap-guangzhou"
  logset_name          = "tf-example"
  topic_name           = "tf-example"
}
`

const testAccClsCloudProductLogTaskUpdate = `
resource "tencentcloud_cls_cloud_product_log_task" "example" {
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
