package tke

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKubernetesClusterRelease() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesClusterReleaseCreate,
		Read:   resourceTencentCloudKubernetesClusterReleaseRead,
		Update: resourceTencentCloudKubernetesClusterReleaseUpdate,
		Delete: resourceTencentCloudKubernetesClusterReleaseDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster ID.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Application name, maximum 63 characters, can only contain lowercase letters, numbers, and the separator \"-\", and must start with a lowercase letter and end with a number or lowercase letter.",
			},

			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Application namespace, obtained from the cluster details namespace.",
			},

			"chart": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Chart name (obtained from the application market) or the download URL of the chart package when installing from a third-party repo, redirect-type chart URLs are not supported, must end with *.tgz.",
			},

			"values": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Custom parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"raw_original": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Custom parameter original value.",
						},
						"values_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Custom parameter value type.",
						},
					},
				},
			},

			"chart_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Chart source, range: tke-market or other, default value: tke-market.",
			},

			"chart_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Chart version.",
			},

			"chart_repo_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Chart repository URL address.",
			},

			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Chart access username.",
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Chart access password.",
			},

			"chart_namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Chart namespace, when ChartFrom is tke-market, ChartNamespace is not empty, value is the Namespace returned by the DescribeProducts interface.",
			},

			"cluster_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cluster type, supports tke, eks, tkeedge, external (registered cluster).",
			},

			// computed
			"cluster_release_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster release ID.",
			},

			"release_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster release status.",
			},
		},
	}
}

func resourceTencentCloudKubernetesClusterReleaseCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_release.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId            = tccommon.GetLogId(tccommon.ContextNil)
		ctx              = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service          = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request          = tkev20180525.NewCreateClusterReleaseRequest()
		response         = tkev20180525.NewCreateClusterReleaseResponse()
		clusterId        string
		name             string
		namespace        string
		clusterReleaseId string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
		name = v.(string)
	}

	if v, ok := d.GetOk("namespace"); ok {
		request.Namespace = helper.String(v.(string))
		namespace = v.(string)
	}

	if v, ok := d.GetOk("chart"); ok {
		request.Chart = helper.String(v.(string))
	}

	if valuesMap, ok := helper.InterfacesHeadMap(d, "values"); ok {
		releaseValues := tkev20180525.ReleaseValues{}
		if v, ok := valuesMap["raw_original"].(string); ok && v != "" {
			releaseValues.RawOriginal = helper.String(v)
		}

		if v, ok := valuesMap["values_type"].(string); ok && v != "" {
			releaseValues.ValuesType = helper.String(v)
		}

		request.Values = &releaseValues
	}

	if v, ok := d.GetOk("chart_from"); ok {
		request.ChartFrom = helper.String(v.(string))
	}

	if v, ok := d.GetOk("chart_version"); ok {
		request.ChartVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("chart_repo_url"); ok {
		request.ChartRepoURL = helper.String(v.(string))
	}

	if v, ok := d.GetOk("username"); ok {
		request.Username = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("chart_namespace"); ok {
		request.ChartNamespace = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_type"); ok {
		request.ClusterType = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().CreateClusterReleaseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Release == nil {
			return resource.NonRetryableError(fmt.Errorf("Create kubernetes cluster release failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create kubernetes cluster release failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.Release.ID == nil {
		return fmt.Errorf("ID is nil.")
	}

	clusterReleaseId = *response.Response.Release.ID
	_ = d.Set("cluster_release_id", clusterReleaseId)
	d.SetId(strings.Join([]string{clusterId, namespace, name}, tccommon.FILED_SP))

	// wait
	reqErr = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		result, e := service.DescribeKubernetesClusterPendingReleaseById(ctx, clusterId, clusterReleaseId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		// get release detail
		if result == nil {
			respData, err := service.DescribeKubernetesClusterReleaseById(ctx, clusterId, namespace, name)
			if err != nil {
				return tccommon.RetryError(e)
			}

			if respData == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe kubernetes cluster release details failed, Response is nil."))
			}

			return nil
		}

		if result.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe kubernetes cluster release failed, Response is nil."))
		}

		if *result.Status == "deployed" || *result.Status == "failed" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Cluster release is still install...Status is %s", *result.Status))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create kubernetes cluster release failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudKubernetesClusterReleaseRead(d, meta)
}

func resourceTencentCloudKubernetesClusterReleaseRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_release.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	namespace := idSplit[1]
	name := idSplit[2]

	respData, err := service.DescribeKubernetesClusterReleaseById(ctx, clusterId, namespace, name)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_kubernetes_cluster_release` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("name", name)
	_ = d.Set("namespace", namespace)

	if respData.ChartFrom != nil {
		_ = d.Set("chart_from", respData.ChartFrom)
	}

	if respData.ChartVersion != nil {
		_ = d.Set("chart_version", respData.ChartVersion)
	}

	if respData.Status != nil {
		_ = d.Set("release_status", respData.Status)
	}

	return nil
}

func resourceTencentCloudKubernetesClusterReleaseUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_release.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	namespace := idSplit[1]
	name := idSplit[2]

	needChange := false
	mutableArgs := []string{"chart", "values", "chart_from", "chart_version", "chart_repo_url", "username", "password", "chart_namespace", "cluster_type"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := tkev20180525.NewUpgradeClusterReleaseRequest()
		if v, ok := d.GetOk("chart"); ok {
			request.Chart = helper.String(v.(string))
		}

		if valuesMap, ok := helper.InterfacesHeadMap(d, "values"); ok {
			releaseValues := tkev20180525.ReleaseValues{}
			if v, ok := valuesMap["raw_original"].(string); ok && v != "" {
				releaseValues.RawOriginal = helper.String(v)
			}

			if v, ok := valuesMap["values_type"].(string); ok && v != "" {
				releaseValues.ValuesType = helper.String(v)
			}

			request.Values = &releaseValues
		}

		if v, ok := d.GetOk("chart_from"); ok {
			request.ChartFrom = helper.String(v.(string))
		}

		if v, ok := d.GetOk("chart_version"); ok {
			request.ChartVersion = helper.String(v.(string))
		}

		if v, ok := d.GetOk("chart_repo_url"); ok {
			request.ChartRepoURL = helper.String(v.(string))
		}

		if v, ok := d.GetOk("username"); ok {
			request.Username = helper.String(v.(string))
		}

		if v, ok := d.GetOk("password"); ok {
			request.Password = helper.String(v.(string))
		}

		if v, ok := d.GetOk("chart_namespace"); ok {
			request.ChartNamespace = helper.String(v.(string))
		}

		if v, ok := d.GetOk("cluster_type"); ok {
			request.ClusterType = helper.String(v.(string))
		}

		request.ClusterId = &clusterId
		request.Namespace = &namespace
		request.Name = &name
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().UpgradeClusterReleaseWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update kubernetes cluster release failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// wait
		waitReq := tkev20180525.NewDescribeClusterReleasesRequest()
		waitReq.ClusterId = &clusterId
		waitReq.Namespace = &namespace
		waitReq.ReleaseName = &name
		reqErr = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DescribeClusterReleasesWithContext(ctx, waitReq)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.ReleaseSet == nil || len(result.Response.ReleaseSet) != 1 {
				return resource.NonRetryableError(fmt.Errorf("Describe kubernetes cluster release failed, Response is nil."))
			}

			release := result.Response.ReleaseSet[0]
			if release.Status == nil {
				return resource.NonRetryableError(fmt.Errorf("Status is nil."))
			}

			if *release.Status == "deployed" || *release.Status == "failed" {
				return nil
			}

			return resource.RetryableError(fmt.Errorf("upgrade release is still running...status is %s", *release.Status))
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s upgrade kubernetes cluster release failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudKubernetesClusterReleaseRead(d, meta)
}

func resourceTencentCloudKubernetesClusterReleaseDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_release.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tkev20180525.NewUninstallClusterReleaseRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	namespace := idSplit[1]
	name := idSplit[2]

	request.ClusterId = &clusterId
	request.Namespace = &namespace
	request.Name = &name
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().UninstallClusterReleaseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete kubernetes cluster release failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// For now no need to perform asynchronous task progress queries
	return nil
}
