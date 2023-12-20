package dbbrain

import (
	"context"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainSecurityAuditLogDownloadUrls() *schema.Resource {
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
				Description: "Service product type, supported values: `mysql` - ApsaraDB for MySQL.",
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
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_security_audit_log_download_urls.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var sagId string
	var asyncReqId int
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("sec_audit_group_id"); ok {
		paramMap["sec_audit_group_id"] = helper.String(v.(string))
		sagId = v.(string)
	}

	if v, _ := d.GetOk("async_request_id"); v != nil {
		paramMap["async_request_id"] = helper.IntUint64(v.(int))
		asyncReqId = v.(int)
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var urls []*string
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var e error
		urls, e = service.DescribeDbbrainSecurityAuditLogDownloadUrlsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if urls != nil {
		_ = d.Set("urls", urls)
	}

	d.SetId(strings.Join([]string{sagId, helper.Int64ToStr(int64(asyncReqId))}, tccommon.FILED_SP))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), urls); e != nil {
			return e
		}
	}
	return nil
}
