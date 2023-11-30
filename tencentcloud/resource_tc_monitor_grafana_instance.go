package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorGrafanaInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorGrafanaInstanceRead,
		Create: resourceTencentCloudMonitorGrafanaInstanceCreate,
		Update: resourceTencentCloudMonitorGrafanaInstanceUpdate,
		Delete: resourceTencentCloudMonitorGrafanaInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Grafana instance id.",
			},

			"root_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Grafana external url which could be accessed by user.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Vpc Id.",
			},

			"subnet_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Computed:    true,
				Description: "Subnet Id array.",
			},

			"grafana_init_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Grafana server admin password.",
			},

			"enable_internet": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Control whether grafana could be accessed by internet.",
			},

			"is_distroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Deprecated:  "It has been deprecated from version 1.81.16.",
				Description: "Whether to clean up completely, the default is false.",
			},

			"is_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to clean up completely, the default is false.",
			},

			"instance_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Grafana instance status, 1: Creating, 2: Running, 6: Stopped.",
			},

			"internet_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Grafana intranet address.",
			},

			"internal_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Grafana public address.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = monitor.NewCreateGrafanaInstanceRequest()
		response *monitor.CreateGrafanaInstanceResponse
	)

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_ids"); ok {
		subnetIdsSet := v.(*schema.Set).List()
		for i := range subnetIdsSet {
			subnetIds := subnetIdsSet[i].(string)
			request.SubnetIds = append(request.SubnetIds, &subnetIds)
		}
	}

	if v, ok := d.GetOk("grafana_init_password"); ok {
		request.GrafanaInitPassword = helper.String(v.(string))
	}

	if v, _ := d.GetOk("enable_internet"); v != nil {
		// Internal account won't open
		request.EnableInternet = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreateGrafanaInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create monitor grafanaInstance failed, reason:%+v", logId, err)
		return err
	}

	grafanaInstanceId := *response.Response.InstanceId

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeMonitorGrafanaInstance(ctx, grafanaInstanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.InstanceStatus == 2 {
			return nil
		}
		if *instance.InstanceStatus == 6 {
			return resource.NonRetryableError(fmt.Errorf("grafanaInstance status is %v, operate failed.", *instance.InstanceStatus))
		}
		return resource.RetryableError(fmt.Errorf("grafanaInstance status is %v, retry...", *instance.InstanceStatus))
	})
	if err != nil {
		return err
	}

	d.SetId(grafanaInstanceId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::monitor:%s:uin/:grafana-instance/%s", region, grafanaInstanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	return resourceTencentCloudMonitorGrafanaInstanceRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	grafanaInstance, err := service.DescribeMonitorGrafanaInstance(ctx, instanceId)

	if err != nil {
		return err
	}

	if grafanaInstance == nil {
		d.SetId("")
		return fmt.Errorf("resource `grafanaInstance` %s does not exist", instanceId)
	}

	if grafanaInstance.InstanceName != nil {
		_ = d.Set("instance_name", grafanaInstance.InstanceName)
	}

	if grafanaInstance.InstanceId != nil {
		_ = d.Set("instance_id", grafanaInstance.InstanceId)
	}

	if grafanaInstance.RootUrl != nil {
		_ = d.Set("root_url", grafanaInstance.RootUrl)
	}

	if grafanaInstance.VpcId != nil {
		_ = d.Set("vpc_id", grafanaInstance.VpcId)
	}

	if grafanaInstance.SubnetIds != nil {
		var subnetIds []string
		for _, v := range grafanaInstance.SubnetIds {
			subnetIds = append(subnetIds, *v)
		}
		_ = d.Set("subnet_ids", subnetIds)
	}

	if grafanaInstance.InternetUrl != nil && *grafanaInstance.InternetUrl != "" {
		_ = d.Set("enable_internet", true)
	} else {
		_ = d.Set("enable_internet", false)
	}

	if grafanaInstance.InternetUrl != nil {
		_ = d.Set("internet_url", grafanaInstance.InternetUrl)
	}

	if grafanaInstance.InternalUrl != nil {
		_ = d.Set("internal_url", grafanaInstance.InternalUrl)
	}

	if grafanaInstance.InstanceStatus != nil {
		_ = d.Set("instance_status", grafanaInstance.InstanceStatus)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "monitor", "grafana-instance", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudMonitorGrafanaInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()

	if d.HasChange("instance_name") {
		request := monitor.NewModifyGrafanaInstanceRequest()
		request.InstanceId = &instanceId
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().ModifyGrafanaInstance(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("vpc_id") {
		return fmt.Errorf("`vpc_id` do not support change now.")
	}

	if d.HasChange("subnet_ids") {
		return fmt.Errorf("`subnet_ids` do not support change now.")
	}

	if d.HasChange("grafana_init_password") {
		return fmt.Errorf("`grafana_init_password` do not support change now.")
	}

	if d.HasChange("enable_internet") {
		request := monitor.NewEnableGrafanaInternetRequest()
		request.InstanceID = &instanceId

		if v, ok := d.GetOk("enable_internet"); ok {
			request.EnableInternet = helper.Bool(v.(bool))
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().EnableGrafanaInternet(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("monitor", "grafana-instance", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudMonitorGrafanaInstanceRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	if err := service.DeleteMonitorGrafanaInstanceById(ctx, instanceId); err != nil {
		return err
	}

	err := resource.Retry(1*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeMonitorGrafanaInstance(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.InstanceStatus == 6 {
			return nil
		}
		if *instance.InstanceStatus == 7 {
			return resource.NonRetryableError(fmt.Errorf("grafanaInstance status is %v, operate failed.", *instance.InstanceStatus))
		}
		return resource.RetryableError(fmt.Errorf("grafanaInstance status is %v, retry...", *instance.InstanceStatus))
	})
	if err != nil {
		return err
	}

	claenFlag := false
	if v, ok := d.GetOk("is_distroy"); ok && v.(bool) {
		claenFlag = true
	}
	if v, ok := d.GetOk("is_destroy"); ok && v.(bool) {
		claenFlag = true
	}
	if claenFlag {
		if err := service.CleanGrafanaInstanceById(ctx, instanceId); err != nil {
			return err
		}

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			instance, errRet := service.DescribeMonitorGrafanaInstance(ctx, instanceId)
			if errRet != nil {
				return retryError(errRet, InternalError)
			}
			if instance == nil {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("grafanaInstance status is %v, retry...", *instance.InstanceStatus))
		})
		if err != nil {
			return err
		}
	}

	return nil
}
