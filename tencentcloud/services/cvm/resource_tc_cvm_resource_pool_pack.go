package cvm

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmResourcePoolPack() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmResourcePoolPackCreate,
		Read:   resourceTencentCloudCvmResourcePoolPackRead,
		Delete: resourceTencentCloudCvmResourcePoolPackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The availability zone where the resource pool pack is located. Format: ap-guangzhou-6.",
			},
			"instance_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance type for the resource pool pack. Only half-machine/full-machine specifications are supported. Format: SA9.96XLARGE1152 (SA9 half-machine).",
			},
			"period": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The period of the resource pool pack in months. Range: 1-60.",
			},
			"resource_pool_pack_type": {
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Resource pool pack type. Options: EXCLUSIVE (exclusive, default), SHARED (shared). Note: Only EXCLUSIVE is supported in the first phase.",
			},
			"auto_placement": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Auto placement switch. Default: true. When enabled, the system will search for suitable pools in pools with this capability enabled when creating instances without specifying a resource pool.",
			},
			"dedicated_resource_pool_pack_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The name of the resource pool pack. Length: 1-60 characters, supports Chinese, English, numbers, hyphens '-', and underscores '_'.",
			},
			"renew_flag": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Auto renewal flag. Options: NOTIFY_AND_AUTO_RENEW (notify and auto renew), NOTIFY_AND_MANUAL_RENEW (notify and manual renew, default), DISABLE_NOTIFY_AND_MANUAL_RENEW (do not notify and manual renew).",
			},

			// Computed attributes
			"dedicated_resource_pack_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The ID of the created resource pool pack.",
			},
			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Resource pool pack status. Values: CREATING (creating), ACTIVE (running), FAILED (creation failed), RETIRED (expired).",
			},
			"instance_family": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Instance family. Format: SA9.",
			},
			"start_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Resource pool pack creation time. Format: YYYY-MM-DDThh:mm:ssZ.",
			},
			"end_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Resource pool pack expiration time. Format: YYYY-MM-DDThh:mm:ssZ.",
			},
		},
	}
}

func resourceTencentCloudCvmResourcePoolPackCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_resource_pool_pack.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request = cvm.NewPurchaseResourcePoolPacksRequest()
		service = CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	// Required fields
	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	request.InstanceCount = helper.IntUint64(1)

	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("period"); ok {
		request.Period = helper.IntUint64(v.(int))
	}

	// Optional fields
	if v, ok := d.GetOk("resource_pool_pack_type"); ok {
		request.ResourcePoolPackType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_placement"); ok {
		request.AutoPlacement = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("dedicated_resource_pool_pack_name"); ok {
		request.DedicatedResourcePoolPackName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("renew_flag"); ok {
		request.RenewFlag = helper.String(v.(string))
	}

	// Create resource pool packs
	packIds, err := service.CreateCvmResourcePoolPacks(ctx, request)
	if err != nil {
		log.Printf("[CRITAL]%s create cvm resource pool packs failed, reason:%+v", logId, err)
		return err
	}

	if len(packIds) < 1 {
		return fmt.Errorf("resource `tencentcloud_cvm_resource_pool_pack` create failed, no pack ID returned")
	}

	// Set resource ID to the first pack ID
	d.SetId(*packIds[0])

	// Wait for resource to be ready with retry logic
	err = resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		pack, e := service.DescribeCvmResourcePoolPackById(ctx, d.Id())
		if e != nil {
			return resource.NonRetryableError(e)
		}
		if pack == nil {
			return resource.RetryableError(fmt.Errorf("resource pool pack %s not found yet", d.Id()))
		}
		if pack.Status != nil && *pack.Status == "CREATING" {
			return resource.RetryableError(fmt.Errorf("resource pool pack %s is still creating", d.Id()))
		}
		if pack.Status != nil && *pack.Status == "FAILED" {
			return resource.NonRetryableError(fmt.Errorf("resource pool pack %s creation failed", d.Id()))
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s wait for resource pool pack ready failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCvmResourcePoolPackRead(d, meta)
}

func resourceTencentCloudCvmResourcePoolPackRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_resource_pool_pack.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	packId := d.Id()

	// Retry logic for eventual consistency
	var pack *cvm.ResourcePoolPack
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCvmResourcePoolPackById(ctx, packId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		pack = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read cvm resource pool pack failed, reason:%+v", logId, err)
		return err
	}

	if pack == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tencentcloud_cvm_resource_pool_pack` %s does not exist", logId, packId)
		return nil
	}

	// Set computed fields
	if pack.DedicatedResourcePackId != nil {
		_ = d.Set("dedicated_resource_pack_id", pack.DedicatedResourcePackId)
	}

	if pack.Zone != nil {
		_ = d.Set("zone", pack.Zone)
	}

	if pack.InstanceType != nil {
		_ = d.Set("instance_type", pack.InstanceType)
	}

	if pack.InstanceFamily != nil {
		_ = d.Set("instance_family", pack.InstanceFamily)
	}

	if pack.ResourcePoolPackType != nil {
		_ = d.Set("resource_pool_pack_type", pack.ResourcePoolPackType)
	}

	if pack.Status != nil {
		_ = d.Set("status", pack.Status)
	}

	if pack.AutoPlacement != nil {
		_ = d.Set("auto_placement", pack.AutoPlacement)
	}

	if pack.DedicatedResourcePackName != nil {
		_ = d.Set("dedicated_resource_pool_pack_name", pack.DedicatedResourcePackName)
	}

	if pack.RenewFlag != nil {
		_ = d.Set("renew_flag", pack.RenewFlag)
	}

	if pack.StartTime != nil {
		_ = d.Set("start_time", pack.StartTime)
	}

	if pack.EndTime != nil {
		_ = d.Set("end_time", pack.EndTime)
	}

	return nil
}

func resourceTencentCloudCvmResourcePoolPackDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_resource_pool_pack.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	packId := d.Id()

	if err := service.DeleteCvmResourcePoolPack(ctx, []*string{&packId}); err != nil {
		log.Printf("[CRITAL]%s delete cvm resource pool pack failed, reason:%+v", logId, err)
		return err
	}

	// Wait for deletion to complete
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		pack, e := service.DescribeCvmResourcePoolPackById(ctx, packId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if pack != nil {
			return resource.RetryableError(fmt.Errorf("resource pool pack %s still exists", packId))
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete cvm resource pool pack failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
