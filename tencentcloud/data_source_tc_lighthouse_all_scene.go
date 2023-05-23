/*
Use this data source to query detailed information of all region lighthouse scene

Example Usage

```hcl
data "tencentcloud_lighthouse_all_scene" "scene" {
  offset = 0
  limit = 20
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLighthouseAllScene() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseAllSceneRead,
		Schema: map[string]*schema.Schema{
			"scene_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of scene IDs.",
			},

			"offset": {
				Optional:    true,
				Default:     0,
				Type:        schema.TypeInt,
				Description: "Offset. Default value is 0.",
			},

			"limit": {
				Optional:    true,
				Default:     20,
				Type:        schema.TypeInt,
				Description: "Number of returned results. Default value is 20. Maximum value is 100.",
			},

			"scene_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of scene info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scene_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use scene Id.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use the scene presentation name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Use scene description.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudLighthouseAllSceneRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_lighthouse_scene.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("scene_ids"); ok {
		sceneIdsSet := v.(*schema.Set).List()
		sceneIds := make([]string, 0)
		for _, sceneId := range sceneIdsSet {
			sceneIds = append(sceneIds, sceneId.(string))
		}
		paramMap["scene_ids"] = sceneIds
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["offset"] = v.(int)
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["limit"] = v.(int)
	}

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var sceneSet []*lighthouse.SceneInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseAllSceneByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		sceneSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(sceneSet))
	tmpList := make([]map[string]interface{}, 0, len(sceneSet))
	for _, scene := range sceneSet {
		ids = append(ids, *scene.SceneId)
		tmpList = append(tmpList, map[string]interface{}{
			"scene_id":     *scene.SceneId,
			"display_name": *scene.DisplayName,
			"description":  *scene.Description,
		})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("scene_set", tmpList)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
