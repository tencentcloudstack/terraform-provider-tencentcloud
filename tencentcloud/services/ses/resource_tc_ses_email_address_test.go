package ses_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSesEmail_address_basic -v
func TestAccTencentCloudSesEmail_address_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckBusiness(t, tcacctest.ACCOUNT_TYPE_SES) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesEmail_address,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ses_email_address.email_address", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ses_email_address.email_address", "email_address", "aaa@iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_ses_email_address.email_address", "email_sender_name", "aaa"),
				),
			},
			{
				ResourceName:      "tencentcloud_ses_email_address.email_address",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSesEmail_address = `

resource "tencentcloud_ses_email_address" "email_address" {
  email_address     = "aaa@iac-tf.cloud"
  email_sender_name = "aaa"
}

`
