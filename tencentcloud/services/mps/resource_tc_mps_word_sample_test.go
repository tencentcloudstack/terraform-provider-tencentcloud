package mps_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsWordSampleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsWordSample,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_word_sample.word_sample", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_word_sample.word_sample", "usages.#", "3"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_word_sample.word_sample", "usages.*", "Recognition.Ocr"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_word_sample.word_sample", "usages.*", "Review.Ocr"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_word_sample.word_sample", "usages.*", "Review.Asr"),
					resource.TestCheckResourceAttr("tencentcloud_mps_word_sample.word_sample", "keyword", "tf_test_kw_1"),
					resource.TestCheckResourceAttr("tencentcloud_mps_word_sample.word_sample", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_word_sample.word_sample", "tags.*", "tags_1"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_word_sample.word_sample", "tags.*", "tags_2"),
				),
			},
			{
				Config: testAccMpsWordSample_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_word_sample.word_sample", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_word_sample.word_sample", "usages.#", "2"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_word_sample.word_sample", "usages.*", "Recognition.Asr"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_word_sample.word_sample", "usages.*", "Review.Ocr"),
					resource.TestCheckResourceAttr("tencentcloud_mps_word_sample.word_sample", "keyword", "tf_test_kw_1"),
					resource.TestCheckResourceAttr("tencentcloud_mps_word_sample.word_sample", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_word_sample.word_sample", "tags.*", "tags_1"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_word_sample.word_sample", "tags.*", "tags_3"),
				),
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
  usages = ["Recognition.Ocr","Review.Ocr","Review.Asr"]
  keyword = "tf_test_kw_1"
  tags = ["tags_1", "tags_2"]
}

`

const testAccMpsWordSample_update = `

resource "tencentcloud_mps_word_sample" "word_sample" {
  usages = ["Recognition.Asr","Review.Ocr"]
  keyword = "tf_test_kw_1"
  tags = ["tags_1", "tags_3"]
}

`
