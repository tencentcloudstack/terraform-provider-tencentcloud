/*
Provides a resource to create a cvm export_images

Example Usage

```hcl
resource "tencentcloud_cvm_export_images" "export_images" {
  bucket_name = "test-bucket-AppId"
  image_ids =
  export_format = "RAW"
  file_name_prefix_list =
  only_export_root_disk = true
  dry_run = false
  role_name = "CVM_QcsRole"
}
```

Import

cvm export_images can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_export_images.export_images export_images_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudCvmExportImages() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmExportImagesCreate,
		Read:   resourceTencentCloudCvmExportImagesRead,
		Delete: resourceTencentCloudCvmExportImagesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "COS bucket name.",
			},

			"image_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of image IDs.",
			},

			"export_format": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Format of the exported image file. Valid values: RAW, QCOW2, VHD and VMDK. Default value: RAW.",
			},

			"file_name_prefix_list": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Prefix list of the name of exported files.",
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
		request    = cvm.NewExportImagesRequest()
		response   = cvm.NewExportImagesResponse()
		imageId    string
		bucketName string
	)
	if v, ok := d.GetOk("bucket_name"); ok {
		bucketName = v.(string)
		request.BucketName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_ids"); ok {
		imageIdsSet := v.(*schema.Set).List()
		for i := range imageIdsSet {
			imageIds := imageIdsSet[i].(string)
			request.ImageIds = append(request.ImageIds, &imageIds)
		}
	}

	if v, ok := d.GetOk("export_format"); ok {
		request.ExportFormat = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_name_prefix_list"); ok {
		fileNamePrefixListSet := v.(*schema.Set).List()
		for i := range fileNamePrefixListSet {
			fileNamePrefixList := fileNamePrefixListSet[i].(string)
			request.FileNamePrefixList = append(request.FileNamePrefixList, &fileNamePrefixList)
		}
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
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm exportImages failed, reason:%+v", logId, err)
		return err
	}

	imageId = *response.Response.ImageId
	d.SetId(strings.Join([]string{imageId, bucketName}, FILED_SP))

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
