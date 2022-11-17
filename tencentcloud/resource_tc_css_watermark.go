/*
Provides a resource to create a css watermark

Example Usage

```hcl
resource "tencentcloud_css_watermark" "watermark" {
  picture_url = "picture_url"
  watermark_name = "watermark_name"
  x_position = 0
  y_position = 0
  width = 0
  height = 0
}

```
Import

css watermark can be imported using the id, e.g.
```
$ terraform import tencentcloud_css_watermark.watermark watermark_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssWatermark() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudCssWatermarkRead,
		Create: resourceTencentCloudCssWatermarkCreate,
		Update: resourceTencentCloudCssWatermarkUpdate,
		Delete: resourceTencentCloudCssWatermarkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"picture_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "watermark url.",
			},

			"watermark_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "watermark name.",
			},

			"x_position": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "x position of the picture.",
			},

			"y_position": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "y position of the picture.",
			},

			"width": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "width of the picture.",
			},

			"height": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "height of the picture.",
			},

			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "status. 0: not used, 1: used.",
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
		response    *css.AddLiveWatermarkResponse
		watermarkId string
	)

	if v, ok := d.GetOk("picture_url"); ok {

		request.PictureUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("watermark_name"); ok {

		request.WatermarkName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("x_position"); ok {
		request.XPosition = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("y_position"); ok {
		request.YPosition = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("width"); ok {
		request.Width = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("height"); ok {
		request.Height = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().AddLiveWatermark(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create css watermark failed, reason:%+v", logId, err)
		return err
	}

	watermarkId = helper.UInt64ToStr(*response.Response.WatermarkId)

	d.SetId(watermarkId)
	return resourceTencentCloudCssWatermarkRead(d, meta)
}

func resourceTencentCloudCssWatermarkRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_watermark.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	watermarkId := d.Id()

	watermark, err := service.DescribeCssWatermark(ctx, watermarkId)

	if err != nil {
		return err
	}

	if watermark == nil {
		d.SetId("")
		return fmt.Errorf("resource `watermark` %s does not exist", watermarkId)
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

	if watermark.Status != nil {
		_ = d.Set("status", watermark.Status)
	}

	return nil
}

func resourceTencentCloudCssWatermarkUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_watermark.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	// ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := css.NewUpdateLiveWatermarkRequest()

	watermarkId := d.Id()

	request.WatermarkId = helper.StrToInt64Point(watermarkId)

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
		if v, ok := d.GetOk("x_position"); ok {
			request.XPosition = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("y_position") {
		if v, ok := d.GetOk("y_position"); ok {
			request.YPosition = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("width") {
		if v, ok := d.GetOk("width"); ok {
			request.Width = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("height") {
		if v, ok := d.GetOk("height"); ok {
			request.Height = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().UpdateLiveWatermark(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create css watermark failed, reason:%+v", logId, err)
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

	if err := service.DeleteCssWatermarkById(ctx, helper.StrToInt64Point(watermarkId)); err != nil {
		return err
	}

	return nil
}
