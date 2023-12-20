package cdb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func DataSourceTencentCloudMysqlProxyCustom() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlProxyCustomRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instanced id.",
			},

			"custom_conf": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "proxy configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "equipment.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "type.",
						},
						"device_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Equipment type.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "number of cores.",
						},
					},
				},
			},

			"weight_rule": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "weight limit.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"less_than": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "division ceiling.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "weight limit.",
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

func dataSourceTencentCloudMysqlProxyCustomRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_proxy_custom.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var proxyCustom *cdb.DescribeProxyCustomConfResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlProxyCustomById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		proxyCustom = result
		return nil
	})
	if err != nil {
		return err
	}

	if proxyCustom.CustomConf != nil {
		customConfigMap := map[string]interface{}{}

		if proxyCustom.CustomConf.Device != nil {
			customConfigMap["device"] = proxyCustom.CustomConf.Device
		}

		if proxyCustom.CustomConf.Type != nil {
			customConfigMap["type"] = proxyCustom.CustomConf.Type
		}

		if proxyCustom.CustomConf.DeviceType != nil {
			customConfigMap["device_type"] = proxyCustom.CustomConf.DeviceType
		}

		if proxyCustom.CustomConf.Memory != nil {
			customConfigMap["memory"] = proxyCustom.CustomConf.Memory
		}

		if proxyCustom.CustomConf.Cpu != nil {
			customConfigMap["cpu"] = proxyCustom.CustomConf.Cpu
		}

		_ = d.Set("custom_conf", customConfigMap)
	}

	if proxyCustom.WeightRule != nil {
		ruleMap := map[string]interface{}{}

		if proxyCustom.WeightRule.LessThan != nil {
			ruleMap["less_than"] = proxyCustom.WeightRule.LessThan
		}

		if proxyCustom.WeightRule.Weight != nil {
			ruleMap["weight"] = proxyCustom.WeightRule.Weight
		}

		_ = d.Set("weight_rule", ruleMap)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
