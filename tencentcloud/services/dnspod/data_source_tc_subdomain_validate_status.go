package dnspod

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSubdomainValidateStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSubdomainValidateStatusRead,
		Schema: map[string]*schema.Schema{
			"domain_zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone domain for which to view the verification status of TXT records.",
			},

			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status. 0: not ready; 1: ready.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSubdomainValidateStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_subdomain_validate_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		domainZone string
	)
	if v, ok := d.GetOk("domain_zone"); ok {
		domainZone = v.(string)
	}
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain_zone"); ok {
		paramMap["DomainZone"] = helper.String(v.(string))
	}

	var status int
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSubdomainValidateStatusByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		status = result
		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(domainZone)
	_ = d.Set("status", status)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
