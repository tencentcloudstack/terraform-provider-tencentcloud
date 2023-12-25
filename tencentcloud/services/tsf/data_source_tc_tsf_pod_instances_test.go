package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfPodInstancesDataSource_basic -v
func TestAccTencentCloudTsfPodInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfPodInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_pod_instances.pod_instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_pod_instances.pod_instances", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_pod_instances.pod_instances", "result.0.total_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_pod_instances.pod_instances", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_pod_instances.pod_instances", "result.0.content.0.created_at"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_pod_instances.pod_instances", "result.0.content.0.ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_pod_instances.pod_instances", "result.0.content.0.node_instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_pod_instances.pod_instances", "result.0.content.0.node_ip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_pod_instances.pod_instances", "result.0.content.0.pod_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_pod_instances.pod_instances", "result.0.content.0.pod_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_pod_instances.pod_instances", "result.0.content.0.ready_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_pod_instances.pod_instances", "result.0.content.0.runtime"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_pod_instances.pod_instances", "result.0.content.0.service_instance_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_pod_instances.pod_instances", "result.0.content.0.status"),
				),
			},
		},
	})
}

const testAccTsfPodInstancesDataSourceVar = `
variable "group_id" {
	default = "` + tcacctest.DefaultTsfContainerGroupId + `"
}

variable "pod_name" {
	default = "` + tcacctest.DefaultTsfpodName + `"
}
`

const testAccTsfPodInstancesDataSource = testAccTsfPodInstancesDataSourceVar + `

data "tencentcloud_tsf_pod_instances" "pod_instances" {
	group_id = var.group_id
	pod_name_list = [var.pod_name]
}

`
