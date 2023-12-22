package ses_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudSesEmailIdentitiesDataSource_basic -v
func TestAccTencentCloudSesEmailIdentitiesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-hongkong")
			tcacctest.AccPreCheckBusiness(t, tcacctest.ACCOUNT_TYPE_SES)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesEmailIdentitiesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ses_email_identities.email_identities"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_email_identities.email_identities", "email_identities.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_email_identities.email_identities", "email_identities.0.current_reputation_level"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_email_identities.email_identities", "email_identities.0.daily_quota"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_email_identities.email_identities", "email_identities.0.identity_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_email_identities.email_identities", "email_identities.0.identity_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_email_identities.email_identities", "email_identities.0.sending_enabled"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_email_identities.email_identities", "max_daily_quota"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_email_identities.email_identities", "max_reputation_level"),
				),
			},
		},
	})
}

const testAccSesEmailIdentitiesDataSource = `

data "tencentcloud_ses_email_identities" "email_identities" {
}

`
