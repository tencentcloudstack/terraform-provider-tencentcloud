/*
Provides a resource to create a cynosdb proxy_end_point

Example Usage

```hcl
resource "tencentcloud_cynosdb_proxy_end_point" "proxy_end_point" {
  cluster_id = "cynosdbmysql-xxxxxxx"
  unique_vpc_id = "无"
  unique_subnet_id = "无"
  connection_pool_type = "SessionConnectionPool"
  open_connection_pool = "yes"
  connection_pool_time_out = 0
  security_group_ids =
  description = "无"
  vip = "无"
  weight_mode = "无"
  auto_add_ro = "无"
  fail_over = "无"
  consistency_type = "无"
  rw_type = "无"
  consistency_time_out = 无
  trans_split = 无
  access_mode = "无"
  instance_weights {
		instance_id = ""
		weight =

  }
}
```

Import

cynosdb proxy_end_point can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_proxy_end_point.proxy_end_point proxy_end_point_id
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

func resourceTencentCloudCynosdbProxyEndPoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbProxyEndPointCreate,
		Read:   resourceTencentCloudCynosdbProxyEndPointRead,
		Update: resourceTencentCloudCynosdbProxyEndPointUpdate,
		Delete: resourceTencentCloudCynosdbProxyEndPointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

			"vip": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "VIP Information.",
			},

			"weight_mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Weight mode: system system allocation, custom customization.",
			},

			"auto_add_ro": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Do you want to automatically add read-only instances? Yes - Yes, no - Do not automatically add.",
			},

			"fail_over": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Enable Failover.",
			},

			"consistency_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Consistency type: event, global, session.",
			},

			"rw_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Read and write attributes: READWRITE, READONLY.",
			},

			"consistency_time_out": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Consistency timeout.",
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

			"instance_weights": {
				Optional:    true,
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
		},
	}
}

func resourceTencentCloudCynosdbProxyEndPointCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy_end_point.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cynosdb.NewCreateProxyEndPointRequest()
		response  = cynosdb.NewCreateProxyEndPointResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
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

	if v, ok := d.GetOk("vip"); ok {
		request.Vip = helper.String(v.(string))
	}

	if v, ok := d.GetOk("weight_mode"); ok {
		request.WeightMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auto_add_ro"); ok {
		request.AutoAddRo = helper.String(v.(string))
	}

	if v, ok := d.GetOk("fail_over"); ok {
		request.FailOver = helper.String(v.(string))
	}

	if v, ok := d.GetOk("consistency_type"); ok {
		request.ConsistencyType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rw_type"); ok {
		request.RwType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("consistency_time_out"); ok {
		request.ConsistencyTimeOut = helper.IntInt64(v.(int))
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

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbProxyEndPointStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbProxyEndPointRead(d, meta)
}

func resourceTencentCloudCynosdbProxyEndPointRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy_end_point.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	proxyEndPointId := d.Id()

	proxyEndPoint, err := service.DescribeCynosdbProxyEndPointById(ctx, clusterId)
	if err != nil {
		return err
	}

	if proxyEndPoint == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbProxyEndPoint` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if proxyEndPoint.ClusterId != nil {
		_ = d.Set("cluster_id", proxyEndPoint.ClusterId)
	}

	if proxyEndPoint.UniqueVpcId != nil {
		_ = d.Set("unique_vpc_id", proxyEndPoint.UniqueVpcId)
	}

	if proxyEndPoint.UniqueSubnetId != nil {
		_ = d.Set("unique_subnet_id", proxyEndPoint.UniqueSubnetId)
	}

	if proxyEndPoint.ConnectionPoolType != nil {
		_ = d.Set("connection_pool_type", proxyEndPoint.ConnectionPoolType)
	}

	if proxyEndPoint.OpenConnectionPool != nil {
		_ = d.Set("open_connection_pool", proxyEndPoint.OpenConnectionPool)
	}

	if proxyEndPoint.ConnectionPoolTimeOut != nil {
		_ = d.Set("connection_pool_time_out", proxyEndPoint.ConnectionPoolTimeOut)
	}

	if proxyEndPoint.SecurityGroupIds != nil {
		_ = d.Set("security_group_ids", proxyEndPoint.SecurityGroupIds)
	}

	if proxyEndPoint.Description != nil {
		_ = d.Set("description", proxyEndPoint.Description)
	}

	if proxyEndPoint.Vip != nil {
		_ = d.Set("vip", proxyEndPoint.Vip)
	}

	if proxyEndPoint.WeightMode != nil {
		_ = d.Set("weight_mode", proxyEndPoint.WeightMode)
	}

	if proxyEndPoint.AutoAddRo != nil {
		_ = d.Set("auto_add_ro", proxyEndPoint.AutoAddRo)
	}

	if proxyEndPoint.FailOver != nil {
		_ = d.Set("fail_over", proxyEndPoint.FailOver)
	}

	if proxyEndPoint.ConsistencyType != nil {
		_ = d.Set("consistency_type", proxyEndPoint.ConsistencyType)
	}

	if proxyEndPoint.RwType != nil {
		_ = d.Set("rw_type", proxyEndPoint.RwType)
	}

	if proxyEndPoint.ConsistencyTimeOut != nil {
		_ = d.Set("consistency_time_out", proxyEndPoint.ConsistencyTimeOut)
	}

	if proxyEndPoint.TransSplit != nil {
		_ = d.Set("trans_split", proxyEndPoint.TransSplit)
	}

	if proxyEndPoint.AccessMode != nil {
		_ = d.Set("access_mode", proxyEndPoint.AccessMode)
	}

	if proxyEndPoint.InstanceWeights != nil {
		instanceWeightsList := []interface{}{}
		for _, instanceWeights := range proxyEndPoint.InstanceWeights {
			instanceWeightsMap := map[string]interface{}{}

			if proxyEndPoint.InstanceWeights.InstanceId != nil {
				instanceWeightsMap["instance_id"] = proxyEndPoint.InstanceWeights.InstanceId
			}

			if proxyEndPoint.InstanceWeights.Weight != nil {
				instanceWeightsMap["weight"] = proxyEndPoint.InstanceWeights.Weight
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

	logId := getLogId(contextNil)

	var (
		switchProxyVpcRequest  = cynosdb.NewSwitchProxyVpcRequest()
		switchProxyVpcResponse = cynosdb.NewSwitchProxyVpcResponse()
	)

	proxyEndPointId := d.Id()

	request.ClusterId = &clusterId

	immutableArgs := []string{"cluster_id", "unique_vpc_id", "unique_subnet_id", "connection_pool_type", "open_connection_pool", "connection_pool_time_out", "security_group_ids", "description", "vip", "weight_mode", "auto_add_ro", "fail_over", "consistency_type", "rw_type", "consistency_time_out", "trans_split", "access_mode", "instance_weights"}

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

	if d.HasChange("connection_pool_type") {
		if v, ok := d.GetOk("connection_pool_type"); ok {
			request.ConnectionPoolType = helper.String(v.(string))
		}
	}

	if d.HasChange("open_connection_pool") {
		if v, ok := d.GetOk("open_connection_pool"); ok {
			request.OpenConnectionPool = helper.String(v.(string))
		}
	}

	if d.HasChange("connection_pool_time_out") {
		if v, ok := d.GetOkExists("connection_pool_time_out"); ok {
			request.ConnectionPoolTimeOut = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("weight_mode") {
		if v, ok := d.GetOk("weight_mode"); ok {
			request.WeightMode = helper.String(v.(string))
		}
	}

	if d.HasChange("auto_add_ro") {
		if v, ok := d.GetOk("auto_add_ro"); ok {
			request.AutoAddRo = helper.String(v.(string))
		}
	}

	if d.HasChange("fail_over") {
		if v, ok := d.GetOk("fail_over"); ok {
			request.FailOver = helper.String(v.(string))
		}
	}

	if d.HasChange("consistency_type") {
		if v, ok := d.GetOk("consistency_type"); ok {
			request.ConsistencyType = helper.String(v.(string))
		}
	}

	if d.HasChange("rw_type") {
		if v, ok := d.GetOk("rw_type"); ok {
			request.RwType = helper.String(v.(string))
		}
	}

	if d.HasChange("consistency_time_out") {
		if v, ok := d.GetOkExists("consistency_time_out"); ok {
			request.ConsistencyTimeOut = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("trans_split") {
		if v, ok := d.GetOkExists("trans_split"); ok {
			request.TransSplit = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("access_mode") {
		if v, ok := d.GetOk("access_mode"); ok {
			request.AccessMode = helper.String(v.(string))
		}
	}

	if d.HasChange("instance_weights") {
		if v, ok := d.GetOk("instance_weights"); ok {
			for _, item := range v.([]interface{}) {
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
		log.Printf("[CRITAL]%s update cynosdb proxyEndPoint failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbProxyEndPointRead(d, meta)
}

func resourceTencentCloudCynosdbProxyEndPointDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_proxy_end_point.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	proxyEndPointId := d.Id()

	if err := service.DeleteCynosdbProxyEndPointById(ctx, clusterId); err != nil {
		return err
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbProxyEndPointStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
