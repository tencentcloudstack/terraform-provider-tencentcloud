/*
Use this data source to query detailed information of mariadb file_download_url

Example Usage

```hcl
data "tencentcloud_mariadb_file_download_url" "file_download_url" {
  instance_id = ""
  file_path = ""
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

func dataSourceTencentCloudMariadbFileDownloadUrl() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbFileDownloadUrlRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"file_path": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Unsigned file path.",
			},

			"pre_signed_url": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Signed download URL.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMariadbFileDownloadUrlRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mariadb_file_download_url.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_path"); ok {
		paramMap["FilePath"] = helper.String(v.(string))
	}

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbFileDownloadUrlByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		preSignedUrl = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(preSignedUrl))
	if preSignedUrl != nil {
		_ = d.Set("pre_signed_url", preSignedUrl)
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
