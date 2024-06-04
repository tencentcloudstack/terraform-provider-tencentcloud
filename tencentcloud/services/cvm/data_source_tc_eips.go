package cvm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudEips() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEipsRead,
		Schema: map[string]*schema.Schema{
			"eip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the EIP to be queried.",
			},

			"eip_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of EIP. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the EIP.",
						},
						"eip_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the EIP.",
						},
						"eip_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the EIP.",
						},
						"eip_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the EIP.",
						},
						"eni_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The eni id to bind with the EIP.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance id to bind with the EIP.",
						},
						"public_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The elastic ip address.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The EIP current status.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of the EIP.",
						},
					},
				},
			},

			"eip_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the EIP to be queried.",
			},

			"public_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The elastic ip address.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The tags of EIP.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudEipsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_eips.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	var filtersList []*vpc.Filter
	filtersMap := map[string]*vpc.Filter{}
	filter := vpc.Filter{}
	name := "address-id"
	filter.Name = &name
	if v, ok := d.GetOk("eip_id"); ok {
		filter.Values = []*string{helper.String(v.(string))}
	}
	filtersMap["Temp0"] = &filter
	if v, ok := filtersMap["Temp0"]; ok && len(v.Values) > 0 {
		filtersList = append(filtersList, v)
	}
	filter2 := vpc.Filter{}
	name2 := "address-name"
	filter2.Name = &name2
	if v, ok := d.GetOk("eip_name"); ok {
		filter2.Values = []*string{helper.String(v.(string))}
	}
	filtersMap["Temp1"] = &filter2
	if v, ok := filtersMap["Temp1"]; ok && len(v.Values) > 0 {
		filtersList = append(filtersList, v)
	}
	filter3 := vpc.Filter{}
	name3 := "public-ip"
	filter3.Name = &name3
	if v, ok := d.GetOk("public_ip"); ok {
		filter3.Values = []*string{helper.String(v.(string))}
	}
	filtersMap["Temp2"] = &filter3
	if v, ok := filtersMap["Temp2"]; ok && len(v.Values) > 0 {
		filtersList = append(filtersList, v)
	}
	paramMap["Filters"] = filtersList

	var respData *vpc.DescribeAddressesResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEipsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	if err := dataSourceTencentCloudEipsReadPostHandleResponse0(ctx, paramMap, respData); err != nil {
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataSourceTencentCloudEipsReadOutputContent(ctx)); e != nil {
			return e
		}
	}

	return nil
}
