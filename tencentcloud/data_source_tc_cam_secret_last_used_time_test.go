package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCamSecretLastUsedTimeDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamSecretLastUsedTimeDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cam_secret_last_used_time.secret_last_used_time")),
			},
		},
	})
}

const testAccCamSecretLastUsedTimeDataSource = `

data "tencentcloud_cam_secret_last_used_time" "secret_last_used_time" {
  secret_id_list = ["AKID9KZRGC1d8VjtqdIkEcnproZFzmW8oUDz"]
  }

`
