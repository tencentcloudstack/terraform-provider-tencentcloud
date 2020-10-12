/*
Provide a resource to create a vod super player config.

Example Usage

```hcl
resource "tencentcloud_vod_adaptive_dynamic_streaming_template" "foo" {
  format                          = "HLS"
  name                            = "tf-adaptive"
  drm_type                        = "SimpleAES"
  disable_higher_video_bitrate    = 0
  disable_higher_video_resolution = 0
  comment                         = "test"

  stream_info {
    video {
      codec   = "libx264"
      fps     = 3
      bitrate = 128
    }
    audio {
      codec       = "libfdk_aac"
      bitrate     = 128
      sample_rate = 32000
    }
    remove_audio = 1
  }
  stream_info {
    video {
      codec   = "libx264"
      fps     = 4
      bitrate = 256
    }
    audio {
      codec       = "libfdk_aac"
      bitrate     = 256
      sample_rate = 44100
    }
    remove_audio = 1
  }
}

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
  resolution_adaptive = "close"
}

resource "tencentcloud_vod_super_player_config" "foo" {
  name                    = "tf-super-player"
  drm_switch              = "ON"
  drm_streaming_info {
    simple_aes_definition = tencentcloud_vod_adaptive_dynamic_streaming_template.foo.id
  }
  image_sprite_definition = tencentcloud_vod_image_sprite_template.foo.id
  resolution_names {
    min_edge_length = 889
    name            = "test1"
  }
  resolution_names {
    min_edge_length = 890
    name            = "test2"
  }
  domain                  = "Default"
  scheme                  = "Default"
  comment                 = "test"
}
```

Import

Vod super player config can be imported using the name, e.g.

```
$ terraform import tencentcloud_vod_super_player_config.foo tf-super-player
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudVodSuperPlayerConfig() *schema.Resource {
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
				ValidateFunc: validateStringLengthInRange(1, 64),
				Description:  "Player configuration name, which can contain up to 64 letters, digits, underscores, and hyphens (such as test_ABC-123) and must be unique under a user.",
			},
			"drm_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "OFF",
				ValidateFunc: validateAllowedStringValue([]string{"ON", "OFF"}),
				Description:  "Switch of DRM-protected adaptive bitstream playback: `ON`: enabled, indicating to play back only output adaptive bitstreams protected by DRM; `OFF`: disabled, indicating to play back unencrypted output adaptive bitstreams. Default value: `OFF`.",
			},
			"adaptive_dynamic_streaming_definition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the unencrypted adaptive bitrate streaming template that allows output, which is required if `drm_switch` is `OFF`.",
			},
			"drm_streaming_info": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				Description: "Content of the DRM-protected adaptive bitrate streaming template that allows output, which is required if `drm_switch` is `ON`.",
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
				ValidateFunc: validateStringLengthInRange(1, 256),
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
	defer logElapsed("resource.tencentcloud_vod_super_player_config.create")()

	var (
		logId   = getLogId(contextNil)
		request = vod.NewCreateSuperPlayerConfigRequest()
	)

	request.Name = helper.String(d.Get("name").(string))
	request.DrmSwitch = helper.String(d.Get("drm_switch").(string))
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
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err = meta.(*TencentCloudClient).apiV3Conn.UseVodClient().CreateSuperPlayerConfig(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return retryError(err)
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
	defer logElapsed("resource.tencentcloud_vod_super_player_config.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		id         = d.Id()
		client     = meta.(*TencentCloudClient).apiV3Conn
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
	_ = d.Set("drm_switch", config.DrmSwitch)
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
	defer logElapsed("resource.tencentcloud_vod_super_player_config.update")()

	var (
		logId   = getLogId(contextNil)
		request = vod.NewModifySuperPlayerConfigRequest()
		id      = d.Id()
	)

	request.Name = &id
	if d.HasChange("drm_switch") {
		request.DrmSwitch = helper.String(d.Get("drm_switch").(string))
	}
	if d.HasChange("adaptive_dynamic_streaming_definition") {
		idUint, _ := strconv.ParseUint(d.Get("adaptive_dynamic_streaming_definition").(string), 0, 64)
		request.AdaptiveDynamicStreamingDefinition = &idUint
	}
	if d.HasChange("drm_streaming_info") {
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
		idUint, _ := strconv.ParseUint(d.Get("image_sprite_definition").(string), 0, 64)
		request.ImageSpriteDefinition = &idUint
	}
	if d.HasChange("resolution_names") {
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
		request.Domain = helper.String(d.Get("domain").(string))
	}
	if d.HasChange("scheme") {
		request.Scheme = helper.String(d.Get("scheme").(string))
	}
	if d.HasChange("comment") {
		request.Comment = helper.String(d.Get("comment").(string))
	}
	if d.HasChange("sub_app_id") {
		request.SubAppId = helper.IntUint64(d.Get("sub_app_id").(int))
	}

	var err error
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err = meta.(*TencentCloudClient).apiV3Conn.UseVodClient().ModifySuperPlayerConfig(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudVodSuperPlayerConfigRead(d, meta)
}

func resourceTencentCloudVodSuperPlayerConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_super_player_config.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	vodService := VodService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	if err := vodService.DeleteSuperPlayerConfig(ctx, id, uint64(d.Get("sub_app_id").(int))); err != nil {
		return err
	}

	return nil
}
