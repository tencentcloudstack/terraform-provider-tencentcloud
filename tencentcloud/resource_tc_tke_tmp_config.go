package tencentcloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudTkeTmpConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTkeTmpConfigCreate,
		Read:   resourceTencentCloudTkeTmpConfigRead,
		Update: resourceTencentCloudTkeTmpConfigUpdate,
		Delete: resourceTencentCloudTkeTmpConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of instance.",
			},
			"cluster_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of cluster.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of cluster.",
			},
			"config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Global configuration.",
			},
			"service_monitors": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration of the service monitors.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name.",
						},
						"config": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Config.",
						},
						"template_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Used for output parameters, if the configuration comes from a template, it is the template id.",
						},
					},
				},
			},
			"pod_monitors": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration of the pod monitors.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name.",
						},
						"config": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Config.",
						},
						"template_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Used for output parameters, if the configuration comes from a template, it is the template id.",
						},
					},
				},
			},
			"raw_jobs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Configuration of the native prometheus job.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name.",
						},
						"config": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Config.",
						},
						"template_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Used for output parameters, if the configuration comes from a template, it is the template id.",
						},
					},
				},
			},
		},
		//compare to console, miss cam_role and running_version and lock_initial_node and security_proof
	}
}

func resourceTencentCloudTkeTmpConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_config.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId    = getLogId(contextNil)
		service  = TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
		configId = d.Id()
	)

	params, err := service.DescribeTkeTmpConfigById(logId, configId)

	if err != nil {
		return err
	}

	if params == nil {
		d.SetId("")
		return fmt.Errorf("resource `prometheus_config` %s does not exist", configId)
	}

	if e := d.Set("config", params.Config); e != nil {
		log.Printf("[CRITAL]%s provider set config fail, reason:%s\n", logId, e.Error())
		return e
	}
	if e := d.Set("service_monitors", flattenPrometheusConfigItems(params.ServiceMonitors)); e != nil {
		log.Printf("[CRITAL]%s provider set service_monitors fail, reason:%s\n", logId, e.Error())
		return e
	}
	if e := d.Set("pod_monitors", flattenPrometheusConfigItems(params.PodMonitors)); e != nil {
		log.Printf("[CRITAL]%s provider set pod_monitors fail, reason:%s\n", logId, e.Error())
		return e
	}
	if e := d.Set("raw_jobs", flattenPrometheusConfigItems(params.RawJobs)); e != nil {
		log.Printf("[CRITAL]%s provider set raw_jobs fail, reason:%s\n", logId, e.Error())
		return e
	}
	return nil
}

func resourceTencentCloudTkeTmpConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_config.create")()
	defer inconsistentCheck(d, meta)()

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	ret, err := service.CreateTkeTmpConfig(d)

	if err != nil {
		return err
	}

	ids := strings.Join([]string{ret.InstanceId, ret.ClusterType, ret.ClusterId}, FILED_SP)
	d.SetId(ids)

	return resourceTencentCloudTkeTmpConfigRead(d, meta)
}

func resourceTencentCloudTkeTmpConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_config.update, Id: %s", d.Id())()
	defer inconsistentCheck(d, meta)()

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}
	if d.HasChange("cluster_id") {
		return fmt.Errorf("`cluster_id` do not support change now.")
	}
	if d.HasChange("cluster_type") {
		return fmt.Errorf("`cluster_type` do not support change now.")
	}

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	if err := service.UpdateTkeTmpConfig(d); err != nil {
		return err
	}

	return resourceTencentCloudTkeTmpConfigRead(d, meta)
}

func resourceTencentCloudTkeTmpConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tke_tmp_config.delete, Id: %s", d.Id())()
	defer inconsistentCheck(d, meta)()

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	if err := service.DeleteTkeTmpConfig(d); err != nil {
		return err
	}

	return nil
}
