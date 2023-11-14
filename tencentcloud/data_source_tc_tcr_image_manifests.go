/*
Use this data source to query detailed information of tcr image_manifests

Example Usage

```hcl
data "tencentcloud_tcr_image_manifests" "image_manifests" {
  registry_id = "tcr-xxx"
  namespace_name = "ns"
  repository_name = "repo"
  image_version = "v1"
    }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTcrImageManifests() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcrImageManifestsRead,
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"namespace_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Namespace name.",
			},

			"repository_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Repository name.",
			},

			"image_version": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Image version name.",
			},

			"manifest": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Image manifests info.",
			},

			"config": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Image config info.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTcrImageManifestsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcr_image_manifests.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("registry_id"); ok {
		paramMap["RegistryId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_name"); ok {
		paramMap["NamespaceName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("repository_name"); ok {
		paramMap["RepositoryName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_version"); ok {
		paramMap["ImageVersion"] = helper.String(v.(string))
	}

	service := TcrService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcrImageManifestsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		manifest = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(manifest))
	if manifest != nil {
		_ = d.Set("manifest", manifest)
	}

	if config != nil {
		_ = d.Set("config", config)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
