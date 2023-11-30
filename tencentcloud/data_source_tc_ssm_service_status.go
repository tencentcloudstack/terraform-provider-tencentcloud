package tencentcloud

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"
)

func dataSourceTencentCloudSsmServiceStatus() *schema.Resource {
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
	defer logElapsed("data_source.tencentcloud_ssm_service_status.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		service       = SsmService{client: meta.(*TencentCloudClient).apiV3Conn}
		serviceStatus *ssm.GetServiceStatusResponseParams
	)

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSsmServiceStatusByFilter(ctx)
		if e != nil {
			return retryError(e)
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
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
