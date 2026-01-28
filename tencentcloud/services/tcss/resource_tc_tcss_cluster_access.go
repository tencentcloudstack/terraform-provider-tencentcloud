package tcss

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcssv20201101 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcss/v20201101"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTcssClusterAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcssClusterAccessCreate,
		Read:   resourceTencentCloudTcssClusterAccessRead,
		Update: resourceTencentCloudTcssClusterAccessUpdate,
		Delete: resourceTencentCloudTcssClusterAccessDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster Id.",
			},

			"switch_on": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable cluster defend status.",
			},

			// computed
			"accessed_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster access status.",
			},

			"defender_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster defender status.",
			},
		},
	}
}

func resourceTencentCloudTcssClusterAccessCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcss_cluster_access.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = tcssv20201101.NewCreateClusterAccessRequest()
		clusterId string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterIDs = append(request.ClusterIDs, helper.String(v.(string)))
		clusterId = v.(string)
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcssV20201101Client().CreateClusterAccessWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create tcss cluster access failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(clusterId)

	// wait
	waitReq := tcssv20201101.NewDescribeUserClusterRequest()
	waitReq.Offset = helper.IntUint64(0)
	waitReq.Limit = helper.IntUint64(1)
	waitReq.Filters = []*tcssv20201101.ComplianceFilters{
		&tcssv20201101.ComplianceFilters{
			Name:       helper.String("ClusterID"),
			Values:     helper.Strings([]string{clusterId}),
			ExactMatch: helper.Bool(true),
		},
	}
	reqErr = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcssV20201101Client().DescribeUserClusterWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), waitReq.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.ClusterInfoList == nil || len(result.Response.ClusterInfoList) < 1 {
			return resource.NonRetryableError(fmt.Errorf("Describe user cluster failed, Response is nil."))
		}

		clusterInfo := result.Response.ClusterInfoList[0]
		if clusterInfo.AccessedStatus == nil {
			return resource.NonRetryableError(fmt.Errorf("AccessedStatus is nil."))
		}

		if *clusterInfo.AccessedStatus == "AccessedException" {
			return resource.NonRetryableError(fmt.Errorf("AccessedStatus is `AccessedException`."))
		}

		if *clusterInfo.AccessedStatus == "AccessedDefended" || *clusterInfo.AccessedStatus == "AccessedInstalled" || *clusterInfo.AccessedStatus == "AccessedPartialDefence" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Cluster access is still running...Accessed status is %s", *clusterInfo.AccessedStatus))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create tcss cluster access failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// set switch
	if v, ok := d.GetOkExists("switch_on"); ok {
		if v.(bool) {
			request := tcssv20201101.NewModifyDefendStatusRequest()
			request.InstanceIDs = helper.Strings([]string{clusterId})
			request.InstanceType = helper.String("Cluster")
			request.IsAll = helper.Bool(false)
			request.SwitchOn = helper.Bool(v.(bool))
			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcssV20201101Client().ModifyDefendStatusWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s modify tcss cluster access failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			// wait
			waitReq := tcssv20201101.NewDescribeUserClusterRequest()
			waitReq.Offset = helper.IntUint64(0)
			waitReq.Limit = helper.IntUint64(1)
			waitReq.Filters = []*tcssv20201101.ComplianceFilters{
				&tcssv20201101.ComplianceFilters{
					Name:       helper.String("ClusterID"),
					Values:     helper.Strings([]string{clusterId}),
					ExactMatch: helper.Bool(true),
				},
			}
			reqErr = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcssV20201101Client().DescribeUserClusterWithContext(ctx, waitReq)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), waitReq.ToJsonString())
				}

				if result == nil || result.Response == nil || result.Response.ClusterInfoList == nil || len(result.Response.ClusterInfoList) < 1 {
					return resource.NonRetryableError(fmt.Errorf("Describe user cluster failed, Response is nil."))
				}

				clusterInfo := result.Response.ClusterInfoList[0]
				if clusterInfo.DefenderStatus == nil {
					return resource.NonRetryableError(fmt.Errorf("DefenderStatus is nil."))
				}

				if *clusterInfo.DefenderStatus == "Defended" || *clusterInfo.DefenderStatus == "PartDefened" {
					return nil
				}

				return resource.RetryableError(fmt.Errorf("Cluster access is still running...Defender status is %s", *clusterInfo.DefenderStatus))
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s modify tcss cluster access failed, reason:%+v", logId, reqErr)
				return reqErr
			}
		}
	}

	return resourceTencentCloudTcssClusterAccessRead(d, meta)
}

func resourceTencentCloudTcssClusterAccessRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcss_cluster_access.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = TcssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		clusterId = d.Id()
	)

	respData, err := service.DescribeTcssClusterAccessById(ctx, clusterId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_tcss_cluster_access` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.ClusterId != nil {
		_ = d.Set("cluster_id", respData.ClusterId)
	}

	if respData.DefenderStatus != nil {
		if *respData.DefenderStatus == "Defended" || *respData.DefenderStatus == "PartDefened" {
			_ = d.Set("switch_on", true)
		} else {
			_ = d.Set("switch_on", false)
		}

		_ = d.Set("defender_status", *respData.DefenderStatus)
	}

	if respData.AccessedStatus != nil {
		_ = d.Set("accessed_status", *respData.AccessedStatus)
	}

	return nil
}

func resourceTencentCloudTcssClusterAccessUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcss_cluster_access.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = tcssv20201101.NewModifyDefendStatusRequest()
		clusterId = d.Id()
	)

	if d.HasChange("switch_on") {
		var switchOn bool
		if v, ok := d.GetOkExists("switch_on"); ok {
			switchOn = v.(bool)
		}

		request.InstanceIDs = helper.Strings([]string{clusterId})
		request.InstanceType = helper.String("Cluster")
		request.IsAll = helper.Bool(false)
		request.SwitchOn = helper.Bool(switchOn)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcssV20201101Client().ModifyDefendStatusWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s modify tcss cluster access failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// wait
		waitReq := tcssv20201101.NewDescribeUserClusterRequest()
		waitReq.Offset = helper.IntUint64(0)
		waitReq.Limit = helper.IntUint64(1)
		waitReq.Filters = []*tcssv20201101.ComplianceFilters{
			&tcssv20201101.ComplianceFilters{
				Name:       helper.String("ClusterID"),
				Values:     helper.Strings([]string{clusterId}),
				ExactMatch: helper.Bool(true),
			},
		}
		reqErr = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcssV20201101Client().DescribeUserClusterWithContext(ctx, waitReq)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), waitReq.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.ClusterInfoList == nil || len(result.Response.ClusterInfoList) < 1 {
				return resource.NonRetryableError(fmt.Errorf("Describe user cluster failed, Response is nil."))
			}

			clusterInfo := result.Response.ClusterInfoList[0]
			if clusterInfo.DefenderStatus == nil {
				return resource.NonRetryableError(fmt.Errorf("DefenderStatus is nil."))
			}

			if switchOn {
				if *clusterInfo.DefenderStatus == "Defended" || *clusterInfo.DefenderStatus == "PartDefened" {
					return nil
				}
			} else {
				if *clusterInfo.DefenderStatus == "UnDefended" {
					return nil
				}
			}

			return resource.RetryableError(fmt.Errorf("Cluster access is still running...Defender status is %s", *clusterInfo.DefenderStatus))
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s modify tcss cluster access failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTcssClusterAccessRead(d, meta)
}

func resourceTencentCloudTcssClusterAccessDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcss_cluster_access.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = tcssv20201101.NewUninstallClusterContainerSecurityRequest()
		clusterId = d.Id()
	)

	request.ClusterIDs = append(request.ClusterIDs, &clusterId)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcssV20201101Client().UninstallClusterContainerSecurityWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete tcss cluster access failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// wait
	waitReq := tcssv20201101.NewDescribeUserClusterRequest()
	waitReq.Offset = helper.IntUint64(0)
	waitReq.Limit = helper.IntUint64(1)
	waitReq.Filters = []*tcssv20201101.ComplianceFilters{
		&tcssv20201101.ComplianceFilters{
			Name:       helper.String("ClusterID"),
			Values:     helper.Strings([]string{clusterId}),
			ExactMatch: helper.Bool(true),
		},
	}
	reqErr = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTcssV20201101Client().DescribeUserClusterWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), waitReq.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.ClusterInfoList == nil || len(result.Response.ClusterInfoList) < 1 {
			return resource.NonRetryableError(fmt.Errorf("Describe user cluster failed, Response is nil."))
		}

		clusterInfo := result.Response.ClusterInfoList[0]
		if clusterInfo.AccessedStatus == nil {
			return resource.NonRetryableError(fmt.Errorf("AccessedStatus is nil."))
		}

		if *clusterInfo.AccessedStatus == "AccessedUninstallException" {
			return resource.NonRetryableError(fmt.Errorf("AccessedStatus is `AccessedUninstallException`."))
		}

		if *clusterInfo.AccessedStatus == "AccessedNone" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Cluster access is still running...Accessed status is %s", *clusterInfo.AccessedStatus))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete tcss cluster access failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
