package tencentcloud

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCvmExportImages() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmExportImagesCreate,
		Read:   resourceTencentCloudCvmExportImagesRead,
		Delete: resourceTencentCloudCvmExportImagesDelete,
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "COS bucket name.",
			},

			"image_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Image ID.",
			},

			"file_name_prefix": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Prefix of exported file.",
			},

			"export_format": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Format of the exported image file. Valid values: RAW, QCOW2, VHD and VMDK. Default value: RAW.",
			},

			"only_export_root_disk": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to export only the system disk.",
			},

			"dry_run": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Check whether the image can be exported.",
			},

			"role_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Role name (Default: CVM_QcsRole). Before exporting the images, make sure the role exists, and it has write permission to COS.",
			},
		},
	}
}

func resourceTencentCloudCvmExportImagesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_export_images.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = cvm.NewExportImagesRequest()
		imageId        string
		bucketName     string
		fileNamePrefix string
	)
	imageId = d.Get("image_id").(string)
	bucketName = d.Get("bucket_name").(string)
	fileNamePrefix = d.Get("file_name_prefix").(string)
	request.ImageIds = []*string{&imageId}
	request.BucketName = helper.String(bucketName)
	request.FileNamePrefixList = []*string{&fileNamePrefix}

	if v, ok := d.GetOk("export_format"); ok {
		request.ExportFormat = helper.String(v.(string))
	}

	if v, _ := d.GetOk("only_export_root_disk"); v != nil {
		request.OnlyExportRootDisk = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("dry_run"); v != nil {
		request.DryRun = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("role_name"); ok {
		request.RoleName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ExportImages(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm exportImages failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(imageId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"NORMAL"}, 20*readRetryTimeout, time.Second, service.CvmSyncImagesStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCvmExportImagesRead(d, meta)
}

func resourceTencentCloudCvmExportImagesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_export_images.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmExportImagesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_export_images.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
