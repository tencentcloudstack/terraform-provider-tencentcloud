package ssm

import (
	"context"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"
)

func DataSourceTencentCloudSsmServiceStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSsmServiceStatusRead,
		Schema: map[string]*schema.Schema{
			"service_enabled": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "True means the service has been activated, false means the service has not been activated yet.",
			},
			"invalid_type": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Service unavailability type: 0-Not purchased, 1-Normal, 2-Service suspended due to arrears, 3-Resource release.",
			},
			"access_key_escrow_enabled": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "True means that the user can already use the key safe hosting function, false means that the user cannot use the key safe hosting function temporarily.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSsmServiceStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ssm_service_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service       = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		serviceStatus *ssm.GetServiceStatusResponseParams
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSsmServiceStatusByFilter(ctx)
		if e != nil {
			return tccommon.RetryError(e)
		}

		serviceStatus = result
		return nil
	})

	if err != nil {
		return err
	}

	if serviceStatus.ServiceEnabled != nil {
		_ = d.Set("service_enabled", serviceStatus.ServiceEnabled)
	}

	if serviceStatus.InvalidType != nil {
		_ = d.Set("invalid_type", serviceStatus.InvalidType)
	}

	if serviceStatus.AccessKeyEscrowEnabled != nil {
		_ = d.Set("access_key_escrow_enabled", serviceStatus.AccessKeyEscrowEnabled)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
