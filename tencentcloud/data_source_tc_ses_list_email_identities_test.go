package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSesListEmailIdentitiesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesListEmailIdentitiesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ses_list_email_identities.list_email_identities")),
			},
		},
	})
}

const testAccSesListEmailIdentitiesDataSource = `

data "tencentcloud_ses_list_email_identities" "list_email_identities" {
      }

`
