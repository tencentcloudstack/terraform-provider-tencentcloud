package bdrc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBdrcDisasterRecoveryProtectGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBdrcDisasterRecoveryProtectGroupCreate,
		Read:   resourceTencentCloudBdrcDisasterRecoveryProtectGroupRead,
		Update: resourceTencentCloudBdrcDisasterRecoveryProtectGroupUpdate,
		Delete: resourceTencentCloudBdrcDisasterRecoveryProtectGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"site_pair_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the disaster recovery site pair to which the protect group belongs.",
			},

			"protect_group_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Product type of the disaster recovery protect group. Valid values: `DISK`, `INSTANCE`, `CFS`.",
			},

			"recovery_point_objective": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Expected RPO of the protect group, in minutes (currently only 15 is supported).",
			},

			"protect_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the protect group, up to 60 characters.",
			},

			"data_direction": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Data replication direction. Valid values: `POSITIVE`, `REVERSE`.",
			},

			"app_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "User AppId.",
			},

			"site_pair_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the disaster recovery site pair.",
			},

			"source_region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source region.",
			},

			"source_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source zone.",
			},

			"source_vpc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source VPC.",
			},

			"target_region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target region.",
			},

			"target_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target zone.",
			},

			"target_vpc": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target VPC.",
			},

			"copy_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Replication technology (SYN sync / ASY async).",
			},

			"disaster_recovery_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Disaster recovery type (CROSS_ZONE / CROSS_REGION / CROSS_CLOUD).",
			},

			"peer_cloud_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Peer cloud name (only returned when DisasterRecoveryType is CROSS_CLOUD).",
			},

			"create_from": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation source (LOCAL local / PEER peer).",
			},

			"life_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Lifecycle state.",
			},

			"account_uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account Uin of the protect group owner.",
			},

			"sub_account_uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sub account Uin of the protect group creator.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"modify_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification time.",
			},

			"bind_protected_resource_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of protected resources bound to the protect group.",
			},

			"error_recovery_point_objective_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of replication pairs whose RPO is abnormal (not synced for more than 15 minutes).",
			},

			"protected_resource_status_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Protected resource status statistics, key is the replication pair status, value is the resource count under that status.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Replication pair status.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource count under this status.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudBdrcDisasterRecoveryProtectGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_disaster_recovery_protect_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = bdrcv20260330.NewCreateDisasterRecoveryProtectGroupRequest()
		response = bdrcv20260330.NewCreateDisasterRecoveryProtectGroupResponse()
	)

	if v, ok := d.GetOk("site_pair_id"); ok {
		request.SitePairId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protect_group_type"); ok {
		request.ProtectGroupType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("recovery_point_objective"); ok {
		request.RecoveryPointObjective = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("protect_group_name"); ok {
		request.ProtectGroupName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("data_direction"); ok {
		request.DataDirection = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client().CreateDisasterRecoveryProtectGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create bdrc disaster_recovery_protect_group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create bdrc disaster_recovery_protect_group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[CRITAL]%s create bdrc disaster_recovery_protect_group current id=%s", logId, d.Id())
	if response.Response.ProtectGroupId == nil {
		return fmt.Errorf("ProtectGroupId is nil.")
	}

	d.SetId(*response.Response.ProtectGroupId)
	return resourceTencentCloudBdrcDisasterRecoveryProtectGroupRead(d, meta)
}

func resourceTencentCloudBdrcDisasterRecoveryProtectGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_disaster_recovery_protect_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = BdrcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	protectGroupId := d.Id()
	protectGroupType := d.Get("protect_group_type").(string)

	respData, err := service.DescribeDisasterRecoveryProtectGroupById(ctx, protectGroupId, protectGroupType)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] bdrc disaster_recovery_protect_group id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.SitePairId != nil {
		_ = d.Set("site_pair_id", respData.SitePairId)
	}

	if respData.ProtectGroupType != nil {
		_ = d.Set("protect_group_type", respData.ProtectGroupType)
	}

	if respData.RecoveryPointObjective != nil {
		_ = d.Set("recovery_point_objective", respData.RecoveryPointObjective)
	}

	if respData.ProtectGroupName != nil {
		_ = d.Set("protect_group_name", respData.ProtectGroupName)
	}

	if respData.DataDirection != nil {
		_ = d.Set("data_direction", respData.DataDirection)
	}

	if respData.AppId != nil {
		_ = d.Set("app_id", respData.AppId)
	}

	if respData.SitePairName != nil {
		_ = d.Set("site_pair_name", respData.SitePairName)
	}

	if respData.SourceRegion != nil {
		_ = d.Set("source_region", respData.SourceRegion)
	}

	if respData.SourceZone != nil {
		_ = d.Set("source_zone", respData.SourceZone)
	}

	if respData.SourceVpc != nil {
		_ = d.Set("source_vpc", respData.SourceVpc)
	}

	if respData.TargetRegion != nil {
		_ = d.Set("target_region", respData.TargetRegion)
	}

	if respData.TargetZone != nil {
		_ = d.Set("target_zone", respData.TargetZone)
	}

	if respData.TargetVpc != nil {
		_ = d.Set("target_vpc", respData.TargetVpc)
	}

	if respData.CopyType != nil {
		_ = d.Set("copy_type", respData.CopyType)
	}

	if respData.DisasterRecoveryType != nil {
		_ = d.Set("disaster_recovery_type", respData.DisasterRecoveryType)
	}

	if respData.PeerCloudName != nil {
		_ = d.Set("peer_cloud_name", respData.PeerCloudName)
	}

	if respData.CreateFrom != nil {
		_ = d.Set("create_from", respData.CreateFrom)
	}

	if respData.LifeState != nil {
		_ = d.Set("life_state", respData.LifeState)
	}

	if respData.AccountUin != nil {
		_ = d.Set("account_uin", respData.AccountUin)
	}

	if respData.SubAccountUin != nil {
		_ = d.Set("sub_account_uin", respData.SubAccountUin)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.ModifyTime != nil {
		_ = d.Set("modify_time", respData.ModifyTime)
	}

	if respData.BindProtectedResourceCount != nil {
		_ = d.Set("bind_protected_resource_count", respData.BindProtectedResourceCount)
	}

	if respData.ErrorRecoveryPointObjectiveCount != nil {
		_ = d.Set("error_recovery_point_objective_count", respData.ErrorRecoveryPointObjectiveCount)
	}

	if respData.ProtectedResourceStatusSet != nil && len(respData.ProtectedResourceStatusSet) > 0 {
		statusSetList := make([]map[string]interface{}, 0, len(respData.ProtectedResourceStatusSet))
		for _, status := range respData.ProtectedResourceStatusSet {
			statusMap := map[string]interface{}{}
			if status.Status != nil {
				statusMap["status"] = status.Status
			}

			if status.Count != nil {
				statusMap["count"] = status.Count
			}

			statusSetList = append(statusSetList, statusMap)
		}

		_ = d.Set("protected_resource_status_set", statusSetList)
	}

	return nil
}

func resourceTencentCloudBdrcDisasterRecoveryProtectGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_disaster_recovery_protect_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	immutableArgs := []string{"site_pair_id", "protect_group_type", "recovery_point_objective", "data_direction"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("bdrc disaster_recovery_protect_group `%s` is immutable, cannot be updated, please recreate.", v)
		}
	}

	if d.HasChange("protect_group_name") {
		request := bdrcv20260330.NewModifyProtectGroupAttributeRequest()
		protectGroupId := d.Id()
		request.ProtectGroupId = helper.String(protectGroupId)
		if v, ok := d.GetOk("protect_group_name"); ok {
			request.ProtectGroupName = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client().ModifyProtectGroupAttributeWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Update bdrc disaster_recovery_protect_group failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update bdrc disaster_recovery_protect_group failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudBdrcDisasterRecoveryProtectGroupRead(d, meta)
}

func resourceTencentCloudBdrcDisasterRecoveryProtectGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_disaster_recovery_protect_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bdrcv20260330.NewDeleteDisasterRecoveryProtectGroupsRequest()
	)

	protectGroupId := d.Id()
	request.ProtectGroups = []*string{helper.String(protectGroupId)}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client().DeleteDisasterRecoveryProtectGroupsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete bdrc disaster_recovery_protect_group failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete bdrc disaster_recovery_protect_group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
