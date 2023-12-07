package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	antiddos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/antiddos/v20200309"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAntiddosBasicDeviceStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAntiddosBasicDeviceStatusRead,
		Schema: map[string]*schema.Schema{
			"ip_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Ip resource list.",
			},

			"id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Named resource transfer ID.",
			},

			"filter_region": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Region Id.",
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Return resources and status, status code: 1- Blocking status 2- Normal status 3- Attack status.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Properties name.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Properties value.",
						},
					},
				},
			},

			"clb_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Properties name.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Properties value.",
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

func dataSourceTencentCloudAntiddosBasicDeviceStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_antiddos_basic_device_status.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("ip_list"); ok {
		ipListSet := v.(*schema.Set).List()
		paramMap["IpList"] = helper.InterfacesStringsPoint(ipListSet)
	}

	if v, ok := d.GetOk("id_list"); ok {
		idListSet := v.(*schema.Set).List()
		paramMap["IdList"] = helper.InterfacesStringsPoint(idListSet)
	}

	if v, ok := d.GetOkExists("filter_region"); ok {
		paramMap["FilterRegion"] = helper.IntUint64(v.(int))
	}

	service := AntiddosService{client: meta.(*TencentCloudClient).apiV3Conn}

	var basicDeviceStatus *antiddos.DescribeBasicDeviceStatusResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeAntiddosBasicDeviceStatusByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		basicDeviceStatus = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0)

	if basicDeviceStatus.Data != nil {
		data := basicDeviceStatus.Data
		dataTmpList := make([]map[string]interface{}, 0, len(data))

		for _, keyValue := range data {
			keyValueMap := map[string]interface{}{}
			if keyValue.Key != nil {
				keyValueMap["key"] = keyValue.Key
			}
			if keyValue.Value != nil {
				keyValueMap["value"] = keyValue.Value
			}
			dataTmpList = append(dataTmpList, keyValueMap)
		}
		tmpList = append(tmpList, dataTmpList...)
		_ = d.Set("data", dataTmpList)
	}

	if basicDeviceStatus.CLBData != nil {
		clbData := basicDeviceStatus.CLBData
		clbDataTmpList := make([]map[string]interface{}, 0, len(clbData))

		for _, keyValue := range clbData {
			keyValueMap := map[string]interface{}{}
			if keyValue.Key != nil {
				keyValueMap["key"] = keyValue.Key
			}
			if keyValue.Value != nil {
				keyValueMap["value"] = keyValue.Value
			}
			clbDataTmpList = append(clbDataTmpList, keyValueMap)
		}
		tmpList = append(tmpList, clbDataTmpList...)
		_ = d.Set("clb_data", clbDataTmpList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
