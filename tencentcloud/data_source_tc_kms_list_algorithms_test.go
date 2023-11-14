package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudKmsListAlgorithmsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsListAlgorithmsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_kms_list_algorithms.list_algorithms")),
			},
		},
	})
}

const testAccKmsListAlgorithmsDataSource = `

data "tencentcloud_kms_list_algorithms" "list_algorithms" {
      }

`
