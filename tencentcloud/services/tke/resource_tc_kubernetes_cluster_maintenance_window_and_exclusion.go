package tke

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKubernetesClusterMaintenanceWindowAndExclusion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesClusterMaintenanceWindowAndExclusionCreate,
		Read:   resourceTencentCloudKubernetesClusterMaintenanceWindowAndExclusionRead,
		Update: resourceTencentCloudKubernetesClusterMaintenanceWindowAndExclusionUpdate,
		Delete: resourceTencentCloudKubernetesClusterMaintenanceWindowAndExclusionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster ID, can be obtained from the container service cluster console (https://console.cloud.tencent.com/tke2/cluster).",
			},

			"maintenance_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Maintenance start time.",
			},

			"duration": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maintenance duration (hours).",
			},

			"day_of_week": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "'Maintenance cycle (which days of the week), supported parameter values are as follows:\n\n- MO: Monday\n- TU: Tuesday\n- WE: Wednesday\n- TH: Thursday\n- FR: Friday\n- SA: Saturday\n- SU: Sunday.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"exclusions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Maintenance exclusions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Maintenance exclusion name.",
						},
						"start_at": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Maintenance exclusion start time.",
						},
						"end_at": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Maintenance exclusion end time.",
						},
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maintenance exclusion ID.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudKubernetesClusterMaintenanceWindowAndExclusionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_maintenance_window_and_exclusion.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = tkev20180525.NewCreateClusterMaintenanceWindowAndExclusionsRequest()
		clusterId string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterID = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("maintenance_time"); ok {
		request.MaintenanceTime = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("duration"); ok {
		request.Duration = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("day_of_week"); ok {
		dayOfWeekSet := v.(*schema.Set).List()
		for i := range dayOfWeekSet {
			dayOfWeek := dayOfWeekSet[i].(string)
			request.DayOfWeek = append(request.DayOfWeek, helper.String(dayOfWeek))
		}
	}

	if v, ok := d.GetOk("exclusions"); ok {
		for _, item := range v.([]interface{}) {
			exclusionsMap := item.(map[string]interface{})
			maintenanceExclusion := tkev20180525.MaintenanceExclusion{}
			if v, ok := exclusionsMap["name"].(string); ok && v != "" {
				maintenanceExclusion.Name = helper.String(v)
			}

			if v, ok := exclusionsMap["start_at"].(string); ok && v != "" {
				maintenanceExclusion.StartAt = helper.String(v)
			}

			if v, ok := exclusionsMap["end_at"].(string); ok && v != "" {
				maintenanceExclusion.EndAt = helper.String(v)
			}

			request.Exclusions = append(request.Exclusions, &maintenanceExclusion)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().CreateClusterMaintenanceWindowAndExclusionsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create kubernetes cluster maintenance window and exclusion failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(clusterId)
	return resourceTencentCloudKubernetesClusterMaintenanceWindowAndExclusionRead(d, meta)
}

func resourceTencentCloudKubernetesClusterMaintenanceWindowAndExclusionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_maintenance_window_and_exclusion.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		clusterId = d.Id()
	)

	respData, err := service.DescribeKubernetesClusterMaintenanceWindowAndExclusionById(ctx, clusterId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_kubernetes_cluster_maintenance_window_and_exclusion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.ClusterID != nil {
		_ = d.Set("cluster_id", respData.ClusterID)
	}

	if respData.MaintenanceTime != nil {
		_ = d.Set("maintenance_time", respData.MaintenanceTime)
	}

	if respData.Duration != nil {
		_ = d.Set("duration", respData.Duration)
	}

	if respData.DayOfWeek != nil {
		_ = d.Set("day_of_week", respData.DayOfWeek)
	}

	if respData.Exclusions != nil {
		exclusionsList := make([]map[string]interface{}, 0, len(respData.Exclusions))
		for _, exclusions := range respData.Exclusions {
			exclusionsMap := map[string]interface{}{}
			if exclusions.Name != nil {
				exclusionsMap["name"] = exclusions.Name
			}

			if exclusions.StartAt != nil {
				exclusionsMap["start_at"] = exclusions.StartAt
			}

			if exclusions.EndAt != nil {
				exclusionsMap["end_at"] = exclusions.EndAt
			}

			if exclusions.ID != nil {
				exclusionsMap["id"] = exclusions.ID
			}

			exclusionsList = append(exclusionsList, exclusionsMap)
		}

		_ = d.Set("exclusions", exclusionsList)
	}

	return nil
}

func resourceTencentCloudKubernetesClusterMaintenanceWindowAndExclusionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_maintenance_window_and_exclusion.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		clusterId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"maintenance_time", "duration", "day_of_week", "exclusions"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := tkev20180525.NewModifyClusterMaintenanceWindowAndExclusionsRequest()
		if v, ok := d.GetOk("maintenance_time"); ok {
			request.MaintenanceTime = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("duration"); ok {
			request.Duration = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("day_of_week"); ok {
			dayOfWeekSet := v.(*schema.Set).List()
			for i := range dayOfWeekSet {
				dayOfWeek := dayOfWeekSet[i].(string)
				request.DayOfWeek = append(request.DayOfWeek, helper.String(dayOfWeek))
			}
		}

		if v, ok := d.GetOk("exclusions"); ok {
			for _, item := range v.([]interface{}) {
				exclusionsMap := item.(map[string]interface{})
				maintenanceExclusion := tkev20180525.MaintenanceExclusion{}
				if v, ok := exclusionsMap["name"].(string); ok && v != "" {
					maintenanceExclusion.Name = helper.String(v)
				}

				if v, ok := exclusionsMap["start_at"].(string); ok && v != "" {
					maintenanceExclusion.StartAt = helper.String(v)
				}

				if v, ok := exclusionsMap["end_at"].(string); ok && v != "" {
					maintenanceExclusion.EndAt = helper.String(v)
				}

				request.Exclusions = append(request.Exclusions, &maintenanceExclusion)
			}
		}

		request.ClusterID = &clusterId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ModifyClusterMaintenanceWindowAndExclusionsWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update kubernetes cluster maintenance window and exclusion failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudKubernetesClusterMaintenanceWindowAndExclusionRead(d, meta)
}

func resourceTencentCloudKubernetesClusterMaintenanceWindowAndExclusionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_maintenance_window_and_exclusion.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = tkev20180525.NewDeleteClusterMaintenanceWindowAndExclusionRequest()
		clusterId = d.Id()
	)

	request.ClusterID = &clusterId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DeleteClusterMaintenanceWindowAndExclusionWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete kubernetes cluster maintenance window and exclusion failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
