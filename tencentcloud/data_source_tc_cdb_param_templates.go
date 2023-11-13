/*
Use this data source to query detailed information of cdb param_templates

Example Usage

```hcl
data "tencentcloud_cdb_param_templates" "param_templates" {
  engine_versions = &lt;nil&gt;
  engine_types = &lt;nil&gt;
  template_names = &lt;nil&gt;
  template_ids = &lt;nil&gt;
  total_count = &lt;nil&gt;
  items {
		template_id = &lt;nil&gt;
		name = &lt;nil&gt;
		description = &lt;nil&gt;
		engine_version = &lt;nil&gt;
		template_type = &lt;nil&gt;

  }
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCdbParamTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbParamTemplatesRead,
		Schema: map[string]*schema.Schema{
			"engine_versions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Engine version, the default is to query all.",
			},

			"engine_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Engine type, the default is to query all.",
			},

			"template_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Template name list, the default is to query all.",
			},

			"template_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Template id list, the default is to query all.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "Number of parameter templates.",
			},

			"items": {
				Type:        schema.TypeList,
				Description: "Parameter template details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": {
							Type:        schema.TypeInt,
							Description: "Parameter template ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Parameter template name.",
						},
						"description": {
							Type:        schema.TypeString,
							Description: "Parameter template description.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Description: "Instance engine version.",
						},
						"template_type": {
							Type:        schema.TypeString,
							Description: "Parameter template type.",
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

func dataSourceTencentCloudCdbParamTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_param_templates.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("engine_versions"); ok {
		engineVersionsSet := v.(*schema.Set).List()
		paramMap["EngineVersions"] = helper.InterfacesStringsPoint(engineVersionsSet)
	}

	if v, ok := d.GetOk("engine_types"); ok {
		engineTypesSet := v.(*schema.Set).List()
		paramMap["EngineTypes"] = helper.InterfacesStringsPoint(engineTypesSet)
	}

	if v, ok := d.GetOk("template_names"); ok {
		templateNamesSet := v.(*schema.Set).List()
		paramMap["TemplateNames"] = helper.InterfacesStringsPoint(templateNamesSet)
	}

	if v, ok := d.GetOk("template_ids"); ok {
		templateIdsSet := v.(*schema.Set).List()
		for i := range templateIdsSet {
			templateIds := templateIdsSet[i].(int)
			paramMap["TemplateIds"] = append(paramMap["TemplateIds"], helper.IntInt64(templateIds))
		}
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("items"); ok {
		itemsSet := v.([]interface{})
		tmpSet := make([]*cdb.ParamTemplateInfo, 0, len(itemsSet))

		for _, item := range itemsSet {
			paramTemplateInfo := cdb.ParamTemplateInfo{}
			paramTemplateInfoMap := item.(map[string]interface{})

			if v, ok := paramTemplateInfoMap["template_id"]; ok {
				paramTemplateInfo.TemplateId = helper.IntInt64(v.(int))
			}
			if v, ok := paramTemplateInfoMap["name"]; ok {
				paramTemplateInfo.Name = helper.String(v.(string))
			}
			if v, ok := paramTemplateInfoMap["description"]; ok {
				paramTemplateInfo.Description = helper.String(v.(string))
			}
			if v, ok := paramTemplateInfoMap["engine_version"]; ok {
				paramTemplateInfo.EngineVersion = helper.String(v.(string))
			}
			if v, ok := paramTemplateInfoMap["template_type"]; ok {
				paramTemplateInfo.TemplateType = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &paramTemplateInfo)
		}
		paramMap["items"] = tmpSet
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var items []*cdb.ParamTemplateInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbParamTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
