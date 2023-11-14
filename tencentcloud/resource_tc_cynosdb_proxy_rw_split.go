/*
Provides a resource to create a cynosdb proxy_rw_split

Example Usage

```hcl
resource "tencentcloud_cynosdb_proxy_rw_split" "proxy_rw_split" {
  cluster_id = "cynosdbmysql-xxxxxxx"
  proxy_group_id = "无"
  consistency_type = "无"
  consistency_time_out = "无"
  weight_mode = "无"
  instance_weights {
		instance_id = ""
		weight =

  }
  fail_over = "无"
  auto_add_ro = "无"
  open_rw = "无"
  rw_type = "无"
  trans_split = 无
  access_mode = "无"
  open_connection_pool = "无"
  connection_pool_type = "无"
  connection_pool_time_out = 无
}
```

Import

cynosdb proxy_rw_split can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_proxy_rw_split.proxy_rw_split proxy_rw_split_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceTencentCloudCynosdbProxyRwSplit() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbProxyRwSplitCreate,
		Read:   resourceTencentCloudCynosdbProxyRwSplitRead,
		Update: resourceTencentCloudCynosdbProxyRwSplitUpdate,
		Delete: resourceTencentCloudCynosdbProxyRwSplitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"proxy_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database Agent Group ID.",
			},

			"consistency_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Consistency type; Eventual - Eventual consistency, session - session consistency, global - global consistency.",
			},

			"consistency_time_out": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Consistency timeout.",
			},

			"weight_mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Reading and writing weight allocation mode; System automatic allocation: system, custom: custom.",
			},

			"instance_weights": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Instance read-only weight.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance Id.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Instance Weight.",
						},
					},
				},
			},

			"fail_over": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Whether to enable failover. After the agent fails, the connection address will be routed to the main instance, with values of yes and no.",
			},

			"auto_add_ro": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Automatically add read-only instances, values: &amp;amp;#39;yes&amp;amp;#39;,&amp;amp;#39; no &amp;amp;#39;.",
			},

			"open_rw": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Do you want to turn on read write separation.",
			},

			"rw_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Read and write types: READWRITE, READONLY.",
			},

			"trans_split": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Transaction splitting.",
			},

			"access_mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Connection mode: nearby, balance.",
			},

			"open_connection_pool": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Open Connection pool: yes, no.",
			},

			"connection_pool_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Connection pool type: SessionConnectionPool.",
			},

			"connection_pool_time_out": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Connection pool time.",
			},
		},
	}
}

func resourceTencentCloudCynosdbProxyRwSplitCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy_rw_split.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudCynosdbProxyRwSplitUpdate(d, meta)
}

func resourceTencentCloudCynosdbProxyRwSplitRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy_rw_split.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	proxyRwSplitId := d.Id()

	proxyRwSplit, err := service.DescribeCynosdbProxyRwSplitById(ctx, instanceId)
	if err != nil {
		return err
	}

	if proxyRwSplit == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbProxyRwSplit` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if proxyRwSplit.ClusterId != nil {
		_ = d.Set("cluster_id", proxyRwSplit.ClusterId)
	}

	if proxyRwSplit.ProxyGroupId != nil {
		_ = d.Set("proxy_group_id", proxyRwSplit.ProxyGroupId)
	}

	if proxyRwSplit.ConsistencyType != nil {
		_ = d.Set("consistency_type", proxyRwSplit.ConsistencyType)
	}

	if proxyRwSplit.ConsistencyTimeOut != nil {
		_ = d.Set("consistency_time_out", proxyRwSplit.ConsistencyTimeOut)
	}

	if proxyRwSplit.WeightMode != nil {
		_ = d.Set("weight_mode", proxyRwSplit.WeightMode)
	}

	if proxyRwSplit.InstanceWeights != nil {
		instanceWeightsList := []interface{}{}
		for _, instanceWeights := range proxyRwSplit.InstanceWeights {
			instanceWeightsMap := map[string]interface{}{}

			if proxyRwSplit.InstanceWeights.InstanceId != nil {
				instanceWeightsMap["instance_id"] = proxyRwSplit.InstanceWeights.InstanceId
			}

			if proxyRwSplit.InstanceWeights.Weight != nil {
				instanceWeightsMap["weight"] = proxyRwSplit.InstanceWeights.Weight
			}

			instanceWeightsList = append(instanceWeightsList, instanceWeightsMap)
		}

		_ = d.Set("instance_weights", instanceWeightsList)

	}

	if proxyRwSplit.FailOver != nil {
		_ = d.Set("fail_over", proxyRwSplit.FailOver)
	}

	if proxyRwSplit.AutoAddRo != nil {
		_ = d.Set("auto_add_ro", proxyRwSplit.AutoAddRo)
	}

	if proxyRwSplit.OpenRw != nil {
		_ = d.Set("open_rw", proxyRwSplit.OpenRw)
	}

	if proxyRwSplit.RwType != nil {
		_ = d.Set("rw_type", proxyRwSplit.RwType)
	}

	if proxyRwSplit.TransSplit != nil {
		_ = d.Set("trans_split", proxyRwSplit.TransSplit)
	}

	if proxyRwSplit.AccessMode != nil {
		_ = d.Set("access_mode", proxyRwSplit.AccessMode)
	}

	if proxyRwSplit.OpenConnectionPool != nil {
		_ = d.Set("open_connection_pool", proxyRwSplit.OpenConnectionPool)
	}

	if proxyRwSplit.ConnectionPoolType != nil {
		_ = d.Set("connection_pool_type", proxyRwSplit.ConnectionPoolType)
	}

	if proxyRwSplit.ConnectionPoolTimeOut != nil {
		_ = d.Set("connection_pool_time_out", proxyRwSplit.ConnectionPoolTimeOut)
	}

	return nil
}

func resourceTencentCloudCynosdbProxyRwSplitUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy_rw_split.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"cluster_id", "proxy_group_id", "consistency_type", "consistency_time_out", "weight_mode", "instance_weights", "fail_over", "auto_add_ro", "open_rw", "rw_type", "trans_split", "access_mode", "open_connection_pool", "connection_pool_type", "connection_pool_time_out"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudCynosdbProxyRwSplitRead(d, meta)
}

func resourceTencentCloudCynosdbProxyRwSplitDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy_rw_split.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
