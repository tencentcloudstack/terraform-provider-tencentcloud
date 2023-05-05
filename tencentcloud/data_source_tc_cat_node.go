/*
Use this data source to query detailed information of cat node

Example Usage

```hcl
data "tencentcloud_cat_node" "node"{
  node_type = 1
  location = 2
  is_ipv6 = false
}
```
*/
package tencentcloud

import (
	"context"
	"log"

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
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Node type 1:IDC,2:LastMile,3:Mobile.",
			},

			"location": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Node area:1=Chinese Mainland,2=Hong Kong, Macao and Taiwan,3=Overseas.",
			},

			"is_ipv6": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "is IPv6.",
			},

			"node_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Node name.",
			},

			"pay_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Payment mode:1=Trial version,2=Paid version.",
			},

			"node_define": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Probe node list.",
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
							Description: "Node ID.",
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node Type;1 = IDC,2 = LastMile,3 = Mobile.",
						},
						"net_service": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Network service provider.",
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
						"ip_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IP type:1 = IPv4,2 = IPv6.",
						},
						"location": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node area:1=Chinese Mainland,2=Hong Kong, Macao and Taiwan,3=Overseas.",
						},
						"code_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "If the node type is base, it is an availability dial test point; if it is blank, it is an advanced dial test point.",
						},
						"node_define_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node status: 1=running, 2=offline.",
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
		paramMap["node_type"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("location"); v != nil {
		paramMap["location"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("is_ipv6"); v != nil {
		paramMap["is_ipv6"] = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("node_name"); ok {
		paramMap["node_name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pay_mode"); ok {
		paramMap["pay_mode"] = helper.IntInt64(v.(int))
	}

	catService := CatService{client: meta.(*TencentCloudClient).apiV3Conn}

	var nodeSets []*cat.NodeDefine
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := catService.DescribeCatNodeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		nodeSets = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Cat nodeSet failed, reason:%+v", logId, err)
		return err
	}

	ids := make([]string, 0, len(nodeSets))
	nodeSetList := make([]map[string]interface{}, 0, len(nodeSets))
	if nodeSets != nil {
		for _, nodeSet := range nodeSets {
			nodeSetMap := map[string]interface{}{}
			if nodeSet.Name != nil {
				nodeSetMap["name"] = nodeSet.Name
			}
			if nodeSet.Code != nil {
				nodeSetMap["code"] = nodeSet.Code
			}
			if nodeSet.Type != nil {
				nodeSetMap["type"] = nodeSet.Type
			}
			if nodeSet.NetService != nil {
				nodeSetMap["net_service"] = nodeSet.NetService
			}
			if nodeSet.District != nil {
				nodeSetMap["district"] = nodeSet.District
			}
			if nodeSet.City != nil {
				nodeSetMap["city"] = nodeSet.City
			}
			if nodeSet.IPType != nil {
				nodeSetMap["ip_type"] = nodeSet.IPType
			}
			if nodeSet.Location != nil {
				nodeSetMap["location"] = nodeSet.Location
			}
			if nodeSet.CodeType != nil {
				nodeSetMap["code_type"] = nodeSet.CodeType
			}
			if nodeSet.NodeDefineStatus != nil {
				nodeSetMap["node_define_status"] = nodeSet.NodeDefineStatus
			}
			ids = append(ids, *nodeSet.Name)
			nodeSetList = append(nodeSetList, nodeSetMap)
		}
		d.SetId(helper.DataResourceIdsHash(ids))
		err = d.Set("node_define", nodeSetList)
		if err != nil {
			return err
		}
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), nodeSetList); e != nil {
			return e
		}
	}

	return nil
}
