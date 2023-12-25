package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApiGroupDataSource_basic -v
func TestAccTencentCloudTsfApiGroupDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApiGroupDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_api_group.api_group"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.acl_mode"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.auth_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.binded_gateway_deploy_groups.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.binded_gateway_deploy_groups.0.application_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.binded_gateway_deploy_groups.0.application_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.binded_gateway_deploy_groups.0.application_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.binded_gateway_deploy_groups.0.cluster_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.binded_gateway_deploy_groups.0.deploy_group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.binded_gateway_deploy_groups.0.deploy_group_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.created_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.gateway_instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.group_context"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.group_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.group_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.service_name_key_position"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_group.api_group", "result.0.content.0.updated_time"),
				),
			},
		},
	})
}

const testAccTsfApiGroupDataSource = `

data "tencentcloud_tsf_api_group" "api_group" {
	group_type = "ms"
	auth_type = "none"
	status = "released"
	order_by = "created_time"
	order_type = 0
	gateway_instance_id = "gw-ins-lvdypq5k"
}

`
