/*
Provides a resource to create a mps content_review_template

Example Usage

```hcl
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
```

Import

mps content_review_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_content_review_template.content_review_template content_review_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMpsContentReviewTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsContentReviewTemplateCreate,
		Read:   resourceTencentCloudMpsContentReviewTemplateRead,
		Update: resourceTencentCloudMpsContentReviewTemplateUpdate,
		Delete: resourceTencentCloudMpsContentReviewTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The name of the content moderation template. Length limit: 64 characters.",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The template description. Length limit: 256 characters.",
			},

			"porn_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Control parameter for porn information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"img_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Control parameter of porn information detection in image.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Switch of a porn information detection in image task. Valid values:&amp;lt;li&amp;gt;ON: Enables a porn information detection in image task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: Disables a porn information detection in image task.&amp;lt;/li&amp;gt;.",
									},
									"label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Filter tag for porn information detection in image. If an audit result contains the selected tag, it will be returned; if the filter tag is empty, all audit results will be returned. Valid values:&amp;lt;li&amp;gt;porn: Porn;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;vulgar: Vulgarity;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;intimacy: Intimacy;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;sexy: Sexiness.&amp;lt;/li&amp;gt;.",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 90 will be used by default. Value range: 0-100.",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 0 will be used by default. Value range: 0-100.",
									},
								},
							},
						},
						"asr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Control parameter of porn information detection in speech.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Switch of a porn information detection in speech task. Valid values:&amp;lt;li&amp;gt;ON: Enables a porn information detection in speech task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: Disables a porn information detection in speech task.&amp;lt;/li&amp;gt;.",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0-100.",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0-100.",
									},
								},
							},
						},
						"ocr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Control parameter of porn information detection in text.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Switch of a porn information detection in text task. Valid values:&amp;lt;li&amp;gt;ON: Enables a porn information detection in text task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: Disables a porn information detection in text task.&amp;lt;/li&amp;gt;.",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0-100.",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0-100.",
									},
								},
							},
						},
					},
				},
			},

			"terrorism_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Control parameter for terrorism information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"img_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The parameters for detecting sensitive information in images.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Whether to detect sensitive information in images. Valid values:&amp;lt;li&amp;gt;ON&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF&amp;lt;/li&amp;gt;.",
									},
									"label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "The filter labels for sensitive information detection in images, which specify the types of sensitive information to return. If this parameter is left empty, the detection results for all labels are returned. Valid values:&amp;lt;li&amp;gt;guns&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;crowd&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;bloody&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;police&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;banners (sensitive flags)&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;militant&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;explosion&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;terrorists&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;scenario (sensitive scenes) &amp;lt;/li&amp;gt;.",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 90 will be used by default. Value range: 0-100.",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 80 will be used by default. Value range: 0-100.",
									},
								},
							},
						},
						"ocr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "The parameters for detecting sensitive information based on OCR.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Whether to detect sensitive information based on OCR. Valid values:&amp;lt;li&amp;gt;ON&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF&amp;lt;/li&amp;gt;.",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0–100.",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0–100.",
									},
								},
							},
						},
					},
				},
			},

			"political_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Control parameter for politically sensitive information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"img_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The parameters for detecting sensitive information in images.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Whether to detect sensitive information in images. Valid values:&amp;lt;li&amp;gt;ON&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF&amp;lt;/li&amp;gt;.",
									},
									"label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "The filter labels for sensitive information detection in images, which specify the types of sensitive information to return. If this parameter is left empty, the detection results for all labels are returned. Valid values:&amp;lt;li&amp;gt;violation_photo (banned icons)&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;politician&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;entertainment (people in the entertainment industry)&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;sport (people in the sports industry)&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;entrepreneur&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;scholar&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;celebrity&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;military (people in military)&amp;lt;/li&amp;gt;.",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 97 will be used by default. Value range: 0-100.",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 95 will be used by default. Value range: 0-100.",
									},
								},
							},
						},
						"asr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The parameters for detecting sensitive information based on ASR.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Whether to detect sensitive information based on ASR. Valid values:&amp;lt;li&amp;gt;ON&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF&amp;lt;/li&amp;gt;.",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0-100.",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0-100.",
									},
								},
							},
						},
						"ocr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The parameters for detecting sensitive information based on OCR.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Whether to detect sensitive information based on OCR. Valid values:&amp;lt;li&amp;gt;ON&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF&amp;lt;/li&amp;gt;.",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0-100.",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0-100.",
									},
								},
							},
						},
					},
				},
			},

			"prohibited_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Control parameter of prohibited information detection. Prohibited information includes:&amp;amp;lt;li&amp;amp;gt;Abusive;&amp;amp;lt;/li&amp;amp;gt;&amp;amp;lt;li&amp;amp;gt;Drug-related.&amp;amp;lt;/li&amp;amp;gt;Note: this parameter is not supported yet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Control parameter of prohibited information detection in speech.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Switch of prohibited information detection in speech task. Valid values:&amp;lt;li&amp;gt;ON: enables prohibited information detection in speech task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: disables prohibited information detection in speech task.&amp;lt;/li&amp;gt;.",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0–100.",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0–100.",
									},
								},
							},
						},
						"ocr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Control parameter of prohibited information detection in text.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Switch of prohibited information detection in text task. Valid values:&amp;lt;li&amp;gt;ON: enables prohibited information detection in text task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: disables prohibited information detection in text task.&amp;lt;/li&amp;gt;.",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0–100.",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0–100.",
									},
								},
							},
						},
					},
				},
			},

			"user_define_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Custom content moderation parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"face_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Control parameter of custom figure audit.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Switch of a custom figure audit task. Valid values:&amp;lt;li&amp;gt;ON: Enables a custom figure audit task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: Disables a custom figure audit task.&amp;lt;/li&amp;gt;.",
									},
									"label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Custom figure filter tag. If an audit result contains the selected tag, it will be returned; if the filter tag is empty, all audit results will be returned. To use the tag filtering feature, you need to add the corresponding tag when adding materials for the custom figure library.There can be up to 10 tags, each with a length limit of 16 characters.",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 97 will be used by default. Value range: 0-100.",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 95 will be used by default. Value range: 0-100.",
									},
								},
							},
						},
						"asr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Control parameter of custom speech audit.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Switch of a custom speech audit task. Valid values:&amp;lt;li&amp;gt;ON: Enables a custom speech audit task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: Disables a custom speech audit task.&amp;lt;/li&amp;gt;.",
									},
									"label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Custom speech filter tag. If an audit result contains the selected tag, it will be returned; if the filter tag is empty, all audit results will be returned. To use the tag filtering feature, you need to add the corresponding tag when adding materials for custom speech keywords.There can be up to 10 tags, each with a length limit of 16 characters.",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0-100.",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0-100.",
									},
								},
							},
						},
						"ocr_review_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Control parameter of custom text audit.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Switch of a custom text audit task. Valid values:&amp;lt;li&amp;gt;ON: Enables a custom text audit task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: Disables a custom text audit task.&amp;lt;/li&amp;gt;.",
									},
									"label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Custom text filter tag. If an audit result contains the selected tag, it will be returned; if the filter tag is empty, all audit results will be returned. To use the tag filtering feature, you need to add the corresponding tag when adding materials for custom text keywords.There can be up to 10 tags, each with a length limit of 16 characters.",
									},
									"block_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0-100.",
									},
									"review_confidence": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0-100.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMpsContentReviewTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_content_review_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mps.NewCreateContentReviewTemplateRequest()
		response   = mps.NewCreateContentReviewTemplateResponse()
		definition int
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "porn_configure"); ok {
		pornConfigureInfo := mps.PornConfigureInfo{}
		if imgReviewInfoMap, ok := helper.InterfaceToMap(dMap, "img_review_info"); ok {
			pornImgReviewTemplateInfo := mps.PornImgReviewTemplateInfo{}
			if v, ok := imgReviewInfoMap["switch"]; ok {
				pornImgReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := imgReviewInfoMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					labelSet := labelSetSet[i].(string)
					pornImgReviewTemplateInfo.LabelSet = append(pornImgReviewTemplateInfo.LabelSet, &labelSet)
				}
			}
			if v, ok := imgReviewInfoMap["block_confidence"]; ok {
				pornImgReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := imgReviewInfoMap["review_confidence"]; ok {
				pornImgReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			pornConfigureInfo.ImgReviewInfo = &pornImgReviewTemplateInfo
		}
		if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
			pornAsrReviewTemplateInfo := mps.PornAsrReviewTemplateInfo{}
			if v, ok := asrReviewInfoMap["switch"]; ok {
				pornAsrReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := asrReviewInfoMap["block_confidence"]; ok {
				pornAsrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := asrReviewInfoMap["review_confidence"]; ok {
				pornAsrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			pornConfigureInfo.AsrReviewInfo = &pornAsrReviewTemplateInfo
		}
		if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
			pornOcrReviewTemplateInfo := mps.PornOcrReviewTemplateInfo{}
			if v, ok := ocrReviewInfoMap["switch"]; ok {
				pornOcrReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
				pornOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
				pornOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			pornConfigureInfo.OcrReviewInfo = &pornOcrReviewTemplateInfo
		}
		request.PornConfigure = &pornConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "terrorism_configure"); ok {
		terrorismConfigureInfo := mps.TerrorismConfigureInfo{}
		if imgReviewInfoMap, ok := helper.InterfaceToMap(dMap, "img_review_info"); ok {
			terrorismImgReviewTemplateInfo := mps.TerrorismImgReviewTemplateInfo{}
			if v, ok := imgReviewInfoMap["switch"]; ok {
				terrorismImgReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := imgReviewInfoMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					labelSet := labelSetSet[i].(string)
					terrorismImgReviewTemplateInfo.LabelSet = append(terrorismImgReviewTemplateInfo.LabelSet, &labelSet)
				}
			}
			if v, ok := imgReviewInfoMap["block_confidence"]; ok {
				terrorismImgReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := imgReviewInfoMap["review_confidence"]; ok {
				terrorismImgReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			terrorismConfigureInfo.ImgReviewInfo = &terrorismImgReviewTemplateInfo
		}
		if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
			terrorismOcrReviewTemplateInfo := mps.TerrorismOcrReviewTemplateInfo{}
			if v, ok := ocrReviewInfoMap["switch"]; ok {
				terrorismOcrReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
				terrorismOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
				terrorismOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			terrorismConfigureInfo.OcrReviewInfo = &terrorismOcrReviewTemplateInfo
		}
		request.TerrorismConfigure = &terrorismConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "political_configure"); ok {
		politicalConfigureInfo := mps.PoliticalConfigureInfo{}
		if imgReviewInfoMap, ok := helper.InterfaceToMap(dMap, "img_review_info"); ok {
			politicalImgReviewTemplateInfo := mps.PoliticalImgReviewTemplateInfo{}
			if v, ok := imgReviewInfoMap["switch"]; ok {
				politicalImgReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := imgReviewInfoMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					labelSet := labelSetSet[i].(string)
					politicalImgReviewTemplateInfo.LabelSet = append(politicalImgReviewTemplateInfo.LabelSet, &labelSet)
				}
			}
			if v, ok := imgReviewInfoMap["block_confidence"]; ok {
				politicalImgReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := imgReviewInfoMap["review_confidence"]; ok {
				politicalImgReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			politicalConfigureInfo.ImgReviewInfo = &politicalImgReviewTemplateInfo
		}
		if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
			politicalAsrReviewTemplateInfo := mps.PoliticalAsrReviewTemplateInfo{}
			if v, ok := asrReviewInfoMap["switch"]; ok {
				politicalAsrReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := asrReviewInfoMap["block_confidence"]; ok {
				politicalAsrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := asrReviewInfoMap["review_confidence"]; ok {
				politicalAsrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			politicalConfigureInfo.AsrReviewInfo = &politicalAsrReviewTemplateInfo
		}
		if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
			politicalOcrReviewTemplateInfo := mps.PoliticalOcrReviewTemplateInfo{}
			if v, ok := ocrReviewInfoMap["switch"]; ok {
				politicalOcrReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
				politicalOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
				politicalOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			politicalConfigureInfo.OcrReviewInfo = &politicalOcrReviewTemplateInfo
		}
		request.PoliticalConfigure = &politicalConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "prohibited_configure"); ok {
		prohibitedConfigureInfo := mps.ProhibitedConfigureInfo{}
		if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
			prohibitedAsrReviewTemplateInfo := mps.ProhibitedAsrReviewTemplateInfo{}
			if v, ok := asrReviewInfoMap["switch"]; ok {
				prohibitedAsrReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := asrReviewInfoMap["block_confidence"]; ok {
				prohibitedAsrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := asrReviewInfoMap["review_confidence"]; ok {
				prohibitedAsrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			prohibitedConfigureInfo.AsrReviewInfo = &prohibitedAsrReviewTemplateInfo
		}
		if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
			prohibitedOcrReviewTemplateInfo := mps.ProhibitedOcrReviewTemplateInfo{}
			if v, ok := ocrReviewInfoMap["switch"]; ok {
				prohibitedOcrReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
				prohibitedOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
				prohibitedOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			prohibitedConfigureInfo.OcrReviewInfo = &prohibitedOcrReviewTemplateInfo
		}
		request.ProhibitedConfigure = &prohibitedConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "user_define_configure"); ok {
		userDefineConfigureInfo := mps.UserDefineConfigureInfo{}
		if faceReviewInfoMap, ok := helper.InterfaceToMap(dMap, "face_review_info"); ok {
			userDefineFaceReviewTemplateInfo := mps.UserDefineFaceReviewTemplateInfo{}
			if v, ok := faceReviewInfoMap["switch"]; ok {
				userDefineFaceReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := faceReviewInfoMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					labelSet := labelSetSet[i].(string)
					userDefineFaceReviewTemplateInfo.LabelSet = append(userDefineFaceReviewTemplateInfo.LabelSet, &labelSet)
				}
			}
			if v, ok := faceReviewInfoMap["block_confidence"]; ok {
				userDefineFaceReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := faceReviewInfoMap["review_confidence"]; ok {
				userDefineFaceReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			userDefineConfigureInfo.FaceReviewInfo = &userDefineFaceReviewTemplateInfo
		}
		if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
			userDefineAsrTextReviewTemplateInfo := mps.UserDefineAsrTextReviewTemplateInfo{}
			if v, ok := asrReviewInfoMap["switch"]; ok {
				userDefineAsrTextReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := asrReviewInfoMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					labelSet := labelSetSet[i].(string)
					userDefineAsrTextReviewTemplateInfo.LabelSet = append(userDefineAsrTextReviewTemplateInfo.LabelSet, &labelSet)
				}
			}
			if v, ok := asrReviewInfoMap["block_confidence"]; ok {
				userDefineAsrTextReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := asrReviewInfoMap["review_confidence"]; ok {
				userDefineAsrTextReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			userDefineConfigureInfo.AsrReviewInfo = &userDefineAsrTextReviewTemplateInfo
		}
		if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
			userDefineOcrTextReviewTemplateInfo := mps.UserDefineOcrTextReviewTemplateInfo{}
			if v, ok := ocrReviewInfoMap["switch"]; ok {
				userDefineOcrTextReviewTemplateInfo.Switch = helper.String(v.(string))
			}
			if v, ok := ocrReviewInfoMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					labelSet := labelSetSet[i].(string)
					userDefineOcrTextReviewTemplateInfo.LabelSet = append(userDefineOcrTextReviewTemplateInfo.LabelSet, &labelSet)
				}
			}
			if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
				userDefineOcrTextReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
			}
			if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
				userDefineOcrTextReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
			}
			userDefineConfigureInfo.OcrReviewInfo = &userDefineOcrTextReviewTemplateInfo
		}
		request.UserDefineConfigure = &userDefineConfigureInfo
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().CreateContentReviewTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps contentReviewTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = *response.Response.Definition
	d.SetId(helper.Int64ToStr(definition))

	return resourceTencentCloudMpsContentReviewTemplateRead(d, meta)
}

func resourceTencentCloudMpsContentReviewTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_content_review_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	contentReviewTemplateId := d.Id()

	contentReviewTemplate, err := service.DescribeMpsContentReviewTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if contentReviewTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsContentReviewTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if contentReviewTemplate.Name != nil {
		_ = d.Set("name", contentReviewTemplate.Name)
	}

	if contentReviewTemplate.Comment != nil {
		_ = d.Set("comment", contentReviewTemplate.Comment)
	}

	if contentReviewTemplate.PornConfigure != nil {
		pornConfigureMap := map[string]interface{}{}

		if contentReviewTemplate.PornConfigure.ImgReviewInfo != nil {
			imgReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.PornConfigure.ImgReviewInfo.Switch != nil {
				imgReviewInfoMap["switch"] = contentReviewTemplate.PornConfigure.ImgReviewInfo.Switch
			}

			if contentReviewTemplate.PornConfigure.ImgReviewInfo.LabelSet != nil {
				imgReviewInfoMap["label_set"] = contentReviewTemplate.PornConfigure.ImgReviewInfo.LabelSet
			}

			if contentReviewTemplate.PornConfigure.ImgReviewInfo.BlockConfidence != nil {
				imgReviewInfoMap["block_confidence"] = contentReviewTemplate.PornConfigure.ImgReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.PornConfigure.ImgReviewInfo.ReviewConfidence != nil {
				imgReviewInfoMap["review_confidence"] = contentReviewTemplate.PornConfigure.ImgReviewInfo.ReviewConfidence
			}

			pornConfigureMap["img_review_info"] = []interface{}{imgReviewInfoMap}
		}

		if contentReviewTemplate.PornConfigure.AsrReviewInfo != nil {
			asrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.PornConfigure.AsrReviewInfo.Switch != nil {
				asrReviewInfoMap["switch"] = contentReviewTemplate.PornConfigure.AsrReviewInfo.Switch
			}

			if contentReviewTemplate.PornConfigure.AsrReviewInfo.BlockConfidence != nil {
				asrReviewInfoMap["block_confidence"] = contentReviewTemplate.PornConfigure.AsrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.PornConfigure.AsrReviewInfo.ReviewConfidence != nil {
				asrReviewInfoMap["review_confidence"] = contentReviewTemplate.PornConfigure.AsrReviewInfo.ReviewConfidence
			}

			pornConfigureMap["asr_review_info"] = []interface{}{asrReviewInfoMap}
		}

		if contentReviewTemplate.PornConfigure.OcrReviewInfo != nil {
			ocrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.PornConfigure.OcrReviewInfo.Switch != nil {
				ocrReviewInfoMap["switch"] = contentReviewTemplate.PornConfigure.OcrReviewInfo.Switch
			}

			if contentReviewTemplate.PornConfigure.OcrReviewInfo.BlockConfidence != nil {
				ocrReviewInfoMap["block_confidence"] = contentReviewTemplate.PornConfigure.OcrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.PornConfigure.OcrReviewInfo.ReviewConfidence != nil {
				ocrReviewInfoMap["review_confidence"] = contentReviewTemplate.PornConfigure.OcrReviewInfo.ReviewConfidence
			}

			pornConfigureMap["ocr_review_info"] = []interface{}{ocrReviewInfoMap}
		}

		_ = d.Set("porn_configure", []interface{}{pornConfigureMap})
	}

	if contentReviewTemplate.TerrorismConfigure != nil {
		terrorismConfigureMap := map[string]interface{}{}

		if contentReviewTemplate.TerrorismConfigure.ImgReviewInfo != nil {
			imgReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.Switch != nil {
				imgReviewInfoMap["switch"] = contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.Switch
			}

			if contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.LabelSet != nil {
				imgReviewInfoMap["label_set"] = contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.LabelSet
			}

			if contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.BlockConfidence != nil {
				imgReviewInfoMap["block_confidence"] = contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.ReviewConfidence != nil {
				imgReviewInfoMap["review_confidence"] = contentReviewTemplate.TerrorismConfigure.ImgReviewInfo.ReviewConfidence
			}

			terrorismConfigureMap["img_review_info"] = []interface{}{imgReviewInfoMap}
		}

		if contentReviewTemplate.TerrorismConfigure.OcrReviewInfo != nil {
			ocrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.TerrorismConfigure.OcrReviewInfo.Switch != nil {
				ocrReviewInfoMap["switch"] = contentReviewTemplate.TerrorismConfigure.OcrReviewInfo.Switch
			}

			if contentReviewTemplate.TerrorismConfigure.OcrReviewInfo.BlockConfidence != nil {
				ocrReviewInfoMap["block_confidence"] = contentReviewTemplate.TerrorismConfigure.OcrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.TerrorismConfigure.OcrReviewInfo.ReviewConfidence != nil {
				ocrReviewInfoMap["review_confidence"] = contentReviewTemplate.TerrorismConfigure.OcrReviewInfo.ReviewConfidence
			}

			terrorismConfigureMap["ocr_review_info"] = []interface{}{ocrReviewInfoMap}
		}

		_ = d.Set("terrorism_configure", []interface{}{terrorismConfigureMap})
	}

	if contentReviewTemplate.PoliticalConfigure != nil {
		politicalConfigureMap := map[string]interface{}{}

		if contentReviewTemplate.PoliticalConfigure.ImgReviewInfo != nil {
			imgReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.Switch != nil {
				imgReviewInfoMap["switch"] = contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.Switch
			}

			if contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.LabelSet != nil {
				imgReviewInfoMap["label_set"] = contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.LabelSet
			}

			if contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.BlockConfidence != nil {
				imgReviewInfoMap["block_confidence"] = contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.ReviewConfidence != nil {
				imgReviewInfoMap["review_confidence"] = contentReviewTemplate.PoliticalConfigure.ImgReviewInfo.ReviewConfidence
			}

			politicalConfigureMap["img_review_info"] = []interface{}{imgReviewInfoMap}
		}

		if contentReviewTemplate.PoliticalConfigure.AsrReviewInfo != nil {
			asrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.PoliticalConfigure.AsrReviewInfo.Switch != nil {
				asrReviewInfoMap["switch"] = contentReviewTemplate.PoliticalConfigure.AsrReviewInfo.Switch
			}

			if contentReviewTemplate.PoliticalConfigure.AsrReviewInfo.BlockConfidence != nil {
				asrReviewInfoMap["block_confidence"] = contentReviewTemplate.PoliticalConfigure.AsrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.PoliticalConfigure.AsrReviewInfo.ReviewConfidence != nil {
				asrReviewInfoMap["review_confidence"] = contentReviewTemplate.PoliticalConfigure.AsrReviewInfo.ReviewConfidence
			}

			politicalConfigureMap["asr_review_info"] = []interface{}{asrReviewInfoMap}
		}

		if contentReviewTemplate.PoliticalConfigure.OcrReviewInfo != nil {
			ocrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.PoliticalConfigure.OcrReviewInfo.Switch != nil {
				ocrReviewInfoMap["switch"] = contentReviewTemplate.PoliticalConfigure.OcrReviewInfo.Switch
			}

			if contentReviewTemplate.PoliticalConfigure.OcrReviewInfo.BlockConfidence != nil {
				ocrReviewInfoMap["block_confidence"] = contentReviewTemplate.PoliticalConfigure.OcrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.PoliticalConfigure.OcrReviewInfo.ReviewConfidence != nil {
				ocrReviewInfoMap["review_confidence"] = contentReviewTemplate.PoliticalConfigure.OcrReviewInfo.ReviewConfidence
			}

			politicalConfigureMap["ocr_review_info"] = []interface{}{ocrReviewInfoMap}
		}

		_ = d.Set("political_configure", []interface{}{politicalConfigureMap})
	}

	if contentReviewTemplate.ProhibitedConfigure != nil {
		prohibitedConfigureMap := map[string]interface{}{}

		if contentReviewTemplate.ProhibitedConfigure.AsrReviewInfo != nil {
			asrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.ProhibitedConfigure.AsrReviewInfo.Switch != nil {
				asrReviewInfoMap["switch"] = contentReviewTemplate.ProhibitedConfigure.AsrReviewInfo.Switch
			}

			if contentReviewTemplate.ProhibitedConfigure.AsrReviewInfo.BlockConfidence != nil {
				asrReviewInfoMap["block_confidence"] = contentReviewTemplate.ProhibitedConfigure.AsrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.ProhibitedConfigure.AsrReviewInfo.ReviewConfidence != nil {
				asrReviewInfoMap["review_confidence"] = contentReviewTemplate.ProhibitedConfigure.AsrReviewInfo.ReviewConfidence
			}

			prohibitedConfigureMap["asr_review_info"] = []interface{}{asrReviewInfoMap}
		}

		if contentReviewTemplate.ProhibitedConfigure.OcrReviewInfo != nil {
			ocrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.ProhibitedConfigure.OcrReviewInfo.Switch != nil {
				ocrReviewInfoMap["switch"] = contentReviewTemplate.ProhibitedConfigure.OcrReviewInfo.Switch
			}

			if contentReviewTemplate.ProhibitedConfigure.OcrReviewInfo.BlockConfidence != nil {
				ocrReviewInfoMap["block_confidence"] = contentReviewTemplate.ProhibitedConfigure.OcrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.ProhibitedConfigure.OcrReviewInfo.ReviewConfidence != nil {
				ocrReviewInfoMap["review_confidence"] = contentReviewTemplate.ProhibitedConfigure.OcrReviewInfo.ReviewConfidence
			}

			prohibitedConfigureMap["ocr_review_info"] = []interface{}{ocrReviewInfoMap}
		}

		_ = d.Set("prohibited_configure", []interface{}{prohibitedConfigureMap})
	}

	if contentReviewTemplate.UserDefineConfigure != nil {
		userDefineConfigureMap := map[string]interface{}{}

		if contentReviewTemplate.UserDefineConfigure.FaceReviewInfo != nil {
			faceReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.Switch != nil {
				faceReviewInfoMap["switch"] = contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.Switch
			}

			if contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.LabelSet != nil {
				faceReviewInfoMap["label_set"] = contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.LabelSet
			}

			if contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.BlockConfidence != nil {
				faceReviewInfoMap["block_confidence"] = contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.ReviewConfidence != nil {
				faceReviewInfoMap["review_confidence"] = contentReviewTemplate.UserDefineConfigure.FaceReviewInfo.ReviewConfidence
			}

			userDefineConfigureMap["face_review_info"] = []interface{}{faceReviewInfoMap}
		}

		if contentReviewTemplate.UserDefineConfigure.AsrReviewInfo != nil {
			asrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.Switch != nil {
				asrReviewInfoMap["switch"] = contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.Switch
			}

			if contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.LabelSet != nil {
				asrReviewInfoMap["label_set"] = contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.LabelSet
			}

			if contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.BlockConfidence != nil {
				asrReviewInfoMap["block_confidence"] = contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.ReviewConfidence != nil {
				asrReviewInfoMap["review_confidence"] = contentReviewTemplate.UserDefineConfigure.AsrReviewInfo.ReviewConfidence
			}

			userDefineConfigureMap["asr_review_info"] = []interface{}{asrReviewInfoMap}
		}

		if contentReviewTemplate.UserDefineConfigure.OcrReviewInfo != nil {
			ocrReviewInfoMap := map[string]interface{}{}

			if contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.Switch != nil {
				ocrReviewInfoMap["switch"] = contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.Switch
			}

			if contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.LabelSet != nil {
				ocrReviewInfoMap["label_set"] = contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.LabelSet
			}

			if contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.BlockConfidence != nil {
				ocrReviewInfoMap["block_confidence"] = contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.BlockConfidence
			}

			if contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.ReviewConfidence != nil {
				ocrReviewInfoMap["review_confidence"] = contentReviewTemplate.UserDefineConfigure.OcrReviewInfo.ReviewConfidence
			}

			userDefineConfigureMap["ocr_review_info"] = []interface{}{ocrReviewInfoMap}
		}

		_ = d.Set("user_define_configure", []interface{}{userDefineConfigureMap})
	}

	return nil
}

func resourceTencentCloudMpsContentReviewTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_content_review_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewModifyContentReviewTemplateRequest()

	contentReviewTemplateId := d.Id()

	request.Definition = &definition

	immutableArgs := []string{"name", "comment", "porn_configure", "terrorism_configure", "political_configure", "prohibited_configure", "user_define_configure"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}
	}

	if d.HasChange("porn_configure") {
		if dMap, ok := helper.InterfacesHeadMap(d, "porn_configure"); ok {
			pornConfigureInfo := mps.PornConfigureInfo{}
			if imgReviewInfoMap, ok := helper.InterfaceToMap(dMap, "img_review_info"); ok {
				pornImgReviewTemplateInfo := mps.PornImgReviewTemplateInfo{}
				if v, ok := imgReviewInfoMap["switch"]; ok {
					pornImgReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := imgReviewInfoMap["label_set"]; ok {
					labelSetSet := v.(*schema.Set).List()
					for i := range labelSetSet {
						labelSet := labelSetSet[i].(string)
						pornImgReviewTemplateInfo.LabelSet = append(pornImgReviewTemplateInfo.LabelSet, &labelSet)
					}
				}
				if v, ok := imgReviewInfoMap["block_confidence"]; ok {
					pornImgReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := imgReviewInfoMap["review_confidence"]; ok {
					pornImgReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				pornConfigureInfo.ImgReviewInfo = &pornImgReviewTemplateInfo
			}
			if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
				pornAsrReviewTemplateInfo := mps.PornAsrReviewTemplateInfo{}
				if v, ok := asrReviewInfoMap["switch"]; ok {
					pornAsrReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := asrReviewInfoMap["block_confidence"]; ok {
					pornAsrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := asrReviewInfoMap["review_confidence"]; ok {
					pornAsrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				pornConfigureInfo.AsrReviewInfo = &pornAsrReviewTemplateInfo
			}
			if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
				pornOcrReviewTemplateInfo := mps.PornOcrReviewTemplateInfo{}
				if v, ok := ocrReviewInfoMap["switch"]; ok {
					pornOcrReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
					pornOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
					pornOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				pornConfigureInfo.OcrReviewInfo = &pornOcrReviewTemplateInfo
			}
			request.PornConfigure = &pornConfigureInfo
		}
	}

	if d.HasChange("terrorism_configure") {
		if dMap, ok := helper.InterfacesHeadMap(d, "terrorism_configure"); ok {
			terrorismConfigureInfo := mps.TerrorismConfigureInfo{}
			if imgReviewInfoMap, ok := helper.InterfaceToMap(dMap, "img_review_info"); ok {
				terrorismImgReviewTemplateInfo := mps.TerrorismImgReviewTemplateInfo{}
				if v, ok := imgReviewInfoMap["switch"]; ok {
					terrorismImgReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := imgReviewInfoMap["label_set"]; ok {
					labelSetSet := v.(*schema.Set).List()
					for i := range labelSetSet {
						labelSet := labelSetSet[i].(string)
						terrorismImgReviewTemplateInfo.LabelSet = append(terrorismImgReviewTemplateInfo.LabelSet, &labelSet)
					}
				}
				if v, ok := imgReviewInfoMap["block_confidence"]; ok {
					terrorismImgReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := imgReviewInfoMap["review_confidence"]; ok {
					terrorismImgReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				terrorismConfigureInfo.ImgReviewInfo = &terrorismImgReviewTemplateInfo
			}
			if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
				terrorismOcrReviewTemplateInfo := mps.TerrorismOcrReviewTemplateInfo{}
				if v, ok := ocrReviewInfoMap["switch"]; ok {
					terrorismOcrReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
					terrorismOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
					terrorismOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				terrorismConfigureInfo.OcrReviewInfo = &terrorismOcrReviewTemplateInfo
			}
			request.TerrorismConfigure = &terrorismConfigureInfo
		}
	}

	if d.HasChange("political_configure") {
		if dMap, ok := helper.InterfacesHeadMap(d, "political_configure"); ok {
			politicalConfigureInfo := mps.PoliticalConfigureInfo{}
			if imgReviewInfoMap, ok := helper.InterfaceToMap(dMap, "img_review_info"); ok {
				politicalImgReviewTemplateInfo := mps.PoliticalImgReviewTemplateInfo{}
				if v, ok := imgReviewInfoMap["switch"]; ok {
					politicalImgReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := imgReviewInfoMap["label_set"]; ok {
					labelSetSet := v.(*schema.Set).List()
					for i := range labelSetSet {
						labelSet := labelSetSet[i].(string)
						politicalImgReviewTemplateInfo.LabelSet = append(politicalImgReviewTemplateInfo.LabelSet, &labelSet)
					}
				}
				if v, ok := imgReviewInfoMap["block_confidence"]; ok {
					politicalImgReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := imgReviewInfoMap["review_confidence"]; ok {
					politicalImgReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				politicalConfigureInfo.ImgReviewInfo = &politicalImgReviewTemplateInfo
			}
			if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
				politicalAsrReviewTemplateInfo := mps.PoliticalAsrReviewTemplateInfo{}
				if v, ok := asrReviewInfoMap["switch"]; ok {
					politicalAsrReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := asrReviewInfoMap["block_confidence"]; ok {
					politicalAsrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := asrReviewInfoMap["review_confidence"]; ok {
					politicalAsrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				politicalConfigureInfo.AsrReviewInfo = &politicalAsrReviewTemplateInfo
			}
			if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
				politicalOcrReviewTemplateInfo := mps.PoliticalOcrReviewTemplateInfo{}
				if v, ok := ocrReviewInfoMap["switch"]; ok {
					politicalOcrReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
					politicalOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
					politicalOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				politicalConfigureInfo.OcrReviewInfo = &politicalOcrReviewTemplateInfo
			}
			request.PoliticalConfigure = &politicalConfigureInfo
		}
	}

	if d.HasChange("prohibited_configure") {
		if dMap, ok := helper.InterfacesHeadMap(d, "prohibited_configure"); ok {
			prohibitedConfigureInfo := mps.ProhibitedConfigureInfo{}
			if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
				prohibitedAsrReviewTemplateInfo := mps.ProhibitedAsrReviewTemplateInfo{}
				if v, ok := asrReviewInfoMap["switch"]; ok {
					prohibitedAsrReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := asrReviewInfoMap["block_confidence"]; ok {
					prohibitedAsrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := asrReviewInfoMap["review_confidence"]; ok {
					prohibitedAsrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				prohibitedConfigureInfo.AsrReviewInfo = &prohibitedAsrReviewTemplateInfo
			}
			if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
				prohibitedOcrReviewTemplateInfo := mps.ProhibitedOcrReviewTemplateInfo{}
				if v, ok := ocrReviewInfoMap["switch"]; ok {
					prohibitedOcrReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
					prohibitedOcrReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
					prohibitedOcrReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				prohibitedConfigureInfo.OcrReviewInfo = &prohibitedOcrReviewTemplateInfo
			}
			request.ProhibitedConfigure = &prohibitedConfigureInfo
		}
	}

	if d.HasChange("user_define_configure") {
		if dMap, ok := helper.InterfacesHeadMap(d, "user_define_configure"); ok {
			userDefineConfigureInfo := mps.UserDefineConfigureInfo{}
			if faceReviewInfoMap, ok := helper.InterfaceToMap(dMap, "face_review_info"); ok {
				userDefineFaceReviewTemplateInfo := mps.UserDefineFaceReviewTemplateInfo{}
				if v, ok := faceReviewInfoMap["switch"]; ok {
					userDefineFaceReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := faceReviewInfoMap["label_set"]; ok {
					labelSetSet := v.(*schema.Set).List()
					for i := range labelSetSet {
						labelSet := labelSetSet[i].(string)
						userDefineFaceReviewTemplateInfo.LabelSet = append(userDefineFaceReviewTemplateInfo.LabelSet, &labelSet)
					}
				}
				if v, ok := faceReviewInfoMap["block_confidence"]; ok {
					userDefineFaceReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := faceReviewInfoMap["review_confidence"]; ok {
					userDefineFaceReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				userDefineConfigureInfo.FaceReviewInfo = &userDefineFaceReviewTemplateInfo
			}
			if asrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "asr_review_info"); ok {
				userDefineAsrTextReviewTemplateInfo := mps.UserDefineAsrTextReviewTemplateInfo{}
				if v, ok := asrReviewInfoMap["switch"]; ok {
					userDefineAsrTextReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := asrReviewInfoMap["label_set"]; ok {
					labelSetSet := v.(*schema.Set).List()
					for i := range labelSetSet {
						labelSet := labelSetSet[i].(string)
						userDefineAsrTextReviewTemplateInfo.LabelSet = append(userDefineAsrTextReviewTemplateInfo.LabelSet, &labelSet)
					}
				}
				if v, ok := asrReviewInfoMap["block_confidence"]; ok {
					userDefineAsrTextReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := asrReviewInfoMap["review_confidence"]; ok {
					userDefineAsrTextReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				userDefineConfigureInfo.AsrReviewInfo = &userDefineAsrTextReviewTemplateInfo
			}
			if ocrReviewInfoMap, ok := helper.InterfaceToMap(dMap, "ocr_review_info"); ok {
				userDefineOcrTextReviewTemplateInfo := mps.UserDefineOcrTextReviewTemplateInfo{}
				if v, ok := ocrReviewInfoMap["switch"]; ok {
					userDefineOcrTextReviewTemplateInfo.Switch = helper.String(v.(string))
				}
				if v, ok := ocrReviewInfoMap["label_set"]; ok {
					labelSetSet := v.(*schema.Set).List()
					for i := range labelSetSet {
						labelSet := labelSetSet[i].(string)
						userDefineOcrTextReviewTemplateInfo.LabelSet = append(userDefineOcrTextReviewTemplateInfo.LabelSet, &labelSet)
					}
				}
				if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
					userDefineOcrTextReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
				}
				if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
					userDefineOcrTextReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
				}
				userDefineConfigureInfo.OcrReviewInfo = &userDefineOcrTextReviewTemplateInfo
			}
			request.UserDefineConfigure = &userDefineConfigureInfo
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ModifyContentReviewTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps contentReviewTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsContentReviewTemplateRead(d, meta)
}

func resourceTencentCloudMpsContentReviewTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_content_review_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}
	contentReviewTemplateId := d.Id()

	if err := service.DeleteMpsContentReviewTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
