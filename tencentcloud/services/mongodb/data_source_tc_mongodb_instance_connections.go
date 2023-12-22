package mongodb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMongodbInstanceConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMongodbInstanceConnectionsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, the format is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page.",
			},

			"clients": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Client connection info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "client connection ip.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "client connection count.",
						},
						"internal_service": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "is internal.",
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

func dataSourceTencentCloudMongodbInstanceConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mongodb_instance_connections.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var clients []*mongodb.ClientConnection

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMongodbInstanceConnectionsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		clients = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(clients))
	tmpList := make([]map[string]interface{}, 0, len(clients))

	if clients != nil {
		for _, clientConnection := range clients {
			clientConnectionMap := map[string]interface{}{}

			if clientConnection.IP != nil {
				clientConnectionMap["ip"] = clientConnection.IP
			}

			if clientConnection.Count != nil {
				clientConnectionMap["count"] = clientConnection.Count
			}

			if clientConnection.InternalService != nil {
				clientConnectionMap["internal_service"] = clientConnection.InternalService
			}

			ids = append(ids, *clientConnection.IP)
			tmpList = append(tmpList, clientConnectionMap)
		}

		_ = d.Set("clients", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
