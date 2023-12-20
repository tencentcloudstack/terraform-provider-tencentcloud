package cat

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCatNode() *schema.Resource {
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
						"task_types": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Computed:    true,
							Description: "The task types supported by the node. `1`: page performance, `2`: file upload, `3`: file download, `4`: port performance, `5`: network quality, `6`: audio and video experience.",
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
	defer tccommon.LogElapsed("data_source.tencentcloud_cat_node.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

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

	catService := CatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var nodeSets []*cat.NodeDefine
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := catService.DescribeCatProbeNodeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		nodeSets = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Cat nodeSet failed, reason:%+v", logId, err)
		return err
	}

	var nodeSetExt []*cat.NodeDefineExt
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, e := catService.DescribeCatNodeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		nodeSetExt = results
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

			for _, node := range nodeSetExt {
				if *node.Code == *nodeSet.Code {
					if node.TaskTypes != nil {
						nodeSetMap["task_types"] = node.TaskTypes
					}
					break
				}
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
		if e := tccommon.WriteToFile(output.(string), nodeSetList); e != nil {
			return e
		}
	}

	return nil
}
