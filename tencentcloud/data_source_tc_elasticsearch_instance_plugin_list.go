/*
Use this data source to query detailed information of elasticsearch instance plugin list

Example Usage

```hcl
data "tencentcloud_elasticsearch_instance_plugin_list" "instance_plugin_list" {
  instance_id = "es-xxxxxx"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudElasticsearchInstancePluginList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudElasticsearchInstancePluginListRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "order field. Valid values: `pluginName`.",
			},

			"order_by_type": {
				Optional: true,
				Type:     schema.TypeString,
				Description: "Order type. Valid values:\n" +
					"- asc: Ascending asc\n" +
					"- desc: Descending Desc.",
			},

			"plugin_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Plugin type. Valid values: `0`: System plugin.",
			},

			"plugin_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Plugin information list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plugin_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plugin name.",
						},
						"plugin_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plugin version.",
						},
						"plugin_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plugin description.",
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "Plugin status. Valid values:\n" +
								"- `-2` has been uninstalled\n" +
								"- `-1` has been installed in\n" +
								"- `0` installation.",
						},
						"removable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the plug-in can be uninstalled.",
						},
						"plugin_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Plugin type. Valid values: `0`: System plugin.",
						},
						"plugin_update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Plugin update time.",
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

func dataSourceTencentCloudElasticsearchInstancePluginListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_elasticsearch_instance_plugin_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("plugin_type"); v != nil {
		paramMap["PluginType"] = helper.IntInt64(v.(int))
	}

	service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}

	var pluginList []*elasticsearch.DescribeInstancePluginInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeElasticsearchInstancePluginListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		pluginList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(pluginList))
	tmpList := make([]map[string]interface{}, 0, len(pluginList))

	if pluginList != nil {
		for _, describeInstancePluginInfo := range pluginList {
			describeInstancePluginInfoMap := map[string]interface{}{}

			if describeInstancePluginInfo.PluginName != nil {
				describeInstancePluginInfoMap["plugin_name"] = describeInstancePluginInfo.PluginName
			}

			if describeInstancePluginInfo.PluginVersion != nil {
				describeInstancePluginInfoMap["plugin_version"] = describeInstancePluginInfo.PluginVersion
			}

			if describeInstancePluginInfo.PluginDesc != nil {
				describeInstancePluginInfoMap["plugin_desc"] = describeInstancePluginInfo.PluginDesc
			}

			if describeInstancePluginInfo.Status != nil {
				describeInstancePluginInfoMap["status"] = describeInstancePluginInfo.Status
			}

			if describeInstancePluginInfo.Removable != nil {
				describeInstancePluginInfoMap["removable"] = describeInstancePluginInfo.Removable
			}

			if describeInstancePluginInfo.PluginType != nil {
				describeInstancePluginInfoMap["plugin_type"] = describeInstancePluginInfo.PluginType
			}

			if describeInstancePluginInfo.PluginUpdateTime != nil {
				describeInstancePluginInfoMap["plugin_update_time"] = describeInstancePluginInfo.PluginUpdateTime
			}

			ids = append(ids, *describeInstancePluginInfo.PluginName)
			tmpList = append(tmpList, describeInstancePluginInfoMap)
		}

		_ = d.Set("plugin_list", tmpList)
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
