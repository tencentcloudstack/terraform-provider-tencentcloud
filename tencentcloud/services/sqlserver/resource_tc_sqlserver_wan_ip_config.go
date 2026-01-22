package sqlserver

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSqlserverWanIpConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverWanIpConfigCreate,
		Read:   resourceTencentCloudSqlserverWanIpConfigRead,
		Update: resourceTencentCloudSqlserverWanIpConfigUpdate,
		Delete: resourceTencentCloudSqlserverWanIpConfigDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID.",
			},

			"ro_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Read only group ID.",
			},

			"enable_wan_ip": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether to open wan ip, true: enable; false: disable.",
			},

			// computed
			"dns_pod_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internet address domain name.",
			},

			"tgw_wan_vport": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "External port number.",
			},

			"ro_group": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Read only group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_pod_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Internet address domain name.",
						},

						"tgw_wan_vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "External port number.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSqlserverWanIpConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_wan_ip_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		instanceId string
		roGroupId  string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("ro_group_id"); ok {
		roGroupId = v.(string)
	}

	if roGroupId != "" {
		d.SetId(strings.Join([]string{instanceId, roGroupId}, tccommon.FILED_SP))
	} else {
		d.SetId(instanceId)
	}

	return resourceTencentCloudSqlserverWanIpConfigUpdate(d, meta)
}

func resourceTencentCloudSqlserverWanIpConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_wan_ip_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()
	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		sqlserverService = SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId       string
		roGroupId        string
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 1 {
		instanceId = idSplit[0]
	} else if len(idSplit) == 2 {
		instanceId = idSplit[0]
		roGroupId = idSplit[1]
	} else {
		return fmt.Errorf("tencentcloud_sqlserver_wan_ip_config id is broken, id is %s", d.Id())
	}

	instance, has, err := sqlserverService.DescribeSqlserverInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if !has {
		d.SetId("")
		return nil
	}

	if instance.DnsPodDomain != nil {
		_ = d.Set("dns_pod_domain", instance.DnsPodDomain)
	}

	if instance.TgwWanVPort != nil {
		_ = d.Set("tgw_wan_vport", instance.TgwWanVPort)
	}

	if roGroupId != "" {
		roGroupList, err := sqlserverService.DescribeReadonlyGroupList(ctx, instanceId)
		if err != nil {
			return err
		}

		if roGroupList == nil {
			d.SetId("")
			log.Printf("[WARN]%s resource `SqlservereReadonlyGroup` [%s] not found, please check if it has been deleted.", logId, d.Id())
			return nil
		}

		tmpList := make([]map[string]interface{}, 0, len(roGroupList))
		for _, v := range roGroupList {
			if v.ReadOnlyGroupId != nil && *v.ReadOnlyGroupId == roGroupId {
				dMap := map[string]interface{}{}
				if v.DnsPodDomain != nil {
					dMap["dns_pod_domain"] = v.DnsPodDomain
				}

				if v.TgwWanVPort != nil {
					dMap["tgw_wan_vport"] = v.TgwWanVPort
				}

				tmpList = append(tmpList, dMap)
			}
		}

		_ = d.Set("ro_group", tmpList)
	}

	return nil
}

func resourceTencentCloudSqlserverWanIpConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_wan_ip_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		instanceId  string
		roGroupId   string
		enableWanIp bool
		flowId      uint64
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 1 {
		instanceId = idSplit[0]
	} else if len(idSplit) == 2 {
		instanceId = idSplit[0]
		roGroupId = idSplit[1]
	} else {
		return fmt.Errorf("tencentcloud_sqlserver_wan_ip_config id is broken, id is %s", d.Id())
	}

	if v, ok := d.GetOkExists("enable_wan_ip"); ok {
		enableWanIp = v.(bool)
	}

	if enableWanIp {
		request := sqlserver.NewModifyOpenWanIpRequest()
		response := sqlserver.NewModifyOpenWanIpResponse()
		request.InstanceId = &instanceId
		if roGroupId != "" {
			request.RoGroupId = &roGroupId
		}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().ModifyOpenWanIpWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Open sqlserver wan ip failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s open sqlserver wan ip failed, reason:%+v", logId, err)
			return err
		}

		if response.Response.FlowId == nil {
			return fmt.Errorf("FlowId is nil.")
		}

		flowId = *response.Response.FlowId
	} else {
		request := sqlserver.NewModifyCloseWanIpRequest()
		response := sqlserver.NewModifyCloseWanIpResponse()
		request.InstanceId = &instanceId
		if roGroupId != "" {
			request.RoGroupId = &roGroupId
		}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().ModifyCloseWanIpWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Close sqlserver wan ip failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s close sqlserver wan ip failed, reason:%+v", logId, err)
			return err
		}

		if response.Response.FlowId == nil {
			return fmt.Errorf("FlowId is nil.")
		}

		flowId = *response.Response.FlowId
	}

	// wait
	flowRequest := sqlserver.NewDescribeFlowStatusRequest()
	flowRequest.FlowId = helper.UInt64Int64(flowId)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSqlserverClient().DescribeFlowStatus(flowRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe FlowStatus failed, Response is nil."))
		}

		if *result.Response.Status == SQLSERVER_TASK_SUCCESS {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Operate sqlserver wan ip status is %d.", *result.Response.Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate sqlserver wan ip failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverWanIpConfigRead(d, meta)
}

func resourceTencentCloudSqlserverWanIpConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_sqlserver_wan_ip_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
