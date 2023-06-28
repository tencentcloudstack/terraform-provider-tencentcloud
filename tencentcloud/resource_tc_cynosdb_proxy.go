/*
Provides a resource to create a cynosdb proxy

Example Usage

```hcl
resource "tencentcloud_cynosdb_proxy" "proxy" {
  cluster_id               = "cynosdbmysql-bws8h88b"
  cpu                      = 2
  mem                      = 4000
  unique_vpc_id            = "vpc-k1t8ickr"
  unique_subnet_id         = "subnet-jdi5xn22"
  connection_pool_type     = "SessionConnectionPool"
  open_connection_pool     = "yes"
  connection_pool_time_out = 30
  security_group_ids       = ["sg-baxfiao5"]
  description              = "desc sample"
  proxy_zones {
    proxy_node_zone  = "ap-guangzhou-7"
    proxy_node_count = 2
  }
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbProxyCreate,
		Read:   resourceTencentCloudCynosdbProxyRead,
		Update: resourceTencentCloudCynosdbProxyUpdate,
		Delete: resourceTencentCloudCynosdbProxyDelete,

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
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Private network ID, which is consistent with the cluster private network ID by default.",
			},
			"unique_subnet_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The private network subnet ID is consistent with the cluster subnet ID by default.",
			},
			"proxy_count": {
				Optional:      true,
				Computed:      true,
				Type:          schema.TypeInt,
				ConflictsWith: []string{"proxy_zones"},
				Description:   "Number of database proxy group nodes. If it is set at the same time as the `proxy_zones` field, the `proxy_zones` parameter shall prevail.",
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
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Security Group ID Array.",
			},
			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description.",
			},
			"proxy_group_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Proxy Group Id.",
			},
			"proxy_zones": {
				Optional:      true,
				Computed:      true,
				Type:          schema.TypeList,
				Description:   "Database node information.",
				ConflictsWith: []string{"proxy_count"},
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

	var (
		logId        = getLogId(contextNil)
		ctx          = context.WithValue(context.TODO(), logIdKey, logId)
		service      = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		request      = cynosdb.NewCreateProxyRequest()
		response     = cynosdb.NewCreateProxyResponse()
		clusterId    string
		proxyGroupId string
		flowId       int64
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("cpu"); ok {
		request.Cpu = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("mem"); ok {
		request.Mem = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("unique_vpc_id"); ok {
		request.UniqueVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("unique_subnet_id"); ok {
		request.UniqueSubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("connection_pool_type"); ok {
		request.ConnectionPoolType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("open_connection_pool"); ok {
		request.OpenConnectionPool = helper.String(v.(string))
	}

	if v, ok := d.GetOk("connection_pool_time_out"); ok {
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
	} else {
		if v, ok = d.GetOk("proxy_count"); ok {
			request.ProxyCount = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().CreateProxy(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("cynosdb proxy not exists")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb proxy failed, reason:%+v", logId, err)
		return err
	}

	flowId = *response.Response.FlowId
	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeFlow(ctx, flowId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}

		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("cynosdb proxy is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb proxy fail, reason:%s\n", logId, err.Error())
		return err
	}

	proxy, err := service.DescribeCynosdbProxyById(ctx, clusterId, "")
	if err != nil {
		return err
	}

	if proxy == nil {
		log.Printf("[WARN]%s resource `CynosdbProxy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	proxyGroupRwInfo := proxy.ProxyGroupInfos[0]
	proxyGroup := proxyGroupRwInfo.ProxyGroup
	if proxyGroup != nil {
		if proxyGroup.ProxyGroupId != nil {
			proxyGroupId = *proxyGroup.ProxyGroupId
		}
	}

	d.SetId(strings.Join([]string{clusterId, proxyGroupId}, FILED_SP))

	return resourceTencentCloudCynosdbProxyRead(d, meta)
}

func resourceTencentCloudCynosdbProxyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	clusterId := idSplit[0]
	proxyGroupId := idSplit[1]

	proxy, err := service.DescribeCynosdbProxyById(ctx, clusterId, proxyGroupId)
	if err != nil {
		return err
	}

	if proxy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbProxy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if proxy != nil {
		proxyGroupRwInfo := proxy.ProxyGroupInfos[0]
		connectionPool := proxyGroupRwInfo.ConnectionPool
		if connectionPool != nil {
			if connectionPool.ConnectionPoolType != nil {
				_ = d.Set("connection_pool_type", connectionPool.ConnectionPoolType)
			}

			if connectionPool.OpenConnectionPool != nil {
				_ = d.Set("open_connection_pool", connectionPool.OpenConnectionPool)
			}

			if connectionPool.ConnectionPoolTimeOut != nil {
				_ = d.Set("connection_pool_time_out", connectionPool.ConnectionPoolTimeOut)
			}
		}

		netAddrInfos := proxyGroupRwInfo.NetAddrInfos
		if netAddrInfos != nil {
			netAddrInfo := netAddrInfos[0]
			if netAddrInfo.Description != nil {
				_ = d.Set("description", netAddrInfo.Description)
			}

			if netAddrInfo.UniqVpcId != nil {
				_ = d.Set("unique_vpc_id", netAddrInfo.UniqVpcId)
			}

			if netAddrInfo.UniqSubnetId != nil {
				_ = d.Set("unique_subnet_id", netAddrInfo.UniqSubnetId)
			}
		}

		proxyGroups := proxyGroupRwInfo.ProxyGroup
		if proxyGroups != nil {
			if proxyGroups.ProxyNodeCount != nil {
				_ = d.Set("proxy_count", proxyGroups.ProxyNodeCount)
			}

			if proxyGroups.ProxyGroupId != nil {
				_ = d.Set("proxy_group_id", proxyGroups.ProxyGroupId)
			}
		}

		proxyNodes := proxyGroupRwInfo.ProxyNodes
		if proxyNodes != nil {
			zoneMap := make(map[string]int)
			for _, v := range proxyNodes {
				if v.Cpu != nil {
					_ = d.Set("cpu", v.Cpu)
				}

				if v.Mem != nil {
					_ = d.Set("mem", v.Mem)
				}

				if v.Zone != nil {
					zone := *v.Zone
					_, ok := zoneMap[zone]
					if ok {
						zoneMap[zone] += 1
					} else {
						zoneMap[zone] = 1
					}
				}
			}

			if zoneMap != nil {
				tmpList := []interface{}{}
				for k, v := range zoneMap {
					tmpMap := make(map[string]interface{})
					tmpMap["proxy_node_zone"] = k
					tmpMap["proxy_node_count"] = v
					tmpList = append(tmpList, tmpMap)
				}
				_ = d.Set("proxy_zones", tmpList)
			}
		}
	}

	return nil
}

func resourceTencentCloudCynosdbProxyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId                  = getLogId(contextNil)
		ctx                    = context.WithValue(context.TODO(), logIdKey, logId)
		service                = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		switchProxyRequest     = cynosdb.NewSwitchProxyVpcRequest()
		modifyProxyDescRequest = cynosdb.NewModifyProxyDescRequest()
		upgradeProxyRequest    = cynosdb.NewUpgradeProxyRequest()
		flowId                 int64
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	clusterId := idSplit[0]
	proxyGroupId := idSplit[1]

	immutableArgs := []string{"cluster_id", "connection_pool_type", "open_connection_pool", "connection_pool_time_out", "security_group_ids"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("unique_vpc_id") || d.HasChange("unique_subnet_id") {
		switchProxyRequest.ClusterId = &clusterId
		switchProxyRequest.ProxyGroupId = &proxyGroupId
		switchProxyRequest.OldIpReserveHours = common.Int64Ptr(1)

		if v, ok := d.GetOk("unique_vpc_id"); ok {
			switchProxyRequest.UniqVpcId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("unique_subnet_id"); ok {
			switchProxyRequest.UniqSubnetId = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().SwitchProxyVpc(switchProxyRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, switchProxyRequest.GetAction(), switchProxyRequest.ToJsonString(), result.ToJsonString())
			}

			flowId = *result.Response.FlowId
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb proxy failed, reason:%+v", logId, err)
			return err
		}

		err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
			ok, err := service.DescribeFlow(ctx, flowId)
			if err != nil {
				if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
					return resource.RetryableError(err)
				} else {
					return resource.NonRetryableError(err)
				}
			}

			if ok {
				return nil
			} else {
				return resource.RetryableError(fmt.Errorf("cynosdb proxy is processing"))
			}
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb proxy fail, reason:%s\n", logId, err.Error())
			return err
		}
	}

	if d.HasChange("description") {
		modifyProxyDescRequest.ClusterId = &clusterId
		modifyProxyDescRequest.ProxyGroupId = &proxyGroupId

		if v, ok := d.GetOk("description"); ok {
			modifyProxyDescRequest.Description = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyProxyDesc(modifyProxyDescRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, switchProxyRequest.GetAction(), switchProxyRequest.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb proxy desc failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("cpu") || d.HasChange("mem") || d.HasChange("proxy_count") || d.HasChange("proxy_zones") {
		upgradeProxyRequest.ClusterId = &clusterId

		if v, ok := d.GetOk("cpu"); ok {
			upgradeProxyRequest.Cpu = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("mem"); ok {
			upgradeProxyRequest.Mem = helper.IntInt64(v.(int))
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
				upgradeProxyRequest.ProxyZones = append(upgradeProxyRequest.ProxyZones, &proxyZone)
			}
		} else {
			if v, ok = d.GetOk("proxy_count"); ok {
				upgradeProxyRequest.ProxyCount = helper.IntInt64(v.(int))
			}
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().UpgradeProxy(upgradeProxyRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, switchProxyRequest.GetAction(), switchProxyRequest.ToJsonString(), result.ToJsonString())
			}

			flowId = *result.Response.FlowId
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s upgrade proxy failed, reason:%+v", logId, err)
			return err
		}

		err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
			ok, err := service.DescribeFlow(ctx, flowId)
			if err != nil {
				if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
					return resource.RetryableError(err)
				} else {
					return resource.NonRetryableError(err)
				}
			}

			if ok {
				return nil
			} else {
				return resource.RetryableError(fmt.Errorf("upgrade proxy is processing"))
			}
		})

		if err != nil {
			log.Printf("[CRITAL]%s upgrade proxy fail, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudCynosdbProxyRead(d, meta)
}

func resourceTencentCloudCynosdbProxyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	clusterId := idSplit[0]

	flowId, err := service.DeleteCynosdbProxyById(ctx, clusterId)
	if err != nil {
		return err
	}

	err = resource.Retry(6*readRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeFlow(ctx, *flowId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}

		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("cynosdb proxy is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s close cynosdb proxy fail, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
