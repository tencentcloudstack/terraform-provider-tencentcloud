package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsContentReviewTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsContentReviewTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_content_review_templates.content_review_templates")),
			},
		},
	})
}

const testAccMpsContentReviewTemplatesDataSource = `

data "tencentcloud_mps_content_review_templates" "content_review_templates" {
  definitions = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  content_review_template_set {
		definition = &lt;nil&gt;
		name = &lt;nil&gt;
		comment = &lt;nil&gt;
		porn_configure {
			img_review_info {
				switch = &lt;nil&gt;
				label_set = &lt;nil&gt;
				block_confidence = &lt;nil&gt;
				review_confidence = &lt;nil&gt;
			}
			asr_review_info {
				switch = &lt;nil&gt;
				block_confidence = &lt;nil&gt;
				review_confidence = &lt;nil&gt;
			}
			ocr_review_info {
				switch = &lt;nil&gt;
				block_confidence = &lt;nil&gt;
				review_confidence = &lt;nil&gt;
			}
		}
		terrorism_configure {
			img_review_info {
				switch = &lt;nil&gt;
				label_set = &lt;nil&gt;
				block_confidence = &lt;nil&gt;
				review_confidence = &lt;nil&gt;
			}
			ocr_review_info {
				switch = &lt;nil&gt;
				block_confidence = &lt;nil&gt;
				review_confidence = &lt;nil&gt;
			}
		}
		political_configure {
			img_review_info {
				switch = &lt;nil&gt;
				label_set = &lt;nil&gt;
				block_confidence = &lt;nil&gt;
				review_confidence = &lt;nil&gt;
			}
			asr_review_info {
				switch = &lt;nil&gt;
				block_confidence = &lt;nil&gt;
				review_confidence = &lt;nil&gt;
			}
			ocr_review_info {
				switch = &lt;nil&gt;
				block_confidence = &lt;nil&gt;
				review_confidence = &lt;nil&gt;
			}
		}
		prohibited_configure {
			asr_review_info {
				switch = &lt;nil&gt;
				block_confidence = &lt;nil&gt;
				review_confidence = &lt;nil&gt;
			}
			ocr_review_info {
				switch = &lt;nil&gt;
				block_confidence = &lt;nil&gt;
				review_confidence = &lt;nil&gt;
			}
		}
		user_define_configure {
			face_review_info {
				switch = &lt;nil&gt;
				label_set = &lt;nil&gt;
				block_confidence = &lt;nil&gt;
				review_confidence = &lt;nil&gt;
			}
			asr_review_info {
				switch = &lt;nil&gt;
				label_set = &lt;nil&gt;
				block_confidence = &lt;nil&gt;
				review_confidence = &lt;nil&gt;
			}
			ocr_review_info {
				switch = &lt;nil&gt;
				label_set = &lt;nil&gt;
				block_confidence = &lt;nil&gt;
				review_confidence = &lt;nil&gt;
			}
		}
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;
		type = &lt;nil&gt;

  }
}

`
