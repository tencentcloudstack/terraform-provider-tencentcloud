package tencentcloud

import (
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func TencentMsyqlBasicInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"period": {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validateAllowedIntValue(MYSQL_AVAILABLE_PERIOD),
		},
		"auto_renew_flag": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
		},
		"engine_version": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validateAllowedStringValue(MYSQL_SUPPORTS_ENGINE),
		},
		"mem_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"volume_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"vpc_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"subnet_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"internet_service": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
		},
		"gtid": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
		},
		"project_id": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"security_groups": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Set: func(v interface{}) int {
				return hashcode.String(v.(string))
			},
		},
		"availability_zone": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"parameters": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"tags": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		// Computed values
		"mysql_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"intranet_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"intranet_port": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"locked": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"task_status": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"internet_host": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"internet_port": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func resourceTencentCloudMysqlInstance() *schema.Resource {
	specialInfo := map[string]*schema.Schema{
		"root_password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validateMysqlPassword,
		},
		"slave_deploy_mode": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
		},
		"first_slave_zone": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"second_slave_zone": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"slave_sync_mode": {
			Type:         schema.TypeInt,
			ValidateFunc: validateAllowedIntValue([]int{0, 1, 2}),
			Optional:     true,
			Default:      0,
		},
	}

	basic := TencentMsyqlBasicInfo()
	for k, v := range basic {
		specialInfo[k] = v
	}
	return &schema.Resource{
		Create: resourceTencentCloudMysqlInstanceCreate,
		Read:   resourceTencentCloudMysqlInstanceRead,
		Update: resourceTencentCloudMysqlInstanceUpdate,
		Delete: resourceTencentCloudMysqlInstanceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: specialInfo,
	}
}

func resourceTencentCloudMysqlInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceTencentCloudMysqlInstanceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceTencentCloudMysqlInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceTencentCloudMysqlInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
