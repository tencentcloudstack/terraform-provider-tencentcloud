package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudVpcResourceDashboard() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudVpcResourceDashboardRead,
		Schema: map[string]*schema.Schema{
			"vpc_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Vpc instance ID, e.g. vpc-f1xjkw1b.",
			},

			"resource_dashboard_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of resource objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC instance ID, such as `vpc-bq4bzxpj`.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet instance ID, such as subnet-bthucmmy.",
						},
						"classic_link": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Classic link.",
						},
						"dcg": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Direct Connect gateway.",
						},
						"pcx": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Peering connection.",
						},
						"ip": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of used IPs except for CVM IP, EIP and network probe IP. The three IP types will be independently counted.",
						},
						"nat": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "NAT gateway.",
						},
						"vpngw": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "VPN gateway.",
						},
						"flow_log": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Flow log.",
						},
						"network_detect": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Network probing.",
						},
						"network_acl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Network ACL.",
						},
						"cvm": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cloud Virtual Machine.",
						},
						"lb": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Load balancer.",
						},
						"cdb": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Relational database.",
						},
						"cmem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TencentDB for Memcached.",
						},
						"cts_db": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cloud time series database.",
						},
						"maria_db": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TencentDB for MariaDB (TDSQL).",
						},
						"sql_server": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TencentDB for SQL Server.",
						},
						"postgres": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TencentDB for PostgreSQL.",
						},
						"nas": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Network attached storage.",
						},
						"greenplumn": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Snova data warehouse.",
						},
						"ckafka": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cloud Kafka (CKafka).",
						},
						"grocery": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Grocery.",
						},
						"hsm": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Data encryption service.",
						},
						"tcaplus": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Game storage - Tcaplus.",
						},
						"cnas": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cnas.",
						},
						"ti_db": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "HTAP database - TiDB.",
						},
						"emr": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "EMR cluster.",
						},
						"seal": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "SEAL.",
						},
						"cfs": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cloud file storage - CFS.",
						},
						"oracle": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Oracle.",
						},
						"elastic_search": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ElasticSearch Service.",
						},
						"t_baas": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Blockchain service.",
						},
						"itop": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Itop.",
						},
						"db_audit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cloud database audit.",
						},
						"cynos_db_postgres": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Enterprise TencentDB - CynosDB for Postgres.",
						},
						"redis": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TencentDB for Redis.",
						},
						"mongo_db": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "TencentDB for MongoDB.",
						},
						"dcdb": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "A distributed cloud database - TencentDB for TDSQL.",
						},
						"cynos_db_mysql": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "An enterprise-grade TencentDB - CynosDB for MySQL.",
						},
						"subnet": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Subnets.",
						},
						"route_table": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Route table.",
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

func dataSourceTencentCloudVpcResourceDashboardRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_vpc_resource_dashboard.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("vpc_ids"); ok {
		vpcIdsSet := v.(*schema.Set).List()
		paramMap["VpcIds"] = helper.InterfacesStringsPoint(vpcIdsSet)
	}

	service := VpcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var resourceDashboardSet []*vpc.ResourceDashboard

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeVpcResourceDashboard(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		resourceDashboardSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(resourceDashboardSet))
	tmpList := make([]map[string]interface{}, 0, len(resourceDashboardSet))

	if resourceDashboardSet != nil {
		for _, resourceDashboard := range resourceDashboardSet {
			resourceDashboardMap := map[string]interface{}{}

			if resourceDashboard.VpcId != nil {
				resourceDashboardMap["vpc_id"] = resourceDashboard.VpcId
			}

			if resourceDashboard.SubnetId != nil {
				resourceDashboardMap["subnet_id"] = resourceDashboard.SubnetId
			}

			if resourceDashboard.Classiclink != nil {
				resourceDashboardMap["classic_link"] = resourceDashboard.Classiclink
			}

			if resourceDashboard.Dcg != nil {
				resourceDashboardMap["dcg"] = resourceDashboard.Dcg
			}

			if resourceDashboard.Pcx != nil {
				resourceDashboardMap["pcx"] = resourceDashboard.Pcx
			}

			if resourceDashboard.Ip != nil {
				resourceDashboardMap["ip"] = resourceDashboard.Ip
			}

			if resourceDashboard.Nat != nil {
				resourceDashboardMap["nat"] = resourceDashboard.Nat
			}

			if resourceDashboard.Vpngw != nil {
				resourceDashboardMap["vpngw"] = resourceDashboard.Vpngw
			}

			if resourceDashboard.FlowLog != nil {
				resourceDashboardMap["flow_log"] = resourceDashboard.FlowLog
			}

			if resourceDashboard.NetworkDetect != nil {
				resourceDashboardMap["network_detect"] = resourceDashboard.NetworkDetect
			}

			if resourceDashboard.NetworkACL != nil {
				resourceDashboardMap["network_acl"] = resourceDashboard.NetworkACL
			}

			if resourceDashboard.CVM != nil {
				resourceDashboardMap["cvm"] = resourceDashboard.CVM
			}

			if resourceDashboard.LB != nil {
				resourceDashboardMap["lb"] = resourceDashboard.LB
			}

			if resourceDashboard.CDB != nil {
				resourceDashboardMap["cdb"] = resourceDashboard.CDB
			}

			if resourceDashboard.Cmem != nil {
				resourceDashboardMap["cmem"] = resourceDashboard.Cmem
			}

			if resourceDashboard.CTSDB != nil {
				resourceDashboardMap["cts_db"] = resourceDashboard.CTSDB
			}

			if resourceDashboard.MariaDB != nil {
				resourceDashboardMap["maria_db"] = resourceDashboard.MariaDB
			}

			if resourceDashboard.SQLServer != nil {
				resourceDashboardMap["sql_server"] = resourceDashboard.SQLServer
			}

			if resourceDashboard.Postgres != nil {
				resourceDashboardMap["postgres"] = resourceDashboard.Postgres
			}

			if resourceDashboard.NAS != nil {
				resourceDashboardMap["nas"] = resourceDashboard.NAS
			}

			if resourceDashboard.Greenplumn != nil {
				resourceDashboardMap["greenplumn"] = resourceDashboard.Greenplumn
			}

			if resourceDashboard.Ckafka != nil {
				resourceDashboardMap["ckafka"] = resourceDashboard.Ckafka
			}

			if resourceDashboard.Grocery != nil {
				resourceDashboardMap["grocery"] = resourceDashboard.Grocery
			}

			if resourceDashboard.HSM != nil {
				resourceDashboardMap["hsm"] = resourceDashboard.HSM
			}

			if resourceDashboard.Tcaplus != nil {
				resourceDashboardMap["tcaplus"] = resourceDashboard.Tcaplus
			}

			if resourceDashboard.Cnas != nil {
				resourceDashboardMap["cnas"] = resourceDashboard.Cnas
			}

			if resourceDashboard.TiDB != nil {
				resourceDashboardMap["ti_db"] = resourceDashboard.TiDB
			}

			if resourceDashboard.Emr != nil {
				resourceDashboardMap["emr"] = resourceDashboard.Emr
			}

			if resourceDashboard.SEAL != nil {
				resourceDashboardMap["seal"] = resourceDashboard.SEAL
			}

			if resourceDashboard.CFS != nil {
				resourceDashboardMap["cfs"] = resourceDashboard.CFS
			}

			if resourceDashboard.Oracle != nil {
				resourceDashboardMap["oracle"] = resourceDashboard.Oracle
			}

			if resourceDashboard.ElasticSearch != nil {
				resourceDashboardMap["elastic_search"] = resourceDashboard.ElasticSearch
			}

			if resourceDashboard.TBaaS != nil {
				resourceDashboardMap["t_baas"] = resourceDashboard.TBaaS
			}

			if resourceDashboard.Itop != nil {
				resourceDashboardMap["itop"] = resourceDashboard.Itop
			}

			if resourceDashboard.DBAudit != nil {
				resourceDashboardMap["db_audit"] = resourceDashboard.DBAudit
			}

			if resourceDashboard.CynosDBPostgres != nil {
				resourceDashboardMap["cynos_db_postgres"] = resourceDashboard.CynosDBPostgres
			}

			if resourceDashboard.Redis != nil {
				resourceDashboardMap["redis"] = resourceDashboard.Redis
			}

			if resourceDashboard.MongoDB != nil {
				resourceDashboardMap["mongo_db"] = resourceDashboard.MongoDB
			}

			if resourceDashboard.DCDB != nil {
				resourceDashboardMap["dcdb"] = resourceDashboard.DCDB
			}

			if resourceDashboard.CynosDBMySQL != nil {
				resourceDashboardMap["cynos_db_mysql"] = resourceDashboard.CynosDBMySQL
			}

			if resourceDashboard.Subnet != nil {
				resourceDashboardMap["subnet"] = resourceDashboard.Subnet
			}

			if resourceDashboard.RouteTable != nil {
				resourceDashboardMap["route_table"] = resourceDashboard.RouteTable
			}

			ids = append(ids, *resourceDashboard.VpcId)
			tmpList = append(tmpList, resourceDashboardMap)
		}

		_ = d.Set("resource_dashboard_set", tmpList)
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
