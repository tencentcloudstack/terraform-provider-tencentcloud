package wedata

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWedataResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWedataResourceGroupCreate,
		Read:   resourceTencentCloudWedataResourceGroupRead,
		Update: resourceTencentCloudWedataResourceGroupUpdate,
		Delete: resourceTencentCloudWedataResourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Resource group name. The name for creating a general resource group must start with a letter, can contain letters, numbers, underscores (_), and up to 64 characters.",
			},

			"type": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Information about the activated resource group.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_group_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Resource group type.\n\n- Schedule --- Scheduling resource group\n- Integration --- Integration resource group  \n- DataService -- Data service resource group.",
						},
						"integration": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Integration resource group, subdivided into real-time resource group and offline resource group (Integration, scheduling, and data service resource groups cannot be purchased simultaneously).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"real_time_data_sync": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Real-time integration resource group.\n\n- i32c (Real-time data synchronization - 16C64G).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"specification": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Resource group specification.",
												},
												"number": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Quantity.",
												},
											},
										},
									},
									"offline_data_sync": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Offline integration resource group.\n\n- integrated (Offline data synchronization - 8C16G)\n- i16 (Offline data synchronization - 8C32G).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"specification": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Resource group specification.",
												},
												"number": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Quantity.",
												},
											},
										},
									},
								},
							},
						},
						"schedule": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Scheduling resource group (Integration, scheduling, and data service resource groups cannot be purchased simultaneously).\n\n- s_test (Test specification)\n- s_small (Basic specification)\n- s_medium (Popular specification)\n- s_large (Professional specification).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"specification": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Resource group specification.",
									},
									"number": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Quantity.",
									},
								},
							},
						},
						"data_service": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Data service resource group (Integration, scheduling, and data service resource groups cannot be purchased simultaneously).\n\n- ds_t (Test specification)\n- ds_s (Basic specification)\n- ds_m (Popular specification)\n- ds_l (Professional specification).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"specification": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Resource group specification.",
									},
									"number": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Quantity.",
									},
								},
							},
						},
					},
				},
			},

			"auto_renew_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether auto-renewal is enabled.",
			},

			"purchase_period": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Purchase duration, in months.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC ID.",
			},

			"subnet": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet.",
			},

			"resource_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Resource purchase region.",
			},

			"associated_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Associated project space project ID.",
			},

			// computed
			"resource_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group ID.",
			},
		},
	}
}

func resourceTencentCloudWedataResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request         = wedatav20250806.NewCreateResourceGroupRequest()
		response        = wedatav20250806.NewCreateResourceGroupResponse()
		resourceGroupId string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if typeMap, ok := helper.InterfacesHeadMap(d, "type"); ok {
		resourceType := wedatav20250806.ResourceType{}
		if v, ok := typeMap["resource_group_type"].(string); ok && v != "" {
			resourceType.ResourceGroupType = helper.String(v)
		}

		if integrationMap, ok := helper.ConvertInterfacesHeadToMap(typeMap["integration"]); ok {
			integrationResource := wedatav20250806.IntegrationResource{}
			if realTimeDataSyncMap, ok := helper.ConvertInterfacesHeadToMap(integrationMap["real_time_data_sync"]); ok {
				resourceGroupSpecification := wedatav20250806.ResourceGroupSpecification{}
				if v, ok := realTimeDataSyncMap["specification"].(string); ok && v != "" {
					resourceGroupSpecification.Specification = helper.String(v)
				}

				if v, ok := realTimeDataSyncMap["number"].(int); ok {
					resourceGroupSpecification.Number = helper.IntInt64(v)
				}

				integrationResource.RealTimeDataSync = &resourceGroupSpecification
			}

			if offlineDataSyncMap, ok := helper.ConvertInterfacesHeadToMap(integrationMap["offline_data_sync"]); ok {
				resourceGroupSpecification2 := wedatav20250806.ResourceGroupSpecification{}
				if v, ok := offlineDataSyncMap["specification"].(string); ok && v != "" {
					resourceGroupSpecification2.Specification = helper.String(v)
				}

				if v, ok := offlineDataSyncMap["number"].(int); ok {
					resourceGroupSpecification2.Number = helper.IntInt64(v)
				}

				integrationResource.OfflineDataSync = &resourceGroupSpecification2
			}

			resourceType.Integration = &integrationResource
		}

		if scheduleMap, ok := helper.ConvertInterfacesHeadToMap(typeMap["schedule"]); ok {
			resourceGroupSpecification3 := wedatav20250806.ResourceGroupSpecification{}
			if v, ok := scheduleMap["specification"].(string); ok && v != "" {
				resourceGroupSpecification3.Specification = helper.String(v)
			}

			if v, ok := scheduleMap["number"].(int); ok {
				resourceGroupSpecification3.Number = helper.IntInt64(v)
			}

			resourceType.Schedule = &resourceGroupSpecification3
		}

		if dataServiceMap, ok := helper.ConvertInterfacesHeadToMap(typeMap["data_service"]); ok {
			resourceGroupSpecification4 := wedatav20250806.ResourceGroupSpecification{}
			if v, ok := dataServiceMap["specification"].(string); ok && v != "" {
				resourceGroupSpecification4.Specification = helper.String(v)
			}

			if v, ok := dataServiceMap["number"].(int); ok {
				resourceGroupSpecification4.Number = helper.IntInt64(v)
			}

			resourceType.DataService = &resourceGroupSpecification4
		}

		request.Type = &resourceType
	}

	if v, ok := d.GetOkExists("auto_renew_enabled"); ok {
		request.AutoRenewEnabled = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("purchase_period"); ok {
		request.PurchasePeriod = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet"); ok {
		request.SubNet = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_region"); ok {
		request.ResourceRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("associated_project_id"); ok {
		request.AssociatedProjectId = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().CreateResourceGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Create wedata resource group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create wedata resource group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if !*response.Response.Data.Status {
		return fmt.Errorf("Create wedata resource group failed, Status is false.")
	}

	if response.Response.Data.ResourceGroupId == nil || *response.Response.Data.ResourceGroupId == "" {
		return fmt.Errorf("ResourceGroupId is nil.")
	}

	resourceGroupId = *response.Response.Data.ResourceGroupId
	d.SetId(resourceGroupId)
	return resourceTencentCloudWedataResourceGroupRead(d, meta)
}

func resourceTencentCloudWedataResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service         = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		resourceGroupId = d.Id()
	)

	respData, err := service.DescribeWedataResourceGroupById(ctx, resourceGroupId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_wedata_resource_group` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	for _, items := range respData {
		if items.Name != nil {
			_ = d.Set("name", items.Name)
		}

		if items.AutoRenewEnabled != nil {
			_ = d.Set("auto_renew_enabled", items.AutoRenewEnabled)
		}

		if items.VpcId != nil {
			_ = d.Set("vpc_id", items.VpcId)
		}

		if items.SubNet != nil {
			_ = d.Set("subnet", items.SubNet)
		}

		if items.Region != nil {
			_ = d.Set("resource_region", items.Region)
		}

		if items.Id != nil {
			_ = d.Set("resource_group_id", items.Id)
		}
	}

	return nil
}

func resourceTencentCloudWedataResourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		resourceGroupId = d.Id()
	)

	if d.HasChange("auto_renew_enabled") || d.HasChange("purchase_period") {
		request := wedatav20250806.NewUpdateResourceGroupRequest()
		if v, ok := d.GetOkExists("auto_renew_enabled"); ok {
			request.AutoRenewEnabled = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("purchase_period"); ok {
			request.PurchasePeriod = helper.IntInt64(v.(int))
		}

		request.Id = &resourceGroupId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().UpdateResourceGroupWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.Data == nil || result.Response.Data.Status == nil {
				return resource.NonRetryableError(fmt.Errorf("Update wedata resource group failed, Response is nil."))
			}

			if !*result.Response.Data.Status {
				return resource.NonRetryableError(fmt.Errorf("Update wedata resource group failed, Status is false."))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update wedata resource group failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudWedataResourceGroupRead(d, meta)
}

func resourceTencentCloudWedataResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_wedata_resource_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request         = wedatav20250806.NewDeleteResourceGroupRequest()
		response        = wedatav20250806.NewDeleteResourceGroupResponse()
		resourceGroupId = d.Id()
	)

	request.Id = &resourceGroupId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWedataV20250806Client().DeleteResourceGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete wedata resource group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete wedata resource group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Data.Status != nil && *response.Response.Data.Status == true {
		return nil
	}

	return fmt.Errorf("Delete wedata resource group failed, Status is false.")
}
