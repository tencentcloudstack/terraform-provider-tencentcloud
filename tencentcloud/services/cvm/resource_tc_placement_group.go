package cvm

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func ResourceTencentCloudPlacementGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPlacementGroupCreate,
		Read:   resourceTencentCloudPlacementGroupRead,
		Update: resourceTencentCloudPlacementGroupUpdate,
		Delete: resourceTencentCloudPlacementGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of the placement group, 1-60 characters in length.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CVM_PLACEMENT_GROUP_TYPE),
				Description:  "Type of the placement group. Valid values: `HOST`, `SW` and `RACK`.",
			},

			// computed
			"cvm_quota_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum number of hosts in the placement group.",
			},
			"current_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of hosts in the placement group.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the placement group.",
			},
		},
	}
}

func resourceTencentCloudPlacementGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_placement_group.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	placementName := d.Get("name").(string)
	placementType := d.Get("type").(string)
	var id string
	var errRet error
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		id, errRet = cvmService.CreatePlacementGroup(ctx, placementName, placementType)
		if errRet != nil {
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(id)

	return resourceTencentCloudPlacementGroupRead(d, meta)
}

func resourceTencentCloudPlacementGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_placement_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	placementId := d.Id()
	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var placement *cvm.DisasterRecoverGroup
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		placement, errRet = cvmService.DescribePlacementGroupById(ctx, placementId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if placement == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", placement.Name)
	_ = d.Set("type", placement.Type)
	_ = d.Set("cvm_quota_total", placement.CvmQuotaTotal)
	_ = d.Set("current_num", placement.CurrentNum)
	_ = d.Set("create_time", placement.CreateTime)

	return nil
}

func resourceTencentCloudPlacementGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_placement_group.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	placementId := d.Id()
	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	if d.HasChange("name") {
		placementName := d.Get("name").(string)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			err := cvmService.ModifyPlacementGroup(ctx, placementId, placementName)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudPlacementGroupRead(d, meta)
}

func resourceTencentCloudPlacementGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_placement_group.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	placementId := d.Id()
	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		err := cvmService.DeletePlacementGroup(ctx, placementId)
		if err != nil {
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
