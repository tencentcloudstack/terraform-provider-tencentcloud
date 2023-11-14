package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsWordSamplesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsWordSamplesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_word_samples.word_samples")),
			},
		},
	})
}

const testAccMpsWordSamplesDataSource = `

data "tencentcloud_mps_word_samples" "word_samples" {
  keywords = 
  usages = 
  }

`
