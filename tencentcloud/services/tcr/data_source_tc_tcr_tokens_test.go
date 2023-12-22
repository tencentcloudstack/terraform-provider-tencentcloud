package tcr_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataTCRTokensNameAll = "data.tencentcloud_tcr_tokens.id_test"

func TestAccTencentCloudTcrTokensData(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTCRTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTCRTokensBasic,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
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

const testAccTencentCloudDataTCRTokensBasic = tcacctest.DefaultTCRInstanceData + `
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
