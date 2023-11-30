package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpcLimits() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcLimitsRead,
		Schema: map[string]*schema.Schema{
			"limit_types": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Quota name. A maximum of 100 quota types can be queried each time.",
			},

			"vpc_limit_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "vpc limit.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limit_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "type of vpc limit.",
						},
						"limit_value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "value of vpc limit.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudVpcLimitsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_limits.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("limit_types"); ok {
		limitTypesSet := v.(*schema.Set).List()
		paramMap["LimitTypes"] = helper.InterfacesStringsPoint(limitTypesSet)
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var vpcLimitSet []*vpc.VpcLimit

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcLimitsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		vpcLimitSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(vpcLimitSet))
	tmpList := make([]map[string]interface{}, 0, len(vpcLimitSet))

	if vpcLimitSet != nil {
		for _, vpcLimit := range vpcLimitSet {
			vpcLimitMap := map[string]interface{}{}

			if vpcLimit.LimitType != nil {
				vpcLimitMap["limit_type"] = vpcLimit.LimitType
			}

			if vpcLimit.LimitValue != nil {
				vpcLimitMap["limit_value"] = vpcLimit.LimitValue
			}

			ids = append(ids, *vpcLimit.LimitType)
			tmpList = append(tmpList, vpcLimitMap)
		}

		_ = d.Set("vpc_limit_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
