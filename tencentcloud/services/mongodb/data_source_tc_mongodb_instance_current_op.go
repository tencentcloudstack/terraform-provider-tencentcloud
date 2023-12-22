package mongodb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMongodbInstanceCurrentOp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMongodbInstanceCurrentOpRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, the format is: cmgo-9d0p6umb.Same as the instance ID displayed in the cloud database console page.",
			},

			"ns": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter condition, the namespace namespace to which the operation belongs, in the format of db.collection.",
			},

			"millisecond_running": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Filter condition, the time that the operation has been executed (unit: millisecond),the result will return the operation that exceeds the set time, the default value is 0,and the value range is [0, 3600000].",
			},

			"op": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter condition, operation type, possible values: none, update, insert, query, command, getmore,remove and killcursors.",
			},

			"replica_set_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "filter condition, shard name.",
			},

			"state": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter condition, node status, possible value: primary, secondary.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Returns the sorted field of the result set, currently supports: MicrosecsRunning/microsecsrunning,the default is ascending sort.",
			},

			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Returns the sorting method of the result set, possible values: ASC/asc or DESC/desc.",
			},

			"current_ops": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "current operation list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"op_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "operation id.",
						},
						"ns": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "operation namespace.",
						},
						"query": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "operation query.",
						},
						"op": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "operation value.",
						},
						"replica_set_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Replication name.",
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "operation state.",
						},
						"operation": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "operation info.",
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node name.",
						},
						"microsecs_running": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "running time(ms).",
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

func dataSourceTencentCloudMongodbInstanceCurrentOpRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mongodb_instance_current_op.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("ns"); ok {
		paramMap["ns"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("millisecond_running"); ok {
		paramMap["millisecond_running"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("op"); ok {
		paramMap["op"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("replica_set_name"); ok {
		paramMap["replica_set_name"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("state"); ok {
		paramMap["state"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["order_by"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["order_by_type"] = helper.String(v.(string))
	}

	service := MongodbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var currentOps []*mongodb.CurrentOp

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMongodbInstanceCurrentOpByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		currentOps = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(currentOps))
	tmpList := make([]map[string]interface{}, 0, len(currentOps))

	if currentOps != nil {
		for _, currentOp := range currentOps {
			currentOpMap := map[string]interface{}{}

			if currentOp.OpId != nil {
				currentOpMap["op_id"] = currentOp.OpId
			}

			if currentOp.Ns != nil {
				currentOpMap["ns"] = currentOp.Ns
			}

			if currentOp.Query != nil {
				currentOpMap["query"] = currentOp.Query
			}

			if currentOp.Op != nil {
				currentOpMap["op"] = currentOp.Op
			}

			if currentOp.ReplicaSetName != nil {
				currentOpMap["replica_set_name"] = currentOp.ReplicaSetName
			}

			if currentOp.State != nil {
				currentOpMap["state"] = currentOp.State
			}

			if currentOp.Operation != nil {
				currentOpMap["operation"] = currentOp.Operation
			}

			if currentOp.NodeName != nil {
				currentOpMap["node_name"] = currentOp.NodeName
			}

			if currentOp.MicrosecsRunning != nil {
				currentOpMap["microsecs_running"] = currentOp.MicrosecsRunning
			}

			ids = append(ids, helper.Int64ToStr(*currentOp.OpId))
			tmpList = append(tmpList, currentOpMap)
		}

		_ = d.Set("current_ops", tmpList)
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
