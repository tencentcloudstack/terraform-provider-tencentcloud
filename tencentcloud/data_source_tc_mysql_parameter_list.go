package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TencentCloudMysqlParameterDetail() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"parameter_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Parameter name.",
		},
		"parameter_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Parameter type.",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Parameter specification description.",
		},
		"current_value": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Current value.",
		},
		"default_value": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Default value.",
		},
		"enum_value": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Enumerated value.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"max": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Maximum value for the parameter.",
		},
		"min": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Minimum value for the parameter.",
		},
		"need_reboot": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Indicates whether reboot is needed to enable the new parameters.",
		},
	}
}

func dataSourceTencentCloudMysqlParameterList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMysqlParameterListRead,
		Schema: map[string]*schema.Schema{
			"mysql_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance ID.",
			},
			"engine_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"5.1", "5.5", "5.6", "5.7", "8.0"}),
				Description:  "The version number of the database engine to use. Supported versions include 5.5/5.6/5.7/8.0.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			"parameter_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of parameters. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: TencentCloudMysqlParameterDetail(),
				},
			},
		},
	}
}

func dataSourceTencentMysqlParameterListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mysql_parameter_list.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
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
	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("parameter_list", parameterList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set parameter list fail, reason:%s\n ", logId, err.Error())
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = writeToFile(output.(string), parameterList); err != nil {
			return err
		}
	}
	return nil
}
