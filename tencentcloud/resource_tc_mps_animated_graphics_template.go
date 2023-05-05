/*
Provides a resource to create a mps animated_graphics_template

Example Usage

```hcl
resource "tencentcloud_mps_animated_graphics_template" "animated_graphics_template" {
  format              = "gif"
  fps                 = 20
  height              = 130
  name                = "terraform-test"
  quality             = 75
  resolution_adaptive = "open"
  width               = 140
}
```

Import

mps animated_graphics_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_animated_graphics_template.animated_graphics_template animated_graphics_template_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMpsAnimatedGraphicsTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsAnimatedGraphicsTemplateCreate,
		Read:   resourceTencentCloudMpsAnimatedGraphicsTemplateRead,
		Update: resourceTencentCloudMpsAnimatedGraphicsTemplateUpdate,
		Delete: resourceTencentCloudMpsAnimatedGraphicsTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"fps": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Frame rate, value range: [1, 30], unit: Hz.",
			},

			"width": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum value of the animation width (or long side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
			},

			"height": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum value of the animation height (or short side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
			},

			"resolution_adaptive": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Adaptive resolution, optional value:open: At this time, Width represents the long side of the video, Height represents the short side of the video.close: At this point, Width represents the width of the video, and Height represents the height of the video.Default value: open.",
			},

			"format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Animation format, the values are gif and webp. Default is gif.",
			},

			"quality": {
				Optional:    true,
				Type:        schema.TypeFloat,
				Description: "Image quality, value range: [1, 100], default value is 75.",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Rotation diagram template name, length limit: 64 characters.",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template description information, length limit: 256 characters.",
			},
		},
	}
}

func resourceTencentCloudMpsAnimatedGraphicsTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_animated_graphics_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mps.NewCreateAnimatedGraphicsTemplateRequest()
		response   = mps.NewCreateAnimatedGraphicsTemplateResponse()
		definition uint64
	)
	if v, ok := d.GetOkExists("fps"); ok {
		request.Fps = helper.IntUint64(v.(int))
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

	if v, ok := d.GetOkExists("quality"); ok {
		request.Quality = helper.Float64(v.(float64))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().CreateAnimatedGraphicsTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps animatedGraphicsTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = *response.Response.Definition
	d.SetId(helper.UInt64ToStr(definition))

	return resourceTencentCloudMpsAnimatedGraphicsTemplateRead(d, meta)
}

func resourceTencentCloudMpsAnimatedGraphicsTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_animated_graphics_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	definition := d.Id()

	animatedGraphicsTemplate, err := service.DescribeMpsAnimatedGraphicsTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if animatedGraphicsTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsAnimatedGraphicsTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if animatedGraphicsTemplate.Fps != nil {
		_ = d.Set("fps", animatedGraphicsTemplate.Fps)
	}

	if animatedGraphicsTemplate.Width != nil {
		_ = d.Set("width", animatedGraphicsTemplate.Width)
	}

	if animatedGraphicsTemplate.Height != nil {
		_ = d.Set("height", animatedGraphicsTemplate.Height)
	}

	if animatedGraphicsTemplate.ResolutionAdaptive != nil {
		_ = d.Set("resolution_adaptive", animatedGraphicsTemplate.ResolutionAdaptive)
	}

	if animatedGraphicsTemplate.Format != nil {
		_ = d.Set("format", animatedGraphicsTemplate.Format)
	}

	if animatedGraphicsTemplate.Quality != nil {
		_ = d.Set("quality", animatedGraphicsTemplate.Quality)
	}

	if animatedGraphicsTemplate.Name != nil {
		_ = d.Set("name", animatedGraphicsTemplate.Name)
	}

	if animatedGraphicsTemplate.Comment != nil {
		_ = d.Set("comment", animatedGraphicsTemplate.Comment)
	}

	return nil
}

func resourceTencentCloudMpsAnimatedGraphicsTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_animated_graphics_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewModifyAnimatedGraphicsTemplateRequest()

	definition := d.Id()

	request.Definition = helper.StrToUint64Point(definition)

	mutableArgs := []string{"fps", "width", "height", "resolution_adaptive", "format", "quality", "name", "comment"}

	needChange := false

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if v, ok := d.GetOkExists("fps"); ok {
			request.Fps = helper.IntUint64(v.(int))
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

		if v, ok := d.GetOkExists("quality"); ok {
			request.Quality = helper.Float64(v.(float64))
		}

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ModifyAnimatedGraphicsTemplate(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mps animatedGraphicsTemplate failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMpsAnimatedGraphicsTemplateRead(d, meta)
}

func resourceTencentCloudMpsAnimatedGraphicsTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_animated_graphics_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}
	definition := d.Id()

	if err := service.DeleteMpsAnimatedGraphicsTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
