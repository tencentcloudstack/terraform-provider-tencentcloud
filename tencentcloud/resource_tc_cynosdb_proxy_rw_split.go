/*
Provides a resource to create a cynosdb proxy_rw_split

Example Usage

```hcl
resource "tencentcloud_cynosdb_proxy_rw_split" "proxy_rw_split" {
  cluster_id           = "cynosdbmysql-cgd2gpwr"
  proxy_group_id       = "cynosdbmysql-proxy-l6zf9t30"
  consistency_type     = "global"
  consistency_time_out = "30"
  weight_mode          = "system"
  instance_weights {
    instance_id = "cynosdbmysql-ins-9810be9i"
    weight      = 0
  }
  fail_over                = "yes"
  auto_add_ro              = "no"
  open_rw                  = "yes"
  rw_type                  = "READWRITE"
  trans_split              = false
  access_mode              = "balance"
  open_connection_pool     = "yes"
  connection_pool_type     = "SessionConnectionPool"
  connection_pool_time_out = 30
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbProxyRwSplit() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbProxyRwSplitCreate,
		Read:   resourceTencentCloudCynosdbProxyRwSplitRead,
		Delete: resourceTencentCloudCynosdbProxyRwSplitDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"proxy_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Database Agent Group ID.",
			},
			"consistency_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Consistency type; Eventual - Eventual consistency, session - session consistency, global - global consistency.",
			},
			"consistency_time_out": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Consistency timeout.",
			},
			"weight_mode": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Reading and writing weight allocation mode; System automatic allocation: system, custom: custom.",
			},
			"instance_weights": {
				Optional:    true,
				ForceNew:    true,
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
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Whether to enable failover. After the agent fails, the connection address will be routed to the main instance, with values of yes and no.",
			},
			"auto_add_ro": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Automatically add read-only instances, values: &amp;#39;yes&amp;#39;,&amp;#39; no &amp;#39;.",
			},
			"open_rw": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Do you want to turn on read write separation.",
			},
			"rw_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Read and write types: READWRITE, READONLY.",
			},
			"trans_split": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Transaction splitting.",
			},
			"access_mode": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Connection mode: nearby, balance.",
			},
			"open_connection_pool": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Open Connection pool: yes, no.",
			},
			"connection_pool_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Connection pool type: SessionConnectionPool.",
			},
			"connection_pool_time_out": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Connection pool time.",
			},
		},
	}
}

func resourceTencentCloudCynosdbProxyRwSplitCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy_rw_split.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId     = getLogId(contextNil)
		request   = cynosdb.NewModifyProxyRwSplitRequest()
		clusterId string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("proxy_group_id"); ok {
		request.ProxyGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("consistency_type"); ok {
		request.ConsistencyType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("consistency_time_out"); ok {
		request.ConsistencyTimeOut = helper.String(v.(string))
	}

	if v, ok := d.GetOk("weight_mode"); ok {
		request.WeightMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_weights"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			proxyInstanceWeight := cynosdb.ProxyInstanceWeight{}
			if v, ok := dMap["instance_id"]; ok {
				proxyInstanceWeight.InstanceId = helper.String(v.(string))
			}
			if v, ok := dMap["weight"]; ok {
				proxyInstanceWeight.Weight = helper.IntInt64(v.(int))
			}
			request.InstanceWeights = append(request.InstanceWeights, &proxyInstanceWeight)
		}
	}

	if v, ok := d.GetOk("fail_over"); ok {
		request.FailOver = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auto_add_ro"); ok {
		request.AutoAddRo = helper.String(v.(string))
	}

	if v, ok := d.GetOk("open_rw"); ok {
		request.OpenRw = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rw_type"); ok {
		request.RwType = helper.String(v.(string))
	}

	if v, _ := d.GetOk("trans_split"); v != nil {
		request.TransSplit = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("access_mode"); ok {
		request.AccessMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("open_connection_pool"); ok {
		request.OpenConnectionPool = helper.String(v.(string))
	}

	if v, ok := d.GetOk("connection_pool_type"); ok {
		request.ConnectionPoolType = helper.String(v.(string))
	}

	if v, _ := d.GetOk("connection_pool_time_out"); v != nil {
		request.ConnectionPoolTimeOut = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyProxyRwSplit(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb proxyRwSplit failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId)

	return resourceTencentCloudCynosdbProxyRwSplitRead(d, meta)
}

func resourceTencentCloudCynosdbProxyRwSplitRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy_rw_split.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbProxyRwSplitDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy_rw_split.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
