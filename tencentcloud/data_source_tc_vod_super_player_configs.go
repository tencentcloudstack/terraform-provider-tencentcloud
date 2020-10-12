/*
Use this data source to query detailed information of Vod super player configs.

Example Usage

```hcl
data "tencentcloud_vod_super_player_configs" "foo" {
  type = "Custom"
  name = "tf-super-player"
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVodSuperPlayerConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVodSuperPlayerConfigsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of super player config.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Config type filter. Valid values: `Preset`: preset template; `Custom`: custom template.",
			},
			"sub_app_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Subapplication ID in VOD. If you need to access a resource in a subapplication, enter the subapplication ID in this field; otherwise, leave it empty.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"config_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of super player configs. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template type filter. Valid values: `Preset`: preset template; `Custom`: custom template.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Player configuration name, which can contain up to 64 letters, digits, underscores, and hyphens (such as test_ABC-123) and must be unique under a user.",
						},
						"drm_switch": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Switch of DRM-protected adaptive bitstream playback: `ON`: enabled, indicating to play back only output adaptive bitstreams protected by DRM; `OFF`: disabled, indicating to play back unencrypted output adaptive bitstreams.",
						},
						"adaptive_dynamic_streaming_definition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the unencrypted adaptive bitrate streaming template that allows output, which is required if `drm_switch` is `OFF`.",
						},
						"drm_streaming_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Content of the DRM-protected adaptive bitrate streaming template that allows output, which is required if `drm_switch` is `ON`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"simple_aes_definition": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the adaptive dynamic streaming template whose protection type is `SimpleAES`.",
									},
								},
							},
						},
						"image_sprite_definition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the image sprite template that allows output.",
						},
						"resolution_names": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Display name of player for substreams with different resolutions. If this parameter is left empty or an empty array, the default configuration will be used: `min_edge_length: 240, name: LD`; `min_edge_length: 480, name: SD`; `min_edge_length: 720, name: HD`; `min_edge_length: 1080, name: FHD`; `min_edge_length: 1440, name: 2K`; `min_edge_length: 2160, name: 4K`; `min_edge_length: 4320, name: 8K`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_edge_length": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Length of video short side in px.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Display name.",
									},
								},
							},
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name used for playback. If it is left empty or set to `Default`, the domain name configured in [Default Distribution Configuration](https://cloud.tencent.com/document/product/266/33373) will be used.",
						},
						"scheme": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheme used for playback. If it is left empty or set to `Default`, the scheme configured in [Default Distribution Configuration](https://cloud.tencent.com/document/product/266/33373) will be used. Other valid values: `HTTP`; `HTTPS`.",
						},
						"comment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template description.",
						},
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
				},
			},
		},
	}
}

func dataSourceTencentCloudVodSuperPlayerConfigsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vod_super_player_configs.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	filter := make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		filter["names"] = []string{v.(string)}
	}
	if v, ok := d.GetOk("type"); ok {
		filter["type"] = v.(string)
	}
	if v, ok := d.GetOk("sub_app_id"); ok {
		filter["sub_appid"] = v.(int)
	}

	vodService := VodService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	configs, err := vodService.DescribeSuperPlayerConfigsByFilter(ctx, filter)
	if err != nil {
		return err
	}

	configsList := make([]map[string]interface{}, 0, len(configs))
	ids := make([]string, 0, len(configs))
	for _, item := range configs {
		configsList = append(configsList, func() map[string]interface{} {
			mapping := map[string]interface{}{
				"type":        item.Type,
				"name":        item.Name,
				"drm_switch":  item.DrmSwitch,
				"domain":      item.Domain,
				"scheme":      item.Scheme,
				"comment":     item.Comment,
				"create_time": item.CreateTime,
				"update_time": item.UpdateTime,
			}
			// workaround for AdaptiveDynamicStreamingDefinition para cuz it's dirty data.
			if *item.DrmSwitch == "OFF" {
				mapping["adaptive_dynamic_streaming_definition"] = strconv.FormatUint(*item.AdaptiveDynamicStreamingDefinition, 10)
			}
			if item.DrmStreamingsInfo != nil && item.DrmStreamingsInfo.SimpleAesDefinition != nil {
				mapping["drm_streaming_info"] = []map[string]interface{}{
					{
						"simple_aes_definition": strconv.FormatUint(*item.DrmStreamingsInfo.SimpleAesDefinition, 10),
					},
				}
			}
			if item.ImageSpriteDefinition != nil {
				mapping["image_sprite_definition"] = strconv.FormatUint(*item.ImageSpriteDefinition, 10)
			}
			mapping["resolution_names"] = func() []map[string]interface{} {
				namesMap := make([]map[string]interface{}, 0, len(item.ResolutionNameSet))
				for _, v := range item.ResolutionNameSet {
					namesMap = append(namesMap, map[string]interface{}{
						"min_edge_length": v.MinEdgeLength,
						"name":            v.Name,
					})
				}
				return namesMap
			}()
			ids = append(ids, *item.Name)
			return mapping
		}())
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("config_list", configsList); e != nil {
		log.Printf("[CRITAL]%s provider set vod super player config list fail, reason:%s ", logId, e.Error())
	}

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), configsList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]", logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
