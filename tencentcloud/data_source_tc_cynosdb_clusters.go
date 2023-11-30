package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbClustersRead,

		Schema: map[string]*schema.Schema{
			"db_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of CynosDB, and available values include `MYSQL`, `POSTGRESQL`.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the cluster to be queried.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the project to be queried.",
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the cluster to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"cluster_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of clusters. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the project.",
						},
						"available_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The available zone of the CynosDB Cluster.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the VPC.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the subnet within this VPC.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port of CynosDB cluster.",
						},
						"db_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of CynosDB, and available values include `MYSQL`.",
						},
						"db_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of CynosDB, which is related to `db_type`. For `MYSQL`, available value is `5.7`.",
						},
						"cluster_limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Storage limit of CynosDB cluster instance, unit in GB.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of CynosDB cluster.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of CynosDB cluster.",
						},
						// payment
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. Default value is `POSTPAID_BY_HOUR`.",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Auto renew flag. Valid values are `0`(MANUAL_RENEW), `1`(AUTO_RENEW). Only works for PREPAID cluster.",
						},
						"cluster_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the Cynosdb cluster.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the CynosDB cluster.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCynosdbClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_clusters.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	params := make(map[string]string)
	if v, ok := d.GetOk("cluster_id"); ok {
		params["ClusterId"] = v.(string)
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		params["ClusterName"] = v.(string)
	}
	if v, ok := d.GetOkExists("project_id"); ok {
		params["ProjectId"] = fmt.Sprintf("%d", v.(int))
	}
	clusterType := ""
	if v, ok := d.GetOk("cluster_type"); ok {
		clusterType = v.(string)
	}

	cynosdbService := CynosdbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var clusters []*cynosdb.CynosdbCluster
	var err error
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		clusters, err = cynosdbService.DescribeClusters(ctx, params)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read cynosdb clusters failed, reason:%s\n ", logId, err.Error())
		return err
	}

	ids := make([]string, 0, len(clusters))
	clusterList := make([]map[string]interface{}, 0, len(clusters))
	for _, cluster := range clusters {
		if clusterType != "" && clusterType != *cluster.DbType {
			continue
		}
		mapping := map[string]interface{}{
			"cluster_id":      cluster.ClusterId,
			"cluster_name":    cluster.ClusterName,
			"cluster_limit":   cluster.StorageLimit,
			"db_type":         cluster.DbType,
			"available_zone":  cluster.Zone,
			"project_id":      cluster.ProjectID,
			"create_time":     cluster.CreateTime,
			"cluster_status":  cluster.Status,
			"auto_renew_flag": cluster.RenewFlag,
			"port":            cluster.Vport,
			"vpc_id":          cluster.VpcId,
			"subnet_id":       cluster.SubnetId,
			"db_version":      cluster.DbVersion,
			"charge_type":     CYNOSDB_CHARGE_TYPE[*cluster.PayMode],
		}

		clusterList = append(clusterList, mapping)
		ids = append(ids, *cluster.ClusterId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("cluster_list", clusterList)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = writeToFile(output.(string), clusterList); err != nil {
			return err
		}
	}

	return nil
}
