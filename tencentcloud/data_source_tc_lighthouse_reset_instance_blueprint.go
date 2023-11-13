/*
Use this data source to query detailed information of lighthouse reset_instance_blueprint

Example Usage

```hcl
data "tencentcloud_lighthouse_reset_instance_blueprint" "reset_instance_blueprint" {
  instance_id = "lhins-123456"
  offset = 0
  limit = 20
  filters {
		name = "blueprint-id"
		values =

  }
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLighthouseResetInstanceBlueprint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseResetInstanceBlueprintRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Offset. Default value is 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Number of returned results. Default value is 20. Maximum value is 100.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter listblueprint-idFilter by image ID.Type: StringRequired: noblueprint-typeFilter by image type.Valid values: APP_OS: application image; PURE_OS: system image; PRIVATE: custom imageType: StringRequired: noplatform-typeFilter by image platform type.Valid values: LINUX_UNIX: Linux or Unix; WINDOWS: WindowsType: StringRequired: noblueprint-nameFilter by image name.Type: StringRequired: noblueprint-stateFilter by image status.Type: StringRequired: noEach request can contain up to 10 Filters and 5 Filter.Values. BlueprintIds and Filters cannot be specified at the same time.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter value of field.",
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

func dataSourceTencentCloudLighthouseResetInstanceBlueprintRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_lighthouse_reset_instance_blueprint.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*lighthouse.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := lighthouse.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var resetInstanceBlueprintSet []*lighthouse.ResetInstanceBlueprint

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseResetInstanceBlueprintByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		resetInstanceBlueprintSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(resetInstanceBlueprintSet))
	tmpList := make([]map[string]interface{}, 0, len(resetInstanceBlueprintSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
