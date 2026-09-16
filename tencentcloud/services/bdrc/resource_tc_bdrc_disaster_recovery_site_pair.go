package bdrc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	bdrcv20260330 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBdrcDisasterRecoverySitePair() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBdrcDisasterRecoverySitePairCreate,
		Read:   resourceTencentCloudBdrcDisasterRecoverySitePairRead,
		Update: resourceTencentCloudBdrcDisasterRecoverySitePairUpdate,
		Delete: resourceTencentCloudBdrcDisasterRecoverySitePairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"disaster_recovery_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Disaster recovery type, CROSS_REGION (cross-region) or CROSS_ZONE (cross-zone).",
			},

			"source_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Production site region.",
			},

			"source_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Production site availability zone.",
			},

			"target_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Disaster recovery site region.",
			},

			"target_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Disaster recovery site availability zone.",
			},

			"source_vpc": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Production site VPC.",
			},

			"target_vpc": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Disaster recovery site VPC.",
			},

			"site_pair_product_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site pair product type, including DISK, CFS, INSTANCE.",
			},

			"site_pair_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Site pair name, max length is 60 characters.",
			},

			"copy_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Replication technology, SYN (synchronous) / ASY (asynchronous).",
			},

			// computed
			"site_pair_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site pair ID.",
			},

			"site_pair_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site pair state.",
			},

			"create_from": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation source, LOCAL (local creation) / PEER (peer creation).",
			},

			"account_uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Master account Uin of the account that created the site pair.",
			},

			"sub_account_uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sub account Uin of the account that created the site pair.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"bind_protect_group_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of bound protection groups.",
			},

			"error_recovery_point_objective_copy_pair_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of copy pair IDs with RPO anomalies.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"protected_resource_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of protected resources grouped by resource type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type (consistent with SitePairType, such as DISK/CFS/INSTANCE).",
						},
						"resource_id_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of protected source resource IDs under this type.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"protected_resource_status_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Status statistics of protected resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Copy pair status.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of resources under this status.",
						},
					},
				},
			},

			"cross_cloud_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Extra information for cross-cloud scenarios (only returned when IsCrossCloud=true).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_cloud_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source cloud name (peer cloud name in cross-cloud).",
						},
						"target_cloud_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target cloud name (local cloud name in cross-cloud).",
						},
						"source_app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Source cloud AppId.",
						},
						"source_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source cloud master account Uin.",
						},
						"source_sub_account_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source cloud sub account Uin.",
						},
						"source_user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source cloud user name.",
						},
						"target_app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Target cloud AppId.",
						},
						"target_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target cloud master account Uin.",
						},
						"target_sub_account_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target cloud sub account Uin.",
						},
						"peer_region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Peer cloud region display name.",
						},
						"peer_zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Peer cloud availability zone display name.",
						},
						"peer_vpc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Peer cloud VPC display name.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudBdrcDisasterRecoverySitePairCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_disaster_recovery_site_pair.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		ctx                 = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request             = bdrcv20260330.NewCreateDisasterRecoverySitePairRequest()
		response            = bdrcv20260330.NewCreateDisasterRecoverySitePairResponse()
		sitePairProductType string
	)

	if v, ok := d.GetOk("disaster_recovery_type"); ok {
		request.DisasterRecoveryType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_region"); ok {
		request.SourceRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_zone"); ok {
		request.SourceZone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_region"); ok {
		request.TargetRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_zone"); ok {
		request.TargetZone = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source_vpc"); ok {
		request.SourceVpc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_vpc"); ok {
		request.TargetVpc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("site_pair_product_type"); ok {
		request.SitePairProductType = helper.String(v.(string))
		sitePairProductType = v.(string)
	}

	if v, ok := d.GetOk("site_pair_name"); ok {
		request.SitePairName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("copy_type"); ok {
		request.CopyType = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client().CreateDisasterRecoverySitePairWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create bdrc disaster_recovery_site_pair failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create bdrc disaster_recovery_site_pair failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create bdrc disaster_recovery_site_pair success, logId=%s, d.Id()=%s", logId, logId, *response.Response.SitePairId)

	if response.Response.SitePairId == nil || *response.Response.SitePairId == "" {
		return fmt.Errorf("Create bdrc disaster_recovery_site_pair failed, SitePairId is empty.")
	}

	sitePairId := *response.Response.SitePairId
	d.SetId(strings.Join([]string{sitePairId, sitePairProductType}, tccommon.FILED_SP))
	return resourceTencentCloudBdrcDisasterRecoverySitePairRead(d, meta)
}

func resourceTencentCloudBdrcDisasterRecoverySitePairRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_disaster_recovery_site_pair.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewBdrcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	parts := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(parts) != 2 {
		return fmt.Errorf("resource id is broken, id: %s (expected `<site_pair_id>%s<site_pair_product_type>`)", d.Id(), tccommon.FILED_SP)
	}

	sitePairId := parts[0]
	sitePairProductType := parts[1]

	respData, err := service.DescribeDisasterRecoverySitePairById(ctx, sitePairId, sitePairProductType)
	if err != nil {
		log.Printf("[CRITAL]%s read bdrc disaster_recovery_site_pair failed, reason:%+v", logId, err)
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] tencentcloud_bdrc_disaster_recovery_site_pair id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.SitePairName != nil {
		_ = d.Set("site_pair_name", respData.SitePairName)
	}

	if respData.DisasterRecoveryType != nil {
		_ = d.Set("disaster_recovery_type", respData.DisasterRecoveryType)
	}

	if respData.SourceRegion != nil {
		_ = d.Set("source_region", respData.SourceRegion)
	}

	if respData.SourceZone != nil {
		_ = d.Set("source_zone", respData.SourceZone)
	}

	if respData.TargetRegion != nil {
		_ = d.Set("target_region", respData.TargetRegion)
	}

	if respData.TargetZone != nil {
		_ = d.Set("target_zone", respData.TargetZone)
	}

	if respData.SourceVpc != nil {
		_ = d.Set("source_vpc", respData.SourceVpc)
	}

	if respData.TargetVpc != nil {
		_ = d.Set("target_vpc", respData.TargetVpc)
	}

	if respData.SitePairType != nil {
		_ = d.Set("site_pair_product_type", respData.SitePairType)
	}

	if respData.CopyType != nil {
		_ = d.Set("copy_type", respData.CopyType)
	}

	if respData.SitePairId != nil {
		_ = d.Set("site_pair_id", respData.SitePairId)
	}

	if respData.SitePairState != nil {
		_ = d.Set("site_pair_state", respData.SitePairState)
	}

	if respData.CreateFrom != nil {
		_ = d.Set("create_from", respData.CreateFrom)
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

	if respData.BindProtectGroupCount != nil {
		_ = d.Set("bind_protect_group_count", respData.BindProtectGroupCount)
	}

	// The four computed list attributes below are always written back to state,
	// even when the API returns nothing. Skipping the write leaves the
	// attribute null in state, which Terraform treats as "not yet populated"
	// and reports as a perpetual diff.
	errorRecoveryPointObjectiveCopyPairSetList := make([]string, 0, len(respData.ErrorRecoveryPointObjectiveCopyPairSet))
	for _, v := range respData.ErrorRecoveryPointObjectiveCopyPairSet {
		if v != nil {
			errorRecoveryPointObjectiveCopyPairSetList = append(errorRecoveryPointObjectiveCopyPairSetList, *v)
		}
	}
	_ = d.Set("error_recovery_point_objective_copy_pair_set", errorRecoveryPointObjectiveCopyPairSetList)

	protectedResourceSetList := make([]map[string]interface{}, 0, len(respData.ProtectedResourceSet))
	for _, protectedResource := range respData.ProtectedResourceSet {
		protectedResourceMap := map[string]interface{}{}
		if protectedResource.ResourceType != nil {
			protectedResourceMap["resource_type"] = protectedResource.ResourceType
		}

		if protectedResource.ResourceIdSet != nil {
			resourceIdSetList := make([]string, 0, len(protectedResource.ResourceIdSet))
			for _, v := range protectedResource.ResourceIdSet {
				if v != nil {
					resourceIdSetList = append(resourceIdSetList, *v)
				}
			}
			protectedResourceMap["resource_id_set"] = resourceIdSetList
		}

		protectedResourceSetList = append(protectedResourceSetList, protectedResourceMap)
	}
	_ = d.Set("protected_resource_set", protectedResourceSetList)

	protectedResourceStatusSetList := make([]map[string]interface{}, 0, len(respData.ProtectedResourceStatusSet))
	for _, protectedResourceStatus := range respData.ProtectedResourceStatusSet {
		protectedResourceStatusMap := map[string]interface{}{}
		if protectedResourceStatus.Status != nil {
			protectedResourceStatusMap["status"] = protectedResourceStatus.Status
		}

		if protectedResourceStatus.Count != nil {
			protectedResourceStatusMap["count"] = protectedResourceStatus.Count
		}

		protectedResourceStatusSetList = append(protectedResourceStatusSetList, protectedResourceStatusMap)
	}
	_ = d.Set("protected_resource_status_set", protectedResourceStatusSetList)

	crossCloudDetailsList := make([]map[string]interface{}, 0, 1)
	if respData.CrossCloudDetails != nil {
		crossCloudDetailsMap := map[string]interface{}{}
		crossCloudDetails := respData.CrossCloudDetails
		if crossCloudDetails.SourceCloudName != nil {
			crossCloudDetailsMap["source_cloud_name"] = crossCloudDetails.SourceCloudName
		}

		if crossCloudDetails.TargetCloudName != nil {
			crossCloudDetailsMap["target_cloud_name"] = crossCloudDetails.TargetCloudName
		}

		if crossCloudDetails.SourceAppId != nil {
			crossCloudDetailsMap["source_app_id"] = crossCloudDetails.SourceAppId
		}

		if crossCloudDetails.SourceUin != nil {
			crossCloudDetailsMap["source_uin"] = crossCloudDetails.SourceUin
		}

		if crossCloudDetails.SourceSubAccountUin != nil {
			crossCloudDetailsMap["source_sub_account_uin"] = crossCloudDetails.SourceSubAccountUin
		}

		if crossCloudDetails.SourceUserName != nil {
			crossCloudDetailsMap["source_user_name"] = crossCloudDetails.SourceUserName
		}

		if crossCloudDetails.TargetAppId != nil {
			crossCloudDetailsMap["target_app_id"] = crossCloudDetails.TargetAppId
		}

		if crossCloudDetails.TargetUin != nil {
			crossCloudDetailsMap["target_uin"] = crossCloudDetails.TargetUin
		}

		if crossCloudDetails.TargetSubAccountUin != nil {
			crossCloudDetailsMap["target_sub_account_uin"] = crossCloudDetails.TargetSubAccountUin
		}

		if crossCloudDetails.PeerRegionName != nil {
			crossCloudDetailsMap["peer_region_name"] = crossCloudDetails.PeerRegionName
		}

		if crossCloudDetails.PeerZoneName != nil {
			crossCloudDetailsMap["peer_zone_name"] = crossCloudDetails.PeerZoneName
		}

		if crossCloudDetails.PeerVpcName != nil {
			crossCloudDetailsMap["peer_vpc_name"] = crossCloudDetails.PeerVpcName
		}

		crossCloudDetailsList = append(crossCloudDetailsList, crossCloudDetailsMap)
	}
	_ = d.Set("cross_cloud_details", crossCloudDetailsList)

	return nil
}

func resourceTencentCloudBdrcDisasterRecoverySitePairUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_disaster_recovery_site_pair.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	parts := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(parts) != 2 {
		return fmt.Errorf("resource id is broken, id: %s (expected `<site_pair_id>%s<site_pair_product_type>`)", d.Id(), tccommon.FILED_SP)
	}

	sitePairId := parts[0]

	immutableArgs := []string{"disaster_recovery_type", "source_region", "source_zone", "target_region", "target_zone", "source_vpc", "target_vpc", "site_pair_product_type", "copy_type"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("Update bdrc disaster_recovery_site_pair failed, %s is immutable, please recreate the resource.", v)
		}
	}

	if d.HasChange("site_pair_name") {
		request := bdrcv20260330.NewModifySitePairAttributeRequest()
		request.SitePairId = helper.String(sitePairId)

		if v, ok := d.GetOk("site_pair_name"); ok {
			request.SitePairName = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client().ModifySitePairAttributeWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify bdrc disaster_recovery_site_pair failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update bdrc disaster_recovery_site_pair failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudBdrcDisasterRecoverySitePairRead(d, meta)
}

func resourceTencentCloudBdrcDisasterRecoverySitePairDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_disaster_recovery_site_pair.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bdrcv20260330.NewDeleteDisasterRecoverySitePairsRequest()
	)

	parts := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(parts) != 2 {
		return fmt.Errorf("resource id is broken, id: %s (expected `<site_pair_id>%s<site_pair_product_type>`)", d.Id(), tccommon.FILED_SP)
	}

	sitePairId := parts[0]

	request.SitePairIds = helper.Strings([]string{sitePairId})
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client().DeleteDisasterRecoverySitePairsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete bdrc disaster_recovery_site_pair failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete bdrc disaster_recovery_site_pair failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
