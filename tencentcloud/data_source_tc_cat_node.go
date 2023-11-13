/*
Use this data source to query detailed information of cat node

Example Usage

```hcl
data "tencentcloud_cat_node" "node" {
  node_type =
  location =
  is_i_pv6 =
  node_name = ""
  pay_mode =
  task_type =
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCatNode() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCatNodeRead,
		Schema: map[string]*schema.Schema{
			"node_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Node type 1:IDC,2:LastMile,3:Mobile.",
			},

			"location": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Node area:1=Chinese Mainland,2=Hong Kong, Macao and Taiwan,3=Overseas.",
			},

			"is_i_pv6": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Is IPv6.",
			},

			"node_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Node name.",
			},

			"pay_mode": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Payment mode:1=Trial version,2=Paid version.",
			},

			"task_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Task type: 1=Page performance, 2=File upload, 3=File download, 4=Port performance, 5=Net quality, 6=Audiovisual experience.",
			},

			"node_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Node List. Note: This field may return null, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node name.",
						},
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node code.",
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node type 1:IDC,2:LastMile,3:Mobile.",
						},
						"net_service": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Internet service provider.",
						},
						"district": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "District.",
						},
						"city": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "City.",
						},
						"i_p_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IP Type. 1=IPv4, 2=IPv6. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"location": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node area:1=Chinese Mainland,2=Hong Kong, Macao and Taiwan,3=Overseas. Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"code_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node type. If it is &amp;#39;base&amp;#39;, it is an availability probe point; if it is empty, it is an advanced probe point. Note: This field may return null, indicating that no valid value was found.",
						},
						"task_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "Task type: 1=Page performance, 2=File upload, 3=File download, 4=Port performance, 5=Net quality, 6=Audiovisual experience. Note: This field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudCatNodeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cat_node.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("node_type"); v != nil {
		paramMap["NodeType"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("location"); v != nil {
		paramMap["Location"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("is_i_pv6"); v != nil {
		paramMap["IsIPv6"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("node_name"); ok {
		paramMap["NodeName"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("pay_mode"); v != nil {
		paramMap["PayMode"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("task_type"); v != nil {
		paramMap["TaskType"] = helper.IntInt64(v.(int))
	}

	service := CatService{client: meta.(*TencentCloudClient).apiV3Conn}

	var nodeSet []*cat.NodeDefineExt

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCatNodeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		nodeSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(nodeSet))
	tmpList := make([]map[string]interface{}, 0, len(nodeSet))

	if nodeSet != nil {
		for _, nodeDefineExt := range nodeSet {
			nodeDefineExtMap := map[string]interface{}{}

			if nodeDefineExt.Name != nil {
				nodeDefineExtMap["name"] = nodeDefineExt.Name
			}

			if nodeDefineExt.Code != nil {
				nodeDefineExtMap["code"] = nodeDefineExt.Code
			}

			if nodeDefineExt.Type != nil {
				nodeDefineExtMap["type"] = nodeDefineExt.Type
			}

			if nodeDefineExt.NetService != nil {
				nodeDefineExtMap["net_service"] = nodeDefineExt.NetService
			}

			if nodeDefineExt.District != nil {
				nodeDefineExtMap["district"] = nodeDefineExt.District
			}

			if nodeDefineExt.City != nil {
				nodeDefineExtMap["city"] = nodeDefineExt.City
			}

			if nodeDefineExt.IPType != nil {
				nodeDefineExtMap["i_p_type"] = nodeDefineExt.IPType
			}

			if nodeDefineExt.Location != nil {
				nodeDefineExtMap["location"] = nodeDefineExt.Location
			}

			if nodeDefineExt.CodeType != nil {
				nodeDefineExtMap["code_type"] = nodeDefineExt.CodeType
			}

			if nodeDefineExt.TaskTypes != nil {
				nodeDefineExtMap["task_types"] = nodeDefineExt.TaskTypes
			}

			ids = append(ids, *nodeDefineExt.Code)
			tmpList = append(tmpList, nodeDefineExtMap)
		}

		_ = d.Set("node_set", tmpList)
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
