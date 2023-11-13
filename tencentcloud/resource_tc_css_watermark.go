/*
Provides a resource to create a css watermark

Example Usage

```hcl
resource "tencentcloud_css_watermark" "watermark" {
  picture_url = &lt;nil&gt;
  watermark_name = &lt;nil&gt;
  x_position = &lt;nil&gt;
  y_position = &lt;nil&gt;
  width = &lt;nil&gt;
  height = &lt;nil&gt;
}
```

Import

css watermark can be imported using the id, e.g.

```
terraform import tencentcloud_css_watermark.watermark watermark_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCssWatermark() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssWatermarkCreate,
		Read:   resourceTencentCloudCssWatermarkRead,
		Update: resourceTencentCloudCssWatermarkUpdate,
		Delete: resourceTencentCloudCssWatermarkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"picture_url": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Watermark url.",
			},

			"watermark_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Watermark name.",
			},

			"x_position": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "X position of the picture.",
			},

			"y_position": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Y position of the picture.",
			},

			"width": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Width of the picture.",
			},

			"height": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Height of the picture.",
			},
		},
	}
}

func resourceTencentCloudCssWatermarkCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_watermark.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = css.NewAddLiveWatermarkRequest()
		response    = css.NewAddLiveWatermarkResponse()
		watermarkId int
	)
	if v, ok := d.GetOk("picture_url"); ok {
		request.PictureUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("watermark_name"); ok {
		request.WatermarkName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("x_position"); ok {
		request.XPosition = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("y_position"); ok {
		request.YPosition = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("width"); ok {
		request.Width = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("height"); ok {
		request.Height = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().AddLiveWatermark(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css watermark failed, reason:%+v", logId, err)
		return err
	}

	watermarkId = *response.Response.WatermarkId
	d.SetId(helper.Int64ToStr(watermarkId))

	return resourceTencentCloudCssWatermarkRead(d, meta)
}

func resourceTencentCloudCssWatermarkRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_watermark.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	watermarkId := d.Id()

	watermark, err := service.DescribeCssWatermarkById(ctx, watermarkId)
	if err != nil {
		return err
	}

	if watermark == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssWatermark` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if watermark.PictureUrl != nil {
		_ = d.Set("picture_url", watermark.PictureUrl)
	}

	if watermark.WatermarkName != nil {
		_ = d.Set("watermark_name", watermark.WatermarkName)
	}

	if watermark.XPosition != nil {
		_ = d.Set("x_position", watermark.XPosition)
	}

	if watermark.YPosition != nil {
		_ = d.Set("y_position", watermark.YPosition)
	}

	if watermark.Width != nil {
		_ = d.Set("width", watermark.Width)
	}

	if watermark.Height != nil {
		_ = d.Set("height", watermark.Height)
	}

	return nil
}

func resourceTencentCloudCssWatermarkUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_watermark.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := css.NewUpdateLiveWatermarkRequest()

	watermarkId := d.Id()

	request.WatermarkId = &watermarkId

	immutableArgs := []string{"picture_url", "watermark_name", "x_position", "y_position", "width", "height"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("picture_url") {
		if v, ok := d.GetOk("picture_url"); ok {
			request.PictureUrl = helper.String(v.(string))
		}
	}

	if d.HasChange("watermark_name") {
		if v, ok := d.GetOk("watermark_name"); ok {
			request.WatermarkName = helper.String(v.(string))
		}
	}

	if d.HasChange("x_position") {
		if v, ok := d.GetOkExists("x_position"); ok {
			request.XPosition = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("y_position") {
		if v, ok := d.GetOkExists("y_position"); ok {
			request.YPosition = helper.IntInt64(v.(int))
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().UpdateLiveWatermark(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update css watermark failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCssWatermarkRead(d, meta)
}

func resourceTencentCloudCssWatermarkDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_watermark.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}
	watermarkId := d.Id()

	if err := service.DeleteCssWatermarkById(ctx, watermarkId); err != nil {
		return err
	}

	return nil
}
