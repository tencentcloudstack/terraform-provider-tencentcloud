package mongodb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func DataSourceTencentCloudMongodbInstanceUrls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMongodbInstanceUrlsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance id.",
			},

			"urls": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Example connection string access address in the form of an instance URI. Contains: URI type and connection string address.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url_type": {
							Type:     schema.TypeString,
							Required: true,
							Description: "Refers to the URI category, including:\n" +
								"	- CLUSTER_ALL: Refers to the main node connected to the library instance through this URI, which can be read and write;\n" +
								"	- CLUSTER_READ_READONLY: Refers to the read-only node connected to the instance through this URI;\n" +
								"	- CLUSTER_READ_SECONDARY: Refers to connecting the instance slave node through this URI;\n" +
								"	- CLUSTER_READ_SECONDARY_AND_READONLY: Refers to the read-only slave node connected to the instance through this URI;\n" +
								"	- CLUSTER_PRIMARY_AND_SECONDARY: This URI connects the instance master node and slave node;\n" +
								"	- MONGOS_ALL: means that each Mongos node is connected through this URI and can be read and write;\n" +
								"	- MONGOS_READ_READONLY: Refers to the read-only node connected to Mongos through this URI;\n" +
								"	- MONGOS_READ_SECONDARY: Refers to the slave node connected to Mongos through this URI;\n" +
								"	- MONGOS_READ_PRIMARY_AND_SECONDARY: refers to the connection between the master node and slave node of Mongos through this URI;\n" +
								"	- MONGOS_READ_SECONDARY_AND_READONLY: refers to the connection between Mongos slave node and read-only node through this URI.",
						},
						"address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Example connection string access address in the form of an instance URI.",
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

func dataSourceTencentCloudMongodbInstanceUrlsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mongodb_instance_urls.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	var respData []*mongodb.DbURL
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMongodbInstanceUrls(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	urlsList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, urls := range respData {
			urlsMap := map[string]interface{}{}

			if urls.URLType != nil {
				urlsMap["url_type"] = urls.URLType
			}

			if urls.Address != nil {
				urlsMap["address"] = urls.Address
			}

			urlsList = append(urlsList, urlsMap)
		}

		_ = d.Set("urls", urlsList)
	}

	d.SetId(instanceId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), urlsList); e != nil {
			return e
		}
	}

	return nil
}
