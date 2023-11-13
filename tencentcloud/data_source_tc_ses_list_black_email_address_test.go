package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSesListBlackEmailAddressDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesListBlackEmailAddressDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ses_list_black_email_address.list_black_email_address")),
			},
		},
	})
}

const testAccSesListBlackEmailAddressDataSource = `

data "tencentcloud_ses_list_black_email_address" "list_black_email_address" {
  start_date = "2020-09-22"
  end_date = "2020-09-23"
  email_address = "xxx@mail.qcloud.com"
  task_i_d = "7000"
  }

`
