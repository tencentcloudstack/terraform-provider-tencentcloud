package tke

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudKubernetesClusterExtraArgsConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesClusterExtraArgsConfigCreate,
		Read:   resourceTencentCloudKubernetesClusterExtraArgsConfigRead,
		Update: resourceTencentCloudKubernetesClusterExtraArgsConfigUpdate,
		Delete: resourceTencentCloudKubernetesClusterExtraArgsConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster ID. Only managed clusters are supported.",
			},

			"kube_apiserver": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Custom args for kube-apiserver, format: [\"k1=v1\", \"k2=v2\"], e.g. [\"max-requests-inflight=500\",\"feature-gates=PodShareProcessNamespace=true\"].",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"kube_controller_manager": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Custom args for kube-controller-manager, format: [\"k1=v1\", \"k2=v2\"].",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"kube_scheduler": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Custom args for kube-scheduler, format: [\"k1=v1\", \"k2=v2\"].",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"etcd": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Custom args for etcd. Only standalone clusters are supported, format: [\"k1=v1\", \"k2=v2\"].",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudKubernetesClusterExtraArgsConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_extra_args_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	clusterId := d.Get("cluster_id").(string)
	d.SetId(clusterId)

	return resourceTencentCloudKubernetesClusterExtraArgsConfigUpdate(d, meta)
}

func resourceTencentCloudKubernetesClusterExtraArgsConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_extra_args_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	clusterId := d.Id()
	_ = d.Set("cluster_id", clusterId)

	extraArgs, err := service.DescribeKubernetesClusterExtraArgsConfig(ctx, clusterId)
	if err != nil {
		return err
	}

	if extraArgs == nil {
		log.Printf("[WARN]%s resource `tencentcloud_kubernetes_cluster_extra_args_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if extraArgs.KubeAPIServer != nil {
		_ = d.Set("kube_apiserver", helper.PStrings(extraArgs.KubeAPIServer))
	}

	if extraArgs.KubeControllerManager != nil {
		_ = d.Set("kube_controller_manager", helper.PStrings(extraArgs.KubeControllerManager))
	}

	if extraArgs.KubeScheduler != nil {
		_ = d.Set("kube_scheduler", helper.PStrings(extraArgs.KubeScheduler))
	}

	if extraArgs.Etcd != nil {
		_ = d.Set("etcd", helper.PStrings(extraArgs.Etcd))
	}

	return nil
}

func resourceTencentCloudKubernetesClusterExtraArgsConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_extra_args_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tke.NewModifyClusterExtraArgsRequest()
	)

	clusterId := d.Id()
	request.ClusterId = &clusterId
	request.ClusterExtraArgs = &tke.ClusterExtraArgs{}

	if v, ok := d.GetOk("kube_apiserver"); ok {
		request.ClusterExtraArgs.KubeAPIServer = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("kube_controller_manager"); ok {
		request.ClusterExtraArgs.KubeControllerManager = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("kube_scheduler"); ok {
		request.ClusterExtraArgs.KubeScheduler = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("etcd"); ok {
		request.ClusterExtraArgs.Etcd = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ModifyClusterExtraArgsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update kubernetes cluster extra args config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	waitTimeout := d.Timeout(schema.TimeoutUpdate)
	if d.IsNewResource() {
		waitTimeout = d.Timeout(schema.TimeoutCreate)
	}

	if err := waitForClusterExtraArgsTaskDone(ctx, meta, clusterId, logId, waitTimeout); err != nil {
		log.Printf("[CRITAL]%s wait for cluster extra args task failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudKubernetesClusterExtraArgsConfigRead(d, meta)
}

func resourceTencentCloudKubernetesClusterExtraArgsConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_extra_args_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func waitForClusterExtraArgsTaskDone(ctx context.Context, meta interface{}, clusterId, logId string, timeout time.Duration) error {
	request := tke.NewDescribeTasksRequest()
	request.Latest = helper.Bool(true)
	request.Filter = []*tke.Filter{
		{
			Name:   helper.String("TaskType"),
			Values: []*string{helper.String("apply_cluster_extra_args")},
		},
		{
			Name:   helper.String("ClusterId"),
			Values: []*string{&clusterId},
		},
	}

	return resource.Retry(timeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DescribeTasksWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.RetryableError(fmt.Errorf("waiting for cluster extra args task, response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if len(result.Response.Tasks) == 0 {
			return resource.RetryableError(fmt.Errorf("waiting for cluster extra args task, no tasks found"))
		}

		task := result.Response.Tasks[0]
		if task.LifeState == nil {
			return resource.RetryableError(fmt.Errorf("waiting for cluster extra args task, LifeState is nil"))
		}

		lifeState := *task.LifeState
		if lifeState == "done" {
			return nil
		}

		if lifeState == "abort" || lifeState == "aborted" {
			lastErr := ""
			if task.LastError != nil {
				lastErr = *task.LastError
			}

			return resource.NonRetryableError(fmt.Errorf("cluster extra args task failed with state: %s, error: %s", lifeState, lastErr))
		}

		return resource.RetryableError(fmt.Errorf("waiting for cluster extra args task, current LifeState: %s", lifeState))
	})
}
