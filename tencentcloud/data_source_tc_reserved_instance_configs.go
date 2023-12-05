package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudReservedInstanceConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudReservedInstanceConfigsRead,

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The available zone that the reserved instance locates at.",
			},
			"duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{31536000, 94608000}),
				Description:  "Validity period of the reserved instance. Valid values are `31536000`(1 year) and `94608000`(3 years).",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of reserved instance.",
			},
			"offering_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by Payment Type. Such as All Upfront.",
			},
			"product_description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by the Platform Description (that is, operating system) for Reserved Instance billing. Shaped like: linux.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"config_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of reserved instance configuration. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Configuration ID of the purchasable reserved instance.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone of the purchasable reserved instance.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type of the reserved instance.",
						},
						"duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Validity period of the reserved instance.",
						},
						"price": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Purchase price of the reserved instance.",
						},
						"currency_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Settlement currency of the reserved instance, which is a standard currency code as listed in ISO 4217.",
						},
						"platform": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Platform of the reserved instance.",
						},
						"offering_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OfferingType of the reserved instance.",
						},
						"usage_price": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "UsagePrice of the reserved instance.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudReservedInstanceConfigsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_reserved_instance_configs.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	filter := make(map[string]string)
	if v, ok := d.GetOk("availability_zone"); ok {
		filter["zone"] = v.(string)
	}
	if v, ok := d.GetOk("duration"); ok {
		filter["duration"] = strconv.Itoa(v.(int))
	}
	if v, ok := d.GetOk("instance_type"); ok {
		filter["instance-type"] = v.(string)
	}
	if v, ok := d.GetOk("offering_type"); ok {
		filter["offering-type"] = v.(string)
	}
	if v, ok := d.GetOk("product_description"); ok {
		filter["product-description"] = v.(string)
	}

	var configs []*cvm.ReservedInstancesOffering
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		configs, errRet = cvmService.DescribeReservedInstanceConfigs(ctx, filter)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	configList := make([]map[string]interface{}, 0, len(configs))
	ids := make([]string, 0, len(configs))
	for _, config := range configs {
		mapping := map[string]interface{}{
			"config_id":         config.ReservedInstancesOfferingId,
			"availability_zone": config.Zone,
			"instance_type":     config.InstanceType,
			"duration":          config.Duration,
			"price":             config.FixedPrice,
			"currency_code":     config.CurrencyCode,
			"platform":          config.ProductDescription,
			"offering_type":     config.OfferingType,
			"usage_price":       config.UsagePrice,
		}
		configList = append(configList, mapping)
		ids = append(ids, *config.ReservedInstancesOfferingId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("config_list", configList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set config list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), configList); err != nil {
			return err
		}
	}
	return nil
}
