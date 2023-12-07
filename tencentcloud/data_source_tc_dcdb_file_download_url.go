package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDcdbFileDownloadUrl() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbFileDownloadUrlRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"shard_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance Shard ID.",
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

func dataSourceTencentCloudDcdbFileDownloadUrlRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dcdb_file_download_url.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var (
		preSignedUrl *string
		instanceId   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("shard_id"); ok {
		paramMap["ShardId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_path"); ok {
		paramMap["FilePath"] = helper.String(v.(string))
	}

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDcdbFileDownloadUrlByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		preSignedUrl = result
		return nil
	})
	if err != nil {
		return err
	}

	if preSignedUrl != nil {
		_ = d.Set("pre_signed_url", preSignedUrl)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), preSignedUrl); e != nil {
			return e
		}
	}
	return nil
}
