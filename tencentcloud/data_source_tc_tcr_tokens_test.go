package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataTCRTokensNameAll = "data.tencentcloud_tcr_tokens.id_test"

func TestAccTencentCloudTCRTokensData(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTCRTokensBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataTCRTokensNameAll, "token_list.0.token_id"),
					resource.TestCheckResourceAttrSet(testDataTCRTokensNameAll, "token_list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataTCRTokensNameAll, "token_list.0.description"),
					resource.TestCheckResourceAttr(testDataTCRTokensNameAll, "token_list.0.enable", "true"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTCRTokensBasic = defaultTCRInstanceData + `
resource "tencentcloud_tcr_token" "mytcr_token" {
  instance_id = local.tcr_id
  description       = "test"
  enable   = true
}

data "tencentcloud_tcr_tokens" "id_test" {
  token_id = tencentcloud_tcr_token.mytcr_token.token_id
  instance_id = local.tcr_id
}
`
