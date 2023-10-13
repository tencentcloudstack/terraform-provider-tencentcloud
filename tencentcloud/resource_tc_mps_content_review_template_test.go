package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsContentReviewTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsContentReviewTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "name", "tf_test_content_review_temp"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "comment", "tf test content review temp"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "porn_configure.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "porn_configure.0.img_review_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "porn_configure.0.img_review_info.0.switch", "ON"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "porn_configure.0.img_review_info.0.label_set.*", "porn"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "porn_configure.0.img_review_info.0.label_set.*", "vulgar"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "porn_configure.0.img_review_info.0.block_confidence", "60"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "porn_configure.0.img_review_info.0.review_confidence", "100"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "porn_configure.0.asr_review_info.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "porn_configure.0.ocr_review_info.#"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "terrorism_configure.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "terrorism_configure.0.img_review_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "terrorism_configure.0.img_review_info.0.switch", "ON"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "terrorism_configure.0.img_review_info.0.label_set.*", "guns"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "terrorism_configure.0.img_review_info.0.label_set.*", "crowd"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "terrorism_configure.0.img_review_info.0.block_confidence", "60"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "terrorism_configure.0.img_review_info.0.review_confidence", "100"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "terrorism_configure.0.ocr_review_info.#"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "political_configure.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "political_configure.0.img_review_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "political_configure.0.img_review_info.0.switch", "ON"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "political_configure.0.img_review_info.0.label_set.*", "violation_photo"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "political_configure.0.img_review_info.0.label_set.*", "politician"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "political_configure.0.img_review_info.0.block_confidence", "60"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "political_configure.0.img_review_info.0.review_confidence", "100"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "political_configure.0.asr_review_info.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "political_configure.0.ocr_review_info.#"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "prohibited_configure.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "prohibited_configure.0.asr_review_info.0.switch", "ON"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "prohibited_configure.0.asr_review_info.0.block_confidence", "60"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "prohibited_configure.0.asr_review_info.0.review_confidence", "100"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "prohibited_configure.0.ocr_review_info.0.switch", "ON"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "prohibited_configure.0.ocr_review_info.0.block_confidence", "60"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "prohibited_configure.0.ocr_review_info.0.review_confidence", "100"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "user_define_configure.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "user_define_configure.0.face_review_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.face_review_info.0.switch", "ON"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.face_review_info.0.label_set.*", "FACE_1"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.face_review_info.0.label_set.*", "FACE_2"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.face_review_info.0.block_confidence", "60"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.face_review_info.0.review_confidence", "100"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "user_define_configure.0.asr_review_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.asr_review_info.0.switch", "ON"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.asr_review_info.0.label_set.*", "VOICE_1"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.asr_review_info.0.label_set.*", "VOICE_2"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.asr_review_info.0.block_confidence", "60"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.asr_review_info.0.review_confidence", "100"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "user_define_configure.0.ocr_review_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.ocr_review_info.0.switch", "ON"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.ocr_review_info.0.label_set.*", "VIDEO_1"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.ocr_review_info.0.label_set.*", "VIDEO_2"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.ocr_review_info.0.block_confidence", "60"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.ocr_review_info.0.review_confidence", "100"),
				),
			},
			{
				Config: testAccMpsContentReviewTemplate_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "name", "tf_test_content_review_temp_changed"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "comment", "tf test content review temp changed"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "porn_configure.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "porn_configure.0.img_review_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "porn_configure.0.img_review_info.0.switch", "ON"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "porn_configure.0.img_review_info.0.label_set.*", "porn"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "porn_configure.0.img_review_info.0.label_set.*", "sexy"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "porn_configure.0.img_review_info.0.block_confidence", "80"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "porn_configure.0.img_review_info.0.review_confidence", "100"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "porn_configure.0.asr_review_info.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "porn_configure.0.ocr_review_info.#"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "terrorism_configure.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "terrorism_configure.0.img_review_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "terrorism_configure.0.img_review_info.0.switch", "ON"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "terrorism_configure.0.img_review_info.0.label_set.*", "guns"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "terrorism_configure.0.img_review_info.0.label_set.*", "bloody"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "terrorism_configure.0.img_review_info.0.block_confidence", "80"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "terrorism_configure.0.img_review_info.0.review_confidence", "100"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "terrorism_configure.0.ocr_review_info.#"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "political_configure.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "political_configure.0.img_review_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "political_configure.0.img_review_info.0.switch", "ON"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "political_configure.0.img_review_info.0.label_set.*", "entertainment"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "political_configure.0.img_review_info.0.label_set.*", "politician"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "political_configure.0.img_review_info.0.block_confidence", "60"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "political_configure.0.img_review_info.0.review_confidence", "90"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "political_configure.0.asr_review_info.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "political_configure.0.ocr_review_info.#"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "prohibited_configure.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "prohibited_configure.0.asr_review_info.0.switch", "ON"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "prohibited_configure.0.asr_review_info.0.block_confidence", "60"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "prohibited_configure.0.asr_review_info.0.review_confidence", "90"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "prohibited_configure.0.ocr_review_info.0.switch", "ON"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "prohibited_configure.0.ocr_review_info.0.block_confidence", "60"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "prohibited_configure.0.ocr_review_info.0.review_confidence", "90"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "user_define_configure.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "user_define_configure.0.face_review_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.face_review_info.0.switch", "ON"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.face_review_info.0.label_set.*", "FACE_1"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.face_review_info.0.label_set.*", "FACE_3"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.face_review_info.0.block_confidence", "80"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.face_review_info.0.review_confidence", "100"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "user_define_configure.0.asr_review_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.asr_review_info.0.switch", "OFF"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.asr_review_info.0.label_set.*", "VOICE_1"),
					resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.asr_review_info.0.label_set.*", "VOICE_2"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.asr_review_info.0.block_confidence", "60"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.asr_review_info.0.review_confidence", "100"),

					resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.template", "user_define_configure.0.ocr_review_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.ocr_review_info.0.switch", "OFF"),
					// resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.ocr_review_info.0.label_set.*", "VIDEO_1"),
					// resource.TestCheckTypeSetElemAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.ocr_review_info.0.label_set.*", "VIDEO_2"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.ocr_review_info.0.block_confidence", "60"),
					resource.TestCheckResourceAttr("tencentcloud_mps_content_review_template.template", "user_define_configure.0.ocr_review_info.0.review_confidence", "100"),
				),
			},
			{
				ResourceName:      "tencentcloud_mps_content_review_template.template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsContentReviewTemplate = `

resource "tencentcloud_mps_content_review_template" "template" {
  name    = "tf_test_content_review_temp"
  comment = "tf test content review temp"
  porn_configure {
    img_review_info {
      switch            = "ON"
      label_set         = ["porn", "vulgar"]
      block_confidence  = 60
      review_confidence = 100
    }
    asr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }
    ocr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }

  }
  terrorism_configure {
    img_review_info {
      switch            = "ON"
      label_set         = ["guns", "crowd"]
      block_confidence  = 60
      review_confidence = 100
    }
    ocr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }

  }
  political_configure {
    img_review_info {
      switch            = "ON"
      label_set         = ["violation_photo", "politician"]
      block_confidence  = 60
      review_confidence = 100
    }
    asr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }
    ocr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }

  }
  prohibited_configure {
    asr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }
    ocr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }

  }
  user_define_configure {
    face_review_info {
      switch            = "ON"
      label_set         = ["FACE_1", "FACE_2"]
      block_confidence  = 60
      review_confidence = 100
    }
    asr_review_info {
      switch            = "ON"
      label_set         = ["VOICE_1", "VOICE_2"]
      block_confidence  = 60
      review_confidence = 100
    }
    ocr_review_info {
      switch            = "ON"
      label_set         = ["VIDEO_1", "VIDEO_2"]
      block_confidence  = 60
      review_confidence = 100
    }
  }
}

`

const testAccMpsContentReviewTemplate_update = `

resource "tencentcloud_mps_content_review_template" "template" {
  name    = "tf_test_content_review_temp_changed"
  comment = "tf test content review temp changed"
  porn_configure {
    img_review_info {
      switch            = "ON"
      label_set         = ["porn", "sexy"]
      block_confidence  = 80
      review_confidence = 100
    }
    asr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }
    ocr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }

  }
  terrorism_configure {
    img_review_info {
      switch            = "ON"
      label_set         = ["guns", "bloody"]
      block_confidence  = 80
      review_confidence = 100
    }
    ocr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }

  }
  political_configure {
    img_review_info {
      switch            = "ON"
      label_set         = ["entertainment", "politician"]
      block_confidence  = 60
      review_confidence = 90
    }
    asr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }
    ocr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }

  }
  prohibited_configure {
    asr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 90
    }
    ocr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 90
    }

  }
  user_define_configure {
    face_review_info {
      switch            = "ON"
      label_set         = ["FACE_1", "FACE_3"]
      block_confidence  = 80
      review_confidence = 100
    }
    asr_review_info {
      switch            = "OFF"
      label_set         = ["VOICE_1", "VOICE_2"]
      block_confidence  = 60
      review_confidence = 100
    }
    ocr_review_info {
      switch            = "OFF"
      label_set         = ["VIDEO_1", "VIDEO_2"]
      block_confidence  = 60
      review_confidence = 100
    }
  }
}

`
