package cdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlProxyAddressConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlProxyAddressConfigCreate,
		Read:   resourceTencentCloudMysqlProxyAddressConfigRead,
		Update: resourceTencentCloudMysqlProxyAddressConfigUpdate,
		Delete: resourceTencentCloudMysqlProxyAddressConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID, such as: cdb-xxxxxxxx.",
			},

			"proxy_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Proxy group ID, such as: proxy-xxxxxxxx. Can be obtained through the DescribeCdbProxyInfo interface.",
			},

			"proxy_address_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Proxy address ID, such as: proxyaddr-xxxxxxxx. Can be obtained through the DescribeCdbProxyInfo interface.",
			},

			"weight_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Weight allocation mode. Valid values: `system` (system auto-allocation), `custom` (custom).",
			},

			"is_kick_out": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable delay elimination. Valid values: `true`, `false`.",
			},

			"min_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Minimum reserved quantity. Minimum value: 0. Note: only valid when IsKickOut is true.",
			},

			"max_delay": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Delay elimination threshold in milliseconds. Value range: [1, 10000].",
			},

			"fail_over": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to enable failover. Valid values: `true`, `false`.",
			},

			"auto_add_ro": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to automatically add read-only instances. Valid values: `true`, `false`.",
			},

			"read_only": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether it is read-only. Valid values: `true`, `false`.",
			},

			"trans_split": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable transaction splitting. Default value: `false`.",
			},

			"connection_pool": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable connection pool. Default is disabled. Note: for MySQL 8.0, the kernel minor version must be >= MySQL 8.0 20230630.",
			},

			"proxy_allocation": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Read/write weight allocation. If WeightMode is `system`, the input weight does not take effect.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Region, such as: ap-guangzhou.",
						},
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Availability zone, such as: ap-guangzhou-2.",
						},
						"proxy_instance": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Proxy instance list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Instance ID.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Weight value.",
									},
								},
							},
						},
					},
				},
			},

			"auto_load_balance": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable adaptive load balancing. Default is disabled.",
			},

			"access_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Access mode. Valid values: `nearby` (nearby access), `balance` (load balancing). Default: `nearby`.",
			},

			"ap_node_as_ro_node": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to treat libra nodes as regular RO nodes.",
			},

			"ap_query_to_other_node": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When libra node fails, whether to forward to other nodes.",
			},
		},
	}
}

func resourceTencentCloudMysqlProxyAddressConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_proxy_address_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId     string
		proxyGroupId   string
		proxyAddressId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("proxy_group_id"); ok {
		proxyGroupId = v.(string)
	}

	if v, ok := d.GetOk("proxy_address_id"); ok {
		proxyAddressId = v.(string)
	}

	d.SetId(strings.Join([]string{instanceId, proxyGroupId, proxyAddressId}, tccommon.FILED_SP))
	return resourceTencentCloudMysqlProxyAddressConfigUpdate(d, meta)
}

func resourceTencentCloudMysqlProxyAddressConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_proxy_address_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	proxyGroupId := idSplit[1]
	proxyAddressId := idSplit[2]

	respData, err := service.DescribeMysqlProxyAddressConfig(ctx, instanceId, proxyGroupId, proxyAddressId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_mysql_proxy_address_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("proxy_group_id", proxyGroupId)

	if respData.ProxyAddressId != nil {
		_ = d.Set("proxy_address_id", respData.ProxyAddressId)
	}

	if respData.WeightMode != nil {
		_ = d.Set("weight_mode", respData.WeightMode)
	}

	if respData.IsKickOut != nil {
		_ = d.Set("is_kick_out", respData.IsKickOut)
	}

	if respData.MinCount != nil {
		_ = d.Set("min_count", respData.MinCount)
	}

	if respData.MaxDelay != nil {
		_ = d.Set("max_delay", respData.MaxDelay)
	}

	if respData.FailOver != nil {
		_ = d.Set("fail_over", respData.FailOver)
	}

	if respData.AutoAddRo != nil {
		_ = d.Set("auto_add_ro", respData.AutoAddRo)
	}

	if respData.ReadOnly != nil {
		_ = d.Set("read_only", respData.ReadOnly)
	}

	if respData.TransSplit != nil {
		_ = d.Set("trans_split", respData.TransSplit)
	}

	if respData.ConnectionPool != nil {
		_ = d.Set("connection_pool", respData.ConnectionPool)
	}

	if respData.AutoLoadBalance != nil {
		_ = d.Set("auto_load_balance", respData.AutoLoadBalance)
	}

	if respData.AccessMode != nil {
		_ = d.Set("access_mode", respData.AccessMode)
	}

	if respData.ApNodeAsRoNode != nil {
		_ = d.Set("ap_node_as_ro_node", respData.ApNodeAsRoNode)
	}

	if respData.ApQueryToOtherNode != nil {
		_ = d.Set("ap_query_to_other_node", respData.ApQueryToOtherNode)
	}

	if len(respData.ProxyAllocation) > 0 {
		proxyAllocationList := make([]map[string]interface{}, 0, len(respData.ProxyAllocation))
		for _, allocation := range respData.ProxyAllocation {
			allocationMap := map[string]interface{}{}
			if allocation.Region != nil {
				allocationMap["region"] = allocation.Region
			}

			if allocation.Zone != nil {
				allocationMap["zone"] = allocation.Zone
			}

			if len(allocation.ProxyInstance) > 0 {
				proxyInstanceList := make([]map[string]interface{}, 0, len(allocation.ProxyInstance))
				for _, inst := range allocation.ProxyInstance {
					instMap := map[string]interface{}{}
					if inst.InstanceId != nil {
						instMap["instance_id"] = inst.InstanceId
					}

					if inst.Weight != nil {
						instMap["weight"] = inst.Weight
					}

					proxyInstanceList = append(proxyInstanceList, instMap)
				}

				allocationMap["proxy_instance"] = proxyInstanceList
			}

			proxyAllocationList = append(proxyAllocationList, allocationMap)
		}

		_ = d.Set("proxy_allocation", proxyAllocationList)
	}

	return nil
}

func resourceTencentCloudMysqlProxyAddressConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_proxy_address_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	proxyGroupId := idSplit[1]
	proxyAddressId := idSplit[2]

	// Validate that the proxy group exists and the proxy address ID is valid before calling AdjustCdbProxyAddress.
	proxyInfo, err := service.DescribeMysqlProxyById(ctx, instanceId, proxyGroupId)
	if err != nil {
		return fmt.Errorf("DescribeCdbProxyInfo failed before update: %w", err)
	}

	if proxyInfo == nil {
		return fmt.Errorf("proxy group [%s] not found in instance [%s]", proxyGroupId, instanceId)
	}

	found := false
	for _, addr := range proxyInfo.ProxyAddress {
		if addr.ProxyAddressId != nil && *addr.ProxyAddressId == proxyAddressId {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("proxy address [%s] not found in proxy group [%s], instance [%s]", proxyAddressId, proxyGroupId, instanceId)
	}

	request := mysql.NewAdjustCdbProxyAddressRequest()
	response := mysql.NewAdjustCdbProxyAddressResponse()
	request.ProxyGroupId = &proxyGroupId
	request.ProxyAddressId = &proxyAddressId

	if v, ok := d.GetOk("weight_mode"); ok {
		request.WeightMode = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("is_kick_out"); ok {
		request.IsKickOut = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("min_count"); ok {
		request.MinCount = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("max_delay"); ok {
		request.MaxDelay = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("fail_over"); ok {
		request.FailOver = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("auto_add_ro"); ok {
		request.AutoAddRo = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("read_only"); ok {
		request.ReadOnly = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("trans_split"); ok {
		request.TransSplit = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("connection_pool"); ok {
		request.ConnectionPool = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("auto_load_balance"); ok {
		request.AutoLoadBalance = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("access_mode"); ok {
		request.AccessMode = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("ap_node_as_ro_node"); ok {
		if v.(bool) {
			request.ApNodeAsRoNode = helper.Bool(true)
		} else {
			request.ApNodeAsRoNode = helper.Bool(false)
		}
	}

	if v, ok := d.GetOkExists("ap_query_to_other_node"); ok {
		if v.(bool) {
			request.ApQueryToOtherNode = helper.Bool(true)
		} else {
			request.ApQueryToOtherNode = helper.Bool(false)
		}
	}

	if v, ok := d.GetOk("proxy_allocation"); ok {
		for _, item := range v.([]interface{}) {
			allocationMap := item.(map[string]interface{})
			allocation := mysql.ProxyAllocation{}
			if v, ok := allocationMap["region"].(string); ok && v != "" {
				allocation.Region = helper.String(v)
			}

			if v, ok := allocationMap["zone"].(string); ok && v != "" {
				allocation.Zone = helper.String(v)
			}

			if v, ok := allocationMap["proxy_instance"]; ok {
				for _, instItem := range v.([]interface{}) {
					instMap := instItem.(map[string]interface{})
					inst := mysql.ProxyInst{}
					if v, ok := instMap["instance_id"].(string); ok && v != "" {
						inst.InstanceId = helper.String(v)
					}

					if v, ok := instMap["weight"].(int); ok {
						inst.Weight = helper.IntUint64(v)
					}

					allocation.ProxyInstance = append(allocation.ProxyInstance, &inst)
				}
			}

			request.ProxyAllocation = append(request.ProxyAllocation, &allocation)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().AdjustCdbProxyAddressWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.AsyncRequestId == nil {
			return resource.NonRetryableError(fmt.Errorf("Update mysql proxy address config failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update mysql proxy address config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Poll DescribeCdbProxyInfo until TaskStatus is empty, indicating the async task has completed.
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
			return resource.RetryableError(fmt.Errorf("%s update mysql proxy address config status is %s", instanceId, taskStatus))
		}

		err = fmt.Errorf("%s update mysql proxy address config status is %s, we won't wait for it finish, it show message:%s\n", instanceId, taskStatus, message)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s update mysql proxy address config fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudMysqlProxyAddressConfigRead(d, meta)
}

func resourceTencentCloudMysqlProxyAddressConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_proxy_address_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
