package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfGatewayAllGroupApisDataSource_basic -v
func TestAccTencentCloudTsfGatewayAllGroupApisDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfGatewayAllGroupApisDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_gateway_all_group_apis.gateway_all_group_apis"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_gateway_all_group_apis.gateway_all_group_apis", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_gateway_all_group_apis.gateway_all_group_apis", "result.0.gateway_deploy_group_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_gateway_all_group_apis.gateway_all_group_apis", "result.0.gateway_deploy_group_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_gateway_all_group_apis.gateway_all_group_apis", "result.0.group_num"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_gateway_all_group_apis.gateway_all_group_apis", "result.0.groups.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_gateway_all_group_apis.gateway_all_group_apis", "result.0.groups.0.gateway_instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_gateway_all_group_apis.gateway_all_group_apis", "result.0.groups.0.gateway_instance_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_gateway_all_group_apis.gateway_all_group_apis", "result.0.groups.0.group_api_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_gateway_all_group_apis.gateway_all_group_apis", "result.0.groups.0.group_apis.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_gateway_all_group_apis.gateway_all_group_apis", "result.0.groups.0.group_apis.0.api_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_gateway_all_group_apis.gateway_all_group_apis", "result.0.groups.0.group_apis.0.method"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_gateway_all_group_apis.gateway_all_group_apis", "result.0.groups.0.group_apis.0.microservice_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_gateway_all_group_apis.gateway_all_group_apis", "result.0.groups.0.group_apis.0.path"),
				),
			},
		},
	})
}

const testAccTsfGatewayAllGroupApisDataSource = `

data "tencentcloud_tsf_gateway_all_group_apis" "gateway_all_group_apis" {
	gateway_deploy_group_id = "group-aeoej4qy"
	search_word = "user"
}

`
