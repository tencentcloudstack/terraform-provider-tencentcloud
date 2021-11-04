/*
Provide a resource to create a VOD super player config.

Example Usage

```hcl
resource "tencentcloud_vod_adaptive_dynamic_streaming_template" "foo" {
  format                          = "HLS"
  name                            = "tf-adaptive"
  drm_type                        = "SimpleAES"
  disable_higher_video_bitrate    = false
  disable_higher_video_resolution = false
  comment                         = "test"

  stream_info {
    video {
      codec               = "libx265"
      fps                 = 4
      bitrate             = 129
      resolution_adaptive = false
      width               = 128
      height              = 128
      fill_type           = "stretch"
    }
    audio {
      codec         = "libmp3lame"
      bitrate       = 129
      sample_rate   = 44100
      audio_channel = "dual"
    }
    remove_audio = false
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
    remove_audio = true
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
  resolution_adaptive = false
}

resource "tencentcloud_vod_super_player_config" "foo" {
  name                    = "tf-super-player"
  drm_switch              = true
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

VOD super player config can be imported using the name, e.g.

```
$ terraform import tencentcloud_vod_super_player_config.foo tf-super-player
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudVodSubApplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodSubApplicationCreate,
		Read:   resourceTencentCloudVodSubApplicationRead,
		Update: resourceTencentCloudVodSubApplicationUpdate,
		Delete: resourceTencentCloudVodSubApplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 40),
				Description:  "Sub application name, which can contain up to 64 letters, digits, underscores, and hyphens (such as test_ABC-123) and must be unique under a user.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sub application description.",
			},
		},
	}
}

func resourceTencentCloudVodSubApplicationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_sub_application.create")()

	var (
		logId      = getLogId(contextNil)
		request    = vod.NewCreateSubAppIdRequest()
		subAppId   *uint64
		subAppName *string
	)

	if v, ok := d.GetOk("name"); ok {
		subAppName = helper.String(v.(string))
		request.Name = subAppName
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := meta.(*TencentCloudClient).apiV3Conn.UseVodClient().CreateSubAppId(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), err.Error())
			return retryError(err)
		}
		subAppId = response.Response.SubAppId
		return nil
	}); err != nil {
		return err
	}

	d.SetId(*subAppName + FILED_SP + helper.UInt64ToStr(*subAppId))

	return resourceTencentCloudVodSubApplicationRead(d, meta)
}

func resourceTencentCloudVodSubApplicationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_sub_application.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId   = getLogId(contextNil)
		ctx     = context.WithValue(context.TODO(), logIdKey, logId)
		client  = meta.(*TencentCloudClient).apiV3Conn
		request = vod.NewDescribeSubAppIdsRequest()
	)

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("scf layer id is borken, id is %s", d.Id())
	}
	subAppName := idSplit[0]
	subAppId := idSplit[1]

	request.Name = &subAppName

	response, err := client.UseVodClient().DescribeSubAppIds(request)
	if err != nil {
		return err
	}
	infoSet := response.Response.SubAppIdInfoSet
	if len(infoSet) == 0 {
		d.SetId("")
		return nil
	}

	for(info in inf)
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

func resourceTencentCloudVodSubApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_sub_application.update")()

	var (
		logId      = getLogId(contextNil)
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

	return nil
}

func resourceTencentCloudVodSubApplicationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_vod_sub_application.delete")()

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
