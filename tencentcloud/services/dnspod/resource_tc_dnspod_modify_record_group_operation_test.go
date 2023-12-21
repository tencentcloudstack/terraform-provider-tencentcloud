package dnspod_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodModifyRecordGroupOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodModifyRecordGroupOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dnspod_modify_record_group_operation.modify_record_group", "domain", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_modify_record_group_operation.modify_record_group", "group_id", "25129"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_modify_record_group_operation.modify_record_group", "record_id", "1363850006"),
					// resource.TestCheckResourceAttr("tencentcloud_dnspod_modify_record_group_operation.modify_record_group", "domain_id", "123"),
				),
			},
		},
	})
}

const testAccDnspodModifyRecordGroupOperation = `

resource "tencentcloud_dnspod_modify_record_group_operation" "modify_record_group" {
  domain = "iac-tf.cloud"
  group_id = "25129"
  record_id = "1363850006"
  # domain_id = "123"
}

`
