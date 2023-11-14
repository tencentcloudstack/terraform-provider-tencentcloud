/*
Provides a resource to create a cvm sync_images

Example Usage

```hcl
resource "tencentcloud_cvm_sync_images" "sync_images" {
  image_ids =
  destination_regions =
  dry_run = false
  image_name = "img-evhmf3fy"
  image_set_required = false
}
```

Import

cvm sync_images can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_sync_images.sync_images sync_images_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCvmSyncImages() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmSyncImagesCreate,
		Read:   resourceTencentCloudCvmSyncImagesRead,
		Delete: resourceTencentCloudCvmSyncImagesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"image_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of image IDs.The specified images must meet the following requirement: the images must be in the `NORMAL` state.",
			},

			"destination_regions": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of destination regions for synchronization. Limits: It must be a valid region. For a custom image, the destination region cannot be the source region. For a shared image, the destination region must be the source region, which indicates to create a copy of the image as a custom image in the same region.",
			},

			"dry_run": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Checks whether image synchronization can be initiated.",
			},

			"image_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Destination image name.",
			},

			"image_set_required": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to return the ID of image created in the destination region.",
			},
		},
	}
}

func resourceTencentCloudCvmSyncImagesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_sync_images.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cvm.NewSyncImagesRequest()
		response = cvm.NewSyncImagesResponse()
		imageId  string
	)
	if v, ok := d.GetOk("image_ids"); ok {
		imageIdsSet := v.(*schema.Set).List()
		for i := range imageIdsSet {
			imageIds := imageIdsSet[i].(string)
			request.ImageIds = append(request.ImageIds, &imageIds)
		}
	}

	if v, ok := d.GetOk("destination_regions"); ok {
		destinationRegionsSet := v.(*schema.Set).List()
		for i := range destinationRegionsSet {
			destinationRegions := destinationRegionsSet[i].(string)
			request.DestinationRegions = append(request.DestinationRegions, &destinationRegions)
		}
	}

	if v, _ := d.GetOk("dry_run"); v != nil {
		request.DryRun = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("image_name"); ok {
		request.ImageName = helper.String(v.(string))
	}

	if v, _ := d.GetOk("image_set_required"); v != nil {
		request.ImageSetRequired = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().SyncImages(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm syncImages failed, reason:%+v", logId, err)
		return err
	}

	imageId = *response.Response.ImageId
	d.SetId(imageId)

	return resourceTencentCloudCvmSyncImagesRead(d, meta)
}

func resourceTencentCloudCvmSyncImagesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_sync_images.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmSyncImagesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_sync_images.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
