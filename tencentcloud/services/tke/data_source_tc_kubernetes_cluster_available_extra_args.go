package tke

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKubernetesClusterAvailableExtraArgs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKubernetesClusterAvailableExtraArgsRead,
		Schema: map[string]*schema.Schema{
			"cluster_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster version, e.g. `1.28.3`.",
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster type. Valid values: `MANAGED_CLUSTER`, `INDEPENDENT_CLUSTER`.",
			},

			// Computed outputs
			"available_extra_args": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Available custom extra arguments for cluster components.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kube_apiserver": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Available custom arguments for kube-apiserver.",
							Elem:        clusterAvailableExtraArgsFlagSchema(),
						},
						"kube_controller_manager": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Available custom arguments for kube-controller-manager.",
							Elem:        clusterAvailableExtraArgsFlagSchema(),
						},
						"kube_scheduler": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Available custom arguments for kube-scheduler.",
							Elem:        clusterAvailableExtraArgsFlagSchema(),
						},
						"kubelet": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Available custom arguments for kubelet.",
							Elem:        clusterAvailableExtraArgsFlagSchema(),
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

func clusterAvailableExtraArgsFlagSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Argument name.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Argument type.",
			},
			"usage": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Argument description.",
			},
			"default": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default value of the argument.",
			},
			"constraint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Valid range or allowed values of the argument.",
			},
		},
	}
}

func dataSourceTencentCloudKubernetesClusterAvailableExtraArgsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kubernetes_cluster_available_extra_args.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	clusterVersion := d.Get("cluster_version").(string)
	clusterType := d.Get("cluster_type").(string)

	var respParams *tke.DescribeClusterAvailableExtraArgsResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClusterAvailableExtraArgs(ctx, clusterVersion, clusterType)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respParams = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	if respParams != nil && respParams.AvailableExtraArgs != nil {
		availableExtraArgsMap := map[string]interface{}{}

		if respParams.AvailableExtraArgs.KubeAPIServer != nil {
			availableExtraArgsMap["kube_apiserver"] = flattenClusterAvailableExtraArgsFlags(respParams.AvailableExtraArgs.KubeAPIServer)
		}

		if respParams.AvailableExtraArgs.KubeControllerManager != nil {
			availableExtraArgsMap["kube_controller_manager"] = flattenClusterAvailableExtraArgsFlags(respParams.AvailableExtraArgs.KubeControllerManager)
		}

		if respParams.AvailableExtraArgs.KubeScheduler != nil {
			availableExtraArgsMap["kube_scheduler"] = flattenClusterAvailableExtraArgsFlags(respParams.AvailableExtraArgs.KubeScheduler)
		}

		if respParams.AvailableExtraArgs.Kubelet != nil {
			availableExtraArgsMap["kubelet"] = flattenClusterAvailableExtraArgsFlags(respParams.AvailableExtraArgs.Kubelet)
		}

		_ = d.Set("available_extra_args", []interface{}{availableExtraArgsMap})
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}

func flattenClusterAvailableExtraArgsFlags(flags []*tke.Flag) []interface{} {
	result := make([]interface{}, 0, len(flags))
	for _, flag := range flags {
		if flag == nil {
			continue
		}

		flagMap := map[string]interface{}{}
		if flag.Name != nil {
			flagMap["name"] = flag.Name
		}

		if flag.Type != nil {
			flagMap["type"] = flag.Type
		}

		if flag.Usage != nil {
			flagMap["usage"] = flag.Usage
		}

		if flag.Default != nil {
			flagMap["default"] = flag.Default
		}

		if flag.Constraint != nil {
			flagMap["constraint"] = flag.Constraint
		}

		result = append(result, flagMap)
	}

	return result
}
