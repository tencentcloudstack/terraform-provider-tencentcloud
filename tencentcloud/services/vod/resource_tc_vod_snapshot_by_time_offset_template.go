package vod

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudVodSnapshotByTimeOffsetTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodSnapshotByTimeOffsetTemplateCreate,
		Read:   resourceTencentCloudVodSnapshotByTimeOffsetTemplateRead,
		Update: resourceTencentCloudVodSnapshotByTimeOffsetTemplateUpdate,
		Delete: resourceTencentCloudVodSnapshotByTimeOffsetTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 64),
				Description:  "Name of a time point screen capturing template. Length limit: 64 characters.",
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
			"format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Image format. Valid values: `jpg`, `png`. Default value: `jpg`.",
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 256),
				Description:  "Template description. Length limit: 256 characters.",
			},
			"sub_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The VOD [application](https://intl.cloud.tencent.com/document/product/266/14574) ID. For customers who activate VOD service from December 25, 2023, if they want to access resources in a VOD application (whether it's the default application or a newly created one), they must fill in this field with the application ID.",
			},
			"fill_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "black",
				Description: "Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the source video. The following fill types are supported: `stretch`: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the source video, which may make the screenshot `shorter` or `longer`; `black`: fill with black. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with black color blocks. `white`: fill with white. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with white color blocks. `gauss`: fill with Gaussian blur. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with Gaussian blur. Default value: `black`.",
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Template type, value range:\n" +
					"- Preset: system preset template;\n" +
					"- Custom: user-defined templates.",
			},
		},
	}
}

func resourceTencentCloudVodSnapshotByTimeOffsetTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_snapshot_by_time_offset_template.create")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = vod.NewCreateSnapshotByTimeOffsetTemplateRequest()
	)

	request.Name = helper.String(d.Get("name").(string))
	request.Width = helper.IntUint64(d.Get("width").(int))
	request.Height = helper.IntUint64(d.Get("height").(int))
	request.ResolutionAdaptive = helper.String(RESOLUTION_ADAPTIVE_TO_STRING[d.Get("resolution_adaptive").(bool)])
	if v, ok := d.GetOk("format"); ok {
		request.Format = helper.String(v.(string))
	}
	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}
	var resourceId string
	if v, ok := d.GetOk("sub_app_id"); ok {
		subAppId := v.(int)
		resourceId += helper.IntToStr(subAppId)
		resourceId += tccommon.FILED_SP
		request.SubAppId = helper.IntUint64(subAppId)
	}
	request.FillType = helper.String(d.Get("fill_type").(string))

	var response *vod.CreateSnapshotByTimeOffsetTemplateResponse
	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().CreateSnapshotByTimeOffsetTemplate(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation" && sdkError.Message == "invalid vod user" {
					return resource.RetryableError(err)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if response == nil || response.Response == nil {
		return fmt.Errorf("for vod snapshot by time offset template creation, response is nil")
	}
	resourceId += strconv.FormatUint(*response.Response.Definition, 10)
	d.SetId(resourceId)

	return resourceTencentCloudVodSnapshotByTimeOffsetTemplateRead(d, meta)
}

func resourceTencentCloudVodSnapshotByTimeOffsetTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_snapshot_by_time_offset_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		subAppId   int
		definition string
		client     = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		vodService = VodService{client: client}
	)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 2 {
		subAppId = helper.StrToInt(idSplit[0])
		definition = idSplit[1]
	} else {
		definition = d.Id()
	}
	// waiting for refreshing cache
	time.Sleep(30 * time.Second)
	template, has, err := vodService.DescribeSnapshotByTimeOffsetTemplatesById(ctx, definition, subAppId)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", template.Name)
	_ = d.Set("type", template.Type)
	_ = d.Set("width", template.Width)
	_ = d.Set("height", template.Height)
	_ = d.Set("resolution_adaptive", *template.ResolutionAdaptive == "open")
	_ = d.Set("format", template.Format)
	_ = d.Set("comment", template.Comment)
	_ = d.Set("fill_type", template.FillType)
	_ = d.Set("create_time", template.CreateTime)
	_ = d.Set("update_time", template.UpdateTime)
	if subAppId != 0 {
		_ = d.Set("sub_app_id", subAppId)
	}

	return nil
}

func resourceTencentCloudVodSnapshotByTimeOffsetTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_snapshot_by_time_offset_template.update")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = vod.NewModifySnapshotByTimeOffsetTemplateRequest()
		subAppId   int
		definition string
		changeFlag = false
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 2 {
		subAppId = helper.StrToInt(idSplit[0])
		definition = idSplit[1]
		request.SubAppId = helper.IntUint64(subAppId)
	} else {
		definition = d.Id()
		if v, ok := d.GetOk("sub_app_id"); ok {
			request.SubAppId = helper.IntUint64(v.(int))
		}
	}

	request.Definition = helper.StrToUint64Point(definition)

	immutableArgs := []string{"sub_app_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		changeFlag = true
		request.Name = helper.String(d.Get("name").(string))
	}
	if d.HasChange("width") || d.HasChange("height") || d.HasChange("resolution_adaptive") {
		changeFlag = true
		request.Width = helper.IntUint64(d.Get("width").(int))
		request.Height = helper.IntUint64(d.Get("height").(int))
		request.ResolutionAdaptive = helper.String(RESOLUTION_ADAPTIVE_TO_STRING[d.Get("resolution_adaptive").(bool)])
	}
	if d.HasChange("format") {
		changeFlag = true
		request.Format = helper.String(d.Get("format").(string))
	}
	if d.HasChange("comment") {
		changeFlag = true
		request.Comment = helper.String(d.Get("comment").(string))
	}
	if d.HasChange("fill_type") {
		changeFlag = true
		request.FillType = helper.String(d.Get("fill_type").(string))
	}

	if changeFlag {
		var err error
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifySnapshotByTimeOffsetTemplate(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return resourceTencentCloudVodSnapshotByTimeOffsetTemplateRead(d, meta)
	}

	return nil
}

func resourceTencentCloudVodSnapshotByTimeOffsetTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_snapshot_by_time_offset_template.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		subAppId   int
		definition string
	)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) == 2 {
		subAppId = helper.StrToInt(idSplit[0])
		definition = idSplit[1]
	} else {
		definition = d.Id()
		if v, ok := d.GetOk("sub_app_id"); ok {
			subAppId = v.(int)
		}
	}
	vodService := VodService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	if err := vodService.DeleteSnapshotByTimeOffsetTemplate(ctx, definition, uint64(subAppId)); err != nil {
		return err
	}

	return nil
}
