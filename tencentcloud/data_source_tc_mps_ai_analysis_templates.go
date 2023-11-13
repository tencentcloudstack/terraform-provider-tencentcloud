/*
Use this data source to query detailed information of mps ai_analysis_templates

Example Usage

```hcl
data "tencentcloud_mps_ai_analysis_templates" "ai_analysis_templates" {
  definitions = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  a_i_analysis_template_set {
		definition = &lt;nil&gt;
		name = &lt;nil&gt;
		comment = &lt;nil&gt;
		classification_configure {
			switch = &lt;nil&gt;
		}
		tag_configure {
			switch = &lt;nil&gt;
		}
		cover_configure {
			switch = &lt;nil&gt;
		}
		frame_tag_configure {
			switch = &lt;nil&gt;
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

func dataSourceTencentCloudMpsAiAnalysisTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsAiAnalysisTemplatesRead,
		Schema: map[string]*schema.Schema{
			"definitions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Ai analysis template uniquely identifies filter conditions, array length limit: 10.",
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

			"a_i_analysis_template_set": {
				Type:        schema.TypeList,
				Description: "Ai analysis template details list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Description: "The unique identifier of the ai analysis template.",
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Ai analysis template name.",
						},
						"comment": {
							Type:        schema.TypeString,
							Description: "The description information of ai analysis template.",
						},
						"classification_configure": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Ai classification task control parameters.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Description: "Ai classification task switch, optional value:ON/OFF.",
									},
								},
							},
						},
						"tag_configure": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Ai tag task control parameters.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Description: "Ai tag task switch, optional value:ON/OFF.",
									},
								},
							},
						},
						"cover_configure": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Ai cover task control parameters.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Description: "Ai cover task switch, optional value:ON/OFF.",
									},
								},
							},
						},
						"frame_tag_configure": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Ai frame tag task control parameters.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Type:        schema.TypeString,
										Description: "Ai frame tag task switch, optional value:ON/OFF.",
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

func dataSourceTencentCloudMpsAiAnalysisTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_ai_analysis_templates.read")()
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

	if v, ok := d.GetOk("a_i_analysis_template_set"); ok {
		aIAnalysisTemplateSetSet := v.([]interface{})
		tmpSet := make([]*mps.AIAnalysisTemplateItem, 0, len(aIAnalysisTemplateSetSet))

		for _, item := range aIAnalysisTemplateSetSet {
			aIAnalysisTemplateItem := mps.AIAnalysisTemplateItem{}
			aIAnalysisTemplateItemMap := item.(map[string]interface{})

			if v, ok := aIAnalysisTemplateItemMap["definition"]; ok {
				aIAnalysisTemplateItem.Definition = helper.IntInt64(v.(int))
			}
			if v, ok := aIAnalysisTemplateItemMap["name"]; ok {
				aIAnalysisTemplateItem.Name = helper.String(v.(string))
			}
			if v, ok := aIAnalysisTemplateItemMap["comment"]; ok {
				aIAnalysisTemplateItem.Comment = helper.String(v.(string))
			}
			if classificationConfigureMap, ok := helper.InterfaceToMap(aIAnalysisTemplateItemMap, "classification_configure"); ok {
				classificationConfigureInfo := mps.ClassificationConfigureInfo{}
				if v, ok := classificationConfigureMap["switch"]; ok {
					classificationConfigureInfo.Switch = helper.String(v.(string))
				}
				aIAnalysisTemplateItem.ClassificationConfigure = &classificationConfigureInfo
			}
			if tagConfigureMap, ok := helper.InterfaceToMap(aIAnalysisTemplateItemMap, "tag_configure"); ok {
				tagConfigureInfo := mps.TagConfigureInfo{}
				if v, ok := tagConfigureMap["switch"]; ok {
					tagConfigureInfo.Switch = helper.String(v.(string))
				}
				aIAnalysisTemplateItem.TagConfigure = &tagConfigureInfo
			}
			if coverConfigureMap, ok := helper.InterfaceToMap(aIAnalysisTemplateItemMap, "cover_configure"); ok {
				coverConfigureInfo := mps.CoverConfigureInfo{}
				if v, ok := coverConfigureMap["switch"]; ok {
					coverConfigureInfo.Switch = helper.String(v.(string))
				}
				aIAnalysisTemplateItem.CoverConfigure = &coverConfigureInfo
			}
			if frameTagConfigureMap, ok := helper.InterfaceToMap(aIAnalysisTemplateItemMap, "frame_tag_configure"); ok {
				frameTagConfigureInfo := mps.FrameTagConfigureInfo{}
				if v, ok := frameTagConfigureMap["switch"]; ok {
					frameTagConfigureInfo.Switch = helper.String(v.(string))
				}
				aIAnalysisTemplateItem.FrameTagConfigure = &frameTagConfigureInfo
			}
			if v, ok := aIAnalysisTemplateItemMap["create_time"]; ok {
				aIAnalysisTemplateItem.CreateTime = helper.String(v.(string))
			}
			if v, ok := aIAnalysisTemplateItemMap["update_time"]; ok {
				aIAnalysisTemplateItem.UpdateTime = helper.String(v.(string))
			}
			if v, ok := aIAnalysisTemplateItemMap["type"]; ok {
				aIAnalysisTemplateItem.Type = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &aIAnalysisTemplateItem)
		}
		paramMap["a_i_analysis_template_set"] = tmpSet
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var aIAnalysisTemplateSet []*mps.AIAnalysisTemplateItem

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsAiAnalysisTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		aIAnalysisTemplateSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(aIAnalysisTemplateSet))
	tmpList := make([]map[string]interface{}, 0, len(aIAnalysisTemplateSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
