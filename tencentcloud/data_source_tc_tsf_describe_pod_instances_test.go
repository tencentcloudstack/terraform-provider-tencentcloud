package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfDescribePodInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDescribePodInstancesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_describe_pod_instances.describe_pod_instances")),
			},
		},
	})
}

const testAccTsfDescribePodInstancesDataSource = `

data "tencentcloud_tsf_describe_pod_instances" "describe_pod_instances" {
  group_id = ""
  pod_name_list = 
  }

`
