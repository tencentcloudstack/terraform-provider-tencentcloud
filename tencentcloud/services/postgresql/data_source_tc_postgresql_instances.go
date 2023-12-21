package postgresql

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPostgresqlInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlInstanceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the postgresql instance to be query.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the postgresql instance to be query.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project ID of the postgresql instance to be query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of postgresql instances. Each element contains the following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the postgresql instance.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the postgresql instance.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pay type of the postgresql instance.",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Auto renew flag.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of the postgresql database engine.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of VPC.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of subnet.",
						},
						"storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Volume size(in GB).",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size(in GB).",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project id, default value is 0.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone.",
						},
						"root_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance root account name, default value is `root`.",
						},
						"public_access_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether to enable the access to an instance from public network or not.",
						},
						"public_access_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Host for public access.",
						},
						"public_access_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port for public access.",
						},
						"private_access_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP address for private access.",
						},
						"private_access_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port for private access.",
						},
						"charset": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Charset of the postgresql instance.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the postgresql instance.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The available tags within this postgresql.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudPostgresqlInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_instances.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	filter := make([]*postgresql.Filter, 0)
	if v, ok := d.GetOk("name"); ok {
		filter = append(filter, &postgresql.Filter{Name: helper.String("db-instance-name"), Values: []*string{helper.String(v.(string))}})
	}
	if v, ok := d.GetOk("id"); ok {
		filter = append(filter, &postgresql.Filter{Name: helper.String("db-instance-id"), Values: []*string{helper.String(v.(string))}})
	}
	if v, ok := d.GetOk("project_id"); ok {
		filter = append(filter, &postgresql.Filter{Name: helper.String("db-project-id"), Values: []*string{helper.String(v.(string))}})
	}

	instanceList, err := service.DescribePostgresqlInstances(ctx, filter)
	if err != nil {
		instanceList, err = service.DescribePostgresqlInstances(ctx, filter)
	}

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceList))
	list := make([]map[string]interface{}, 0, len(instanceList))
	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := &TagService{client: tcClient}

	for _, v := range instanceList {
		listItem := make(map[string]interface{})
		listItem["id"] = v.DBInstanceId
		listItem["name"] = v.DBInstanceName
		listItem["auto_renew_flag"] = v.AutoRenew
		listItem["project_id"] = v.ProjectId
		listItem["storage"] = v.DBInstanceStorage
		listItem["memory"] = v.DBInstanceMemory
		listItem["availability_zone"] = v.Zone
		listItem["create_time"] = v.CreateTime
		listItem["vpc_id"] = v.VpcId
		listItem["subnet_id"] = v.SubnetId
		listItem["engine_version"] = v.DBVersion
		listItem["public_access_switch"] = false
		listItem["charset"] = v.DBCharset
		listItem["public_access_host"] = ""

		// rootUser
		accounts, outErr := service.DescribeRootUser(ctx, *v.DBInstanceId)
		if outErr != nil {
			return outErr
		}
		if len(accounts) > 0 {
			listItem["root_user"] = accounts[0].UserName
		}

		for _, netInfo := range v.DBInstanceNetInfo {
			if *netInfo.NetType == "public" {
				if *netInfo.Status == "opened" || *netInfo.Status == "1" {
					listItem["public_access_switch"] = true
				}
				listItem["public_access_host"] = netInfo.Address
				listItem["public_access_port"] = netInfo.Port
			}
			if (*netInfo.NetType == "private" || *netInfo.NetType == "inner") && *netInfo.Ip != "" {
				listItem["private_access_ip"] = netInfo.Ip
				listItem["private_access_port"] = netInfo.Port
			}
		}

		if *v.PayType == POSTGRESQL_PAYTYPE_PREPAID || *v.PayType == COMMON_PAYTYPE_PREPAID {
			listItem["charge_type"] = COMMON_PAYTYPE_PREPAID
		} else {
			listItem["charge_type"] = COMMON_PAYTYPE_POSTPAID
		}

		//the describe list API is delayed with argument `tag`
		tagList, err := tagService.DescribeResourceTags(ctx, "postgres", "DBInstanceId", tcClient.Region, *v.DBInstanceId)
		if err != nil {
			return err
		}

		listItem["tags"] = tagList

		list = append(list, listItem)
		ids = append(ids, *v.DBInstanceId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("instance_list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}

	return nil
}
