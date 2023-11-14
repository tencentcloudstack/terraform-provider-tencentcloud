/*
Use this data source to query detailed information of lighthouse scene

Example Usage

```hcl
data "tencentcloud_lighthouse_scene" "scene" {
  scene_ids =
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

func dataSourceTencentCloudLighthouseScene() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseSceneRead,
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
				Type:        schema.TypeInt,
				Description: "Offset. Default value is 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Number of returned results. Default value is 20. Maximum value is 100.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudLighthouseSceneRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_lighthouse_scene.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("scene_ids"); ok {
		sceneIdsSet := v.(*schema.Set).List()
		paramMap["SceneIds"] = helper.InterfacesStringsPoint(sceneIdsSet)
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var sceneSet []*lighthouse.Scene

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseSceneByFilter(ctx, paramMap)
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

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
