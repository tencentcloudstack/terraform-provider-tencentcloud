package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudGaapCheckProxyCreate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapCheckProxyCreateRead,
		Schema: map[string]*schema.Schema{
			"access_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The access (acceleration) area of the proxy. The value can be obtained through the interface DescribeAccessRegionsByDestRegion.",
			},

			"real_server_region": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The origin area of the proxy. The value can be obtained through the interface DescribeDestRegions.",
			},

			"bandwidth": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The upper limit of proxy bandwidth, in Mbps.",
			},

			"concurrent": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The upper limit of chanproxynel concurrency, representing the number of simultaneous online connections, in tens of thousands.",
			},

			"group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "If creating a proxy under a proxy group, you need to fill in the ID of the proxy group.",
			},

			"ip_address_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "IP version, can be taken as IPv4 or IPv6, with a default value of IPv4.",
			},

			"network_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Network type, can take values &amp;#39;normal&amp;#39;, &amp;#39;cn2&amp;#39;, default value normal.",
			},

			"package_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Channel package type. Thunder represents the standard proxy group, Accelerator represents the game accelerator proxy, and CrossBorder represents the cross-border proxy.",
			},

			"check_flag": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Query whether the proxy with the given configuration can be created, 1 can be created, 0 cannot be created.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudGaapCheckProxyCreateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_check_proxy_create.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("access_region"); ok {
		paramMap["AccessRegion"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("real_server_region"); ok {
		paramMap["RealServerRegion"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("bandwidth"); v != nil {
		paramMap["Bandwidth"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("concurrent"); v != nil {
		paramMap["Concurrent"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ip_address_version"); ok {
		paramMap["IPAddressVersion"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("network_type"); ok {
		paramMap["NetworkType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("package_type"); ok {
		paramMap["PackageType"] = helper.String(v.(string))
	}

	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}

	var checkFlag *uint64
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapCheckProxyCreate(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		checkFlag = result
		return nil
	})
	if err != nil {
		return err
	}

	result := map[string]interface{}{}

	if checkFlag != nil {
		_ = d.Set("check_flag", *checkFlag)
		result["check_flag"] = *checkFlag
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
