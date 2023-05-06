/*
Provide a datasource to query default mysql parameters.

Example Usage

```hcl
resource "tencentcloud_mysql_default_params" "mysql_57" {
	db_version = "5.7"
}
```


*/
package tencentcloud

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func datasourceTencentCloudMysqlDefaultParams() *schema.Resource {
	return &schema.Resource{
		Read: datasourceTencentCloudMysqlDefaultParamsRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"db_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "MySQL database version.",
			},
			//"template_type": {
			//	Type:        schema.TypeString,
			//	Optional:    true,
			//	Description: "",
			//},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used for save results.",
			},
			"param_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of param detail.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Param current value.",
						},
						"default": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Param default value.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Param description.",
						},
						"enum_value": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Params available values if type of param is enum.",
						},
						"max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Param maximum value if type of param is integer.",
						},
						"min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Param minimum value if type of param is integer.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Param key name.",
						},
						"need_reboot": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates weather the database instance reboot if param modified.",
						},
						"param_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of param.",
						},
					},
				},
			},
		},
	}
}

func datasourceTencentCloudMysqlDefaultParamsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("datasource.tencentcloud_mysql_default_params.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	service := MysqlService{client: client}

	engineVersion := d.Get("db_version").(string)

	params, err := service.DescribeDefaultParameters(ctx, engineVersion)

	if err != nil {
		d.SetId("")
		return err
	}

	d.SetId(fmt.Sprintf("mysql-params-%d", rand.Intn(1000)))

	paramList := make([]map[string]interface{}, 0, len(params))

	for i := range params {
		item := params[i]
		param := map[string]interface{}{
			"current_value": item.CurrentValue,
			"default":       item.Default,
			"description":   item.Description,
			"max":           item.Max,
			"min":           item.Min,
			"name":          item.Name,
			"need_reboot":   item.NeedReboot,
			"param_type":    item.ParamType,
		}

		if item.EnumValue != nil {
			param["enum_value"] = helper.StringsInterfaces(item.EnumValue)
		}
		paramList = append(paramList, param)
	}

	if len(paramList) > 0 {
		err = d.Set("param_list", paramList)
		if err != nil {
			return err
		}
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), params); e != nil {
			return e
		}
	}

	return nil
}
