/*
Provides a resource to update elasticsearch plugins

Example Usage

```hcl
resource "tencentcloud_elasticsearch_update_plugins_operation" "update_plugins_operation" {
  instance_id = "es-xxxxxx"
  install_plugin_list = ["analysis-pinyin"]
  force_restart = false
  force_update = true
}
```
*/
package tencentcloud

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudElasticsearchUpdatePluginsOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudElasticsearchUpdatePluginsOperationCreate,
		Read:   resourceTencentCloudElasticsearchUpdatePluginsOperationRead,
		Delete: resourceTencentCloudElasticsearchUpdatePluginsOperationDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"install_plugin_list": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of plugins that need to be installed.",
			},

			"remove_plugin_list": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of plugins that need to be uninstalled.",
			},

			"force_restart": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to force a restart. Default is false.",
			},

			"force_update": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to reinstall, default value false.",
			},

			"plugin_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Plugin type. 0: system plugin.",
			},
		},
	}
}

func resourceTencentCloudElasticsearchUpdatePluginsOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_update_plugins_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = elasticsearch.NewUpdatePluginsRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("install_plugin_list"); ok {
		installPluginListSet := v.(*schema.Set).List()
		for i := range installPluginListSet {
			installPluginList := installPluginListSet[i].(string)
			request.InstallPluginList = append(request.InstallPluginList, &installPluginList)
		}
	}

	if v, ok := d.GetOk("remove_plugin_list"); ok {
		removePluginListSet := v.(*schema.Set).List()
		for i := range removePluginListSet {
			removePluginList := removePluginListSet[i].(string)
			request.RemovePluginList = append(request.RemovePluginList, &removePluginList)
		}
	}

	if v, ok := d.GetOkExists("force_restart"); ok {
		request.ForceRestart = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("force_update"); ok {
		request.ForceUpdate = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("plugin_type"); ok {
		request.PluginType = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseEsClient().UpdatePlugins(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate elasticsearch UpdatePluginsOperation failed, reason:%+v", logId, err)
		return err
	}

	time.Sleep(2 * time.Second)
	elasticsearchService := ElasticsearchService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	conf := BuildStateChangeConf([]string{}, []string{"1"}, 10*readRetryTimeout, time.Second, elasticsearchService.ElasticsearchInstanceRefreshFunc(instanceId, []string{}))
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	d.SetId(instanceId)

	return resourceTencentCloudElasticsearchUpdatePluginsOperationRead(d, meta)
}

func resourceTencentCloudElasticsearchUpdatePluginsOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_update_plugins_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudElasticsearchUpdatePluginsOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_elasticsearch_update_plugins_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
