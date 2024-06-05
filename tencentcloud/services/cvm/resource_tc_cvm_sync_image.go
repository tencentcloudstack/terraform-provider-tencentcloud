package cvm

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmSyncImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmSyncImageCreate,
		Read:   resourceTencentCloudCvmSyncImageRead,
		Delete: resourceTencentCloudCvmSyncImageDelete,
		Schema: map[string]*schema.Schema{
			"destination_regions": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Description: "List of destination regions for synchronization. Limits: It must be a valid region. For a custom image, the destination region cannot be the source region. For a shared image, the destination region must be the source region, which indicates to create a copy of the image as a custom image in the same region.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"dry_run": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Checks whether image synchronization can be initiated.",
			},

			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Image ID. The specified image must meet the following requirement: the images must be in the `NORMAL` state.",
			},

			"image_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Destination image name.",
			},

			"image_set_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to return the ID of image created in the destination region.",
			},
		},
	}
}

func resourceTencentCloudCvmSyncImageCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_sync_image.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		imageId string
	)
	var (
		request  = cvm.NewSyncImagesRequest()
		response = cvm.NewSyncImagesResponse()
	)

	if v, ok := d.GetOk("image_id"); ok {
		imageId = v.(string)
	}

	request.ImageIds = []*string{helper.String(imageId)}

	if v, ok := d.GetOk("destination_regions"); ok {
		destinationRegionsSet := v.(*schema.Set).List()
		for i := range destinationRegionsSet {
			destinationRegions := destinationRegionsSet[i].(string)
			request.DestinationRegions = append(request.DestinationRegions, helper.String(destinationRegions))
		}
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request.DryRun = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("image_name"); ok {
		request.ImageName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("image_set_required"); ok {
		request.ImageSetRequired = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().SyncImagesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cvm sync image failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	if err := resourceTencentCloudCvmSyncImageCreatePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	d.SetId(imageId)

	return resourceTencentCloudCvmSyncImageRead(d, meta)
}

func resourceTencentCloudCvmSyncImageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_sync_image.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmSyncImageDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_sync_image.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
