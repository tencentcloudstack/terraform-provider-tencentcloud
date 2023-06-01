package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfPodInstancesDataSource_basic -v
func TestAccTencentCloudTsfPodInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfPodInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_pod_instances.pod_instances"),
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

const testAccTsfPodInstancesDataSource = `

data "tencentcloud_tsf_pod_instances" "pod_instances" {
	group_id = "group-ynd95rea"
	pod_name_list = ["keep-terraform-6f8f977688-zvphm"]
}

`
