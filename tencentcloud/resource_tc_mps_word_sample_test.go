package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsWordSampleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsWordSample,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_word_sample.word_sample", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_word_sample.word_sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsWordSample = `

resource "tencentcloud_mps_word_sample" "word_sample" {
  usages = 
  words {
		keyword = ""
		tags = 

  }
}

`
