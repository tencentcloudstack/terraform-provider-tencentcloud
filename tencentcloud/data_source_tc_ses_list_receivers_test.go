package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSesListReceiversDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesListReceiversDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ses_list_receivers.list_receivers")),
			},
		},
	})
}

const testAccSesListReceiversDataSource = `

data "tencentcloud_ses_list_receivers" "list_receivers" {
  status = 1
  key_word = "xxx"
  }

`
