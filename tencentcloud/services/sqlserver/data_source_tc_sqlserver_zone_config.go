package sqlserver

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTencentCloudSqlserverZoneConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentSqlserverZoneConfigRead,
		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			"zone_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of availability zones. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alphabet ID of availability zone.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number ID of availability zone.",
						},
						"specinfo_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of specinfo configurations for the specific availability zone. Each element contains the following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"spec_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Instance specification ID.",
									},
									"machine_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Model ID.",
									},
									"db_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database version information. Valid values: `2008R2 (SQL Server 2008 Enterprise)`, `2012SP3 (SQL Server 2012 Enterprise)`, `2016SP1 (SQL Server 2016 Enterprise)`, `201602 (SQL Server 2016 Standard)`, `2017 (SQL Server 2017 Enterprise)`.",
									},
									"db_version_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Version name corresponding to the `db_version` field.",
									},
									"memory": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Memory size in GB.",
									},
									"cpu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of CPU cores.",
									},
									"min_storage_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum disk size under this specification in GB.",
									},
									"max_storage_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum disk size under this specification in GB.",
									},
									"qps": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "QPS of this specification.",
									},
									"charge_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Billing mode under this specification. Valid values are `POSTPAID_BY_HOUR`, `PREPAID` and `ALL`. `ALL` means both POSTPAID_BY_HOUR and PREPAID.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentSqlserverZoneConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencent_sqlserver_zone_config.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sqlserverService := SqlserverService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	// get zoneinfo
	zoneInfoList, err := sqlserverService.DescribeZones(ctx)
	if err != nil {
		return fmt.Errorf("api[DescribeZones]fail, return %s", err.Error())
	}
	zoneSet := make(map[string]map[string]interface{})
	for _, zoneInfo := range zoneInfoList {
		zoneSetInfo := make(map[string]interface{}, 1)
		zoneSetInfo["id"] = zoneInfo.ZoneId
		zoneSet[*zoneInfo.Zone] = zoneSetInfo
	}

	var zoneList []interface{}
	for k, v := range zoneSet {
		var zoneListItem = make(map[string]interface{})
		zoneListItem["availability_zone"] = k
		zoneListItem["zone_id"] = v["id"]

		// get specinfo for each zone
		specinfoList, err := sqlserverService.DescribeProductConfig(ctx, k)
		if err != nil {
			return fmt.Errorf("api[DescribeProductConfig]fail, return %s", err.Error())
		}
		var specinfoConfigs []interface{}
		for _, specinfoItem := range specinfoList {
			var specinfoConfig = make(map[string]interface{})
			specinfoConfig["spec_id"] = specinfoItem.SpecId
			specinfoConfig["machine_type"] = specinfoItem.MachineType
			specinfoConfig["db_version"] = specinfoItem.Version
			specinfoConfig["db_version_name"] = specinfoItem.VersionName
			specinfoConfig["memory"] = specinfoItem.Memory
			specinfoConfig["cpu"] = specinfoItem.CPU
			specinfoConfig["min_storage_size"] = specinfoItem.MinStorage
			specinfoConfig["max_storage_size"] = specinfoItem.MaxStorage
			specinfoConfig["qps"] = specinfoItem.QPS
			specinfoConfig["charge_type"] = SQLSERVER_CHARGE_TYPE_NAME[*specinfoItem.PayModeStatus]

			specinfoConfigs = append(specinfoConfigs, specinfoConfig)
		}
		zoneListItem["specinfo_list"] = specinfoConfigs
		zoneList = append(zoneList, zoneListItem)
	}

	// set zone_list
	if err := d.Set("zone_list", zoneList); err != nil {
		return fmt.Errorf("[CRITAL]%s provider set zone_list fail, reason:%s\n ", logId, err.Error())
	}

	d.SetId("zone_config")

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), zoneList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
		}

	}
	return nil
}
