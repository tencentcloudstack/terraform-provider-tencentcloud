package vod

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudVodSuperPlayerConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodSuperPlayerConfigCreate,
		Read:   resourceTencentCloudVodSuperPlayerConfigRead,
		Update: resourceTencentCloudVodSuperPlayerConfigUpdate,
		Delete: resourceTencentCloudVodSuperPlayerConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 64),
				Description:  "Player configuration name, which can contain up to 64 letters, digits, underscores, and hyphens (such as test_ABC-123) and must be unique under a user.",
			},
			"drm_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Switch of DRM-protected adaptive bitstream playback: `true`: enabled, indicating to play back only output adaptive bitstreams protected by DRM; `false`: disabled, indicating to play back unencrypted output adaptive bitstreams. Default value: `false`.",
			},
			"adaptive_dynamic_streaming_definition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the unencrypted adaptive bitrate streaming template that allows output, which is required if `drm_switch` is `false`.",
			},
			"drm_streaming_info": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				Description: "Content of the DRM-protected adaptive bitrate streaming template that allows output, which is required if `drm_switch` is `true`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"simple_aes_definition": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the adaptive dynamic streaming template whose protection type is `SimpleAES`.",
						},
					},
				},
			},
			"image_sprite_definition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the image sprite template that allows output.",
			},
			"resolution_names": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Display name of player for substreams with different resolutions. If this parameter is left empty or an empty array, the default configuration will be used: `min_edge_length: 240, name: LD`; `min_edge_length: 480, name: SD`; `min_edge_length: 720, name: HD`; `min_edge_length: 1080, name: FHD`; `min_edge_length: 1440, name: 2K`; `min_edge_length: 2160, name: 4K`; `min_edge_length: 4320, name: 8K`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_edge_length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Length of video short side in px.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Display name.",
						},
					},
				},
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Default",
				Description: "Domain name used for playback. If it is left empty or set to `Default`, the domain name configured in [Default Distribution Configuration](https://cloud.tencent.com/document/product/266/33373) will be used. `Default` by default.",
			},
			"scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Default",
				Description: "Scheme used for playback. If it is left empty or set to `Default`, the scheme configured in [Default Distribution Configuration](https://cloud.tencent.com/document/product/266/33373) will be used. Other valid values: `HTTP`; `HTTPS`.",
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

func resourceTencentCloudVodSuperPlayerConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_super_player_config.create")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = vod.NewCreateSuperPlayerConfigRequest()
	)

	request.Name = helper.String(d.Get("name").(string))
	request.DrmSwitch = helper.String(DRM_SWITCH_TO_STRING[d.Get("drm_switch").(bool)])
	if v, ok := d.GetOk("adaptive_dynamic_streaming_definition"); ok {
		idUint, _ := strconv.ParseUint(v.(string), 0, 64)
		request.AdaptiveDynamicStreamingDefinition = &idUint
	}
	if v, ok := d.GetOk("drm_streaming_info"); ok {
		request.DrmStreamingsInfo = &vod.DrmStreamingsInfo{
			SimpleAesDefinition: func(value interface{}) *uint64 {
				vv := value.([]interface{})
				vvv := vv[0].(map[string]interface{})
				idUint, _ := strconv.ParseUint(vvv["simple_aes_definition"].(string), 0, 64)
				return &idUint
			}(v),
		}
	}
	if v, ok := d.GetOk("image_sprite_definition"); ok {
		idUint, _ := strconv.ParseUint(v.(string), 0, 64)
		request.ImageSpriteDefinition = &idUint
	}
	resolutionNames := []*vod.ResolutionNameInfo{
		{
			MinEdgeLength: helper.Uint64(240),
			Name:          helper.String("LD"),
		},
		{
			MinEdgeLength: helper.Uint64(480),
			Name:          helper.String("SD"),
		},
		{
			MinEdgeLength: helper.Uint64(720),
			Name:          helper.String("HD"),
		},
		{
			MinEdgeLength: helper.Uint64(1080),
			Name:          helper.String("FHD"),
		},
		{
			MinEdgeLength: helper.Uint64(1440),
			Name:          helper.String("2K"),
		},
		{
			MinEdgeLength: helper.Uint64(2160),
			Name:          helper.String("4K"),
		},
		{
			MinEdgeLength: helper.Uint64(4320),
			Name:          helper.String("8K"),
		},
	}
	if v, ok := d.GetOk("resolution_names"); ok {
		resolutionNames = make([]*vod.ResolutionNameInfo, 0, len(v.([]interface{})))
		for _, item := range v.([]interface{}) {
			itemV := item.(map[string]interface{})
			resolutionNames = append(resolutionNames, &vod.ResolutionNameInfo{
				MinEdgeLength: helper.IntUint64(itemV["min_edge_length"].(int)),
				Name:          helper.String(itemV["name"].(string)),
			})
		}
	}
	request.ResolutionNames = resolutionNames
	request.Domain = helper.String(d.Get("domain").(string))
	request.Scheme = helper.String(d.Get("scheme").(string))
	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}
	if v, ok := d.GetOk("sub_app_id"); ok {
		request.SubAppId = helper.IntUint64(v.(int))
	}

	var err error
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().CreateSuperPlayerConfig(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return tccommon.RetryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(d.Get("name").(string))

	return resourceTencentCloudVodSuperPlayerConfigRead(d, meta)
}

func resourceTencentCloudVodSuperPlayerConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_super_player_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		id         = d.Id()
		client     = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		vodService = VodService{client: client}
	)
	config, has, err := vodService.DescribeSuperPlayerConfigsById(ctx, id)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", config.Name)
	_ = d.Set("drm_switch", *config.DrmSwitch == "ON")
	// workaround for AdaptiveDynamicStreamingDefinition para cuz it's dirty data.
	if *config.DrmSwitch == "OFF" {
		_ = d.Set("adaptive_dynamic_streaming_definition", strconv.FormatUint(*config.AdaptiveDynamicStreamingDefinition, 10))
	}
	if config.DrmStreamingsInfo != nil && config.DrmStreamingsInfo.SimpleAesDefinition != nil {
		_ = d.Set("drm_streaming_info", []map[string]interface{}{
			{
				"simple_aes_definition": strconv.FormatUint(*config.DrmStreamingsInfo.SimpleAesDefinition, 10),
			},
		})
	}
	if config.ImageSpriteDefinition != nil {
		_ = d.Set("image_sprite_definition", strconv.FormatUint(*config.ImageSpriteDefinition, 10))
	}
	_ = d.Set("resolution_names", func() []map[string]interface{} {
		namesMap := make([]map[string]interface{}, 0, len(config.ResolutionNameSet))
		for _, v := range config.ResolutionNameSet {
			namesMap = append(namesMap, map[string]interface{}{
				"min_edge_length": v.MinEdgeLength,
				"name":            v.Name,
			})
		}
		return namesMap
	}())
	_ = d.Set("domain", config.Domain)
	_ = d.Set("scheme", config.Scheme)
	_ = d.Set("comment", config.Comment)
	_ = d.Set("create_time", config.CreateTime)
	_ = d.Set("update_time", config.UpdateTime)

	return nil
}

func resourceTencentCloudVodSuperPlayerConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_super_player_config.update")()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		request    = vod.NewModifySuperPlayerConfigRequest()
		id         = d.Id()
		changeFlag = false
	)

	request.Name = &id
	if d.HasChange("drm_switch") {
		changeFlag = true
		request.DrmSwitch = helper.String(DRM_SWITCH_TO_STRING[d.Get("drm_switch").(bool)])
	}
	if d.HasChange("adaptive_dynamic_streaming_definition") {
		changeFlag = true
		idUint, _ := strconv.ParseUint(d.Get("adaptive_dynamic_streaming_definition").(string), 0, 64)
		request.AdaptiveDynamicStreamingDefinition = &idUint
	}
	if d.HasChange("drm_streaming_info") {
		changeFlag = true
		request.DrmStreamingsInfo = func() *vod.DrmStreamingsInfoForUpdate {
			if v, ok := d.GetOk("drm_streaming_info"); !ok {
				return nil
			} else {
				return &vod.DrmStreamingsInfoForUpdate{
					SimpleAesDefinition: func(value interface{}) *uint64 {
						vv := value.([]interface{})
						vvv := vv[0].(map[string]interface{})
						idUint, _ := strconv.ParseUint(vvv["simple_aes_definition"].(string), 0, 64)
						return &idUint
					}(v),
				}
			}
		}()
	}
	if d.HasChange("image_sprite_definition") {
		changeFlag = true
		idUint, _ := strconv.ParseUint(d.Get("image_sprite_definition").(string), 0, 64)
		request.ImageSpriteDefinition = &idUint
	}
	if d.HasChange("resolution_names") {
		changeFlag = true
		v := d.Get("resolution_names")
		resolutionNames := make([]*vod.ResolutionNameInfo, 0, len(v.([]interface{})))
		for _, item := range v.([]interface{}) {
			itemV := item.(map[string]interface{})
			resolutionNames = append(resolutionNames, &vod.ResolutionNameInfo{
				MinEdgeLength: helper.IntUint64(itemV["min_edge_length"].(int)),
				Name:          helper.String(itemV["name"].(string)),
			})
		}
		request.ResolutionNames = resolutionNames
	}
	if d.HasChange("domain") {
		changeFlag = true
		request.Domain = helper.String(d.Get("domain").(string))
	}
	if d.HasChange("scheme") {
		changeFlag = true
		request.Scheme = helper.String(d.Get("scheme").(string))
	}
	if d.HasChange("comment") {
		changeFlag = true
		request.Comment = helper.String(d.Get("comment").(string))
	}
	if d.HasChange("sub_app_id") {
		changeFlag = true
		request.SubAppId = helper.IntUint64(d.Get("sub_app_id").(int))
	}

	if changeFlag {
		var err error
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			_, err = meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifySuperPlayerConfig(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}

		return resourceTencentCloudVodSuperPlayerConfigRead(d, meta)
	}

	return nil
}

func resourceTencentCloudVodSuperPlayerConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_super_player_config.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	id := d.Id()
	vodService := VodService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	if err := vodService.DeleteSuperPlayerConfig(ctx, id, uint64(d.Get("sub_app_id").(int))); err != nil {
		return err
	}

	return nil
}
