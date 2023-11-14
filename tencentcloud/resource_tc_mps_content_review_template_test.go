package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_content_review_template.content_review_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_content_review_template.content_review_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsContentReviewTemplate = `

resource "tencentcloud_mps_content_review_template" "content_review_template" {
  name = ""
  comment = ""
  porn_configure {
		img_review_info {
			switch = ""
			label_set = 
			block_confidence = 
			review_confidence = 
		}
		asr_review_info {
			switch = ""
			block_confidence = 
			review_confidence = 
		}
		ocr_review_info {
			switch = ""
			block_confidence = 
			review_confidence = 
		}

  }
  terrorism_configure {
		img_review_info {
			switch = ""
			label_set = 
			block_confidence = 
			review_confidence = 
		}
		ocr_review_info {
			switch = ""
			block_confidence = 
			review_confidence = 
		}

  }
  political_configure {
		img_review_info {
			switch = ""
			label_set = 
			block_confidence = 
			review_confidence = 
		}
		asr_review_info {
			switch = ""
			block_confidence = 
			review_confidence = 
		}
		ocr_review_info {
			switch = ""
			block_confidence = 
			review_confidence = 
		}

  }
  prohibited_configure {
		asr_review_info {
			switch = ""
			block_confidence = 
			review_confidence = 
		}
		ocr_review_info {
			switch = ""
			block_confidence = 
			review_confidence = 
		}

  }
  user_define_configure {
		face_review_info {
			switch = ""
			label_set = 
			block_confidence = 
			review_confidence = 
		}
		asr_review_info {
			switch = ""
			label_set = 
			block_confidence = 
			review_confidence = 
		}
		ocr_review_info {
			switch = ""
			label_set = 
			block_confidence = 
			review_confidence = 
		}

  }
}

`
