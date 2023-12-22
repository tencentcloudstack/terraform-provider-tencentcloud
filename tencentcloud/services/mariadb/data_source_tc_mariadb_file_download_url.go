package mariadb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMariadbFileDownloadUrl() *schema.Resource {
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
	defer tccommon.LogElapsed("data_source.tencentcloud_mariadb_file_download_url.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service        = MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		responseParams *mariadb.DescribeFileDownloadUrlResponseParams
		instanceId     string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("file_path"); ok {
		paramMap["FilePath"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbFileDownloadUrlByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		responseParams = result
		return nil
	})

	if err != nil {
		return err
	}

	if responseParams.PreSignedUrl != nil {
		_ = d.Set("pre_signed_url", responseParams.PreSignedUrl)
	}

	d.SetId(instanceId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
