/*
Use this data source to query detailed information of cynosdb proxy_version

Example Usage

```hcl
data "tencentcloud_cynosdb_proxy_version" "proxy_version" {
  cluster_id = "cynosdbmysql-xxxxxxx"
  proxy_group_id = "无"
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

func dataSourceTencentCloudCynosdbProxyVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbProxyVersionRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"proxy_group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database Agent Group ID.",
			},

			"support_proxy_versions": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Supported Database Agent Version Collection Note: This field may return null, indicating that a valid value cannot be obtained.",
			},

			"current_proxy_version": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Current proxy version number note: This field may return null, indicating that a valid value cannot be obtained.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCynosdbProxyVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_proxy_version.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_group_id"); ok {
		paramMap["ProxyGroupId"] = helper.String(v.(string))
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var supportProxyVersions []*string

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbProxyVersionByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		supportProxyVersions = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(supportProxyVersions))
	if supportProxyVersions != nil {
		_ = d.Set("support_proxy_versions", supportProxyVersions)
	}

	if currentProxyVersion != nil {
		_ = d.Set("current_proxy_version", currentProxyVersion)
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
