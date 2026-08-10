package tcaplusdb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tcaplusdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb/v20190823"
)

func ResourceTencentCloudTcaplusCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcaplusClusterCreate,
		Read:   resourceTencentCloudTcaplusClusterRead,
		Update: resourceTencentCloudTcaplusClusterUpdate,
		Delete: resourceTencentCloudTcaplusClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"idl_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(TCAPLUS_IDL_TYPES),
				Description:  "IDL type of the TcaplusDB cluster. Valid values: `PROTO` and `TDR`.",
			},
			"cluster_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 30),
				Description:  "Name of the TcaplusDB cluster. Name length should be between 1 and 30.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "VPC id of the TcaplusDB cluster.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subnet id of the TcaplusDB cluster.",
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) < 12 || len(value) > 16 {
						errors = append(errors, fmt.Errorf("invalid password, length should between 12 and 16"))
						return
					}
					var match = make(map[string]bool)
					for i := 0; i < len(value); i++ {
						if len(match) >= 2 {
							break
						}
						if value[i] >= '0' && value[i] <= '9' {
							match["number"] = true
							continue
						}
						if value[i] >= 'a' && value[i] <= 'z' {
							match["low"] = true
							continue
						}
						if value[i] >= 'A' && value[i] <= 'Z' {
							match["up"] = true
							continue
						}
					}
					if len(match) < 2 {
						errors = append(errors, fmt.Errorf("invalid password, a-z and 0-9 and A-Z must contain"))
					}
					return
				},
				Description: "Password of the TcaplusDB cluster. Password length should be between 12 and 16. The password must be a *mix* of uppercase letters (A-Z), lowercase *letters* (a-z) and *numbers* (0-9).",
			},
			"old_password_expire_last": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3600,
				ValidateFunc: tccommon.ValidateIntegerMin(300),
				Description:  "Expiration time of old password after password update, unit: second.",
			},
			"cluster_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Cluster type of the TcaplusDB cluster. `1`: shared cluster, `2`: dedicated cluster. This parameter is only valid for CreateCluster API and cannot be modified once set.",
			},
			"resource_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tag_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Description: "Resource tags of the TcaplusDB cluster. This parameter is only valid for CreateCluster API and cannot be modified once set. Note: This field is write-only and will not be refreshed on Read because the DescribeClusters API does not return cluster-level tags.",
			},
			"server_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"machine_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"machine_num": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"server_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_rate": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_rate": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"read_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"write_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: "Dedicated server machine list of the TcaplusDB cluster. Only valid when `cluster_type` is `2` (dedicated cluster). For creation, each element exposes `machine_type` and `machine_num`. This parameter is only valid for CreateCluster API and cannot be modified once set.",
			},
			"proxy_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"machine_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"machine_num": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"proxy_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"process_speed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"average_process_delay": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"slow_process_speed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: "Dedicated proxy machine list of the TcaplusDB cluster. Only valid when `cluster_type` is `2` (dedicated cluster). For creation, each element exposes `machine_type` and `machine_num`. This parameter is only valid for CreateCluster API and cannot be modified once set.",
			},

			// Computed values.
			"network_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Network type of the TcaplusDB cluster.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the TcaplusDB cluster.",
			},
			"password_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Password status of the TcaplusDB cluster. Valid values: `unmodifiable`, `modifiable`. `unmodifiable`. which means the password can not be changed in this moment; `modifiable`, which means the password can be changed in this moment.",
			},
			"api_access_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access ID of the TcaplusDB cluster.For TcaplusDB SDK connect.",
			},
			"api_access_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Access IP of the TcaplusDB cluster.For TcaplusDB SDK connect.",
			},
			"api_access_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Access port of the TcaplusDB cluster.For TcaplusDB SDK connect.",
			},
			"old_password_expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration time of the old password. If `password_status` is `unmodifiable`, it means the old password has not yet expired.",
			},
		},
	}
}

func resourceTencentCloudTcaplusClusterCreate(d *schema.ResourceData, meta interface{}) error {

	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_cluster.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		idlType     = d.Get("idl_type").(string)
		clusterName = d.Get("cluster_name").(string)
		vpcId       = d.Get("vpc_id").(string)
		subnetId    = d.Get("subnet_id").(string)
		password    = d.Get("password").(string)
		clusterType = int64(d.Get("cluster_type").(int))
	)

	var resourceTags []*tcaplusdb.TagInfoUnit
	if tags, ok := d.Get("resource_tags").([]interface{}); ok && len(tags) > 0 {
		for _, tag := range tags {
			tagMap := tag.(map[string]interface{})
			tagKey := tagMap["tag_key"].(string)
			tagValue := tagMap["tag_value"].(string)
			resourceTags = append(resourceTags, &tcaplusdb.TagInfoUnit{
				TagKey:   &tagKey,
				TagValue: &tagValue,
			})
		}
	}

	var serverList []*tcaplusdb.MachineInfo
	if servers, ok := d.Get("server_list").([]interface{}); ok && len(servers) > 0 {
		for _, server := range servers {
			serverMap := server.(map[string]interface{})
			machineType := serverMap["machine_type"].(string)
			machineNum := int64(serverMap["machine_num"].(int))
			serverList = append(serverList, &tcaplusdb.MachineInfo{
				MachineType: &machineType,
				MachineNum:  &machineNum,
			})
		}
	}

	var proxyList []*tcaplusdb.MachineInfo
	if proxies, ok := d.Get("proxy_list").([]interface{}); ok && len(proxies) > 0 {
		for _, proxy := range proxies {
			proxyMap := proxy.(map[string]interface{})
			machineType := proxyMap["machine_type"].(string)
			machineNum := int64(proxyMap["machine_num"].(int))
			proxyList = append(proxyList, &tcaplusdb.MachineInfo{
				MachineType: &machineType,
				MachineNum:  &machineNum,
			})
		}
	}

	var clusterId string
	var inErr, outErr error

	outErr = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		clusterId, inErr = tcaplusService.CreateCluster(ctx, idlType, clusterName, vpcId, subnetId, password, resourceTags, serverList, proxyList, clusterType)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		return outErr
	}
	log.Printf("[CRITAL]%s tcaplus_cluster clusterId=%s", logId, clusterId)
	d.SetId(clusterId)
	time.Sleep(3 * time.Second)
	return resourceTencentCloudTcaplusClusterRead(d, meta)
}

func resourceTencentCloudTcaplusClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_cluster.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	clusterInfo, has, err := tcaplusService.DescribeCluster(ctx, d.Id())
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			clusterInfo, has, err = tcaplusService.DescribeCluster(ctx, d.Id())
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	_ = d.Set("idl_type", clusterInfo.IdlType)
	_ = d.Set("cluster_name", clusterInfo.ClusterName)
	_ = d.Set("vpc_id", clusterInfo.VpcId)
	_ = d.Set("subnet_id", clusterInfo.SubnetId)
	_ = d.Set("network_type", clusterInfo.NetworkType)
	_ = d.Set("create_time", clusterInfo.CreatedTime)
	_ = d.Set("password_status", clusterInfo.PasswordStatus)
	_ = d.Set("api_access_id", clusterInfo.ApiAccessId)
	_ = d.Set("api_access_ip", clusterInfo.ApiAccessIp)
	_ = d.Set("api_access_port", clusterInfo.ApiAccessPort)

	if clusterInfo.OldPasswordExpireTime == nil || *clusterInfo.OldPasswordExpireTime == "" {
		_ = d.Set("old_password_expire_time", "-")
	} else {
		_ = d.Set("old_password_expire_time", clusterInfo.OldPasswordExpireTime)
	}

	if clusterInfo.ClusterType != nil {
		_ = d.Set("cluster_type", clusterInfo.ClusterType)
	}

	if clusterInfo.ServerList != nil {
		serverList := make([]map[string]interface{}, 0, len(clusterInfo.ServerList))
		for _, server := range clusterInfo.ServerList {
			serverMap := map[string]interface{}{}
			if server.ServerUid != nil {
				serverMap["server_uid"] = *server.ServerUid
			}
			if server.MachineType != nil {
				serverMap["machine_type"] = *server.MachineType
			}
			if server.MemoryRate != nil {
				serverMap["memory_rate"] = *server.MemoryRate
			}
			if server.DiskRate != nil {
				serverMap["disk_rate"] = *server.DiskRate
			}
			if server.ReadNum != nil {
				serverMap["read_num"] = *server.ReadNum
			}
			if server.WriteNum != nil {
				serverMap["write_num"] = *server.WriteNum
			}
			if server.Version != nil {
				serverMap["version"] = *server.Version
			}
			serverList = append(serverList, serverMap)
		}
		_ = d.Set("server_list", serverList)
	}

	if clusterInfo.ProxyList != nil {
		proxyList := make([]map[string]interface{}, 0, len(clusterInfo.ProxyList))
		for _, proxy := range clusterInfo.ProxyList {
			proxyMap := map[string]interface{}{}
			if proxy.ProxyUid != nil {
				proxyMap["proxy_uid"] = *proxy.ProxyUid
			}
			if proxy.MachineType != nil {
				proxyMap["machine_type"] = *proxy.MachineType
			}
			if proxy.ProcessSpeed != nil {
				proxyMap["process_speed"] = *proxy.ProcessSpeed
			}
			if proxy.AverageProcessDelay != nil {
				proxyMap["average_process_delay"] = *proxy.AverageProcessDelay
			}
			if proxy.SlowProcessSpeed != nil {
				proxyMap["slow_process_speed"] = *proxy.SlowProcessSpeed
			}
			if proxy.Version != nil {
				proxyMap["version"] = *proxy.Version
			}
			proxyList = append(proxyList, proxyMap)
		}
		_ = d.Set("proxy_list", proxyList)
	}

	return nil
}

func resourceTencentCloudTcaplusClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_cluster.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	d.Partial(true)

	immutableArgs := []string{"cluster_type", "resource_tags", "server_list", "proxy_list"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("tcaplus_cluster argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cluster_name") {
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			err := tcaplusService.ModifyClusterName(ctx, d.Id(), d.Get("cluster_name").(string))
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("password") {
		oldPwd, newPwd := d.GetChange("password")
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			err := tcaplusService.ModifyClusterPassword(ctx, d.Id(),
				oldPwd.(string),
				newPwd.(string),
				int64(d.Get("old_password_expire_last").(int)))

			if sdkerr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "FailedOperation.OldPasswordInUse" {
					err = fmt.Errorf("[TencentCloudSDKError] Code=FailedOperation.OldPasswordInUse,`password_status` is unmodifiable now, can modify after `old_password_expire_time`")
					return resource.NonRetryableError(err)
				}
			}
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceTencentCloudTcaplusClusterRead(d, meta)
}

func resourceTencentCloudTcaplusClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_cluster.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	_, err := tcaplusService.DeleteCluster(ctx, d.Id())

	if err != nil {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			_, err = tcaplusService.DeleteCluster(ctx, d.Id())
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	_, has, err := tcaplusService.DescribeCluster(ctx, d.Id())
	if err != nil || has {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			_, has, err = tcaplusService.DescribeCluster(ctx, d.Id())
			if err != nil {
				return tccommon.RetryError(err)
			}

			if has {
				err = fmt.Errorf("delete cluster fail, cluster still exist from sdk DescribeClusters")
				return resource.RetryableError(err)
			}

			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		return nil
	} else {
		return errors.New("delete cluster fail, cluster still exist from sdk DescribeClusters")
	}
}
