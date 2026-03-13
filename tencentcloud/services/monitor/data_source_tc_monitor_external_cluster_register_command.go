package monitor

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func DataSourceTencentCloudMonitorExternalClusterRegisterCommand() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorExternalClusterRegisterCommandRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "TMP instance ID.",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "External cluster ID.",
			},

			"command": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Register command, contains YAML configuration for Kubernetes cluster installation.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMonitorExternalClusterRegisterCommandRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_external_cluster_register_command.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.Background(), tccommon.LogIdKey, logId)
		service    = MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId string
		clusterId  string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	var command *string
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeExternalClusterRegisterCommandById(ctx, instanceId, clusterId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		command = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	_ = d.Set("command", command)

	// Set data source ID as composite key
	d.SetId(fmt.Sprintf("%s#%s", instanceId, clusterId))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
