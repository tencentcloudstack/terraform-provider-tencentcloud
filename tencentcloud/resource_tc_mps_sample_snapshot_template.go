/*
Provides a resource to create a mps sample_snapshot_template

Example Usage

```hcl
resource "tencentcloud_mps_sample_snapshot_template" "sample_snapshot_template" {
  sample_type = &lt;nil&gt;
  sample_interval = &lt;nil&gt;
  name = &lt;nil&gt;
  width = 0
  height = 0
  resolution_adaptive = "open"
  format = "jpg"
  comment = &lt;nil&gt;
  fill_type = "black"
}
```

Import

mps sample_snapshot_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_sample_snapshot_template.sample_snapshot_template sample_snapshot_template_id
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

func resourceTencentCloudMpsSampleSnapshotTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsSampleSnapshotTemplateCreate,
		Read:   resourceTencentCloudMpsSampleSnapshotTemplateRead,
		Update: resourceTencentCloudMpsSampleSnapshotTemplateUpdate,
		Delete: resourceTencentCloudMpsSampleSnapshotTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sample_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Sampling snapshot type, optional value:Percent/Time.",
			},

			"sample_interval": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Sampling interval.When SampleType is Percent, specify the percentage of the sampling interval.When SampleType is Time, specify the sampling interval time in seconds.",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sample snapshot template name, length limit: 64 characters.",
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

func resourceTencentCloudMpsSampleSnapshotTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_sample_snapshot_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mps.NewCreateSampleSnapshotTemplateRequest()
		response   = mps.NewCreateSampleSnapshotTemplateResponse()
		definition uint64
	)
	if v, ok := d.GetOk("sample_type"); ok {
		request.SampleType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("sample_interval"); ok {
		request.SampleInterval = helper.IntUint64(v.(int))
	}

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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().CreateSampleSnapshotTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps sampleSnapshotTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = *response.Response.Definition
	d.SetId(helper.UInt64ToStr(definition))

	return resourceTencentCloudMpsSampleSnapshotTemplateRead(d, meta)
}

func resourceTencentCloudMpsSampleSnapshotTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_sample_snapshot_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	definition := d.Id()

	sampleSnapshotTemplate, err := service.DescribeMpsSampleSnapshotTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if sampleSnapshotTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsSampleSnapshotTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if sampleSnapshotTemplate.SampleType != nil {
		_ = d.Set("sample_type", sampleSnapshotTemplate.SampleType)
	}

	if sampleSnapshotTemplate.SampleInterval != nil {
		_ = d.Set("sample_interval", sampleSnapshotTemplate.SampleInterval)
	}

	if sampleSnapshotTemplate.Name != nil {
		_ = d.Set("name", sampleSnapshotTemplate.Name)
	}

	if sampleSnapshotTemplate.Width != nil {
		_ = d.Set("width", sampleSnapshotTemplate.Width)
	}

	if sampleSnapshotTemplate.Height != nil {
		_ = d.Set("height", sampleSnapshotTemplate.Height)
	}

	if sampleSnapshotTemplate.ResolutionAdaptive != nil {
		_ = d.Set("resolution_adaptive", sampleSnapshotTemplate.ResolutionAdaptive)
	}

	if sampleSnapshotTemplate.Format != nil {
		_ = d.Set("format", sampleSnapshotTemplate.Format)
	}

	if sampleSnapshotTemplate.Comment != nil {
		_ = d.Set("comment", sampleSnapshotTemplate.Comment)
	}

	if sampleSnapshotTemplate.FillType != nil {
		_ = d.Set("fill_type", sampleSnapshotTemplate.FillType)
	}

	return nil
}

func resourceTencentCloudMpsSampleSnapshotTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_sample_snapshot_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewModifySampleSnapshotTemplateRequest()

	definition := d.Id()

	request.Definition = helper.StrToUint64Point(definition)

	mutableArgs := []string{"sample_type", "sample_interval", "name", "width", "height", "resolution_adaptive", "format", "comment", "fill_type"}

	needChange := false

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if v, ok := d.GetOk("sample_type"); ok {
			request.SampleType = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("sample_interval"); ok {
			request.SampleInterval = helper.IntUint64(v.(int))
		}

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
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ModifySampleSnapshotTemplate(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mps sampleSnapshotTemplate failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMpsSampleSnapshotTemplateRead(d, meta)
}

func resourceTencentCloudMpsSampleSnapshotTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_sample_snapshot_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}
	definition := d.Id()

	if err := service.DeleteMpsSampleSnapshotTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
