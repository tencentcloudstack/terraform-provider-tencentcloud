package ses_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudSesReceiversDataSource_basic -v
func TestAccTencentCloudSesReceiversDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-hongkong")
			tcacctest.AccPreCheckBusiness(t, tcacctest.ACCOUNT_TYPE_SES)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesReceiversDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ses_receivers.receivers"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_receivers.receivers", "data.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_receivers.receivers", "data.0.count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_receivers.receivers", "data.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_receivers.receivers", "data.0.desc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_receivers.receivers", "data.0.receiver_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_receivers.receivers", "data.0.receivers_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_receivers.receivers", "data.0.receivers_status"),
				),
			},
		},
	})
}

const testAccSesReceiversDataSource = `

data "tencentcloud_ses_receivers" "receivers" {
  status   = 3
  key_word = "keep"
}

`
