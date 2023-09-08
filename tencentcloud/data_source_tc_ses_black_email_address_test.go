package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudSesBlackEmailAddressDataSource_basic -v
func TestAccTencentCloudSesBlackEmailAddressDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-hongkong")
			testAccPreCheckBusiness(t, ACCOUNT_TYPE_SES)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesBlackEmailAddressDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_ses_black_email_address.black_email_address"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_black_email_address.black_email_address", "black_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_black_email_address.black_email_address", "black_list.0.bounce_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_black_email_address.black_email_address", "black_list.0.email_address"),
				),
			},
		},
	})
}

const testAccSesBlackEmailAddressDataSource = `

data "tencentcloud_ses_black_email_address" "black_email_address" {
  start_date = "2020-09-22"
  end_date = "2023-09-23"
  email_address = "terraform-tf@gmail.com"
}

`
