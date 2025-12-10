package cvm

import (
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmSyncImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmSyncImageCreate,
		Read:   resourceTencentCloudCvmSyncImageRead,
		Delete: resourceTencentCloudCvmSyncImageDelete,

		Schema: map[string]*schema.Schema{
			"image_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Image ID. The specified image must meet the following requirement: the images must be in the `NORMAL` state.",
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

			"encrypt": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to synchronize as an encrypted custom image. Default value is `false`. Synchronization to an encrypted custom image is only supported within the same region.",
			},

			"kms_key_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "KMS key ID used when synchronizing to an encrypted custom image. This parameter is valid only synchronizing to an encrypted image. If KmsKeyId is not specified, the default CBS cloud product KMS key is used.",
			},

			"image_set": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image ID.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region of the image.",
						},
					},
				},
				Description: "ID of the image created in the destination region.",
			},
		},
	}
}

func resourceTencentCloudCvmSyncImageCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_sync_image.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := cvm.NewSyncImagesRequest()
	response := cvm.NewSyncImagesResponse()
	imageId := d.Get("image_id").(string)
	request.ImageIds = []*string{&imageId}

	if v, ok := d.GetOk("destination_regions"); ok {
		destinationRegionsSet := v.(*schema.Set).List()
		for i := range destinationRegionsSet {
			destinationRegions := destinationRegionsSet[i].(string)
			request.DestinationRegions = append(request.DestinationRegions, &destinationRegions)
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

	if v, ok := d.GetOkExists("encrypt"); ok {
		request.Encrypt = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		request.KmsKeyId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().SyncImages(request)
		if e != nil {
			return tccommon.RetryError(e)
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

	if response == nil || response.Response == nil || response.Response.ImageSet == nil {
		err = fmt.Errorf("Response is nil")
		return err
	}

	d.SetId(imageId)

	imageSetList := []interface{}{}
	for _, image := range response.Response.ImageSet {
		imageMap := map[string]interface{}{}

		if image.ImageId != nil {
			imageMap["image_id"] = image.ImageId
		}

		if image.Region != nil {
			imageMap["region"] = image.Region
		}

		imageSetList = append(imageSetList, imageMap)
	}

	_ = d.Set("image_set", imageSetList)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"NORMAL"}, 20*tccommon.ReadRetryTimeout, time.Second, service.CvmSyncImagesStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

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
