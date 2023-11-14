/*
Use this data source to query detailed information of mps content_review_templates

Example Usage

```hcl
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
				Description: "The content review template uniquely identifies the filter condition, and the array length limit: 50.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page offset, default: 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Return the number of records, default value: 10, maximum value: 100.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template type filter condition, optional value:Preset: system preset template.Custom: user-defined template.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "Total number of records matching filter condition.",
			},

			"content_review_template_set": {
				Type:        schema.TypeList,
				Description: "Content review template details list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Description: "The unique identifier of the content review template.",
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Content review template name, length limit: 64 characters.",
						},
						"comment": {
							Type:        schema.TypeString,
							Description: "Content review template description information, length limit: 256 characters.",
						},
						"porn_configure": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Control parameters for porn image.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"img_review_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "Porn image Identification Control Parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Description: "Porn screen task switch, optional value:ON/OFF.",
												},
												"label_set": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Description: "Porn image filter label, if the review result contains the selected label, the result will be returned. If the filter label is empty, all the review results will be returned. The optional value is:porn, vulgar, intimacy, sexy.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 90 points. Value range: 0~100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 0. Value range: 0~100.",
												},
											},
										},
									},
									"asr_review_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "Voice pornography control parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Description: "Voice pornography task switch, optional value:ON/OFF.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.",
												},
											},
										},
									},
									"ocr_review_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "Ocr pornography control parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Description: "Ocr pornography task switch, optional value:ON/OFF.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.",
												},
											},
										},
									},
								},
							},
						},
						"terrorism_configure": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Control parameters for unsafe information.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"img_review_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "Terrorism image task control parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Description: "Terrorism image task switch, optional value:ON/OFF.",
												},
												"label_set": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Description: "Terrorism image filter tag, if the review result contains the selected tag, the result will be returned, if the filter tag is empty, all the review results will be returned, the optional value is:guns, crowd, bloody, police, banners, militant, explosion, terrorists, scenario.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 90 points. Value range: 0~100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 80 points. Value range: 0~100.",
												},
											},
										},
									},
									"ocr_review_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "Ocr terrorism task Control Parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Description: "Ocr terrorism image task switch, optional value:ON/OFF.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.",
												},
											},
										},
									},
								},
							},
						},
						"political_configure": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Political control parameters.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"img_review_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "Political image control parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Description: "Political image task switch, optional value:ON/OFF.",
												},
												"label_set": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Description: "Political image filter tag, if the review result contains the selected tag, the result will be returned, if the filter tag is empty, all the review results will be returned, the optional value is:violation_photo, politician, entertainment, sport, entrepreneur, scholar, celebrity, military.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 97 points. Value range: 0~100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 95 points. Value range: 0~100.",
												},
											},
										},
									},
									"asr_review_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "Political asr control parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Description: "Political asr task switch, optional value:ON/OFF.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.",
												},
											},
										},
									},
									"ocr_review_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "Political ocr control parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Description: "Political ocr task switch, optional value:ON/OFF.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.",
												},
											},
										},
									},
								},
							},
						},
						"prohibited_configure": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Prohibited control parameters. Prohibited content includes:abuse, drug-related violations.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"asr_review_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "Voice Prohibition Control Parameters.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Description: "Voice Prohibition task switch, optional value:ON/OFF.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.",
												},
											},
										},
									},
									"ocr_review_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "Ocr Prohibition Control Parameters.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Description: "Ocr Prohibition task switch, optional value:ON/OFF.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.",
												},
											},
										},
									},
								},
							},
						},
						"user_define_configure": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "User-Defined Content Moderation Control Parameters.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"face_review_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "User-defined face review control parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Description: "User-defined face review task switch, optional value:ON/OFF.",
												},
												"label_set": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Description: "User-defined face review tags, the review result contains the selected tag and returns the result, if the filter tag is empty, all review results are returned. If you want to use the tag filtering function, when adding a face library, you need to add the corresponding character tag.The maximum number of tags is 10, and the length of each tag is up to 16 characters.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.",
												},
											},
										},
									},
									"asr_review_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "User-defined asr text review control parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Description: "User-defined asr review task switch, optional value:ON/OFF.",
												},
												"label_set": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Description: "User-defined asr tags, the review result contains the selected tag and returns the result, if the filter tag is empty, all review results are returned. If you want to use the tag filtering function, when adding a asr library, you need to add the corresponding character tag.The maximum number of tags is 10, and the length of each tag is up to 16 characters.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.",
												},
											},
										},
									},
									"ocr_review_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Description: "User-defined ocr text review control parameters.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"switch": {
													Type:        schema.TypeString,
													Description: "User-defined ocr text review task switch, optional value:ON/OFF.",
												},
												"label_set": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Description: "User-defined ocr tags, the review result contains the selected tag and returns the result, if the filter tag is empty, all review results are returned. If you want to use the tag filtering function, when adding a ocr library, you need to add the corresponding character tag.The maximum number of tags is 10, and the length of each tag is up to 16 characters.",
												},
												"block_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.",
												},
												"review_confidence": {
													Type:        schema.TypeInt,
													Description: "The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.",
												},
											},
										},
									},
								},
							},
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Template creation time, in [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
						},
						"update_time": {
							Type:        schema.TypeString,
							Description: "Template last modified time, using [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Template type, optional value:Preset: system preset template.Custom: user-defined template.Note: This field may return null, indicating that no valid value can be obtained.",
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

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("content_review_template_set"); ok {
		contentReviewTemplateSetSet := v.([]interface{})
		tmpSet := make([]*mps.ContentReviewTemplateItem, 0, len(contentReviewTemplateSetSet))

		for _, item := range contentReviewTemplateSetSet {
			contentReviewTemplateItem := mps.ContentReviewTemplateItem{}
			contentReviewTemplateItemMap := item.(map[string]interface{})

			if v, ok := contentReviewTemplateItemMap["definition"]; ok {
				contentReviewTemplateItem.Definition = helper.IntInt64(v.(int))
			}
			if v, ok := contentReviewTemplateItemMap["name"]; ok {
				contentReviewTemplateItem.Name = helper.String(v.(string))
			}
			if v, ok := contentReviewTemplateItemMap["comment"]; ok {
				contentReviewTemplateItem.Comment = helper.String(v.(string))
			}
			if pornConfigureMap, ok := helper.InterfaceToMap(contentReviewTemplateItemMap, "porn_configure"); ok {
				pornConfigureInfo := mps.PornConfigureInfo{}
				if imgReviewInfoMap, ok := helper.InterfaceToMap(pornConfigureMap, "img_review_info"); ok {
					pornImgReviewTemplateInfo := mps.PornImgReviewTemplateInfo{}
					if v, ok := imgReviewInfoMap["switch"]; ok {
						pornImgReviewTemplateInfo.Switch = helper.String(v.(string))
					}
					if v, ok := imgReviewInfoMap["label_set"]; ok {
						labelSetSet := v.(*schema.Set).List()
						pornImgReviewTemplateInfo.LabelSet = helper.InterfacesStringsPoint(labelSetSet)
					}
					if v, ok := imgReviewInfoMap["block_confidence"]; ok {
						pornImgReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
					}
					if v, ok := imgReviewInfoMap["review_confidence"]; ok {
						pornImgReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
					}
					pornConfigureInfo.ImgReviewInfo = &pornImgReviewTemplateInfo
				}
				if asrReviewInfoMap, ok := helper.InterfaceToMap(pornConfigureMap, "asr_review_info"); ok {
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
				if ocrReviewInfoMap, ok := helper.InterfaceToMap(pornConfigureMap, "ocr_review_info"); ok {
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
				contentReviewTemplateItem.PornConfigure = &pornConfigureInfo
			}
			if terrorismConfigureMap, ok := helper.InterfaceToMap(contentReviewTemplateItemMap, "terrorism_configure"); ok {
				terrorismConfigureInfo := mps.TerrorismConfigureInfo{}
				if imgReviewInfoMap, ok := helper.InterfaceToMap(terrorismConfigureMap, "img_review_info"); ok {
					terrorismImgReviewTemplateInfo := mps.TerrorismImgReviewTemplateInfo{}
					if v, ok := imgReviewInfoMap["switch"]; ok {
						terrorismImgReviewTemplateInfo.Switch = helper.String(v.(string))
					}
					if v, ok := imgReviewInfoMap["label_set"]; ok {
						labelSetSet := v.(*schema.Set).List()
						terrorismImgReviewTemplateInfo.LabelSet = helper.InterfacesStringsPoint(labelSetSet)
					}
					if v, ok := imgReviewInfoMap["block_confidence"]; ok {
						terrorismImgReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
					}
					if v, ok := imgReviewInfoMap["review_confidence"]; ok {
						terrorismImgReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
					}
					terrorismConfigureInfo.ImgReviewInfo = &terrorismImgReviewTemplateInfo
				}
				if ocrReviewInfoMap, ok := helper.InterfaceToMap(terrorismConfigureMap, "ocr_review_info"); ok {
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
				contentReviewTemplateItem.TerrorismConfigure = &terrorismConfigureInfo
			}
			if politicalConfigureMap, ok := helper.InterfaceToMap(contentReviewTemplateItemMap, "political_configure"); ok {
				politicalConfigureInfo := mps.PoliticalConfigureInfo{}
				if imgReviewInfoMap, ok := helper.InterfaceToMap(politicalConfigureMap, "img_review_info"); ok {
					politicalImgReviewTemplateInfo := mps.PoliticalImgReviewTemplateInfo{}
					if v, ok := imgReviewInfoMap["switch"]; ok {
						politicalImgReviewTemplateInfo.Switch = helper.String(v.(string))
					}
					if v, ok := imgReviewInfoMap["label_set"]; ok {
						labelSetSet := v.(*schema.Set).List()
						politicalImgReviewTemplateInfo.LabelSet = helper.InterfacesStringsPoint(labelSetSet)
					}
					if v, ok := imgReviewInfoMap["block_confidence"]; ok {
						politicalImgReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
					}
					if v, ok := imgReviewInfoMap["review_confidence"]; ok {
						politicalImgReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
					}
					politicalConfigureInfo.ImgReviewInfo = &politicalImgReviewTemplateInfo
				}
				if asrReviewInfoMap, ok := helper.InterfaceToMap(politicalConfigureMap, "asr_review_info"); ok {
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
				if ocrReviewInfoMap, ok := helper.InterfaceToMap(politicalConfigureMap, "ocr_review_info"); ok {
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
				contentReviewTemplateItem.PoliticalConfigure = &politicalConfigureInfo
			}
			if prohibitedConfigureMap, ok := helper.InterfaceToMap(contentReviewTemplateItemMap, "prohibited_configure"); ok {
				prohibitedConfigureInfo := mps.ProhibitedConfigureInfo{}
				if asrReviewInfoMap, ok := helper.InterfaceToMap(prohibitedConfigureMap, "asr_review_info"); ok {
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
				if ocrReviewInfoMap, ok := helper.InterfaceToMap(prohibitedConfigureMap, "ocr_review_info"); ok {
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
				contentReviewTemplateItem.ProhibitedConfigure = &prohibitedConfigureInfo
			}
			if userDefineConfigureMap, ok := helper.InterfaceToMap(contentReviewTemplateItemMap, "user_define_configure"); ok {
				userDefineConfigureInfo := mps.UserDefineConfigureInfo{}
				if faceReviewInfoMap, ok := helper.InterfaceToMap(userDefineConfigureMap, "face_review_info"); ok {
					userDefineFaceReviewTemplateInfo := mps.UserDefineFaceReviewTemplateInfo{}
					if v, ok := faceReviewInfoMap["switch"]; ok {
						userDefineFaceReviewTemplateInfo.Switch = helper.String(v.(string))
					}
					if v, ok := faceReviewInfoMap["label_set"]; ok {
						labelSetSet := v.(*schema.Set).List()
						userDefineFaceReviewTemplateInfo.LabelSet = helper.InterfacesStringsPoint(labelSetSet)
					}
					if v, ok := faceReviewInfoMap["block_confidence"]; ok {
						userDefineFaceReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
					}
					if v, ok := faceReviewInfoMap["review_confidence"]; ok {
						userDefineFaceReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
					}
					userDefineConfigureInfo.FaceReviewInfo = &userDefineFaceReviewTemplateInfo
				}
				if asrReviewInfoMap, ok := helper.InterfaceToMap(userDefineConfigureMap, "asr_review_info"); ok {
					userDefineAsrTextReviewTemplateInfo := mps.UserDefineAsrTextReviewTemplateInfo{}
					if v, ok := asrReviewInfoMap["switch"]; ok {
						userDefineAsrTextReviewTemplateInfo.Switch = helper.String(v.(string))
					}
					if v, ok := asrReviewInfoMap["label_set"]; ok {
						labelSetSet := v.(*schema.Set).List()
						userDefineAsrTextReviewTemplateInfo.LabelSet = helper.InterfacesStringsPoint(labelSetSet)
					}
					if v, ok := asrReviewInfoMap["block_confidence"]; ok {
						userDefineAsrTextReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
					}
					if v, ok := asrReviewInfoMap["review_confidence"]; ok {
						userDefineAsrTextReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
					}
					userDefineConfigureInfo.AsrReviewInfo = &userDefineAsrTextReviewTemplateInfo
				}
				if ocrReviewInfoMap, ok := helper.InterfaceToMap(userDefineConfigureMap, "ocr_review_info"); ok {
					userDefineOcrTextReviewTemplateInfo := mps.UserDefineOcrTextReviewTemplateInfo{}
					if v, ok := ocrReviewInfoMap["switch"]; ok {
						userDefineOcrTextReviewTemplateInfo.Switch = helper.String(v.(string))
					}
					if v, ok := ocrReviewInfoMap["label_set"]; ok {
						labelSetSet := v.(*schema.Set).List()
						userDefineOcrTextReviewTemplateInfo.LabelSet = helper.InterfacesStringsPoint(labelSetSet)
					}
					if v, ok := ocrReviewInfoMap["block_confidence"]; ok {
						userDefineOcrTextReviewTemplateInfo.BlockConfidence = helper.IntInt64(v.(int))
					}
					if v, ok := ocrReviewInfoMap["review_confidence"]; ok {
						userDefineOcrTextReviewTemplateInfo.ReviewConfidence = helper.IntInt64(v.(int))
					}
					userDefineConfigureInfo.OcrReviewInfo = &userDefineOcrTextReviewTemplateInfo
				}
				contentReviewTemplateItem.UserDefineConfigure = &userDefineConfigureInfo
			}
			if v, ok := contentReviewTemplateItemMap["create_time"]; ok {
				contentReviewTemplateItem.CreateTime = helper.String(v.(string))
			}
			if v, ok := contentReviewTemplateItemMap["update_time"]; ok {
				contentReviewTemplateItem.UpdateTime = helper.String(v.(string))
			}
			if v, ok := contentReviewTemplateItemMap["type"]; ok {
				contentReviewTemplateItem.Type = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &contentReviewTemplateItem)
		}
		paramMap["content_review_template_set"] = tmpSet
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

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
