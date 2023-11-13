/*
Provides a resource to create a cynosdb proxy

Example Usage

```hcl
resource "tencentcloud_cynosdb_proxy" "proxy" {
  cluster_id = "cynosdbmysql-xxxxxxx"
  cpu = 2
  mem = 4000
  unique_vpc_id = "无"
  unique_subnet_id = "无"
  proxy_count = 2
  connection_pool_type = "SessionConnectionPool"
  open_connection_pool = "yes"
  connection_pool_time_out = 0
  security_group_ids =
  description = "无"
  proxy_zones {
		proxy_node_zone = ""
		proxy_node_count =

  }
}
```

Import

cynosdb proxy can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_proxy.proxy proxy_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCynosdbProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbProxyCreate,
		Read:   resourceTencentCloudCynosdbProxyRead,
		Update: resourceTencentCloudCynosdbProxyUpdate,
		Delete: resourceTencentCloudCynosdbProxyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"cpu": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Number of CPU cores.",
			},

			"mem": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Memory.",
			},

			"unique_vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Private network ID, which is consistent with the cluster private network ID by default.",
			},

			"unique_subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The private network subnet ID is consistent with the cluster subnet ID by default.",
			},

			"proxy_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Number of database proxy group nodes.",
			},

			"connection_pool_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Connection pool type: SessionConnectionPool (session level Connection pool).",
			},

			"open_connection_pool": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Whether to enable Connection pool, yes - enable, no - do not enable.",
			},

			"connection_pool_time_out": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Connection pool threshold: unit (second).",
			},

			"security_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security Group ID Array.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description.",
			},

			"proxy_zones": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Database node information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_node_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Proxy node availability zone.",
						},
						"proxy_node_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of proxy nodes.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCynosdbProxyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cynosdb.NewCreateProxyRequest()
		response  = cynosdb.NewCreateProxyResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("cpu"); ok {
		request.Cpu = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("mem"); ok {
		request.Mem = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("unique_vpc_id"); ok {
		request.UniqueVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("unique_subnet_id"); ok {
		request.UniqueSubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("proxy_count"); ok {
		request.ProxyCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("connection_pool_type"); ok {
		request.ConnectionPoolType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("open_connection_pool"); ok {
		request.OpenConnectionPool = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("connection_pool_time_out"); ok {
		request.ConnectionPoolTimeOut = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_zones"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			proxyZone := cynosdb.ProxyZone{}
			if v, ok := dMap["proxy_node_zone"]; ok {
				proxyZone.ProxyNodeZone = helper.String(v.(string))
			}
			if v, ok := dMap["proxy_node_count"]; ok {
				proxyZone.ProxyNodeCount = helper.IntInt64(v.(int))
			}
			request.ProxyZones = append(request.ProxyZones, &proxyZone)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().CreateProxy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb proxy failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbProxyStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbProxyRead(d, meta)
}

func resourceTencentCloudCynosdbProxyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	proxyId := d.Id()

	proxy, err := service.DescribeCynosdbProxyById(ctx, clusterId)
	if err != nil {
		return err
	}

	if proxy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbProxy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if proxy.ClusterId != nil {
		_ = d.Set("cluster_id", proxy.ClusterId)
	}

	if proxy.Cpu != nil {
		_ = d.Set("cpu", proxy.Cpu)
	}

	if proxy.Mem != nil {
		_ = d.Set("mem", proxy.Mem)
	}

	if proxy.UniqueVpcId != nil {
		_ = d.Set("unique_vpc_id", proxy.UniqueVpcId)
	}

	if proxy.UniqueSubnetId != nil {
		_ = d.Set("unique_subnet_id", proxy.UniqueSubnetId)
	}

	if proxy.ProxyCount != nil {
		_ = d.Set("proxy_count", proxy.ProxyCount)
	}

	if proxy.ConnectionPoolType != nil {
		_ = d.Set("connection_pool_type", proxy.ConnectionPoolType)
	}

	if proxy.OpenConnectionPool != nil {
		_ = d.Set("open_connection_pool", proxy.OpenConnectionPool)
	}

	if proxy.ConnectionPoolTimeOut != nil {
		_ = d.Set("connection_pool_time_out", proxy.ConnectionPoolTimeOut)
	}

	if proxy.SecurityGroupIds != nil {
		_ = d.Set("security_group_ids", proxy.SecurityGroupIds)
	}

	if proxy.Description != nil {
		_ = d.Set("description", proxy.Description)
	}

	if proxy.ProxyZones != nil {
		proxyZonesList := []interface{}{}
		for _, proxyZones := range proxy.ProxyZones {
			proxyZonesMap := map[string]interface{}{}

			if proxy.ProxyZones.ProxyNodeZone != nil {
				proxyZonesMap["proxy_node_zone"] = proxy.ProxyZones.ProxyNodeZone
			}

			if proxy.ProxyZones.ProxyNodeCount != nil {
				proxyZonesMap["proxy_node_count"] = proxy.ProxyZones.ProxyNodeCount
			}

			proxyZonesList = append(proxyZonesList, proxyZonesMap)
		}

		_ = d.Set("proxy_zones", proxyZonesList)

	}

	return nil
}

func resourceTencentCloudCynosdbProxyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		switchProxyVpcRequest  = cynosdb.NewSwitchProxyVpcRequest()
		switchProxyVpcResponse = cynosdb.NewSwitchProxyVpcResponse()
	)

	proxyId := d.Id()

	request.ClusterId = &clusterId

	immutableArgs := []string{"cluster_id", "cpu", "mem", "unique_vpc_id", "unique_subnet_id", "proxy_count", "connection_pool_type", "open_connection_pool", "connection_pool_time_out", "security_group_ids", "description", "proxy_zones"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cluster_id") {
		if v, ok := d.GetOk("cluster_id"); ok {
			request.ClusterId = helper.String(v.(string))
		}
	}

	if d.HasChange("cpu") {
		if v, ok := d.GetOkExists("cpu"); ok {
			request.Cpu = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("mem") {
		if v, ok := d.GetOkExists("mem"); ok {
			request.Mem = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("proxy_count") {
		if v, ok := d.GetOkExists("proxy_count"); ok {
			request.ProxyCount = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("proxy_zones") {
		if v, ok := d.GetOk("proxy_zones"); ok {
			for _, item := range v.([]interface{}) {
				proxyZone := cynosdb.ProxyZone{}
				if v, ok := dMap["proxy_node_zone"]; ok {
					proxyZone.ProxyNodeZone = helper.String(v.(string))
				}
				if v, ok := dMap["proxy_node_count"]; ok {
					proxyZone.ProxyNodeCount = helper.IntInt64(v.(int))
				}
				request.ProxyZones = append(request.ProxyZones, &proxyZone)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().SwitchProxyVpc(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb proxy failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbProxyRead(d, meta)
}

func resourceTencentCloudCynosdbProxyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	proxyId := d.Id()

	if err := service.DeleteCynosdbProxyById(ctx, clusterId); err != nil {
		return err
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbProxyStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
