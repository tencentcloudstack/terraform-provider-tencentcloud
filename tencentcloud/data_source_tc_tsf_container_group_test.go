package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfDContainerGroupDataSource_basic -v
func TestAccTencentCloudTsfDContainerGroupDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfContainerGroupDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_container_group.container_group"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_container_group.container_group", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_container_group.container_group", "result.0.total_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_container_group.container_group", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_container_group.container_group", "result.0.content.0.alias"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_container_group.container_group", "result.0.content.0.cluster_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_container_group.container_group", "result.0.content.0.cluster_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_container_group.container_group", "result.0.content.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_container_group.container_group", "result.0.content.0.group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_container_group.container_group", "result.0.content.0.group_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_container_group.container_group", "result.0.content.0.kube_inject_enable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_container_group.container_group", "result.0.content.0.namespace_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_container_group.container_group", "result.0.content.0.namespace_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_container_group.container_group", "result.0.content.0.updated_time"),
				),
			},
		},
	})
}

const testAccTsfContainerGroupDataSourceVar = `
variable "application_id" {
	default = "` + defaultTsfApplicationId + `"
}
variable "cluster_id" {
	default = "` + defaultTsfClustId + `"
}
variable "namespace_id" {
	default = "` + defaultNamespaceId + `"
}
`

const testAccTsfContainerGroupDataSource = testAccTsfContainerGroupDataSourceVar + `

data "tencentcloud_tsf_container_group" "container_group" {
	application_id = var.application_id
	search_word = "keep"
	order_by = "createTime"
	order_type = 0
	cluster_id = var.cluster_id
	namespace_id = var.namespace_id
}

`
