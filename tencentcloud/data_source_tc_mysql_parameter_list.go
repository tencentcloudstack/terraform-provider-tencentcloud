package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func TencentCloudMysqlParameterDetail() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"parameter_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"parameter_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"current_value": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"default_value": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"enum_value": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"max": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"min": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"need_reboot": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func dataSourceTencentCloudMysqlParameterList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMysqlParameterListRead,
		Schema: map[string]*schema.Schema{
			"mysql_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"5.1", "5.5", "5.6", "5.7"}),
			},
			"result_output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parameter_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: TencentCloudMysqlParameterDetail(),
				},
			},
		},
	}
}

func dataSourceTencentMysqlParameterListRead(d *schema.ResourceData, meta interface{}) error {
	defer LogElapsed("data_source.tencentcloud_mysql_parameter_list.read")()

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	mysqlService := MysqlService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var parameterDetails []*cdb.ParameterDetail
	var err error
	instanceIdString := ""
	engineVersionString := ""
	if instanceId, ok := d.GetOk("mysql_id"); ok {
		instanceIdString = instanceId.(string)
		parameterDetails, err = mysqlService.DescribeInstanceParameters(ctx, instanceIdString)
	} else if engineVersion, ok := d.GetOk("engine_version"); ok {
		engineVersionString = engineVersion.(string)
		parameterDetails, err = mysqlService.DescribeDefaultParameters(ctx, engineVersionString)
	} else {
		return fmt.Errorf("mysql_id and engine_version cannot be empty at the same time")
	}
	if err != nil {
		return fmt.Errorf("api[DescribeParameters]fail, return %s", err.Error())
	}

	parameterList := make([]map[string]interface{}, 0, len(parameterDetails))
	for _, item := range parameterDetails {
		mapping := map[string]interface{}{
			"parameter_name": *item.Name,
			"parameter_type": *item.ParamType,
			"description":    *item.Description,
			"current_value":  *item.CurrentValue,
			"default_value":  *item.Default,
			"enum_value":     item.EnumValue,
			"max":            *item.Max,
			"min":            *item.Min,
			"need_reboot":    *item.NeedReboot,
		}
		parameterList = append(parameterList, mapping)
	}
	ids := make([]string, 3)
	ids[0] = "DescribeParameter"
	ids[1] = instanceIdString
	ids[2] = engineVersionString
	d.SetId(dataResourceIdsHash(ids))
	err = d.Set("parameter_list", parameterList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set parameter list fail, reason:%s\n ", logId, err.Error())
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		writeToFile(output.(string), parameterList)
	}
	return nil
}
