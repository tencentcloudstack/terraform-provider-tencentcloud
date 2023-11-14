/*
Use this data source to query detailed information of mps ai_recognition_templates

Example Usage

```hcl
data "tencentcloud_mps_ai_recognition_templates" "ai_recognition_templates" {
  definitions = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  a_i_recognition_template_set {
		definition = &lt;nil&gt;
		name = &lt;nil&gt;
		comment = &lt;nil&gt;
		face_configure {
			switch = &lt;nil&gt;
			score =
			default_library_label_set = &lt;nil&gt;
			user_define_library_label_set = &lt;nil&gt;
			face_library = "All"
		}
		ocr_full_text_configure {
			switch = &lt;nil&gt;
		}
		ocr_words_configure {
			switch = &lt;nil&gt;
			label_set = &lt;nil&gt;
		}
		asr_full_text_configure {
			switch = &lt;nil&gt;
			subtitle_format = &lt;nil&gt;
		}
		asr_words_configure {
			switch = &lt;nil&gt;
			label_set = &lt;nil&gt;
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

func dataSourceTencentCloudMpsAiRecognitionTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsAiRecognitionTemplatesRead,
		Schema: map[string]*schema.Schema{
			"definitions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Ai recognition template uniquely identifies filter conditions, array length limit: 10.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page offset, default: 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Return the number of records, default value: 10, maximum value: 50.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template type filter condition, return all if not filled, optional value:Preset: system preset template.Custom: user-defined template.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "Total number of records matching filter condition.",
			},

			"a_i_recognition_template_set": {
				Type:        schema.TypeList,
				Description: "Ai recognition template details list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Description: "The unique identifier of the ai recognition template.",
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Ai recognition template name.",
						},
						"comment": {
							Type:        schema.TypeString,
							Description: "The description information of ai recognition template.",
						},
						"face_configure": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Face recognition control parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Description: "Ai face recognition task switch, optional value:ON/OFF.",
									},
									"score": {
										Type:        schema.TypeFloat,
										Description: "Face recognition filter score, when the recognition result reaches the score above, the recognition result will be returned. The default is 95 points. Value range: 0 - 100.",
									},
									"default_library_label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "Default face filter tag, specify the tag of the default face that needs to be returned. If not filled or empty, all default face results will be returned. Label optional value:entertainment, sport, politician.",
									},
									"user_define_library_label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "User-defined face filter tag, specify the tag of the user-defined face that needs to be returned. If not filled or empty, all custom face results will be returned.The maximum number of tags is 100, and the length of each tag is up to 16 characters.",
									},
									"face_library": {
										Type:        schema.TypeString,
										Description: "Face library selection, optional value:Default, UserDefine, AllDefault value: All, use the system default face library and user-defined face library.",
									},
								},
							},
						},
						"ocr_full_text_configure": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Ocr full text control parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Description: "Ocr full text recognition task switch, optional value:ON/OFF.",
									},
								},
							},
						},
						"ocr_words_configure": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Ocr words recognition control parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Description: "Ocr words recognition task switch, optional value:ON/OFF.",
									},
									"label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "Keyword filter label, specify the label of the keyword to be returned. If not filled or empty, all results will be returned.The maximum number of tags is 10, and the length of each tag is up to 16 characters.",
									},
								},
							},
						},
						"asr_full_text_configure": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Asr full text recognition control parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Description: "Asr full text recognition task switch, optional value:ON/OFF.",
									},
									"subtitle_format": {
										Type:        schema.TypeString,
										Description: "Generated subtitle file format, if left blank or blank string means no subtitle file will be generated, optional value:vtt: Generate WebVTT subtitle files.",
									},
								},
							},
						},
						"asr_words_configure": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Asr word recognition control parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Description: "Asr word recognition task switch, optional value:ON/OFF.",
									},
									"label_set": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "Keyword filter label, specify the label of the keyword to be returned. If not filled or empty, all results will be returned.The maximum number of tags is 10, and the length of each tag is up to 16 characters.",
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

func dataSourceTencentCloudMpsAiRecognitionTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_ai_recognition_templates.read")()
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

	if v, ok := d.GetOk("a_i_recognition_template_set"); ok {
		aIRecognitionTemplateSetSet := v.([]interface{})
		tmpSet := make([]*mps.AIRecognitionTemplateItem, 0, len(aIRecognitionTemplateSetSet))

		for _, item := range aIRecognitionTemplateSetSet {
			aIRecognitionTemplateItem := mps.AIRecognitionTemplateItem{}
			aIRecognitionTemplateItemMap := item.(map[string]interface{})

			if v, ok := aIRecognitionTemplateItemMap["definition"]; ok {
				aIRecognitionTemplateItem.Definition = helper.IntInt64(v.(int))
			}
			if v, ok := aIRecognitionTemplateItemMap["name"]; ok {
				aIRecognitionTemplateItem.Name = helper.String(v.(string))
			}
			if v, ok := aIRecognitionTemplateItemMap["comment"]; ok {
				aIRecognitionTemplateItem.Comment = helper.String(v.(string))
			}
			if faceConfigureMap, ok := helper.InterfaceToMap(aIRecognitionTemplateItemMap, "face_configure"); ok {
				faceConfigureInfo := mps.FaceConfigureInfo{}
				if v, ok := faceConfigureMap["switch"]; ok {
					faceConfigureInfo.Switch = helper.String(v.(string))
				}
				if v, ok := faceConfigureMap["score"]; ok {
					faceConfigureInfo.Score = helper.Float64(v.(float64))
				}
				if v, ok := faceConfigureMap["default_library_label_set"]; ok {
					defaultLibraryLabelSetSet := v.(*schema.Set).List()
					faceConfigureInfo.DefaultLibraryLabelSet = helper.InterfacesStringsPoint(defaultLibraryLabelSetSet)
				}
				if v, ok := faceConfigureMap["user_define_library_label_set"]; ok {
					userDefineLibraryLabelSetSet := v.(*schema.Set).List()
					faceConfigureInfo.UserDefineLibraryLabelSet = helper.InterfacesStringsPoint(userDefineLibraryLabelSetSet)
				}
				if v, ok := faceConfigureMap["face_library"]; ok {
					faceConfigureInfo.FaceLibrary = helper.String(v.(string))
				}
				aIRecognitionTemplateItem.FaceConfigure = &faceConfigureInfo
			}
			if ocrFullTextConfigureMap, ok := helper.InterfaceToMap(aIRecognitionTemplateItemMap, "ocr_full_text_configure"); ok {
				ocrFullTextConfigureInfo := mps.OcrFullTextConfigureInfo{}
				if v, ok := ocrFullTextConfigureMap["switch"]; ok {
					ocrFullTextConfigureInfo.Switch = helper.String(v.(string))
				}
				aIRecognitionTemplateItem.OcrFullTextConfigure = &ocrFullTextConfigureInfo
			}
			if ocrWordsConfigureMap, ok := helper.InterfaceToMap(aIRecognitionTemplateItemMap, "ocr_words_configure"); ok {
				ocrWordsConfigureInfo := mps.OcrWordsConfigureInfo{}
				if v, ok := ocrWordsConfigureMap["switch"]; ok {
					ocrWordsConfigureInfo.Switch = helper.String(v.(string))
				}
				if v, ok := ocrWordsConfigureMap["label_set"]; ok {
					labelSetSet := v.(*schema.Set).List()
					ocrWordsConfigureInfo.LabelSet = helper.InterfacesStringsPoint(labelSetSet)
				}
				aIRecognitionTemplateItem.OcrWordsConfigure = &ocrWordsConfigureInfo
			}
			if asrFullTextConfigureMap, ok := helper.InterfaceToMap(aIRecognitionTemplateItemMap, "asr_full_text_configure"); ok {
				asrFullTextConfigureInfo := mps.AsrFullTextConfigureInfo{}
				if v, ok := asrFullTextConfigureMap["switch"]; ok {
					asrFullTextConfigureInfo.Switch = helper.String(v.(string))
				}
				if v, ok := asrFullTextConfigureMap["subtitle_format"]; ok {
					asrFullTextConfigureInfo.SubtitleFormat = helper.String(v.(string))
				}
				aIRecognitionTemplateItem.AsrFullTextConfigure = &asrFullTextConfigureInfo
			}
			if asrWordsConfigureMap, ok := helper.InterfaceToMap(aIRecognitionTemplateItemMap, "asr_words_configure"); ok {
				asrWordsConfigureInfo := mps.AsrWordsConfigureInfo{}
				if v, ok := asrWordsConfigureMap["switch"]; ok {
					asrWordsConfigureInfo.Switch = helper.String(v.(string))
				}
				if v, ok := asrWordsConfigureMap["label_set"]; ok {
					labelSetSet := v.(*schema.Set).List()
					asrWordsConfigureInfo.LabelSet = helper.InterfacesStringsPoint(labelSetSet)
				}
				aIRecognitionTemplateItem.AsrWordsConfigure = &asrWordsConfigureInfo
			}
			if v, ok := aIRecognitionTemplateItemMap["create_time"]; ok {
				aIRecognitionTemplateItem.CreateTime = helper.String(v.(string))
			}
			if v, ok := aIRecognitionTemplateItemMap["update_time"]; ok {
				aIRecognitionTemplateItem.UpdateTime = helper.String(v.(string))
			}
			if v, ok := aIRecognitionTemplateItemMap["type"]; ok {
				aIRecognitionTemplateItem.Type = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &aIRecognitionTemplateItem)
		}
		paramMap["a_i_recognition_template_set"] = tmpSet
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var aIRecognitionTemplateSet []*mps.AIRecognitionTemplateItem

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsAiRecognitionTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		aIRecognitionTemplateSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(aIRecognitionTemplateSet))
	tmpList := make([]map[string]interface{}, 0, len(aIRecognitionTemplateSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
