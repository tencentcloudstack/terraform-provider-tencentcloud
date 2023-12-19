package cdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlProxy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlProxyCreate,
		Read:   resourceTencentCloudMysqlProxyRead,
		Update: resourceTencentCloudMysqlProxyUpdate,
		Delete: resourceTencentCloudMysqlProxyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"uniq_vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Vpc id.",
			},

			"uniq_subnet_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Subnet id.",
			},

			"proxy_node_custom": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Node specification configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of nodes.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of CPU cores.",
						},
						"mem": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Memory size.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region.",
						},
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Zone.",
						},
					},
				},
			},

			"security_group": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security group.",
			},

			"desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Describe.",
			},

			"connection_pool_limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Connection Pool Threshold.",
			},

			"vip": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "IP address.",
			},

			"vport": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Port.",
			},

			"proxy_version": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The current version of the database agent. No need to fill in when creating.",
			},

			"upgrade_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Upgrade time: nowTime (upgrade completed) timeWindow (instance maintenance time), Required when modifying the agent version, No need to fill in when creating.",
			},

			"proxy_group_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Proxy group id.",
			},

			"proxy_address_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Proxy address id.",
			},
		},
	}
}

func resourceTencentCloudMysqlProxyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_proxy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request    = mysql.NewCreateCdbProxyRequest()
		response   = mysql.NewCreateCdbProxyResponse()
		instanceId string
		vpcId      string
		subnetId   string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_vpc_id"); ok {
		vpcId = v.(string)
		request.UniqVpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_subnet_id"); ok {
		subnetId = v.(string)
		request.UniqSubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_node_custom"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			proxyNodeCustom := mysql.ProxyNodeCustom{}
			if v, ok := dMap["node_count"]; ok {
				proxyNodeCustom.NodeCount = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["cpu"]; ok {
				proxyNodeCustom.Cpu = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["mem"]; ok {
				proxyNodeCustom.Mem = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["region"]; ok {
				proxyNodeCustom.Region = helper.String(v.(string))
			}
			if v, ok := dMap["zone"]; ok {
				proxyNodeCustom.Zone = helper.String(v.(string))
			}
			request.ProxyNodeCustom = append(request.ProxyNodeCustom, &proxyNodeCustom)
		}
	}

	if v, ok := d.GetOk("security_group"); ok {
		securityGroupSet := v.(*schema.Set).List()
		for i := range securityGroupSet {
			securityGroup := securityGroupSet[i].(string)
			request.SecurityGroup = append(request.SecurityGroup, &securityGroup)
		}
	}

	if v, ok := d.GetOk("desc"); ok {
		request.Desc = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("connection_pool_limit"); ok {
		request.ConnectionPoolLimit = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().CreateCdbProxy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql proxy failed, reason:%+v", logId, err)
		return err
	}

	asyncRequestId := *response.Response.AsyncRequestId
	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
			return nil
		}
		if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("%s create mysql proxy status is %s", instanceId, taskStatus))
		}
		err = fmt.Errorf("%s create mysql proxy status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mysql proxy fail, reason:%s\n ", logId, err.Error())
		return err
	}

	proxy, err := service.DescribeMysqlProxyById(ctx, instanceId, "")
	if err != nil {
		return err
	}

	proxyAddressId := *proxy.ProxyAddress[0].ProxyAddressId
	proxyGroupId := *proxy.ProxyGroupId
	d.SetId(instanceId + tccommon.FILED_SP + proxyGroupId + tccommon.FILED_SP + proxyAddressId)

	ip := d.Get("vip").(string)
	port := d.Get("vport").(int)
	if ip != "" || port > 0 {
		err := service.ModifyCdbProxyAddressVipAndVPort(ctx, proxyGroupId, proxyAddressId, vpcId, subnetId, ip, uint64(port))
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudMysqlProxyRead(d, meta)
}

func resourceTencentCloudMysqlProxyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_proxy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", d.Id())
	}

	instanceId := items[0]
	proxyGroupId := items[1]
	proxyAddressId := items[2]

	proxy, err := service.DescribeMysqlProxyById(ctx, instanceId, proxyGroupId)
	if err != nil {
		return err
	}

	if proxy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlProxy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	proxyAddress := proxy.ProxyAddress[0]
	if proxyAddress.UniqVpcId != nil {
		_ = d.Set("uniq_vpc_id", proxyAddress.UniqVpcId)
	}

	if proxyAddress.UniqSubnetId != nil {
		_ = d.Set("uniq_subnet_id", proxyAddress.UniqSubnetId)
	}

	if proxy.ProxyNode != nil {
		proxyNodeCustomList := []interface{}{}
		for _, proxyNodeCustom := range proxy.ProxyNode {
			proxyNodeCustomMap := map[string]interface{}{}

			proxyNodeCustomMap["node_count"] = len(proxy.ProxyNode)

			if proxyNodeCustom.Cpu != nil {
				proxyNodeCustomMap["cpu"] = proxyNodeCustom.Cpu
			}

			if proxyNodeCustom.Mem != nil {
				proxyNodeCustomMap["mem"] = proxyNodeCustom.Mem
			}

			if proxyNodeCustom.Region != nil {
				proxyNodeCustomMap["region"] = proxyNodeCustom.Region
			}

			if proxyNodeCustom.Zone != nil {
				proxyNodeCustomMap["zone"] = proxyNodeCustom.Zone
			}

			proxyNodeCustomList = append(proxyNodeCustomList, proxyNodeCustomMap)
		}

		_ = d.Set("proxy_node_custom", proxyNodeCustomList)

	}

	securityGroups, err := service.DescribeDBSecurityGroups(ctx, proxyAddressId)
	if err != nil {
		return err
	}
	if len(securityGroups) > 0 {
		_ = d.Set("security_group", securityGroups)
	}

	if proxyAddress.Desc != nil {
		_ = d.Set("desc", proxyAddress.Desc)
	}

	if proxy.ConnectionPoolLimit != nil {
		_ = d.Set("connection_pool_limit", proxy.ConnectionPoolLimit)
	}

	if proxyAddress.Vip != nil {
		_ = d.Set("vip", proxyAddress.Vip)
	}

	if proxyAddress.VPort != nil {
		_ = d.Set("vport", proxyAddress.VPort)
	}

	_ = d.Set("proxy_group_id", proxyGroupId)

	_ = d.Set("proxy_address_id", proxyAddressId)

	if proxy.ProxyVersion != nil {
		_ = d.Set("proxy_version", proxy.ProxyVersion)
	}

	return nil
}

func resourceTencentCloudMysqlProxyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_proxy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", d.Id())
	}

	instanceId := items[0]
	proxyGroupId := items[1]
	proxyAddressId := items[2]

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	proxy, err := service.DescribeMysqlProxyById(ctx, instanceId, proxyGroupId)
	if err != nil {
		return err
	}

	if proxy == nil || proxy.ProxyGroupId == nil {
		return fmt.Errorf("Instance `%s` proxy does not exist", instanceId)
	}

	immutableArgs := []string{"instance_id", "security_group"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("proxy_node_custom") {
		request := mysql.NewAdjustCdbProxyRequest()
		response := mysql.NewAdjustCdbProxyResponse()

		request.InstanceId = &instanceId
		request.ProxyGroupId = proxy.ProxyGroupId

		if v, ok := d.GetOk("proxy_node_custom"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				proxyNodeCustom := mysql.ProxyNodeCustom{}
				if v, ok := dMap["node_count"]; ok {
					proxyNodeCustom.NodeCount = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["cpu"]; ok {
					proxyNodeCustom.Cpu = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["mem"]; ok {
					proxyNodeCustom.Mem = helper.IntUint64(v.(int))
				}
				if v, ok := dMap["region"]; ok {
					proxyNodeCustom.Region = helper.String(v.(string))
				}
				if v, ok := dMap["zone"]; ok {
					proxyNodeCustom.Zone = helper.String(v.(string))
				}
				request.ProxyNodeCustom = append(request.ProxyNodeCustom, &proxyNodeCustom)
			}
			request.ReloadBalance = helper.String("auto")
			request.UpgradeTime = helper.String("nowTime")
		}

		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().AdjustCdbProxy(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			response = result
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mysql proxy failed, reason:%+v", logId, err)
			return err
		}

		asyncRequestId := *response.Response.AsyncRequestId
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			taskStatus, message, err := service.DescribeAsyncRequestInfo(ctx, asyncRequestId)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("%s update mysql proxy status is %s", instanceId, taskStatus))
			}
			err = fmt.Errorf("%s update mysql proxy status is %s,we won't wait for it finish ,it show message:%s", instanceId, taskStatus, message)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s update mysql proxy fail, reason:%s\n ", logId, err.Error())
			return err
		}
	}

	if d.HasChange("connection_pool_limit") {
		connectionPoolLimit := d.Get("connection_pool_limit")
		request := mysql.NewModifyCdbProxyParamRequest()
		request.InstanceId = &instanceId
		request.ProxyGroupId = proxy.ProxyGroupId
		request.ConnectionPoolLimit = helper.IntUint64(connectionPoolLimit.(int))

		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().ModifyCdbProxyParam(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mysql proxy failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("uniq_vpc_id") || d.HasChange("uniq_subnet_id") || d.HasChange("vip") || d.HasChange("vport") {
		vpcId := d.Get("uniq_vpc_id").(string)
		subnetId := d.Get("uniq_subnet_id").(string)
		ip := d.Get("vip").(string)
		port := uint64(d.Get("vport").(int))

		err := service.ModifyCdbProxyAddressVipAndVPort(ctx, proxyGroupId, proxyAddressId, vpcId, subnetId, ip, port)
		if err != nil {
			return err
		}
	}

	if d.HasChange("desc") {
		desc := d.Get("desc").(string)
		err := service.ModifyCdbProxyAddressDesc(ctx, proxyGroupId, proxyAddressId, desc)
		if err != nil {
			return err
		}
	}

	if d.HasChange("proxy_version") {
		upgradeTime := ""
		if v, ok := d.GetOk("upgrade_time"); ok {
			upgradeTime = v.(string)
		} else {
			return fmt.Errorf("The parameter `upgrade_time` must be filled in when modifying the proxy version")
		}

		oldProxyVersion, proxyVersion := d.GetChange("proxy_version")
		err := service.UpgradeCDBProxyVersion(ctx, instanceId, proxyGroupId, oldProxyVersion.(string), proxyVersion.(string), upgradeTime)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudMysqlProxyRead(d, meta)
}

func resourceTencentCloudMysqlProxyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_proxy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) != 3 {
		return fmt.Errorf("invalid ID %s", d.Id())
	}

	instanceId := items[0]
	proxyGroupId := items[1]
	// proxyAddressId := items[2]

	if err := service.DeleteMysqlProxyById(ctx, instanceId); err != nil {
		return err
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		proxy, err := service.DescribeMysqlProxyById(ctx, instanceId, proxyGroupId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if proxy == nil {
			return nil
		}
		if proxy != nil {
			return resource.RetryableError(fmt.Errorf("%s delete mysql proxy status is %s", instanceId, *proxy.Status))
		}
		err = fmt.Errorf("%s delete mysql proxy status is %s,we won't wait for it finish", instanceId, *proxy.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete mysql proxy fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return nil
}
