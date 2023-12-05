package tencentcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCssPadTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCssPadTemplatesRead,
		Schema: map[string]*schema.Schema{
			"templates": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Live pad template information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Template id.",
						},
						"template_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template name.",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pad content.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template create time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template modify time.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description info.",
						},
						"wait_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Stream interruption waiting time.Value range: 0-30000.Unit: milliseconds.",
						},
						"max_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum pad duration.Value range: 0 - positive infinity.Unit: milliseconds.",
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Pad content type: 1: Image, 2: Video. Default value: 1.",
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

func dataSourceTencentCloudCssPadTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_css_pad_templates.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	var templates []*css.PadTemplate
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCssPadTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		templates = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(templates))
	tmpList := make([]map[string]interface{}, 0, len(templates))
	if templates != nil {
		for _, padTemplate := range templates {
			padTemplateMap := map[string]interface{}{}

			if padTemplate.TemplateId != nil {
				padTemplateMap["template_id"] = padTemplate.TemplateId
			}

			if padTemplate.TemplateName != nil {
				padTemplateMap["template_name"] = padTemplate.TemplateName
			}

			if padTemplate.Url != nil {
				padTemplateMap["url"] = padTemplate.Url
			}

			if padTemplate.CreateTime != nil {
				padTemplateMap["create_time"] = padTemplate.CreateTime
			}

			if padTemplate.UpdateTime != nil {
				padTemplateMap["update_time"] = padTemplate.UpdateTime
			}

			if padTemplate.Description != nil {
				padTemplateMap["description"] = padTemplate.Description
			}

			if padTemplate.WaitDuration != nil {
				padTemplateMap["wait_duration"] = padTemplate.WaitDuration
			}

			if padTemplate.MaxDuration != nil {
				padTemplateMap["max_duration"] = padTemplate.MaxDuration
			}

			if padTemplate.Type != nil {
				padTemplateMap["type"] = padTemplate.Type
			}

			ids = append(ids, fmt.Sprintf("%d", *padTemplate.TemplateId))
			tmpList = append(tmpList, padTemplateMap)
		}

		_ = d.Set("templates", tmpList)
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
