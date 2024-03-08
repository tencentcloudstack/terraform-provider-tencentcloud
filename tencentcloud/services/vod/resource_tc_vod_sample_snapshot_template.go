package vod

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVodSampleSnapshotTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodSampleSnapshotTemplateCreate,
		Read:   resourceTencentCloudVodSampleSnapshotTemplateRead,
		Update: resourceTencentCloudVodSampleSnapshotTemplateUpdate,
		Delete: resourceTencentCloudVodSampleSnapshotTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"sample_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Sampled screencapturing type. Valid values: Percent: by percent. Time: by time interval.",
			},

			"sample_interval": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Sampling interval. If `SampleType` is `Percent`, sampling will be performed at an interval of the specified percentage. If `SampleType` is `Time`, sampling will be performed at the specified time interval in seconds.",
			},

			"sub_app_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The VOD [application](https://intl.cloud.tencent.com/document/product/266/14574) ID. For customers who activate VOD service from December 25, 2023, if they want to access resources in a VOD application (whether it's the default application or a newly created one), they must fill in this field with the application ID.",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Name of a sampled screencapturing template. Length limit: 64 characters.",
			},

			"width": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maximum value of the width (or long side) of a screenshot in px. Value range: 0 and [128, 4,096]. If both `Width` and `Height` are 0, the resolution will be the same as that of the source video; If `Width` is 0, but `Height` is not 0, `Width` will be proportionally scaled; If `Width` is not 0, but `Height` is 0, `Height` will be proportionally scaled; If both `Width` and `Height` are not 0, the custom resolution will be used.Default value: 0.",
			},

			"height": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Maximum value of the height (or short side) of a screenshot in px. Value range: 0 and [128, 4,096]. If both `Width` and `Height` are 0, the resolution will be the same as that of the source video; If `Width` is 0, but `Height` is not 0, `Width` will be proportionally scaled; If `Width` is not 0, but `Height` is 0, `Height` will be proportionally scaled; If both `Width` and `Height` are not 0, the custom resolution will be used.Default value: 0.",
			},

			"resolution_adaptive": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Resolution adaption. Valid values: open: enabled. In this case, `Width` represents the long side of a video, while `Height` the short side; close: disabled. In this case, `Width` represents the width of a video, while `Height` the height.Default value: open.",
			},

			"format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Image format. Valid values: jpg, png. Default value: jpg.",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template description. Length limit: 256 characters.",
			},

			"fill_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fill type. Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the source video. The following fill types are supported:  stretch: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the source video, which may make the screenshot shorter or longer; black: fill with black. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with black color blocks. white: fill with white. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with white color blocks. gauss: fill with Gaussian blur. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with Gaussian blur.Default value: black.",
			},
		},
	}
}

func resourceTencentCloudVodSampleSnapshotTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_sample_snapshot_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = vod.NewCreateSampleSnapshotTemplateRequest()
		response = vod.NewCreateSampleSnapshotTemplateResponse()
		subAppId string
	)

	if v, ok := d.GetOk("sample_type"); ok {
		request.SampleType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("sample_interval"); ok {
		request.SampleInterval = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("sub_app_id"); ok {
		subAppId = helper.IntToStr(v.(int))
		request.SubAppId = helper.IntUint64(v.(int))
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().CreateSampleSnapshotTemplate(request)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation" && sdkError.Message == "invalid vod user" {
					return resource.RetryableError(e)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return resource.NonRetryableError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vod sampleSnapshotTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition := *response.Response.Definition

	d.SetId(subAppId + tccommon.FILED_SP + helper.UInt64ToStr(definition))

	return resourceTencentCloudVodSampleSnapshotTemplateRead(d, meta)
}

func resourceTencentCloudVodSampleSnapshotTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_sample_snapshot_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("sample snapshot id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]
	sampleSnapshotTemplate, err := service.DescribeVodSampleSnapshotTemplateById(ctx, helper.StrToUInt64(subAppId), helper.StrToUInt64(definition))
	if err != nil {
		return err
	}

	if sampleSnapshotTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VodSampleSnapshotTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("sub_app_id", helper.StrToInt(subAppId))
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

func resourceTencentCloudVodSampleSnapshotTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_sample_snapshot_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vod.NewModifySampleSnapshotTemplateRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("sample snapshot id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	immutableArgs := []string{"sub_app_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.Definition = helper.StrToUint64Point(definition)
	request.SubAppId = helper.StrToUint64Point(subAppId)

	if d.HasChange("sample_type") || d.HasChange("sample_interval") || d.HasChange("name") || d.HasChange("width") || d.HasChange("height") || d.HasChange("resolution_adaptive") || d.HasChange("format") || d.HasChange("comment") || d.HasChange("fill_type") {
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
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifySampleSnapshotTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vod sampleSnapshotTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVodSampleSnapshotTemplateRead(d, meta)
}

func resourceTencentCloudVodSampleSnapshotTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_sample_snapshot_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("sample snapshot id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	if err := service.DeleteVodSampleSnapshotTemplateById(ctx, helper.StrToUInt64(subAppId), helper.StrToUInt64(definition)); err != nil {
		return err
	}

	return nil
}
