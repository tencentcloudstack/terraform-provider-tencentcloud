package tencentcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudMongodbShardingInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongodbShardingInstanceCreate,
		Read:   resourceMongodbShardingInstanceRead,
		Update: resourceMongodbShardingInstanceUpdate,
		Delete: resourceMongodbShardingInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(2, 60),
			},
			"shard_quantity": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(2, 20),
			},
			"nodes_per_shard": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(3, 5),
			},
			"memory": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerMin(4),
			},
			"volumn": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerMin(100),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateFunc:
			},
			"machine_type": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateFunc:
			},
			"available_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			// Computed
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vport": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMongodbShardingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceMongodbShardingInstanceRead(d, meta)
}

func resourceMongodbShardingInstanceRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceMongodbShardingInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceMongodbShardingInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
