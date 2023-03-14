/*
Provides a resource to create a mps snapshot_by_timeoffset_template

Example Usage

```hcl
resource "tencentcloud_mps_snapshot_by_timeoffset_template" "snapshot_by_timeoffset_template" {
  fill_type           = "stretch"
  format              = "jpg"
  height              = 128
  name                = "terraform-test"
  resolution_adaptive = "open"
  width               = 140
}
```

Import

mps snapshot_by_timeoffset_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_snapshot_by_timeoffset_template.snapshot_by_timeoffset_template snapshot_by_timeoffset_template_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMpsSnapshotByTimeoffsetTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsSnapshotByTimeoffsetTemplateCreate,
		Read:   resourceTencentCloudMpsSnapshotByTimeoffsetTemplateRead,
		Update: resourceTencentCloudMpsSnapshotByTimeoffsetTemplateUpdate,
		Delete: resourceTencentCloudMpsSnapshotByTimeoffsetTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Snapshot by timeoffset template name, length limit: 64 characters.",
			},

			"width": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum value of the snapshot width (or long side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
			},

			"height": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum value of the snapshot height (or short side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
			},

			"resolution_adaptive": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Adaptive resolution, optional value:open: At this time, Width represents the long side of the video, Height represents the short side of the video.close: At this point, Width represents the width of the video, and Height represents the height of the video.Default value: open.",
			},

			"format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Image format, the value can be jpg, png, webp. Default is jpg.",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template description information, length limit: 256 characters.",
			},

			"fill_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filling type, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling. Optional filling type:stretch: Stretching, stretching each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched.black: Leave black, keep the video aspect ratio unchanged, and fill the rest of the edge with black.white: Leave blank, keep the aspect ratio of the video, and fill the rest of the edge with white.gauss: Gaussian blur, keep the aspect ratio of the video unchanged, and use Gaussian blur for the rest of the edge.Default value: black.",
			},
		},
	}
}

func resourceTencentCloudMpsSnapshotByTimeoffsetTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_snapshot_by_timeoffset_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mps.NewCreateSnapshotByTimeOffsetTemplateRequest()
		response   = mps.NewCreateSnapshotByTimeOffsetTemplateResponse()
		definition uint64
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("width"); ok {
		request.Width = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("height"); ok {
		request.Height = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("resolution_adaptive"); ok {
		request.ResolutionAdaptive = helper.String(v.(string))
	}

	if v, ok := d.GetOk("format"); ok {
		request.Format = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if v, ok := d.GetOk("fill_type"); ok {
		request.FillType = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().CreateSnapshotByTimeOffsetTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps snapshotByTimeoffsetTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = *response.Response.Definition
	d.SetId(helper.UInt64ToStr(definition))

	return resourceTencentCloudMpsSnapshotByTimeoffsetTemplateRead(d, meta)
}

func resourceTencentCloudMpsSnapshotByTimeoffsetTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_snapshot_by_timeoffset_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	definition := d.Id()

	snapshotByTimeoffsetTemplate, err := service.DescribeMpsSnapshotByTimeoffsetTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if snapshotByTimeoffsetTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsSnapshotByTimeoffsetTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if snapshotByTimeoffsetTemplate.Name != nil {
		_ = d.Set("name", snapshotByTimeoffsetTemplate.Name)
	}

	if snapshotByTimeoffsetTemplate.Width != nil {
		_ = d.Set("width", snapshotByTimeoffsetTemplate.Width)
	}

	if snapshotByTimeoffsetTemplate.Height != nil {
		_ = d.Set("height", snapshotByTimeoffsetTemplate.Height)
	}

	if snapshotByTimeoffsetTemplate.ResolutionAdaptive != nil {
		_ = d.Set("resolution_adaptive", snapshotByTimeoffsetTemplate.ResolutionAdaptive)
	}

	if snapshotByTimeoffsetTemplate.Format != nil {
		_ = d.Set("format", snapshotByTimeoffsetTemplate.Format)
	}

	if snapshotByTimeoffsetTemplate.Comment != nil {
		_ = d.Set("comment", snapshotByTimeoffsetTemplate.Comment)
	}

	if snapshotByTimeoffsetTemplate.FillType != nil {
		_ = d.Set("fill_type", snapshotByTimeoffsetTemplate.FillType)
	}

	return nil
}

func resourceTencentCloudMpsSnapshotByTimeoffsetTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_snapshot_by_timeoffset_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	needChange := false

	request := mps.NewModifySnapshotByTimeOffsetTemplateRequest()

	definition := d.Id()

	request.Definition = helper.StrToUint64Point(definition)

	mutableArgs := []string{"name", "width", "height", "resolution_adaptive", "format", "comment", "fill_type"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("width"); ok {
			request.Width = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("height"); ok {
			request.Height = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("resolution_adaptive"); ok {
			request.ResolutionAdaptive = helper.String(v.(string))
		}

		if v, ok := d.GetOk("format"); ok {
			request.Format = helper.String(v.(string))
		}

		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}

		if v, ok := d.GetOk("fill_type"); ok {
			request.FillType = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ModifySnapshotByTimeOffsetTemplate(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mps snapshotByTimeoffsetTemplate failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMpsSnapshotByTimeoffsetTemplateRead(d, meta)
}

func resourceTencentCloudMpsSnapshotByTimeoffsetTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_snapshot_by_timeoffset_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}
	definition := d.Id()

	if err := service.DeleteMpsSnapshotByTimeoffsetTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
