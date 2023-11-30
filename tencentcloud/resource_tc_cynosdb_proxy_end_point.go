package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCynosdbProxyEndPoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbProxyEndPointCreate,
		Read:   resourceTencentCloudCynosdbProxyEndPointRead,
		Update: resourceTencentCloudCynosdbProxyEndPointUpdate,
		Delete: resourceTencentCloudCynosdbProxyEndPointDelete,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"unique_vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Private network ID, which is consistent with the cluster private network ID by default.",
			},
			"unique_subnet_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The private network subnet ID is consistent with the cluster subnet ID by default.",
			},
			"connection_pool_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Connection pool type: SessionConnectionPool (session level Connection pool).",
			},
			"open_connection_pool": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Whether to enable Connection pool, yes - enable, no - do not enable.",
			},
			"connection_pool_time_out": {
				Optional:    true,
				Computed:    true,
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
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Description.",
			},
			"vip": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "VIP Information.",
			},
			"vport": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Port Information.",
			},
			"weight_mode": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Weight mode: system system allocation, custom customization.",
			},
			"auto_add_ro": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Do you want to automatically add read-only instances? Yes - Yes, no - Do not automatically add.",
			},
			"fail_over": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Enable Failover. yes or no.",
			},
			"consistency_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Consistency type: event, global, session.",
			},
			"rw_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Read and write attributes: READWRITE, READONLY.",
			},
			"consistency_time_out": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Consistency timeout.",
			},
			"trans_split": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Transaction splitting.",
			},
			"access_mode": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Connection mode: nearby, balance.",
			},
			"instance_weights": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Instance Weight.",
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
			"instance_group_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Instance Group ID.",
			},
			"proxy_group_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Proxy Group ID.",
			},
		},
	}
}

func resourceTencentCloudCynosdbProxyEndPointCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy_end_point.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId                 = getLogId(contextNil)
		ctx                   = context.WithValue(context.TODO(), logIdKey, logId)
		service               = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		request               = cynosdb.NewCreateProxyEndPointRequest()
		response              = cynosdb.NewCreateProxyEndPointResponse()
		modifyVipVportRequest = cynosdb.NewModifyVipVportRequest()
		clusterId             string
		vip                   *string
		vPort                 *int64
		proxyGroupId          string
		instanceGroupId       string
		flowId                int64
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("unique_vpc_id"); ok {
		request.UniqueVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("unique_subnet_id"); ok {
		request.UniqueSubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("open_connection_pool"); ok {
		if v.(string) == STATUS_YES {
			request.OpenConnectionPool = helper.String(v.(string))

			if v, ok := d.GetOk("connection_pool_type"); ok {
				request.ConnectionPoolType = helper.String(v.(string))
			}

			if v, ok := d.GetOkExists("connection_pool_time_out"); ok {
				request.ConnectionPoolTimeOut = helper.IntInt64(v.(int))
			}

		} else if v.(string) == STATUS_NO {
			if _, ok := d.GetOk("connection_pool_type"); ok {
				return fmt.Errorf("connection_pool_type invalid parameter")
			}

			if _, ok := d.GetOkExists("connection_pool_time_out"); ok {
				return fmt.Errorf("connection_pool_time_out invalid parameter")
			}
		} else {
			return fmt.Errorf("open_connection_pool invalid parameter")
		}
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

	if v, ok := d.GetOk("vip"); ok {
		request.Vip = helper.String(v.(string))
		vip = helper.String(v.(string))
		if v, ok := d.GetOk("vport"); ok {
			vPort = helper.IntInt64(v.(int))
		}
	}

	if v, ok := d.GetOk("weight_mode"); ok {
		request.WeightMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auto_add_ro"); ok {
		request.AutoAddRo = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rw_type"); ok {
		if v.(string) == RW_TYPE {
			request.RwType = helper.String(v.(string))

			if v, ok := d.GetOk("consistency_type"); ok {
				request.ConsistencyType = helper.String(v.(string))
			}

			if v, ok := d.GetOkExists("consistency_time_out"); ok {
				request.ConsistencyTimeOut = helper.IntInt64(v.(int))
			}

			if v, ok := d.GetOk("fail_over"); ok {
				request.FailOver = helper.String(v.(string))
			}

		} else if v.(string) == RO_TYPE {
			if _, ok := d.GetOk("consistency_type"); ok {
				return fmt.Errorf("consistency_type invalid parameter")
			}

			if _, ok := d.GetOkExists("consistency_time_out"); ok {
				return fmt.Errorf("consistency_time_out invalid parameter")
			}

			if _, ok := d.GetOk("fail_over"); ok {
				return fmt.Errorf("fail_over invalid parameter")
			}
		} else {
			return fmt.Errorf("rw_type invalid parameter")
		}
	}

	if v, ok := d.GetOkExists("trans_split"); ok {
		request.TransSplit = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("access_mode"); ok {
		request.AccessMode = helper.String(v.(string))
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().CreateProxyEndPoint(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb proxyEndPoint failed, reason:%+v", logId, err)
		return err
	}

	proxyGroupId = *response.Response.ProxyGroupId
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
			return resource.RetryableError(fmt.Errorf("create cynosdb proxyEndPoint is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cynosdb proxyEndPoint fail, reason:%s\n", logId, err.Error())
		return err
	}

	proxyEndPoint, err := service.DescribeCynosdbProxyEndPointById(ctx, clusterId, proxyGroupId)
	if err != nil {
		return err
	}

	if proxyEndPoint == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbProxyEndPoint` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	instanceGroupId = *proxyEndPoint.NetAddrInfos[0].InstanceGroupId
	sourceVip := *proxyEndPoint.NetAddrInfos[0].Vip
	sourceVport := *proxyEndPoint.NetAddrInfos[0].Vport
	if (vip != nil && vPort != nil) && (*vip != sourceVip && *vPort != sourceVport) {
		modifyVipVportRequest.ClusterId = &clusterId
		modifyVipVportRequest.InstanceGrpId = &instanceGroupId
		modifyVipVportRequest.Vip = vip
		modifyVipVportRequest.Vport = vPort
		modifyVipVportRequest.OldIpReserveHours = helper.IntInt64(0)

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyVipVport(modifyVipVportRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyVipVportRequest.GetAction(), modifyVipVportRequest.ToJsonString(), result.ToJsonString())
			}

			flowId = *result.Response.FlowId
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s create cynosdb proxyEndPoint rw split vip vport failed, reason:%+v", logId, err)
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
				return resource.RetryableError(fmt.Errorf("create cynosdb proxyEndPoint rw split vip vport is processing"))
			}
		})

		if err != nil {
			log.Printf("[CRITAL]%s create cynosdb proxyEndPoint rw split vip vport fail, reason:%s\n", logId, err.Error())
			return err
		}
	}

	d.SetId(strings.Join([]string{clusterId, proxyGroupId, instanceGroupId}, FILED_SP))

	return resourceTencentCloudCynosdbProxyEndPointRead(d, meta)
}

func resourceTencentCloudCynosdbProxyEndPointRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy_end_point.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	clusterId := idSplit[0]
	proxyGroupId := idSplit[1]

	proxyEndPoint, err := service.DescribeCynosdbProxyEndPointById(ctx, clusterId, proxyGroupId)
	if err != nil {
		return err
	}

	if proxyEndPoint == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbProxyEndPoint` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if proxyEndPoint.ProxyGroup.ClusterId != nil {
		_ = d.Set("cluster_id", proxyEndPoint.ProxyGroup.ClusterId)
	}

	if proxyEndPoint.ProxyGroup.ProxyGroupId != nil {
		_ = d.Set("proxy_group_id", proxyEndPoint.ProxyGroup.ProxyGroupId)
	}

	if proxyEndPoint.NetAddrInfos[0].UniqVpcId != nil {
		_ = d.Set("unique_vpc_id", proxyEndPoint.NetAddrInfos[0].UniqVpcId)
	}

	if proxyEndPoint.NetAddrInfos[0].UniqSubnetId != nil {
		_ = d.Set("unique_subnet_id", proxyEndPoint.NetAddrInfos[0].UniqSubnetId)
	}

	if proxyEndPoint.ConnectionPool.ConnectionPoolType != nil {
		_ = d.Set("connection_pool_type", proxyEndPoint.ConnectionPool.ConnectionPoolType)
	}

	if proxyEndPoint.ConnectionPool.OpenConnectionPool != nil {
		if *proxyEndPoint.ConnectionPool.OpenConnectionPool == STATUS_YES {
			_ = d.Set("open_connection_pool", proxyEndPoint.ConnectionPool.OpenConnectionPool)
		} else {
			_ = d.Set("open_connection_pool", STATUS_NO)
		}
	}

	if proxyEndPoint.ConnectionPool.ConnectionPoolTimeOut != nil {
		_ = d.Set("connection_pool_time_out", proxyEndPoint.ConnectionPool.ConnectionPoolTimeOut)
	}

	if proxyEndPoint.NetAddrInfos[0].Description != nil {
		_ = d.Set("description", proxyEndPoint.NetAddrInfos[0].Description)
	}

	if proxyEndPoint.NetAddrInfos[0].Vip != nil {
		_ = d.Set("vip", proxyEndPoint.NetAddrInfos[0].Vip)
	}

	if proxyEndPoint.NetAddrInfos[0].Vport != nil {
		_ = d.Set("vport", proxyEndPoint.NetAddrInfos[0].Vport)
	}

	if proxyEndPoint.NetAddrInfos[0].InstanceGroupId != nil {
		_ = d.Set("instance_group_id", proxyEndPoint.NetAddrInfos[0].InstanceGroupId)
	}

	if proxyEndPoint.ProxyGroupRwInfo.WeightMode != nil {
		_ = d.Set("weight_mode", proxyEndPoint.ProxyGroupRwInfo.WeightMode)
	}

	if proxyEndPoint.ProxyGroupRwInfo.AutoAddRo != nil {
		_ = d.Set("auto_add_ro", proxyEndPoint.ProxyGroupRwInfo.AutoAddRo)
	}

	if proxyEndPoint.ProxyGroupRwInfo.FailOver != nil {
		_ = d.Set("fail_over", proxyEndPoint.ProxyGroupRwInfo.FailOver)
	}

	if proxyEndPoint.ProxyGroupRwInfo.ConsistencyType != nil {
		_ = d.Set("consistency_type", proxyEndPoint.ProxyGroupRwInfo.ConsistencyType)
	}

	if proxyEndPoint.ProxyGroupRwInfo.RwType != nil {
		_ = d.Set("rw_type", proxyEndPoint.ProxyGroupRwInfo.RwType)
	}

	if proxyEndPoint.ProxyGroupRwInfo.ConsistencyTimeOut != nil {
		_ = d.Set("consistency_time_out", proxyEndPoint.ProxyGroupRwInfo.ConsistencyTimeOut)
	}

	if proxyEndPoint.ProxyGroupRwInfo.TransSplit != nil {
		_ = d.Set("trans_split", proxyEndPoint.ProxyGroupRwInfo.TransSplit)
	}

	if proxyEndPoint.ProxyGroupRwInfo.AccessMode != nil {
		_ = d.Set("access_mode", proxyEndPoint.ProxyGroupRwInfo.AccessMode)
	}

	if proxyEndPoint.ProxyGroupRwInfo.InstanceWeights != nil {
		instanceWeightsList := []interface{}{}
		for _, instanceWeight := range proxyEndPoint.ProxyGroupRwInfo.InstanceWeights {
			instanceWeightsMap := map[string]interface{}{}

			if instanceWeight.InstanceId != nil {
				instanceWeightsMap["instance_id"] = instanceWeight.InstanceId
			}

			if instanceWeight.Weight != nil {
				instanceWeightsMap["weight"] = instanceWeight.Weight
			}

			instanceWeightsList = append(instanceWeightsList, instanceWeightsMap)
		}

		_ = d.Set("instance_weights", instanceWeightsList)

	}

	return nil
}

func resourceTencentCloudCynosdbProxyEndPointUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy_end_point.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId                     = getLogId(contextNil)
		ctx                       = context.WithValue(context.TODO(), logIdKey, logId)
		service                   = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		switchProxyVpcRequest     = cynosdb.NewSwitchProxyVpcRequest()
		modifyProxyRwSplitRequest = cynosdb.NewModifyProxyRwSplitRequest()
		modifyProxyDescRequest    = cynosdb.NewModifyProxyDescRequest()
		modifyVipVportRequest     = cynosdb.NewModifyVipVportRequest()
		flowId                    int64
	)

	immutableArgs := []string{"cluster_id", "security_group_ids"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	clusterId := idSplit[0]
	proxyGroupId := idSplit[1]
	instanceGroupId := idSplit[2]

	rwSplit := false
	rwSplitMutableArgs := []string{"connection_pool_type", "open_connection_pool", "connection_pool_time_out", "weight_mode", "auto_add_ro", "fail_over", "consistency_type", "rw_type", "consistency_time_out", "trans_split", "access_mode", "instance_weights"}
	for _, v := range rwSplitMutableArgs {
		if d.HasChange(v) {
			rwSplit = true
			break
		}
	}

	if rwSplit {
		modifyProxyRwSplitRequest.ClusterId = &clusterId
		modifyProxyRwSplitRequest.ProxyGroupId = &proxyGroupId

		if d.HasChange("open_connection_pool") {
			if v, ok := d.GetOk("open_connection_pool"); ok {
				if v.(string) == STATUS_YES {
					modifyProxyRwSplitRequest.OpenConnectionPool = helper.String(v.(string))

					if v, ok := d.GetOk("connection_pool_type"); ok {
						modifyProxyRwSplitRequest.ConnectionPoolType = helper.String(v.(string))
					}

					if v, ok := d.GetOkExists("connection_pool_time_out"); ok {
						modifyProxyRwSplitRequest.ConnectionPoolTimeOut = helper.IntInt64(v.(int))
					}

				} else if v.(string) == STATUS_NO {
					modifyProxyRwSplitRequest.OpenConnectionPool = helper.String(v.(string))

				} else {
					return fmt.Errorf("open_connection_pool invalid parameter")
				}
			}
		}

		if d.HasChange("rw_type") {
			if v, ok := d.GetOk("rw_type"); ok {
				if v.(string) == RW_TYPE {
					modifyProxyRwSplitRequest.RwType = helper.String(v.(string))

					if v, ok := d.GetOk("fail_over"); ok {
						modifyProxyRwSplitRequest.FailOver = helper.String(v.(string))
					}

					if v, ok := d.GetOk("consistency_type"); ok {
						modifyProxyRwSplitRequest.ConsistencyType = helper.String(v.(string))
					}

					if v, ok := d.GetOkExists("consistency_time_out"); ok {
						modifyProxyRwSplitRequest.ConsistencyTimeOut = helper.String(v.(string))
					}

				} else if v.(string) == RO_TYPE {
					modifyProxyRwSplitRequest.RwType = helper.String(v.(string))

				} else {
					return fmt.Errorf("rw_type invalid parameter")
				}
			}
		}

		if v, ok := d.GetOk("weight_mode"); ok {
			modifyProxyRwSplitRequest.WeightMode = helper.String(v.(string))
		}

		if v, ok := d.GetOk("auto_add_ro"); ok {
			modifyProxyRwSplitRequest.AutoAddRo = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("trans_split"); ok {
			modifyProxyRwSplitRequest.TransSplit = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOk("access_mode"); ok {
			modifyProxyRwSplitRequest.AccessMode = helper.String(v.(string))
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
				modifyProxyRwSplitRequest.InstanceWeights = append(modifyProxyRwSplitRequest.InstanceWeights, &proxyInstanceWeight)
			}
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyProxyRwSplit(modifyProxyRwSplitRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyProxyRwSplitRequest.GetAction(), modifyProxyRwSplitRequest.ToJsonString(), result.ToJsonString())
			}

			flowId = *result.Response.FlowId
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb proxyEndPoint rw split, reason:%+v", logId, err)
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
				return resource.RetryableError(fmt.Errorf("cynosdb proxyEndPoint rw split is processing"))
			}
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb proxyEndPoint rw split fail, reason:%s\n", logId, err.Error())
			return err
		}
	}

	if d.HasChange("unique_vpc_id") || d.HasChange("unique_subnet_id") {
		switchProxyVpcRequest.ClusterId = &clusterId
		switchProxyVpcRequest.ProxyGroupId = &proxyGroupId
		switchProxyVpcRequest.OldIpReserveHours = helper.IntInt64(1)

		if v, ok := d.GetOk("unique_vpc_id"); ok {
			switchProxyVpcRequest.UniqVpcId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("unique_subnet_id"); ok {
			switchProxyVpcRequest.UniqSubnetId = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().SwitchProxyVpc(switchProxyVpcRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, switchProxyVpcRequest.GetAction(), switchProxyVpcRequest.ToJsonString(), result.ToJsonString())
			}

			flowId = *result.Response.FlowId
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb proxy vpc failed, reason:%+v", logId, err)
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
				return resource.RetryableError(fmt.Errorf("cynosdb proxy vpc is processing"))
			}
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb proxy vpc fail, reason:%s\n", logId, err.Error())
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
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyProxyDescRequest.GetAction(), modifyProxyDescRequest.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb proxyEndPoint rw split desc failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("vip") || d.HasChange("vport") {
		modifyVipVportRequest.ClusterId = &clusterId
		modifyVipVportRequest.InstanceGrpId = &instanceGroupId
		modifyVipVportRequest.OldIpReserveHours = helper.IntInt64(0)

		if v, ok := d.GetOk("vip"); ok {
			modifyVipVportRequest.Vip = helper.String(v.(string))
		}

		if v, ok := d.GetOk("vport"); ok {
			modifyVipVportRequest.Vport = helper.IntInt64(v.(int))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().ModifyVipVport(modifyVipVportRequest)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyVipVportRequest.GetAction(), modifyVipVportRequest.ToJsonString(), result.ToJsonString())
			}

			flowId = *result.Response.FlowId
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb proxyEndPoint rw split vip vport failed, reason:%+v", logId, err)
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
				return resource.RetryableError(fmt.Errorf("update cynosdb proxyEndPoint rw split vip vport is processing"))
			}
		})

		if err != nil {
			log.Printf("[CRITAL]%s update cynosdb proxyEndPoint rw split vip vport fail, reason:%s\n", logId, err.Error())
			return err
		}
	}

	return resourceTencentCloudCynosdbProxyEndPointRead(d, meta)
}

func resourceTencentCloudCynosdbProxyEndPointDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy_end_point.delete")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		service = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	clusterId := idSplit[0]
	proxyGroupId := idSplit[1]

	flowId, err := service.DeleteCynosdbProxyEndPointById(ctx, clusterId, proxyGroupId)
	if err != nil {
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
			return resource.RetryableError(fmt.Errorf("delete cynosdb proxyEndPoint is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cynosdb proxyEndPoint fail, reason:%s\n", logId, err.Error())
		return err
	}

	return nil
}
