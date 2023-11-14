/*
Use this data source to query detailed information of mps content_review_templates

Example Usage

```hcl
data "tencentcloud_mps_content_review_templates" "content_review_templates" {
  definitions =
  type = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMpsContentReviewTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsContentReviewTemplatesRead,
		Schema: map[string]*schema.Schema{
			"definitions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "The IDs of the content moderation templates to query. Array length limit: 50.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The filter for querying templates. If this parameter is left empty, both preset and custom templates are returned. Valid values:* Preset* Custom.",
			},

			"content_review_template_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of content audit template details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unique ID of a content audit template.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of a content audit template. Length limit: 64 characters.",
						},
						"comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of a content audit template. Length limit: 256 characters.",
						},
						"porn_configure": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Porn information detection control parameter.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"img_review_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Control parameter of porn information detection in image.Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Switch of a porn information detection in image task. Valid values:&amp;lt;li&amp;gt;ON: Enables a porn information detection in image task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: Disables a porn information detection in image task.&amp;lt;/li&amp;gt;.",
												},
												"label_set": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Filter tag for porn information detection in image. If an audit result contains the selected tag, it will be returned; if the filter tag is empty, all audit results will be returned. Valid values:&amp;lt;li&amp;gt;porn: Porn;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;vulgar: Vulgarity;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;intimacy: Intimacy;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;sexy: Sexiness.&amp;lt;/li&amp;gt;.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 90 will be used by default. Value range: 0-100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 0 will be used by default. Value range: 0-100.",
												},
											},
										},
									},
									"asr_review_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Control parameter of porn information detection in speech.Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Switch of a porn information detection in speech task. Valid values:&amp;lt;li&amp;gt;ON: Enables a porn information detection in speech task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: Disables a porn information detection in speech task.&amp;lt;/li&amp;gt;.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0-100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0-100.",
												},
											},
										},
									},
									"ocr_review_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Control parameter of porn information detection in text.Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Switch of a porn information detection in text task. Valid values:&amp;lt;li&amp;gt;ON: Enables a porn information detection in text task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: Disables a porn information detection in text task.&amp;lt;/li&amp;gt;.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0-100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0-100.",
												},
											},
										},
									},
								},
							},
						},
						"terrorism_configure": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The parameters for detecting sensitive information.Note: This field may return `null`, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"img_review_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The parameters for detecting sensitive information in images.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Whether to detect sensitive information in images. Valid values:&amp;lt;li&amp;gt;ON&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF&amp;lt;/li&amp;gt;.",
												},
												"label_set": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "The filter labels for sensitive information detection in images, which specify the types of sensitive information to return. If this parameter is left empty, the detection results for all labels are returned. Valid values:&amp;lt;li&amp;gt;guns&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;crowd&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;bloody&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;police&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;banners (sensitive flags)&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;militant&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;explosion&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;terrorists&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;scenario (sensitive scenes) &amp;lt;/li&amp;gt;.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 90 will be used by default. Value range: 0-100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 80 will be used by default. Value range: 0-100.",
												},
											},
										},
									},
									"ocr_review_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The parameters for detecting sensitive information based on OCR.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Whether to detect sensitive information based on OCR. Valid values:&amp;lt;li&amp;gt;ON&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF&amp;lt;/li&amp;gt;.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0–100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0–100.",
												},
											},
										},
									},
								},
							},
						},
						"political_configure": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The parameters for detecting sensitive information.Note: This field may return `null`, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"img_review_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The parameters for detecting sensitive information in images.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Whether to detect sensitive information in images. Valid values:&amp;lt;li&amp;gt;ON&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF&amp;lt;/li&amp;gt;.",
												},
												"label_set": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "The filter labels for sensitive information detection in images, which specify the types of sensitive information to return. If this parameter is left empty, the detection results for all labels are returned. Valid values:&amp;lt;li&amp;gt;violation_photo (banned icons)&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;politician&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;entertainment (people in the entertainment industry)&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;sport (people in the sports industry)&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;entrepreneur&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;scholar&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;celebrity&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;military (people in military)&amp;lt;/li&amp;gt;.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 97 will be used by default. Value range: 0-100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 95 will be used by default. Value range: 0-100.",
												},
											},
										},
									},
									"asr_review_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The parameters for detecting sensitive information based on ASR.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Whether to detect sensitive information based on ASR. Valid values:&amp;lt;li&amp;gt;ON&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF&amp;lt;/li&amp;gt;.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0-100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0-100.",
												},
											},
										},
									},
									"ocr_review_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The parameters for detecting sensitive information based on OCR.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Whether to detect sensitive information based on OCR. Valid values:&amp;lt;li&amp;gt;ON&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF&amp;lt;/li&amp;gt;.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0-100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0-100.",
												},
											},
										},
									},
								},
							},
						},
						"prohibited_configure": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Control parameter of prohibited information detection. Prohibited information includes:&amp;lt;li&amp;gt;Abusive;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;Drug-related.&amp;lt;/li&amp;gt;Note: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"asr_review_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Control parameter of prohibited information detection in speech.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Switch of prohibited information detection in speech task. Valid values:&amp;lt;li&amp;gt;ON: enables prohibited information detection in speech task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: disables prohibited information detection in speech task.&amp;lt;/li&amp;gt;.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0–100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0–100.",
												},
											},
										},
									},
									"ocr_review_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Control parameter of prohibited information detection in text.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Switch of prohibited information detection in text task. Valid values:&amp;lt;li&amp;gt;ON: enables prohibited information detection in text task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: disables prohibited information detection in text task.&amp;lt;/li&amp;gt;.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0–100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0–100.",
												},
											},
										},
									},
								},
							},
						},
						"user_define_configure": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Custom content audit control parameter.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"face_review_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Control parameter of custom figure audit.Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Switch of a custom figure audit task. Valid values:&amp;lt;li&amp;gt;ON: Enables a custom figure audit task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: Disables a custom figure audit task.&amp;lt;/li&amp;gt;.",
												},
												"label_set": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Custom figure filter tag. If an audit result contains the selected tag, it will be returned; if the filter tag is empty, all audit results will be returned. To use the tag filtering feature, you need to add the corresponding tag when adding materials for the custom figure library.There can be up to 10 tags, each with a length limit of 16 characters.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 97 will be used by default. Value range: 0-100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 95 will be used by default. Value range: 0-100.",
												},
											},
										},
									},
									"asr_review_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Control parameter of custom speech audit.Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Switch of a custom speech audit task. Valid values:&amp;lt;li&amp;gt;ON: Enables a custom speech audit task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: Disables a custom speech audit task.&amp;lt;/li&amp;gt;.",
												},
												"label_set": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Custom speech filter tag. If an audit result contains the selected tag, it will be returned; if the filter tag is empty, all audit results will be returned. To use the tag filtering feature, you need to add the corresponding tag when adding materials for custom speech keywords.There can be up to 10 tags, each with a length limit of 16 characters.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0-100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0-100.",
												},
											},
										},
									},
									"ocr_review_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Control parameter of custom text audit.Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Switch of a custom text audit task. Valid values:&amp;lt;li&amp;gt;ON: Enables a custom text audit task;&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;OFF: Disables a custom text audit task.&amp;lt;/li&amp;gt;.",
												},
												"label_set": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Custom text filter tag. If an audit result contains the selected tag, it will be returned; if the filter tag is empty, all audit results will be returned. To use the tag filtering feature, you need to add the corresponding tag when adding materials for custom text keywords.There can be up to 10 tags, each with a length limit of 16 characters.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for violation. If this score is reached or exceeded during intelligent audit, it will be deemed that a suspected violation has occurred. If this parameter is left empty, 100 will be used by default. Value range: 0-100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Threshold score for human audit. If this score is reached or exceeded during intelligent audit, human audit will be considered necessary. If this parameter is left empty, 75 will be used by default. Value range: 0-100.",
												},
											},
										},
									},
								},
							},
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of a template in [ISO date format](https://intl.cloud.tencent.com/document/product/266/11732?from_cn_redirect=1#iso-.E6.97.A5.E6.9C.9F.E6.A0.BC.E5.BC.8F).",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last modified time of a template in [ISO date format](https://intl.cloud.tencent.com/document/product/266/11732?from_cn_redirect=1#iso-.E6.97.A5.E6.9C.9F.E6.A0.BC.E5.BC.8F).",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The template type. Valid values:* Preset* CustomNote: This field may return `null`, indicating that no valid values can be obtained.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMpsContentReviewTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_content_review_templates.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("definitions"); ok {
		definitionsSet := v.(*schema.Set).List()
		for i := range definitionsSet {
			definitions := definitionsSet[i].(int)
			paramMap["Definitions"] = append(paramMap["Definitions"], helper.IntInt64(definitions))
		}
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var contentReviewTemplateSet []*mps.ContentReviewTemplateItem

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsContentReviewTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		contentReviewTemplateSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(contentReviewTemplateSet))
	tmpList := make([]map[string]interface{}, 0, len(contentReviewTemplateSet))

	if contentReviewTemplateSet != nil {
		for _, contentReviewTemplateItem := range contentReviewTemplateSet {
			contentReviewTemplateItemMap := map[string]interface{}{}

			if contentReviewTemplateItem.Definition != nil {
				contentReviewTemplateItemMap["definition"] = contentReviewTemplateItem.Definition
			}

			if contentReviewTemplateItem.Name != nil {
				contentReviewTemplateItemMap["name"] = contentReviewTemplateItem.Name
			}

			if contentReviewTemplateItem.Comment != nil {
				contentReviewTemplateItemMap["comment"] = contentReviewTemplateItem.Comment
			}

			if contentReviewTemplateItem.PornConfigure != nil {
				pornConfigureMap := map[string]interface{}{}

				if contentReviewTemplateItem.PornConfigure.ImgReviewInfo != nil {
					imgReviewInfoMap := map[string]interface{}{}

					if contentReviewTemplateItem.PornConfigure.ImgReviewInfo.Switch != nil {
						imgReviewInfoMap["switch"] = contentReviewTemplateItem.PornConfigure.ImgReviewInfo.Switch
					}

					if contentReviewTemplateItem.PornConfigure.ImgReviewInfo.LabelSet != nil {
						imgReviewInfoMap["label_set"] = contentReviewTemplateItem.PornConfigure.ImgReviewInfo.LabelSet
					}

					if contentReviewTemplateItem.PornConfigure.ImgReviewInfo.BlockConfidence != nil {
						imgReviewInfoMap["block_confidence"] = contentReviewTemplateItem.PornConfigure.ImgReviewInfo.BlockConfidence
					}

					if contentReviewTemplateItem.PornConfigure.ImgReviewInfo.ReviewConfidence != nil {
						imgReviewInfoMap["review_confidence"] = contentReviewTemplateItem.PornConfigure.ImgReviewInfo.ReviewConfidence
					}

					pornConfigureMap["img_review_info"] = []interface{}{imgReviewInfoMap}
				}

				if contentReviewTemplateItem.PornConfigure.AsrReviewInfo != nil {
					asrReviewInfoMap := map[string]interface{}{}

					if contentReviewTemplateItem.PornConfigure.AsrReviewInfo.Switch != nil {
						asrReviewInfoMap["switch"] = contentReviewTemplateItem.PornConfigure.AsrReviewInfo.Switch
					}

					if contentReviewTemplateItem.PornConfigure.AsrReviewInfo.BlockConfidence != nil {
						asrReviewInfoMap["block_confidence"] = contentReviewTemplateItem.PornConfigure.AsrReviewInfo.BlockConfidence
					}

					if contentReviewTemplateItem.PornConfigure.AsrReviewInfo.ReviewConfidence != nil {
						asrReviewInfoMap["review_confidence"] = contentReviewTemplateItem.PornConfigure.AsrReviewInfo.ReviewConfidence
					}

					pornConfigureMap["asr_review_info"] = []interface{}{asrReviewInfoMap}
				}

				if contentReviewTemplateItem.PornConfigure.OcrReviewInfo != nil {
					ocrReviewInfoMap := map[string]interface{}{}

					if contentReviewTemplateItem.PornConfigure.OcrReviewInfo.Switch != nil {
						ocrReviewInfoMap["switch"] = contentReviewTemplateItem.PornConfigure.OcrReviewInfo.Switch
					}

					if contentReviewTemplateItem.PornConfigure.OcrReviewInfo.BlockConfidence != nil {
						ocrReviewInfoMap["block_confidence"] = contentReviewTemplateItem.PornConfigure.OcrReviewInfo.BlockConfidence
					}

					if contentReviewTemplateItem.PornConfigure.OcrReviewInfo.ReviewConfidence != nil {
						ocrReviewInfoMap["review_confidence"] = contentReviewTemplateItem.PornConfigure.OcrReviewInfo.ReviewConfidence
					}

					pornConfigureMap["ocr_review_info"] = []interface{}{ocrReviewInfoMap}
				}

				contentReviewTemplateItemMap["porn_configure"] = []interface{}{pornConfigureMap}
			}

			if contentReviewTemplateItem.TerrorismConfigure != nil {
				terrorismConfigureMap := map[string]interface{}{}

				if contentReviewTemplateItem.TerrorismConfigure.ImgReviewInfo != nil {
					imgReviewInfoMap := map[string]interface{}{}

					if contentReviewTemplateItem.TerrorismConfigure.ImgReviewInfo.Switch != nil {
						imgReviewInfoMap["switch"] = contentReviewTemplateItem.TerrorismConfigure.ImgReviewInfo.Switch
					}

					if contentReviewTemplateItem.TerrorismConfigure.ImgReviewInfo.LabelSet != nil {
						imgReviewInfoMap["label_set"] = contentReviewTemplateItem.TerrorismConfigure.ImgReviewInfo.LabelSet
					}

					if contentReviewTemplateItem.TerrorismConfigure.ImgReviewInfo.BlockConfidence != nil {
						imgReviewInfoMap["block_confidence"] = contentReviewTemplateItem.TerrorismConfigure.ImgReviewInfo.BlockConfidence
					}

					if contentReviewTemplateItem.TerrorismConfigure.ImgReviewInfo.ReviewConfidence != nil {
						imgReviewInfoMap["review_confidence"] = contentReviewTemplateItem.TerrorismConfigure.ImgReviewInfo.ReviewConfidence
					}

					terrorismConfigureMap["img_review_info"] = []interface{}{imgReviewInfoMap}
				}

				if contentReviewTemplateItem.TerrorismConfigure.OcrReviewInfo != nil {
					ocrReviewInfoMap := map[string]interface{}{}

					if contentReviewTemplateItem.TerrorismConfigure.OcrReviewInfo.Switch != nil {
						ocrReviewInfoMap["switch"] = contentReviewTemplateItem.TerrorismConfigure.OcrReviewInfo.Switch
					}

					if contentReviewTemplateItem.TerrorismConfigure.OcrReviewInfo.BlockConfidence != nil {
						ocrReviewInfoMap["block_confidence"] = contentReviewTemplateItem.TerrorismConfigure.OcrReviewInfo.BlockConfidence
					}

					if contentReviewTemplateItem.TerrorismConfigure.OcrReviewInfo.ReviewConfidence != nil {
						ocrReviewInfoMap["review_confidence"] = contentReviewTemplateItem.TerrorismConfigure.OcrReviewInfo.ReviewConfidence
					}

					terrorismConfigureMap["ocr_review_info"] = []interface{}{ocrReviewInfoMap}
				}

				contentReviewTemplateItemMap["terrorism_configure"] = []interface{}{terrorismConfigureMap}
			}

			if contentReviewTemplateItem.PoliticalConfigure != nil {
				politicalConfigureMap := map[string]interface{}{}

				if contentReviewTemplateItem.PoliticalConfigure.ImgReviewInfo != nil {
					imgReviewInfoMap := map[string]interface{}{}

					if contentReviewTemplateItem.PoliticalConfigure.ImgReviewInfo.Switch != nil {
						imgReviewInfoMap["switch"] = contentReviewTemplateItem.PoliticalConfigure.ImgReviewInfo.Switch
					}

					if contentReviewTemplateItem.PoliticalConfigure.ImgReviewInfo.LabelSet != nil {
						imgReviewInfoMap["label_set"] = contentReviewTemplateItem.PoliticalConfigure.ImgReviewInfo.LabelSet
					}

					if contentReviewTemplateItem.PoliticalConfigure.ImgReviewInfo.BlockConfidence != nil {
						imgReviewInfoMap["block_confidence"] = contentReviewTemplateItem.PoliticalConfigure.ImgReviewInfo.BlockConfidence
					}

					if contentReviewTemplateItem.PoliticalConfigure.ImgReviewInfo.ReviewConfidence != nil {
						imgReviewInfoMap["review_confidence"] = contentReviewTemplateItem.PoliticalConfigure.ImgReviewInfo.ReviewConfidence
					}

					politicalConfigureMap["img_review_info"] = []interface{}{imgReviewInfoMap}
				}

				if contentReviewTemplateItem.PoliticalConfigure.AsrReviewInfo != nil {
					asrReviewInfoMap := map[string]interface{}{}

					if contentReviewTemplateItem.PoliticalConfigure.AsrReviewInfo.Switch != nil {
						asrReviewInfoMap["switch"] = contentReviewTemplateItem.PoliticalConfigure.AsrReviewInfo.Switch
					}

					if contentReviewTemplateItem.PoliticalConfigure.AsrReviewInfo.BlockConfidence != nil {
						asrReviewInfoMap["block_confidence"] = contentReviewTemplateItem.PoliticalConfigure.AsrReviewInfo.BlockConfidence
					}

					if contentReviewTemplateItem.PoliticalConfigure.AsrReviewInfo.ReviewConfidence != nil {
						asrReviewInfoMap["review_confidence"] = contentReviewTemplateItem.PoliticalConfigure.AsrReviewInfo.ReviewConfidence
					}

					politicalConfigureMap["asr_review_info"] = []interface{}{asrReviewInfoMap}
				}

				if contentReviewTemplateItem.PoliticalConfigure.OcrReviewInfo != nil {
					ocrReviewInfoMap := map[string]interface{}{}

					if contentReviewTemplateItem.PoliticalConfigure.OcrReviewInfo.Switch != nil {
						ocrReviewInfoMap["switch"] = contentReviewTemplateItem.PoliticalConfigure.OcrReviewInfo.Switch
					}

					if contentReviewTemplateItem.PoliticalConfigure.OcrReviewInfo.BlockConfidence != nil {
						ocrReviewInfoMap["block_confidence"] = contentReviewTemplateItem.PoliticalConfigure.OcrReviewInfo.BlockConfidence
					}

					if contentReviewTemplateItem.PoliticalConfigure.OcrReviewInfo.ReviewConfidence != nil {
						ocrReviewInfoMap["review_confidence"] = contentReviewTemplateItem.PoliticalConfigure.OcrReviewInfo.ReviewConfidence
					}

					politicalConfigureMap["ocr_review_info"] = []interface{}{ocrReviewInfoMap}
				}

				contentReviewTemplateItemMap["political_configure"] = []interface{}{politicalConfigureMap}
			}

			if contentReviewTemplateItem.ProhibitedConfigure != nil {
				prohibitedConfigureMap := map[string]interface{}{}

				if contentReviewTemplateItem.ProhibitedConfigure.AsrReviewInfo != nil {
					asrReviewInfoMap := map[string]interface{}{}

					if contentReviewTemplateItem.ProhibitedConfigure.AsrReviewInfo.Switch != nil {
						asrReviewInfoMap["switch"] = contentReviewTemplateItem.ProhibitedConfigure.AsrReviewInfo.Switch
					}

					if contentReviewTemplateItem.ProhibitedConfigure.AsrReviewInfo.BlockConfidence != nil {
						asrReviewInfoMap["block_confidence"] = contentReviewTemplateItem.ProhibitedConfigure.AsrReviewInfo.BlockConfidence
					}

					if contentReviewTemplateItem.ProhibitedConfigure.AsrReviewInfo.ReviewConfidence != nil {
						asrReviewInfoMap["review_confidence"] = contentReviewTemplateItem.ProhibitedConfigure.AsrReviewInfo.ReviewConfidence
					}

					prohibitedConfigureMap["asr_review_info"] = []interface{}{asrReviewInfoMap}
				}

				if contentReviewTemplateItem.ProhibitedConfigure.OcrReviewInfo != nil {
					ocrReviewInfoMap := map[string]interface{}{}

					if contentReviewTemplateItem.ProhibitedConfigure.OcrReviewInfo.Switch != nil {
						ocrReviewInfoMap["switch"] = contentReviewTemplateItem.ProhibitedConfigure.OcrReviewInfo.Switch
					}

					if contentReviewTemplateItem.ProhibitedConfigure.OcrReviewInfo.BlockConfidence != nil {
						ocrReviewInfoMap["block_confidence"] = contentReviewTemplateItem.ProhibitedConfigure.OcrReviewInfo.BlockConfidence
					}

					if contentReviewTemplateItem.ProhibitedConfigure.OcrReviewInfo.ReviewConfidence != nil {
						ocrReviewInfoMap["review_confidence"] = contentReviewTemplateItem.ProhibitedConfigure.OcrReviewInfo.ReviewConfidence
					}

					prohibitedConfigureMap["ocr_review_info"] = []interface{}{ocrReviewInfoMap}
				}

				contentReviewTemplateItemMap["prohibited_configure"] = []interface{}{prohibitedConfigureMap}
			}

			if contentReviewTemplateItem.UserDefineConfigure != nil {
				userDefineConfigureMap := map[string]interface{}{}

				if contentReviewTemplateItem.UserDefineConfigure.FaceReviewInfo != nil {
					faceReviewInfoMap := map[string]interface{}{}

					if contentReviewTemplateItem.UserDefineConfigure.FaceReviewInfo.Switch != nil {
						faceReviewInfoMap["switch"] = contentReviewTemplateItem.UserDefineConfigure.FaceReviewInfo.Switch
					}

					if contentReviewTemplateItem.UserDefineConfigure.FaceReviewInfo.LabelSet != nil {
						faceReviewInfoMap["label_set"] = contentReviewTemplateItem.UserDefineConfigure.FaceReviewInfo.LabelSet
					}

					if contentReviewTemplateItem.UserDefineConfigure.FaceReviewInfo.BlockConfidence != nil {
						faceReviewInfoMap["block_confidence"] = contentReviewTemplateItem.UserDefineConfigure.FaceReviewInfo.BlockConfidence
					}

					if contentReviewTemplateItem.UserDefineConfigure.FaceReviewInfo.ReviewConfidence != nil {
						faceReviewInfoMap["review_confidence"] = contentReviewTemplateItem.UserDefineConfigure.FaceReviewInfo.ReviewConfidence
					}

					userDefineConfigureMap["face_review_info"] = []interface{}{faceReviewInfoMap}
				}

				if contentReviewTemplateItem.UserDefineConfigure.AsrReviewInfo != nil {
					asrReviewInfoMap := map[string]interface{}{}

					if contentReviewTemplateItem.UserDefineConfigure.AsrReviewInfo.Switch != nil {
						asrReviewInfoMap["switch"] = contentReviewTemplateItem.UserDefineConfigure.AsrReviewInfo.Switch
					}

					if contentReviewTemplateItem.UserDefineConfigure.AsrReviewInfo.LabelSet != nil {
						asrReviewInfoMap["label_set"] = contentReviewTemplateItem.UserDefineConfigure.AsrReviewInfo.LabelSet
					}

					if contentReviewTemplateItem.UserDefineConfigure.AsrReviewInfo.BlockConfidence != nil {
						asrReviewInfoMap["block_confidence"] = contentReviewTemplateItem.UserDefineConfigure.AsrReviewInfo.BlockConfidence
					}

					if contentReviewTemplateItem.UserDefineConfigure.AsrReviewInfo.ReviewConfidence != nil {
						asrReviewInfoMap["review_confidence"] = contentReviewTemplateItem.UserDefineConfigure.AsrReviewInfo.ReviewConfidence
					}

					userDefineConfigureMap["asr_review_info"] = []interface{}{asrReviewInfoMap}
				}

				if contentReviewTemplateItem.UserDefineConfigure.OcrReviewInfo != nil {
					ocrReviewInfoMap := map[string]interface{}{}

					if contentReviewTemplateItem.UserDefineConfigure.OcrReviewInfo.Switch != nil {
						ocrReviewInfoMap["switch"] = contentReviewTemplateItem.UserDefineConfigure.OcrReviewInfo.Switch
					}

					if contentReviewTemplateItem.UserDefineConfigure.OcrReviewInfo.LabelSet != nil {
						ocrReviewInfoMap["label_set"] = contentReviewTemplateItem.UserDefineConfigure.OcrReviewInfo.LabelSet
					}

					if contentReviewTemplateItem.UserDefineConfigure.OcrReviewInfo.BlockConfidence != nil {
						ocrReviewInfoMap["block_confidence"] = contentReviewTemplateItem.UserDefineConfigure.OcrReviewInfo.BlockConfidence
					}

					if contentReviewTemplateItem.UserDefineConfigure.OcrReviewInfo.ReviewConfidence != nil {
						ocrReviewInfoMap["review_confidence"] = contentReviewTemplateItem.UserDefineConfigure.OcrReviewInfo.ReviewConfidence
					}

					userDefineConfigureMap["ocr_review_info"] = []interface{}{ocrReviewInfoMap}
				}

				contentReviewTemplateItemMap["user_define_configure"] = []interface{}{userDefineConfigureMap}
			}

			if contentReviewTemplateItem.CreateTime != nil {
				contentReviewTemplateItemMap["create_time"] = contentReviewTemplateItem.CreateTime
			}

			if contentReviewTemplateItem.UpdateTime != nil {
				contentReviewTemplateItemMap["update_time"] = contentReviewTemplateItem.UpdateTime
			}

			if contentReviewTemplateItem.Type != nil {
				contentReviewTemplateItemMap["type"] = contentReviewTemplateItem.Type
			}

			ids = append(ids, *contentReviewTemplateItem.Definition)
			tmpList = append(tmpList, contentReviewTemplateItemMap)
		}

		_ = d.Set("content_review_template_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
