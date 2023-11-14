/*
Use this data source to query detailed information of dbbrain security_audit_log_download_urls

Example Usage

```hcl
data "tencentcloud_dbbrain_security_audit_log_download_urls" "security_audit_log_download_urls" {
  sec_audit_group_id = ""
  async_request_id =
  product = ""
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

func dataSourceTencentCloudDbbrainSecurityAuditLogDownloadUrls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainSecurityAuditLogDownloadUrlsRead,
		Schema: map[string]*schema.Schema{
			"sec_audit_group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Security audit group Id.",
			},

			"async_request_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Asynchronous task ID.",
			},

			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values： mysql - ApsaraDB for MySQL.",
			},

			"urls": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of COS links to export results. When the result set is large, it may be divided into multiple urls for download.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainSecurityAuditLogDownloadUrlsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_security_audit_log_download_urls.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("sec_audit_group_id"); ok {
		paramMap["SecAuditGroupId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("async_request_id"); v != nil {
		paramMap["AsyncRequestId"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	var urls []*string

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainSecurityAuditLogDownloadUrlsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		urls = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(urls))
	if urls != nil {
		_ = d.Set("urls", urls)
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
