package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodRecordGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodRecordGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record_group.record_group", "domain", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record_group.record_group", "group_name", "group_demo"),
				),
			},
			{
				Config: testAccDnspodRecordGroupUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record_group.record_group", "domain", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_record_group.record_group", "group_name", "group_demo2"),
				),
			},
			{
				ResourceName:      "tencentcloud_dnspod_record_group.record_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDnspodRecordGroup = `

resource "tencentcloud_dnspod_record_group" "record_group" {
  domain = "iac-tf.cloud"
  group_name = "group_demo"
}

`

const testAccDnspodRecordGroupUp = `

resource "tencentcloud_dnspod_record_group" "record_group" {
  domain = "iac-tf.cloud"
  group_name = "group_demo2"
}

`
