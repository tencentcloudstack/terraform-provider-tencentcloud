package cvm

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCvmImportImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmImportImageCreate,
		Read:   resourceTencentCloudCvmImportImageRead,
		Delete: resourceTencentCloudCvmImportImageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"architecture": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "OS architecture of the image to be imported, `x86_64` or `i386`.",
			},

			"os_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "OS type of the image to be imported. You can call `DescribeImportImageOs` to obtain the list of supported operating systems.",
			},

			"os_version": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "OS version of the image to be imported. You can call `DescribeImportImageOs` to obtain the list of supported operating systems.",
			},

			"image_url": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Address on COS where the image to be imported is stored.",
			},

			"image_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Image name.",
			},

			"image_description": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Image description.",
			},

			"dry_run": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Dry run to check the parameters without performing the operation.",
			},

			"force": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to force import the image.",
			},

			"tag_specification": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Tag description list. This parameter is used to bind a tag to a custom image.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Resource type. Valid values: instance (CVM), host (CDH), image (for image), and keypair (for key). Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"tags": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Tag pairs Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag key.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Tag value.",
									},
								},
							},
						},
					},
				},
			},

			"license_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The license type used to activate the OS after importing an image. Valid values: TencentCloud: Tencent Cloud official license BYOL: Bring Your Own License.",
			},

			"boot_mode": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Boot mode.",
			},
		},
	}
}

func resourceTencentCloudCvmImportImageCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_import_image.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = cvm.NewImportImageRequest()
		imageUrl string
	)
	if v, ok := d.GetOk("architecture"); ok {
		request.Architecture = helper.String(v.(string))
	}

	if v, ok := d.GetOk("os_type"); ok {
		request.OsType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("os_version"); ok {
		request.OsVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_url"); ok {
		imageUrl = v.(string)
		request.ImageUrl = helper.String(imageUrl)
	}

	if v, ok := d.GetOk("image_name"); ok {
		request.ImageName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_description"); ok {
		request.ImageDescription = helper.String(v.(string))
	}

	if v, _ := d.GetOk("dry_run"); v != nil {
		request.DryRun = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("force"); v != nil {
		request.Force = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("tag_specification"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tagSpecification := cvm.TagSpecification{}
			if v, ok := dMap["resource_type"]; ok {
				tagSpecification.ResourceType = helper.String(v.(string))
			}
			if v, ok := dMap["tags"]; ok {
				for _, item := range v.([]interface{}) {
					tagsMap := item.(map[string]interface{})
					tag := cvm.Tag{}
					if v, ok := tagsMap["key"]; ok {
						tag.Key = helper.String(v.(string))
					}
					if v, ok := tagsMap["value"]; ok {
						tag.Value = helper.String(v.(string))
					}
					tagSpecification.Tags = append(tagSpecification.Tags, &tag)
				}
			}
			request.TagSpecification = append(request.TagSpecification, &tagSpecification)
		}
	}

	if v, ok := d.GetOk("license_type"); ok {
		request.LicenseType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("boot_mode"); ok {
		request.BootMode = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ImportImage(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cvm importImage failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(imageUrl)

	return resourceTencentCloudCvmImportImageRead(d, meta)
}

func resourceTencentCloudCvmImportImageRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_import_image.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCvmImportImageDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cvm_import_image.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
