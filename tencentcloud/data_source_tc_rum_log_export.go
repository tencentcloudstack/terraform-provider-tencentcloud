/*
Use this data source to query detailed information of rum log_export

Example Usage

```hcl
data "tencentcloud_rum_log_export" "log_export" {
  name = "log"
  start_time = "1692594840000"
  query = "id:123 AND type:&quot;log&quot;"
  end_time = "1692609240000"
  i_d = 1
  fields =
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudRumLogExport() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRumLogExportRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Export flag name.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start timestamp, in milliseconds.",
			},

			"query": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Log Query syntax statement.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End timestamp, in milliseconds.",
			},

			"i_d": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"fields": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Log fields.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Return result.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudRumLogExportRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_rum_log_export.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		paramMap["Name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("query"); ok {
		paramMap["Query"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("i_d"); v != nil {
		paramMap["ID"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("fields"); ok {
		fieldsSet := v.(*schema.Set).List()
		paramMap["Fields"] = helper.InterfacesStringsPoint(fieldsSet)
	}

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRumLogExportByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		_ = d.Set("result", result)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
