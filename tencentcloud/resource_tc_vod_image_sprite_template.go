/*
Provide a resource to create a VOD image sprite template.

Example Usage

```hcl
resource "tencentcloud_vod_image_sprite_template" "foo" {
  sample_type         = "Percent"
  sample_interval     = 10
  row_count           = 3
  column_count        = 3
  name                = "tf-sprite"
  comment             = "test"
  fill_type           = "stretch"
  width               = 128
  height              = 128
  resolution_adaptive = false
}
```

Import

VOD image sprite template can be imported using the id, e.g.

```
$ terraform import tencentcloud_vod_image_sprite_template.foo 51156
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudVodImageSpriteTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodImageSpriteTemplateCreate,
		Read:   resourceTencentCloudVodImageSpriteTemplateRead,
		Update: resourceTencentCloudVodImageSpriteTemplateUpdate,
		Delete: resourceTencentCloudVodImageSpriteTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"sample_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"Percent", "Time"}),
				Description:  "Sampling type. Valid values: `Percent`, `Time`. `Percent`: by percent. `Time`: by time interval.",
			},
			"sample_interval": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Sampling interval. If `sample_type` is `Percent`, sampling will be performed at an interval of the specified percentage. If `sample_type` is `Time`, sampling will be performed at the specified time interval in seconds.",
			},
			"row_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Subimage row count of an image sprite.",
			},
			"column_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Subimage column count of an image sprite.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 64),
				Description:  "Name of a time point screen capturing template. Length limit: 64 characters.",
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(1, 256),
				Description:  "Template description. Length limit: 256 characters.",
			},
			"fill_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "black",
				Description: "Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the source video. The following fill types are supported: `stretch`: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the source video, which may make the screenshot shorter or longer; `black`: fill with black. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with black color blocks. Default value: `black`.",
			},
			"width": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Maximum value of the `width` (or long side) of a screenshot in px. Value range: 0 and [128, 4,096]. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, width will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used. Default value: `0`.",
			},
			"height": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Maximum value of the `height` (or short side) of a screenshot in px. Value range: 0 and [128, 4,096]. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, `width` will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used. Default value: `0`.",
			},
			"resolution_adaptive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Resolution adaption. Valid values: `true`,`false`. `true`: enabled. In this case, `width` represents the long side of a video, while `height` the short side; `false`: disabled. In this case, `width` represents the width of a video, while `height` the height. Default value: `true`.",
			},
			"sub_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subapplication ID in VOD. If you need to access a resource in a subapplication, enter the subapplication ID in this field; otherwise, leave it empty.",
			},
			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of template in ISO date format.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified time of template in ISO date format.",
			},
		},
	}
}

func resourceTencentCloudVodImageSpriteTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_image_sprite_template.create")()

	var (
		logId   = getLogId(contextNil)
		request = vod.NewCreateImageSpriteTemplateRequest()
	)

	request.SampleType = helper.String(d.Get("sample_type").(string))
	request.SampleInterval = helper.IntUint64(d.Get("sample_interval").(int))
	request.RowCount = helper.IntUint64(d.Get("row_count").(int))
	request.ColumnCount = helper.IntUint64(d.Get("column_count").(int))
	request.Name = helper.String(d.Get("name").(string))
	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}
	request.FillType = helper.String((d.Get("fill_type")).(string))
	request.Width = helper.IntUint64(d.Get("width").(int))
	request.Height = helper.IntUint64(d.Get("height").(int))
	request.ResolutionAdaptive = helper.String(RESOLUTION_ADAPTIVE_TO_STRING[d.Get("resolution_adaptive").(bool)])
	if v, ok := d.GetOk("sub_app_id"); ok {
		request.SubAppId = helper.IntUint64(v.(int))
	}

	var response *vod.CreateImageSpriteTemplateResponse
	var err error
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = meta.(*TencentCloudClient).apiV3Conn.UseVodClient().CreateImageSpriteTemplate(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if response == nil || response.Response == nil {
		return fmt.Errorf("for image sprite template creation, response is nil")
	}
	d.SetId(strconv.FormatUint(*response.Response.Definition, 10))

	return resourceTencentCloudVodImageSpriteTemplateRead(d, meta)
}

func resourceTencentCloudVodImageSpriteTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_image_sprite_template.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		id         = d.Id()
		subAppId   = d.Get("sub_app_id").(int)
		client     = meta.(*TencentCloudClient).apiV3Conn
		vodService = VodService{client: client}
	)
	// waiting for refreshing cache
	time.Sleep(30 * time.Second)
	template, has, err := vodService.DescribeImageSpriteTemplatesById(ctx, id, subAppId)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("sample_type", template.SampleType)
	_ = d.Set("sample_interval", template.SampleInterval)
	_ = d.Set("row_count", template.RowCount)
	_ = d.Set("column_count", template.ColumnCount)
	_ = d.Set("name", template.Name)
	_ = d.Set("comment", template.Comment)
	_ = d.Set("fill_type", template.FillType)
	_ = d.Set("width", template.Width)
	_ = d.Set("height", template.Height)
	_ = d.Set("resolution_adaptive", *template.ResolutionAdaptive == "open")
	_ = d.Set("create_time", template.CreateTime)
	_ = d.Set("update_time", template.UpdateTime)

	return nil
}

func resourceTencentCloudVodImageSpriteTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_image_sprite_template.update")()

	var (
		logId      = getLogId(contextNil)
		request    = vod.NewModifyImageSpriteTemplateRequest()
		id         = d.Id()
		changeFlag = false
	)

	idUint, _ := strconv.ParseUint(id, 0, 64)
	request.Definition = &idUint
	if d.HasChange("sample_type") {
		changeFlag = true
		request.SampleType = helper.String(d.Get("sample_type").(string))
	}
	if d.HasChange("sample_interval") {
		changeFlag = true
		request.SampleInterval = helper.IntUint64(d.Get("sample_interval").(int))
	}
	if d.HasChange("row_count") {
		changeFlag = true
		request.RowCount = helper.IntUint64(d.Get("row_count").(int))
	}
	if d.HasChange("column_count") {
		changeFlag = true
		request.ColumnCount = helper.IntUint64(d.Get("column_count").(int))
	}
	if d.HasChange("name") {
		changeFlag = true
		request.Name = helper.String(d.Get("name").(string))
	}
	if d.HasChange("comment") {
		changeFlag = true
		request.Comment = helper.String(d.Get("comment").(string))
	}
	if d.HasChange("fill_type") {
		changeFlag = true
		request.FillType = helper.String(d.Get("fill_type").(string))
	}
	if d.HasChange("width") || d.HasChange("height") || d.HasChange("resolution_adaptive") {
		changeFlag = true
		request.Width = helper.IntUint64(d.Get("width").(int))
		request.Height = helper.IntUint64(d.Get("height").(int))
		request.ResolutionAdaptive = helper.String(RESOLUTION_ADAPTIVE_TO_STRING[d.Get("resolution_adaptive").(bool)])
	}
	if d.HasChange("sub_app_id") {
		changeFlag = true
		request.SubAppId = helper.IntUint64(d.Get("sub_app_id").(int))
	}

	if changeFlag {
		var err error
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err = meta.(*TencentCloudClient).apiV3Conn.UseVodClient().ModifyImageSpriteTemplate(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
				return retryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return resourceTencentCloudVodImageSpriteTemplateRead(d, meta)
	}

	return nil
}

func resourceTencentCloudVodImageSpriteTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_image_sprite_template.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	vodService := VodService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	if err := vodService.DeleteImageSpriteTemplate(ctx, id, uint64(d.Get("sub_app_id").(int))); err != nil {
		return err
	}

	return nil
}
