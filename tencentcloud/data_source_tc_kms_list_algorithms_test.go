package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudKmsListAlgorithmsDataSource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_kms_list_algorithms.example"),
				),
			},
		},
	})
}

const testAccKmsListAlgorithmsDataSource = `
data "tencentcloud_kms_list_algorithms" "example" {}
`
