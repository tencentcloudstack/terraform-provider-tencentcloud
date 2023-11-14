/*
Provides a resource to create a live snapshot_template

Example Usage

```hcl
resource "tencentcloud_live_snapshot_template" "snapshot_template" {
  template_name = ""
  cos_app_id =
  cos_bucket = ""
  cos_region = ""
  description = ""
  snapshot_interval =
  width =
  height =
  porn_flag =
  cos_prefix = ""
  cos_file_name = ""
}
```

Import

live snapshot_template can be imported using the id, e.g.

```
terraform import tencentcloud_live_snapshot_template.snapshot_template snapshot_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudLiveSnapshotTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveSnapshotTemplateCreate,
		Read:   resourceTencentCloudLiveSnapshotTemplateRead,
		Update: resourceTencentCloudLiveSnapshotTemplateUpdate,
		Delete: resourceTencentCloudLiveSnapshotTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template name.Maximum length: 255 bytes.Only letters, digits, underscores, and hyphens can be contained.",
			},

			"cos_app_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "COS application ID.",
			},

			"cos_bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "COS bucket name.Note: the value of `CosBucket` cannot contain `-[appid]`.",
			},

			"cos_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "COS region.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description.Maximum length: 1,024 bytes.Only letters, digits, underscores, and hyphens can be contained.",
			},

			"snapshot_interval": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Screencapturing interval (s). Default value: 10Value range: 2-300.",
			},

			"width": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Screenshot width. Default value: `0` (original width)Value range: 0-3000.",
			},

			"height": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Screenshot height. Default value: `0` (original height)Value range: 0-2000.",
			},

			"porn_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to enable porn detection. 0: no, 1: yes. Default value: 0.",
			},

			"cos_prefix": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "COS Bucket folder prefix.If no value is entered, the default value`/{Year}-{Month}-{Day}`will be used.",
			},

			"cos_file_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "COS filename.If no value is entered, the default value `{StreamID}-screenshot-{Hour}-{Minute}-{Second}-{Width}x{Height}{Ext}`will be used.",
			},
		},
	}
}

func resourceTencentCloudLiveSnapshotTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_snapshot_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = live.NewCreateLiveSnapshotTemplateRequest()
		response   = live.NewCreateLiveSnapshotTemplateResponse()
		templateId int
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().CreateLiveSnapshotTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create live snapshotTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(helper.Int64ToStr(templateId))

	return resourceTencentCloudLiveSnapshotTemplateRead(d, meta)
}

func resourceTencentCloudLiveSnapshotTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_snapshot_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	snapshotTemplateId := d.Id()

	snapshotTemplate, err := service.DescribeLiveSnapshotTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if snapshotTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LiveSnapshotTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

func resourceTencentCloudLiveSnapshotTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_snapshot_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := live.NewModifyLiveSnapshotTemplateRequest()

	snapshotTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"template_name", "cos_app_id", "cos_bucket", "cos_region", "description", "snapshot_interval", "width", "height", "porn_flag", "cos_prefix", "cos_file_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("template_name") {
		if v, ok := d.GetOk("template_name"); ok {
			request.TemplateName = helper.String(v.(string))
		}
	}

	if d.HasChange("cos_app_id") {
		if v, ok := d.GetOkExists("cos_app_id"); ok {
			request.CosAppId = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("cos_bucket") {
		if v, ok := d.GetOk("cos_bucket"); ok {
			request.CosBucket = helper.String(v.(string))
		}
	}

	if d.HasChange("cos_region") {
		if v, ok := d.GetOk("cos_region"); ok {
			request.CosRegion = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("snapshot_interval") {
		if v, ok := d.GetOkExists("snapshot_interval"); ok {
			request.SnapshotInterval = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("width") {
		if v, ok := d.GetOkExists("width"); ok {
			request.Width = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("height") {
		if v, ok := d.GetOkExists("height"); ok {
			request.Height = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("porn_flag") {
		if v, ok := d.GetOkExists("porn_flag"); ok {
			request.PornFlag = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("cos_prefix") {
		if v, ok := d.GetOk("cos_prefix"); ok {
			request.CosPrefix = helper.String(v.(string))
		}
	}

	if d.HasChange("cos_file_name") {
		if v, ok := d.GetOk("cos_file_name"); ok {
			request.CosFileName = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().ModifyLiveSnapshotTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update live snapshotTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLiveSnapshotTemplateRead(d, meta)
}

func resourceTencentCloudLiveSnapshotTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_snapshot_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}
	snapshotTemplateId := d.Id()

	if err := service.DeleteLiveSnapshotTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
