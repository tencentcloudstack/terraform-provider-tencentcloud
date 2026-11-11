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

func ResourceTencentCloudBdrcInstanceCopyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBdrcInstanceCopyPairCreate,
		Read:   resourceTencentCloudBdrcInstanceCopyPairRead,
		Update: resourceTencentCloudBdrcInstanceCopyPairUpdate,
		Delete: resourceTencentCloudBdrcInstanceCopyPairDelete,
		Schema: map[string]*schema.Schema{
			"protect_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Protect group ID.",
			},

			"create_target_instance_parameters": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				MinItems:    1,
				Description: "Target CVM creation parameters list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Source CVM ID.",
						},
						"instance_charge_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance billing mode.",
						},
						"placement": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "Instance placement.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Zone ID.",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Project ID.",
									},
									"host_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Dedicated host ID.",
									},
									"host_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Dedicated host ID list.",
									},
									"project_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Project name.",
									},
								},
							},
						},
						"image_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Image ID.",
						},
						"system_disk": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "System disk.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Cloud disk type.",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Cloud disk size in GB.",
									},
								},
							},
						},
						"instance_charge_prepaid": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Prepaid billing settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"period": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Purchase duration in months.",
									},
									"renew_flag": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Auto renewal flag.",
									},
								},
							},
						},
						"instance_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instance type.",
						},
						"data_disks": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Data disk list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Cloud disk type.",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Cloud disk size in GB.",
									},
									"delete_with_instance": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Delete with instance.",
									},
								},
							},
						},
						"virtual_private_cloud": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "VPC config.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "VPC ID.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Subnet ID.",
									},
									"subnet_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Subnet name.",
									},
									"as_vpc_gateway": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Used as public network gateway.",
									},
									"private_ip_addresses": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Private IP address list.",
									},
									"vpc_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "VPC name.",
									},
									"ipv6_address_count": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Number of IPv6 addresses.",
									},
								},
							},
						},
						"internet_accessible": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Public bandwidth config.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"internet_charge_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Network billing type.",
									},
									"internet_max_bandwidth_out": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Public network outbound bandwidth cap in Mbps.",
									},
									"public_ip_assigned": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to assign a public IP.",
									},
									"internet_service_provider": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Network service provider.",
									},
								},
							},
						},
						"instance_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instance display name.",
						},
						"login_settings": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Login settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"password": {
										Type:        schema.TypeString,
										Optional:    true,
										Sensitive:   true,
										Description: "Instance login password.",
									},
									"key_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Key ID list.",
									},
									"keep_image_login": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Keep image login settings.",
									},
								},
							},
						},
						"enhanced_service": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Enhanced service config.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_service": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Security service.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to enable.",
												},
											},
										},
									},
									"monitor_service": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Monitor service.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to enable.",
												},
											},
										},
									},
									"automation_service": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Automation service.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to enable.",
												},
											},
										},
									},
									"basic_service": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Basic service.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to enable.",
												},
											},
										},
									},
								},
							},
						},
						"spot_price": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Spot instance max bid.",
						},
						"host_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Instance hostname.",
						},
						"user_data": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User data for the instance.",
						},
						"disaster_recover_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Placement group ID list.",
						},
						"stopped_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Shutdown billing mode.",
						},
						"copy_pair_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Copy pair ID for drill scenarios.",
						},
						"recovery_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Recovery time point for drill scenarios.",
						},
					},
				},
			},

			"instance_copy_pair_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Copy pair name.",
			},

			"recovery_point_objective": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "User-desired RPO in minutes.",
			},

			"delete_target_resource": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to delete the disaster-recovery site disk on destroy.",
			},

			// computed
			"copy_pair_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Copy pair ID.",
			},
			"copy_pair_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Copy pair name from the describe response.",
			},
			"copy_pair_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Copy pair state.",
			},
			"copy_pair_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Copy pair type.",
			},
			"site_pair_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site pair ID.",
			},
			"site_pair_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site pair name.",
			},
			"protect_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Protect group name.",
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
			"source_resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source resource ID.",
			},
			"target_resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target resource ID.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance ID.",
			},
			"instance_copy_pair_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CVM copy pair ID.",
			},
			"percent": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Replication progress percent.",
			},
			"latest_protection_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Latest protection time.",
			},
			"data_direction": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data direction.",
			},
			"create_from": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation source.",
			},
			"disaster_recovery_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Disaster recovery type.",
			},
			"peer_cloud_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Peer cloud name.",
			},
			"rollbacking": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether rollback is in progress.",
			},
			"rollback_percent": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Rollback progress.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},
			"account_uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account Uin.",
			},
			"sub_account_uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sub-account Uin.",
			},
			"drill_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Drill group ID.",
			},
			"protection_time_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Protection time points.",
			},
			"disk_copy_pair_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Disk copy pair list for CVM.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"copy_pair_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk copy pair ID.",
						},
						"copy_pair_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Disk copy pair name.",
						},
						"source_resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source resource ID.",
						},
						"target_resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target resource ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
					},
				},
			},
			"deferred_create": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether deferred creation mode.",
			},
			"target_cvm_created": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether target CVM is actually created.",
			},
			"cvm_create_params": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CVM creation params JSON string.",
			},
		},
	}
}

func resourceTencentCloudBdrcInstanceCopyPairCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_instance_copy_pair.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = bdrcv20260330.NewCreateInstanceCopyPairRequest()
		response = bdrcv20260330.NewCreateInstanceCopyPairResponse()
	)

	if v, ok := d.GetOk("protect_group_id"); ok {
		request.ProtectGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("create_target_instance_parameters"); ok {
		instanceModels := buildCreateInstanceModels(v.([]interface{}))
		request.CreateTargetInstanceParameters = instanceModels
	}

	if v, ok := d.GetOk("instance_copy_pair_name"); ok {
		request.InstanceCopyPairName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("recovery_point_objective"); ok {
		request.RecoveryPointObjective = helper.IntInt64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client().CreateInstanceCopyPairWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create bdrc_instance_copy_pair failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create bdrc_instance_copy_pair failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.CopyPairIds == nil || len(response.Response.CopyPairIds) == 0 {
		log.Printf("[CRITAL]%s bdrc_instance_copy_pair create returned empty CopyPairIds, d.Id()=%s", logId, d.Id())
		return fmt.Errorf("Create bdrc_instance_copy_pair failed, CopyPairIds is empty.")
	}

	copyPairId := *response.Response.CopyPairIds[0]
	d.SetId(copyPairId)

	pollErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		service := NewBdrcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		copyPair, e := service.DescribeBdrcInstanceCopyPairById(ctx, copyPairId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if copyPair == nil {
			return resource.RetryableError(fmt.Errorf("bdrc_instance_copy_pair %s not found yet, keep polling.", copyPairId))
		}

		if copyPair.CopyPairState != nil && *copyPair.CopyPairState != "INIT" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("bdrc_instance_copy_pair %s still in INIT state, keep polling.", copyPairId))
	})

	if pollErr != nil {
		log.Printf("[CRITAL]%s bdrc_instance_copy_pair %s async create polling failed, reason:%+v", logId, copyPairId, pollErr)
		return pollErr
	}

	return resourceTencentCloudBdrcInstanceCopyPairRead(d, meta)
}

func resourceTencentCloudBdrcInstanceCopyPairRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_instance_copy_pair.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewBdrcService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	respData, err := service.DescribeBdrcInstanceCopyPairById(ctx, d.Id())
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[CRUD] bdrc_instance_copy_pair id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if respData.CopyPairId != nil {
		_ = d.Set("copy_pair_id", respData.CopyPairId)
	}

	if respData.CopyPairName != nil {
		_ = d.Set("copy_pair_name", respData.CopyPairName)
	}

	if respData.CopyPairState != nil {
		_ = d.Set("copy_pair_state", respData.CopyPairState)
	}

	if respData.CopyPairType != nil {
		_ = d.Set("copy_pair_type", respData.CopyPairType)
	}

	if respData.SitePairId != nil {
		_ = d.Set("site_pair_id", respData.SitePairId)
	}

	if respData.SitePairName != nil {
		_ = d.Set("site_pair_name", respData.SitePairName)
	}

	if respData.ProtectGroupId != nil {
		_ = d.Set("protect_group_id", respData.ProtectGroupId)
	}

	if respData.ProtectGroupName != nil {
		_ = d.Set("protect_group_name", respData.ProtectGroupName)
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

	if respData.SourceResourceId != nil {
		_ = d.Set("source_resource_id", respData.SourceResourceId)
	}

	if respData.TargetResourceId != nil {
		_ = d.Set("target_resource_id", respData.TargetResourceId)
	}

	if respData.InstanceId != nil {
		_ = d.Set("instance_id", respData.InstanceId)
	}

	if respData.InstanceCopyPairId != nil {
		_ = d.Set("instance_copy_pair_id", respData.InstanceCopyPairId)
	}

	if respData.Percent != nil {
		_ = d.Set("percent", respData.Percent)
	}

	if respData.LatestProtectionTime != nil {
		_ = d.Set("latest_protection_time", respData.LatestProtectionTime)
	}

	if respData.RecoveryPointObjective != nil {
		_ = d.Set("recovery_point_objective", respData.RecoveryPointObjective)
	}

	if respData.DataDirection != nil {
		_ = d.Set("data_direction", respData.DataDirection)
	}

	if respData.CreateFrom != nil {
		_ = d.Set("create_from", respData.CreateFrom)
	}

	if respData.DisasterRecoveryType != nil {
		_ = d.Set("disaster_recovery_type", respData.DisasterRecoveryType)
	}

	if respData.PeerCloudName != nil {
		_ = d.Set("peer_cloud_name", respData.PeerCloudName)
	}

	if respData.Rollbacking != nil {
		_ = d.Set("rollbacking", respData.Rollbacking)
	}

	if respData.RollbackPercent != nil {
		_ = d.Set("rollback_percent", respData.RollbackPercent)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.AccountUin != nil {
		_ = d.Set("account_uin", respData.AccountUin)
	}

	if respData.SubAccountUin != nil {
		_ = d.Set("sub_account_uin", respData.SubAccountUin)
	}

	if respData.DrillGroupId != nil {
		_ = d.Set("drill_group_id", respData.DrillGroupId)
	}

	if respData.ProtectionTimeSet != nil {
		_ = d.Set("protection_time_set", helper.StringsInterfaces(respData.ProtectionTimeSet))
	} else {
		_ = d.Set("protection_time_set", []string{})
	}

	if respData.DiskCopyPairSet != nil {
		diskCopyPairSetList := make([]map[string]interface{}, 0, len(respData.DiskCopyPairSet))
		for _, diskCopyPair := range respData.DiskCopyPairSet {
			diskCopyPairMap := map[string]interface{}{}
			if diskCopyPair.CopyPairId != nil {
				diskCopyPairMap["copy_pair_id"] = diskCopyPair.CopyPairId
			}

			if diskCopyPair.CopyPairName != nil {
				diskCopyPairMap["copy_pair_name"] = diskCopyPair.CopyPairName
			}

			if diskCopyPair.SourceResourceId != nil {
				diskCopyPairMap["source_resource_id"] = diskCopyPair.SourceResourceId
			}

			if diskCopyPair.TargetResourceId != nil {
				diskCopyPairMap["target_resource_id"] = diskCopyPair.TargetResourceId
			}

			if diskCopyPair.CreateTime != nil {
				diskCopyPairMap["create_time"] = diskCopyPair.CreateTime
			}

			diskCopyPairSetList = append(diskCopyPairSetList, diskCopyPairMap)
		}

		_ = d.Set("disk_copy_pair_set", diskCopyPairSetList)
	}

	if respData.DeferredCreate != nil {
		_ = d.Set("deferred_create", respData.DeferredCreate)
	}

	if respData.TargetCvmCreated != nil {
		_ = d.Set("target_cvm_created", respData.TargetCvmCreated)
	}

	if respData.CvmCreateParams != nil {
		_ = d.Set("cvm_create_params", respData.CvmCreateParams)
	}

	return nil
}

func resourceTencentCloudBdrcInstanceCopyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_instance_copy_pair.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	immutableArgs := []string{"protect_group_id", "create_target_instance_parameters", "recovery_point_objective", "delete_target_resource"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("bdrc_instance_copy_pair argument `%s` is immutable, it can only be changed via recreation.", v)
		}
	}

	if d.HasChange("instance_copy_pair_name") {
		request := bdrcv20260330.NewModifyCopyPairAttributeRequest()
		request.CopyPairId = helper.String(d.Id())
		request.CopyPairType = helper.String("INSTANCE")

		if v, ok := d.GetOk("instance_copy_pair_name"); ok {
			request.CopyPairName = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client().ModifyCopyPairAttributeWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify bdrc_instance_copy_pair attribute failed, Response is nil."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update bdrc_instance_copy_pair failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudBdrcInstanceCopyPairRead(d, meta)
}

func resourceTencentCloudBdrcInstanceCopyPairDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_bdrc_instance_copy_pair.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = bdrcv20260330.NewDeleteCopyPairsRequest()
	)

	request.CopyPairIds = []*string{helper.String(d.Id())}
	request.CopyPairType = helper.String("INSTANCE")

	if v, ok := d.GetOkExists("delete_target_resource"); ok {
		request.DeleteTargetResource = helper.Bool(v.(bool))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBdrcV20260330Client().DeleteCopyPairsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete bdrc_instance_copy_pair failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete bdrc_instance_copy_pair failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func buildCreateInstanceModels(list []interface{}) []*bdrcv20260330.CreateInstanceModel {
	if len(list) == 0 {
		return nil
	}

	models := make([]*bdrcv20260330.CreateInstanceModel, 0, len(list))
	for _, item := range list {
		modelMap := item.(map[string]interface{})
		model := bdrcv20260330.CreateInstanceModel{}

		if v, ok := modelMap["source_instance_id"].(string); ok && v != "" {
			model.SourceInstanceId = helper.String(v)
		}

		if v, ok := modelMap["instance_charge_type"].(string); ok && v != "" {
			model.InstanceChargeType = helper.String(v)
		}

		if v, ok := modelMap["placement"].([]interface{}); ok && len(v) > 0 {
			model.Placement = buildPlacement(v[0].(map[string]interface{}))
		}

		if v, ok := modelMap["image_id"].(string); ok && v != "" {
			model.ImageId = helper.String(v)
		}

		if v, ok := modelMap["system_disk"].([]interface{}); ok && len(v) > 0 {
			model.SystemDisk = buildDiskModel(v[0].(map[string]interface{}))
		}

		if v, ok := modelMap["instance_charge_prepaid"].([]interface{}); ok && len(v) > 0 {
			model.InstanceChargePrepaid = buildInstanceChargePrepaid(v[0].(map[string]interface{}))
		}

		if v, ok := modelMap["instance_type"].(string); ok && v != "" {
			model.InstanceType = helper.String(v)
		}

		if v, ok := modelMap["data_disks"].([]interface{}); ok && len(v) > 0 {
			model.DataDisks = buildDiskModels(v)
		}

		if v, ok := modelMap["virtual_private_cloud"].([]interface{}); ok && len(v) > 0 {
			model.VirtualPrivateCloud = buildVirtualPrivateCloud(v[0].(map[string]interface{}))
		}

		if v, ok := modelMap["internet_accessible"].([]interface{}); ok && len(v) > 0 {
			model.InternetAccessible = buildInternetAccessible(v[0].(map[string]interface{}))
		}

		if v, ok := modelMap["instance_name"].(string); ok && v != "" {
			model.InstanceName = helper.String(v)
		}

		if v, ok := modelMap["login_settings"].([]interface{}); ok && len(v) > 0 {
			model.LoginSettings = buildLoginSettings(v[0].(map[string]interface{}))
		}

		if v, ok := modelMap["enhanced_service"].([]interface{}); ok && len(v) > 0 {
			model.EnhancedService = buildEnhancedService(v[0].(map[string]interface{}))
		}

		if v, ok := modelMap["spot_price"].(string); ok && v != "" {
			model.SpotPrice = helper.String(v)
		}

		if v, ok := modelMap["host_name"].(string); ok && v != "" {
			model.HostName = helper.String(v)
		}

		if v, ok := modelMap["user_data"].(string); ok && v != "" {
			model.UserData = helper.String(v)
		}

		if v, ok := modelMap["disaster_recover_group_ids"].([]interface{}); ok && len(v) > 0 {
			model.DisasterRecoverGroupIds = helper.StringsStringsPoint(stringsInterfacesToStrings(v))
		}

		if v, ok := modelMap["stopped_mode"].(string); ok && v != "" {
			model.StoppedMode = helper.String(v)
		}

		if v, ok := modelMap["copy_pair_id"].(string); ok && v != "" {
			model.CopyPairId = helper.String(v)
		}

		if v, ok := modelMap["recovery_time"].(string); ok && v != "" {
			model.RecoveryTime = helper.String(v)
		}

		models = append(models, &model)
	}

	return models
}

func buildPlacement(m map[string]interface{}) *bdrcv20260330.Placement {
	placement := bdrcv20260330.Placement{}

	if v, ok := m["zone"].(string); ok && v != "" {
		placement.Zone = helper.String(v)
	}

	if v, ok := m["project_id"].(int); ok && v != 0 {
		placement.ProjectId = helper.IntInt64(v)
	}

	if v, ok := m["host_id"].(string); ok && v != "" {
		placement.HostId = helper.String(v)
	}

	if v, ok := m["host_ids"].([]interface{}); ok && len(v) > 0 {
		placement.HostIds = helper.StringsStringsPoint(stringsInterfacesToStrings(v))
	}

	if v, ok := m["project_name"].(string); ok && v != "" {
		placement.ProjectName = helper.String(v)
	}

	return &placement
}

func buildDiskModel(m map[string]interface{}) *bdrcv20260330.DiskModel {
	disk := bdrcv20260330.DiskModel{}

	if v, ok := m["disk_type"].(string); ok && v != "" {
		disk.DiskType = helper.String(v)
	}

	if v, ok := m["disk_size"].(int); ok && v != 0 {
		disk.DiskSize = helper.IntInt64(v)
	}

	if v, ok := m["delete_with_instance"].(bool); ok {
		disk.DeleteWithInstance = helper.Bool(v)
	}

	return &disk
}

func buildDiskModels(list []interface{}) []*bdrcv20260330.DiskModel {
	if len(list) == 0 {
		return nil
	}

	disks := make([]*bdrcv20260330.DiskModel, 0, len(list))
	for _, item := range list {
		disks = append(disks, buildDiskModel(item.(map[string]interface{})))
	}

	return disks
}

func buildInstanceChargePrepaid(m map[string]interface{}) *bdrcv20260330.InstanceChargePrepaid {
	prepaid := bdrcv20260330.InstanceChargePrepaid{}

	if v, ok := m["period"].(int); ok && v != 0 {
		prepaid.Period = helper.IntInt64(v)
	}

	if v, ok := m["renew_flag"].(string); ok && v != "" {
		prepaid.RenewFlag = helper.String(v)
	}

	return &prepaid
}

func buildVirtualPrivateCloud(m map[string]interface{}) *bdrcv20260330.VirtualPrivateCloud {
	vpc := bdrcv20260330.VirtualPrivateCloud{}

	if v, ok := m["vpc_id"].(string); ok && v != "" {
		vpc.VpcId = helper.String(v)
	}

	if v, ok := m["subnet_id"].(string); ok && v != "" {
		vpc.SubnetId = helper.String(v)
	}

	if v, ok := m["subnet_name"].(string); ok && v != "" {
		vpc.SubnetName = helper.String(v)
	}

	if v, ok := m["as_vpc_gateway"].(bool); ok {
		vpc.AsVpcGateway = helper.Bool(v)
	}

	if v, ok := m["private_ip_addresses"].([]interface{}); ok && len(v) > 0 {
		vpc.PrivateIpAddresses = helper.StringsStringsPoint(stringsInterfacesToStrings(v))
	}

	if v, ok := m["vpc_name"].(string); ok && v != "" {
		vpc.VpcName = helper.String(v)
	}

	if v, ok := m["ipv6_address_count"].(int); ok && v != 0 {
		vpc.Ipv6AddressCount = helper.IntInt64(v)
	}

	return &vpc
}

func buildInternetAccessible(m map[string]interface{}) *bdrcv20260330.InternetAccessible {
	accessible := bdrcv20260330.InternetAccessible{}

	if v, ok := m["internet_charge_type"].(string); ok && v != "" {
		accessible.InternetChargeType = helper.String(v)
	}

	if v, ok := m["internet_max_bandwidth_out"].(int); ok && v != 0 {
		accessible.InternetMaxBandwidthOut = helper.IntInt64(v)
	}

	if v, ok := m["public_ip_assigned"].(bool); ok {
		accessible.PublicIpAssigned = helper.Bool(v)
	}

	if v, ok := m["internet_service_provider"].(string); ok && v != "" {
		accessible.InternetServiceProvider = helper.String(v)
	}

	return &accessible
}

func buildLoginSettings(m map[string]interface{}) *bdrcv20260330.LoginSettings {
	settings := bdrcv20260330.LoginSettings{}

	if v, ok := m["password"].(string); ok && v != "" {
		settings.Password = helper.String(v)
	}

	if v, ok := m["key_ids"].([]interface{}); ok && len(v) > 0 {
		settings.KeyIds = helper.StringsStringsPoint(stringsInterfacesToStrings(v))
	}

	if v, ok := m["keep_image_login"].(string); ok && v != "" {
		settings.KeepImageLogin = helper.String(v)
	}

	return &settings
}

func buildEnhancedService(m map[string]interface{}) *bdrcv20260330.EnhancedService {
	service := bdrcv20260330.EnhancedService{}

	if v, ok := m["security_service"].([]interface{}); ok && len(v) > 0 {
		securityMap := v[0].(map[string]interface{})
		enabled := bdrcv20260330.RunSecurityServiceEnabled{}
		if val, ok := securityMap["enabled"].(bool); ok {
			enabled.Enabled = helper.Bool(val)
		}
		service.SecurityService = &enabled
	}

	if v, ok := m["monitor_service"].([]interface{}); ok && len(v) > 0 {
		monitorMap := v[0].(map[string]interface{})
		enabled := bdrcv20260330.RunSecurityServiceEnabled{}
		if val, ok := monitorMap["enabled"].(bool); ok {
			enabled.Enabled = helper.Bool(val)
		}
		service.MonitorService = &enabled
	}

	if v, ok := m["automation_service"].([]interface{}); ok && len(v) > 0 {
		automationMap := v[0].(map[string]interface{})
		enabled := bdrcv20260330.AutomationServiceEnabled{}
		if val, ok := automationMap["enabled"].(bool); ok {
			enabled.Enabled = helper.Bool(val)
		}
		service.AutomationService = &enabled
	}

	if v, ok := m["basic_service"].([]interface{}); ok && len(v) > 0 {
		basicMap := v[0].(map[string]interface{})
		enabled := bdrcv20260330.BasicServicesSettings{}
		if val, ok := basicMap["enabled"].(bool); ok {
			enabled.Enabled = helper.Bool(val)
		}
		service.BasicService = &enabled
	}

	return &service
}

func stringsInterfacesToStrings(list []interface{}) []string {
	result := make([]string, 0, len(list))
	for _, v := range list {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}

	return result
}
