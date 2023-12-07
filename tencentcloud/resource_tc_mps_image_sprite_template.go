package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMpsImageSpriteTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsImageSpriteTemplateCreate,
		Read:   resourceTencentCloudMpsImageSpriteTemplateRead,
		Update: resourceTencentCloudMpsImageSpriteTemplateUpdate,
		Delete: resourceTencentCloudMpsImageSpriteTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sample_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Sampling type, optional value:Percent/Time.",
			},

			"sample_interval": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Sampling interval.When SampleType is Percent, specify the percentage of the sampling interval.When SampleType is Time, specify the sampling interval time in seconds.",
			},

			"row_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The number of rows in the small image in the sprite.",
			},

			"column_count": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The number of columns in the small image in the sprite.",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Image sprite template name, length limit: 64 characters.",
			},

			"width": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum value of the width (or long side) of the small image in the sprite image, value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
			},

			"height": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum value of the height (or short side) of the small image in the sprite image, value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
			},

			"resolution_adaptive": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Adaptive resolution, optional value:open: At this time, Width represents the long side of the video, Height represents the short side of the video.close: At this point, Width represents the width of the video, and Height represents the height of the video.Default value: open.",
			},

			"fill_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filling type, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling. Optional filling type:stretch: Stretching, stretching each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched.black: Leave black, keep the video aspect ratio unchanged, and fill the rest of the edge with black.Default value: black.",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template description information, length limit: 256 characters.",
			},

			"format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Image format, the value can be jpg, png, webp. Default is jpg.",
			},
		},
	}
}

func resourceTencentCloudMpsImageSpriteTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_image_sprite_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mps.NewCreateImageSpriteTemplateRequest()
		response   = mps.NewCreateImageSpriteTemplateResponse()
		definition uint64
	)
	if v, ok := d.GetOk("sample_type"); ok {
		request.SampleType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("sample_interval"); ok {
		request.SampleInterval = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("row_count"); ok {
		request.RowCount = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("column_count"); ok {
		request.ColumnCount = helper.IntUint64(v.(int))
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

	if v, ok := d.GetOk("fill_type"); ok {
		request.FillType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if v, ok := d.GetOk("format"); ok {
		request.Format = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().CreateImageSpriteTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps imageSpriteTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = *response.Response.Definition
	d.SetId(helper.UInt64ToStr(definition))

	return resourceTencentCloudMpsImageSpriteTemplateRead(d, meta)
}

func resourceTencentCloudMpsImageSpriteTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_image_sprite_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	definition := d.Id()

	imageSpriteTemplate, err := service.DescribeMpsImageSpriteTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if imageSpriteTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsImageSpriteTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if imageSpriteTemplate.SampleType != nil {
		_ = d.Set("sample_type", imageSpriteTemplate.SampleType)
	}

	if imageSpriteTemplate.SampleInterval != nil {
		_ = d.Set("sample_interval", imageSpriteTemplate.SampleInterval)
	}

	if imageSpriteTemplate.RowCount != nil {
		_ = d.Set("row_count", imageSpriteTemplate.RowCount)
	}

	if imageSpriteTemplate.ColumnCount != nil {
		_ = d.Set("column_count", imageSpriteTemplate.ColumnCount)
	}

	if imageSpriteTemplate.Name != nil {
		_ = d.Set("name", imageSpriteTemplate.Name)
	}

	if imageSpriteTemplate.Width != nil {
		_ = d.Set("width", imageSpriteTemplate.Width)
	}

	if imageSpriteTemplate.Height != nil {
		_ = d.Set("height", imageSpriteTemplate.Height)
	}

	if imageSpriteTemplate.ResolutionAdaptive != nil {
		_ = d.Set("resolution_adaptive", imageSpriteTemplate.ResolutionAdaptive)
	}

	if imageSpriteTemplate.FillType != nil {
		_ = d.Set("fill_type", imageSpriteTemplate.FillType)
	}

	if imageSpriteTemplate.Comment != nil {
		_ = d.Set("comment", imageSpriteTemplate.Comment)
	}

	if imageSpriteTemplate.Format != nil {
		_ = d.Set("format", imageSpriteTemplate.Format)
	}

	return nil
}

func resourceTencentCloudMpsImageSpriteTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_image_sprite_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewModifyImageSpriteTemplateRequest()

	definition := d.Id()

	request.Definition = helper.StrToUint64Point(definition)

	mutableArgs := []string{"sample_type", "sample_interval", "row_count", "column_count", "name", "width", "height", "resolution_adaptive", "fill_type", "comment", "format"}

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

		if v, ok := d.GetOkExists("row_count"); ok {
			request.RowCount = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("column_count"); ok {
			request.ColumnCount = helper.IntUint64(v.(int))
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

		if v, ok := d.GetOk("fill_type"); ok {
			request.FillType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}

		if v, ok := d.GetOk("format"); ok {
			request.Format = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ModifyImageSpriteTemplate(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mps imageSpriteTemplate failed, reason:%+v", logId, err)
			return err
		}

	}

	return resourceTencentCloudMpsImageSpriteTemplateRead(d, meta)
}

func resourceTencentCloudMpsImageSpriteTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_image_sprite_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}
	definition := d.Id()

	if err := service.DeleteMpsImageSpriteTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
