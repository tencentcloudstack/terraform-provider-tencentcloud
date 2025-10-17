package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataDataSources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataDataSourcesRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Data source name.",
			},

			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Data source display name.",
			},

			"type": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Data source type: enumeration values.\n\n- MYSQL\n- TENCENT_MYSQL\n- POSTGRE\n- ORACLE\n- SQLSERVER\n- FTP\n- HIVE\n- HUDI\n- HDFS\n- ICEBERG\n- KAFKA\n- HBASE\n- SPARK\n- VIRTUAL\n- TBASE\n- DB2\n- DM\n- GAUSSDB\n- GBASE\n- IMPALA\n- ES\n- TENCENT_ES\n- GREENPLUM\n- PHOENIX\n- SAP_HANA\n- SFTP\n- OCEANBASE\n- CLICKHOUSE\n- KUDU\n- VERTICA\n- REDIS\n- COS\n- DLC\n- DORIS\n- CKAFKA\n- S3\n- TDSQL\n- TDSQL_MYSQL\n- MONGODB\n- TENCENT_MONGODB\n- REST_API\n- SuperSQL\n- PRESTO\n- TiDB\n- StarRocks\n- Trino\n- Kyuubi\n- TCHOUSE_X\n- TCHOUSE_P\n- TCHOUSE_C\n- TCHOUSE_D\n- INFLUXDB\n- BIG_QUERY\n- SSH\n- BLOB.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"creator": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Creator.",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Data source list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Belonging project ID.",
						},
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Data source ID.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data source type: enumeration values.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data source name.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data source display name, for visual viewing.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data source description information.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Belonging project name.",
						},
						"create_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data source creator.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time.",
						},
						"modify_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modifier.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modification time.",
						},
						"prod_con_properties": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data source configuration information, stored in JSON KV format, varies by data source type.",
						},
						"dev_con_properties": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Same as params, contains data for development data source.",
						},
						"category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Data source category:\n\n- DB - custom source\n- CLUSTER - system source.",
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

func dataSourceTencentCloudWedataDataSourcesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_data_sources.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(nil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		projectId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		paramMap["Name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("display_name"); ok {
		paramMap["DisplayName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		typeList := []*string{}
		typeSet := v.(*schema.Set).List()
		for i := range typeSet {
			tmpType := typeSet[i].(string)
			typeList = append(typeList, helper.String(tmpType))
		}

		paramMap["Type"] = typeList
	}

	if v, ok := d.GetOk("creator"); ok {
		paramMap["Creator"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.DataSource
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataDataSourcesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	itemsList := make([]map[string]interface{}, 0, len(respData))
	for _, items := range respData {
		itemsMap := map[string]interface{}{}
		if items.ProjectId != nil {
			itemsMap["project_id"] = items.ProjectId
		}

		if items.Id != nil {
			itemsMap["id"] = items.Id
		}

		if items.Type != nil {
			itemsMap["type"] = items.Type
		}

		if items.Name != nil {
			itemsMap["name"] = items.Name
		}

		if items.DisplayName != nil {
			itemsMap["display_name"] = items.DisplayName
		}

		if items.Description != nil {
			itemsMap["description"] = items.Description
		}

		if items.ProjectName != nil {
			itemsMap["project_name"] = items.ProjectName
		}

		if items.CreateUser != nil {
			itemsMap["create_user"] = items.CreateUser
		}

		if items.CreateTime != nil {
			itemsMap["create_time"] = items.CreateTime
		}

		if items.ModifyUser != nil {
			itemsMap["modify_user"] = items.ModifyUser
		}

		if items.ModifyTime != nil {
			itemsMap["modify_time"] = items.ModifyTime
		}

		if items.ProdConProperties != nil {
			itemsMap["prod_con_properties"] = items.ProdConProperties
		}

		if items.DevConProperties != nil {
			itemsMap["dev_con_properties"] = items.DevConProperties
		}

		if items.Category != nil {
			itemsMap["category"] = items.Category
		}

		itemsList = append(itemsList, itemsMap)
	}

	_ = d.Set("items", itemsList)

	d.SetId(projectId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), itemsList); e != nil {
			return e
		}
	}

	return nil
}
