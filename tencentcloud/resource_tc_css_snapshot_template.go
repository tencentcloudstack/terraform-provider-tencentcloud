/*
Provides a resource to create a css snapshot_template

Example Usage

```hcl
resource "tencentcloud_css_snapshot_template" "snapshot_template" {
    cos_app_id        = 1308919341
    cos_bucket        = "keep-bucket"
    cos_region        = "ap-guangzhou"
    description       = "snapshot template"
    height            = 0
    porn_flag         = 0
    snapshot_interval = 2
    template_name     = "tf-snapshot-template"
    width             = 0
}
```

Import

css snapshot_template can be imported using the id, e.g.

```
terraform import tencentcloud_css_snapshot_template.snapshot_template templateId
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssSnapshotTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssSnapshotTemplateCreate,
		Read:   resourceTencentCloudCssSnapshotTemplateRead,
		Update: resourceTencentCloudCssSnapshotTemplateUpdate,
		Delete: resourceTencentCloudCssSnapshotTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template name. Maximum length: 255 bytes. Only Chinese, English, numbers, `_`, `-` are supported.",
			},

			"cos_app_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Cos application ID.",
			},

			"cos_bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cos bucket name. Note: The CosBucket parameter value cannot include the - [appid] part.",
			},

			"cos_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cos region.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description information. Maximum length: 1024 bytes. Only `Chinese`, `English`, `numbers`, `_`, `-` are supported.",
			},

			"snapshot_interval": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Screenshot interval, unit: s, default: 10s. Range: 2s~300s.",
			},

			"width": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Screenshot width. Default: 0 (original width). Range: 0-3000.",
			},

			"height": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Screenshot height. Default: 0 (original height). Range: 0-2000.",
			},

			"porn_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether porn is enabled, 0: not enabled, 1: enabled. Default: 0.",
			},

			"cos_prefix": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cos Bucket folder prefix. If it is empty, set according to the default value /{Year}-{Month}-{Day}/.",
			},

			"cos_file_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cos file name. If it is empty, set according to the default value {StreamID}-screenshot-{Hour}-{Minute}-{Second}-{Width}x{Height}{Ext}.",
			},
		},
	}
}

func resourceTencentCloudCssSnapshotTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_snapshot_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = css.NewCreateLiveSnapshotTemplateRequest()
		response   = css.NewCreateLiveSnapshotTemplateResponse()
		templateId int64
	)
	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("cos_app_id"); ok {
		request.CosAppId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cos_bucket"); ok {
		request.CosBucket = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cos_region"); ok {
		request.CosRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("snapshot_interval"); ok {
		request.SnapshotInterval = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("width"); ok {
		request.Width = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("height"); ok {
		request.Height = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("porn_flag"); ok {
		request.PornFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cos_prefix"); ok {
		request.CosPrefix = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cos_file_name"); ok {
		request.CosFileName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().CreateLiveSnapshotTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css snapshotTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(helper.Int64ToStr(templateId))

	return resourceTencentCloudCssSnapshotTemplateRead(d, meta)
}

func resourceTencentCloudCssSnapshotTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_snapshot_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	templateId := d.Id()
	templateIdInt64, _ := strconv.ParseInt(templateId, 10, 64)

	snapshotTemplate, err := service.DescribeCssSnapshotTemplateById(ctx, templateIdInt64)
	if err != nil {
		return err
	}

	if snapshotTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssSnapshotTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if snapshotTemplate.TemplateName != nil {
		_ = d.Set("template_name", snapshotTemplate.TemplateName)
	}

	if snapshotTemplate.CosAppId != nil {
		_ = d.Set("cos_app_id", snapshotTemplate.CosAppId)
	}

	if snapshotTemplate.CosBucket != nil {
		_ = d.Set("cos_bucket", snapshotTemplate.CosBucket)
	}

	if snapshotTemplate.CosRegion != nil {
		_ = d.Set("cos_region", snapshotTemplate.CosRegion)
	}

	if snapshotTemplate.Description != nil {
		_ = d.Set("description", snapshotTemplate.Description)
	}

	if snapshotTemplate.SnapshotInterval != nil {
		_ = d.Set("snapshot_interval", snapshotTemplate.SnapshotInterval)
	}

	if snapshotTemplate.Width != nil {
		_ = d.Set("width", snapshotTemplate.Width)
	}

	if snapshotTemplate.Height != nil {
		_ = d.Set("height", snapshotTemplate.Height)
	}

	if snapshotTemplate.PornFlag != nil {
		_ = d.Set("porn_flag", snapshotTemplate.PornFlag)
	}

	if snapshotTemplate.CosPrefix != nil {
		_ = d.Set("cos_prefix", snapshotTemplate.CosPrefix)
	}

	if snapshotTemplate.CosFileName != nil {
		_ = d.Set("cos_file_name", snapshotTemplate.CosFileName)
	}

	return nil
}

func resourceTencentCloudCssSnapshotTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_snapshot_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := css.NewModifyLiveSnapshotTemplateRequest()

	templateId := d.Id()
	templateIdInt64, _ := strconv.ParseInt(templateId, 10, 64)

	request.TemplateId = &templateIdInt64

	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("cos_app_id"); ok {
		request.CosAppId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cos_bucket"); ok {
		request.CosBucket = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cos_region"); ok {
		request.CosRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("snapshot_interval"); ok {
		request.SnapshotInterval = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("width"); ok {
		request.Width = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("height"); ok {
		request.Height = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("porn_flag"); ok {
		request.PornFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cos_prefix"); ok {
		request.CosPrefix = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cos_file_name"); ok {
		request.CosFileName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().ModifyLiveSnapshotTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update css snapshotTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCssSnapshotTemplateRead(d, meta)
}

func resourceTencentCloudCssSnapshotTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_snapshot_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}
	templateId := d.Id()
	templateIdInt64, _ := strconv.ParseInt(templateId, 10, 64)

	if err := service.DeleteCssSnapshotTemplateById(ctx, templateIdInt64); err != nil {
		return err
	}

	return nil
}
