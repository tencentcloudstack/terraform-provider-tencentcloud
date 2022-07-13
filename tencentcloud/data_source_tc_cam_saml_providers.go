/*
Use this data source to query detailed information of CAM SAML providers

Example Usage

```hcl
data "tencentcloud_cam_saml_providers" "foo" {
  name = "cam-test-provider"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCamSAMLProviders() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCamSAMLProvidersRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the CAM SAML provider to be queried.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the CAM SAML provider.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"provider_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of CAM SAML providers. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CAM SAML provider.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of CAM SAML provider.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the CAM SAML provider.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last modify time of the CAM SAML provider.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCamSAMLProvidersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cam_saml_providers.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	params := make(map[string]interface{})
	if v, ok := d.GetOk("name"); ok {
		params["name"] = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		params["description"] = v.(string)
	}

	camService := CamService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var providers []*cam.SAMLProviderInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := camService.DescribeSAMLProvidersByFilter(ctx, params)
		if e != nil {
			return retryError(e)
		}
		providers = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CAM groups failed, reason:%s\n", logId, err.Error())
		return err
	}
	providerList := make([]map[string]interface{}, 0, len(providers))
	ids := make([]string, 0, len(providers))
	for _, provider := range providers {
		mapping := map[string]interface{}{
			"name":        *provider.Name,
			"description": *provider.Description,
			"create_time": *provider.CreateTime,
			"modify_time": *provider.ModifyTime,
		}
		providerList = append(providerList, mapping)
		ids = append(ids, *provider.Name)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("provider_list", providerList); e != nil {
		log.Printf("[CRITAL]%s provider set provider list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), providerList); e != nil {
			return e
		}
	}

	return nil
}
