package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTencentCloudEipNetworkAccountType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEipNetworkAccountTypeRead,
		Schema: map[string]*schema.Schema{
			"network_account_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The network type of the user account, STANDARD is a standard user, LEGACY is a traditional user.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudEipNetworkAccountTypeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_eip_network_account_type.read")()
	defer inconsistentCheck(d, meta)()

	var networkAccountType *string
	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEipNetworkAccountType(ctx)
		if e != nil {
			return retryError(e)
		}
		networkAccountType = result
		return nil
	})
	if err != nil {
		return err
	}

	if networkAccountType != nil {
		_ = d.Set("network_account_type", networkAccountType)
	}

	d.SetId(*networkAccountType)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), *networkAccountType); e != nil {
			return e
		}
	}
	return nil
}
