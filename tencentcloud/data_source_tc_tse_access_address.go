/*
Use this data source to query detailed information of tse access_address

Example Usage

```hcl
data "tencentcloud_tse_access_address" "access_address" {
  instance_id = "ins-xxxxxx"
  vpc_id = "vpc-xxxxxx"
  subnet_id = "subnet-xxxxxx"
  workload = "pushgateway"
  engine_region = "ap-guangzhou"
                  tags = {
    "createdBy" = "terraform"
  }
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

func dataSourceTencentCloudTseAccessAddress() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseAccessAddressRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Engine instance Id.",
			},

			"vpc_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "VPC ID.",
			},

			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Subnet ID.",
			},

			"workload": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Name of other engine components（pushgateway、polaris-limiter）.",
			},

			"engine_region": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Deploy region.",
			},

			"intranet_address": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Intranet access address.",
			},

			"internet_address": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Public access address.",
			},

			"env_address_infos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Apollo Multi-environment public ip address.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"env_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Env name.",
						},
						"enable_config_internet": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable the config public network.",
						},
						"config_internet_service_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Config public network ip.",
						},
						"config_intranet_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Config Intranet access addressNote: This field may return null, indicating that a valid value is not available.",
						},
						"enable_config_intranet": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to enable the config Intranet clbNote: This field may return null, indicating that a valid value is not available.",
						},
						"internet_band_width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Client public network bandwidthNote: This field may return null, indicating that a valid value is not available.",
						},
					},
				},
			},

			"console_internet_address": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Console public network access addressNote: This field may return null, indicating that a valid value is not available.",
			},

			"console_intranet_address": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Console Intranet access addressNote: This field may return null, indicating that a valid value is not available.",
			},

			"internet_band_width": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Client public network bandwidthNote: This field may return null, indicating that a valid value is not available.",
			},

			"console_internet_band_width": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Console public network bandwidthNote: This field may return null, indicating that a valid value is not available.",
			},

			"limiter_address_infos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Access IP address of the Polaris traffic limiting server nodeNote: This field may return null, indicating that a valid value is not available.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"intranet_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC access IP address listNote: This field may return null, indicating that a valid value is not available.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTseAccessAddressRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tse_access_address.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		paramMap["VpcId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		paramMap["SubnetId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("workload"); ok {
		paramMap["Workload"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("engine_region"); ok {
		paramMap["EngineRegion"] = helper.String(v.(string))
	}

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTseAccessAddressByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		intranetAddress = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(intranetAddress))
	if intranetAddress != nil {
		_ = d.Set("intranet_address", intranetAddress)
	}

	if internetAddress != nil {
		_ = d.Set("internet_address", internetAddress)
	}

	if envAddressInfos != nil {
		for _, envAddressInfo := range envAddressInfos {
			envAddressInfoMap := map[string]interface{}{}

			if envAddressInfo.EnvName != nil {
				envAddressInfoMap["env_name"] = envAddressInfo.EnvName
			}

			if envAddressInfo.EnableConfigInternet != nil {
				envAddressInfoMap["enable_config_internet"] = envAddressInfo.EnableConfigInternet
			}

			if envAddressInfo.ConfigInternetServiceIp != nil {
				envAddressInfoMap["config_internet_service_ip"] = envAddressInfo.ConfigInternetServiceIp
			}

			if envAddressInfo.ConfigIntranetAddress != nil {
				envAddressInfoMap["config_intranet_address"] = envAddressInfo.ConfigIntranetAddress
			}

			if envAddressInfo.EnableConfigIntranet != nil {
				envAddressInfoMap["enable_config_intranet"] = envAddressInfo.EnableConfigIntranet
			}

			if envAddressInfo.InternetBandWidth != nil {
				envAddressInfoMap["internet_band_width"] = envAddressInfo.InternetBandWidth
			}

			ids = append(ids, *envAddressInfo.InstanceId)
			tmpList = append(tmpList, envAddressInfoMap)
		}

		_ = d.Set("env_address_infos", tmpList)
	}

	if consoleInternetAddress != nil {
		_ = d.Set("console_internet_address", consoleInternetAddress)
	}

	if consoleIntranetAddress != nil {
		_ = d.Set("console_intranet_address", consoleIntranetAddress)
	}

	if internetBandWidth != nil {
		_ = d.Set("internet_band_width", internetBandWidth)
	}

	if consoleInternetBandWidth != nil {
		_ = d.Set("console_internet_band_width", consoleInternetBandWidth)
	}

	if limiterAddressInfos != nil {
		for _, polarisLimiterAddress := range limiterAddressInfos {
			polarisLimiterAddressMap := map[string]interface{}{}

			if polarisLimiterAddress.IntranetAddress != nil {
				polarisLimiterAddressMap["intranet_address"] = polarisLimiterAddress.IntranetAddress
			}

			ids = append(ids, *polarisLimiterAddress.InstanceId)
			tmpList = append(tmpList, polarisLimiterAddressMap)
		}

		_ = d.Set("limiter_address_infos", tmpList)
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
