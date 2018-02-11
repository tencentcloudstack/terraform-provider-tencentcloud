package tencentcloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	tencentCloudApiDescribeInstanceTypeConfigsParaLimitMaxFiltersNumber       = 10
	tencentCloudApiDescribeInstanceTypeConfigsParaLimitMaxFilterValuessNumber = 1
)

const (
	tencentCloudApiInstanceTypeFamilyS1  = "S1"
	tencentCloudApiInstanceTypeFamilyS2  = "S2"
	tencentCloudApiInstanceTypeFamilyS3  = "S3"
	tencentCloudApiInstanceTypeFamilySN2 = "SN2"

	tencentCloudApiInstanceTypeFamilyM1 = "M1"
	tencentCloudApiInstanceTypeFamilyM2 = "M2"

	tencentCloudApiInstanceTypeFamilyI1 = "I1"
	tencentCloudApiInstanceTypeFamilyI2 = "I2"

	tencentCloudApiInstanceTypeFamilyC2  = "C2"
	tencentCloudApiInstanceTypeFamilyC3  = "C3"
	tencentCloudApiInstanceTypeFamilyCN2 = "CN2"

	tencentCloudApiInstanceTypeFamilyFX2 = "FX2"
	tencentCloudApiInstanceTypeFamilyGA2 = "GA2"
	tencentCloudApiInstanceTypeFamilyGN2 = "GN2"
)

var (
	availableInstanceTypeFamilies = []string{
		tencentCloudApiInstanceTypeFamilyS1,
		tencentCloudApiInstanceTypeFamilyS2,
		tencentCloudApiInstanceTypeFamilyS3,
		tencentCloudApiInstanceTypeFamilySN2,

		tencentCloudApiInstanceTypeFamilyM1,
		tencentCloudApiInstanceTypeFamilyM2,

		tencentCloudApiInstanceTypeFamilyI1,
		tencentCloudApiInstanceTypeFamilyI2,

		tencentCloudApiInstanceTypeFamilyC2,
		tencentCloudApiInstanceTypeFamilyC3,
		tencentCloudApiInstanceTypeFamilyCN2,

		tencentCloudApiInstanceTypeFamilyFX2,
		tencentCloudApiInstanceTypeFamilyGA2,
		tencentCloudApiInstanceTypeFamilyGN2,
	}
)

func dataSourceInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceTencentCloudFiltersSchema(),

			"cpu_core_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"memory_size": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},

			// Computed values.
			"instance_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_core_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"family": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*TencentCloudClient).commonConn
	filters, filtersOk := d.GetOk("filter")
	cpuCoreCount, cpuCoreCountOk := d.GetOk("cpu_core_count")
	memorySizeCount, memorySizeCountOk := d.GetOk("memory_size")
	if !filtersOk && !cpuCoreCountOk && !memorySizeCountOk {
		return fmt.Errorf("One of filter, cpu_core_count, memory_size must be assigned")
	}
	params := map[string]string{
		"Version": "2017-03-12",
		"Action":  "DescribeInstanceTypeConfigs",
	}

	if filtersOk {
		filterList := filters.(*schema.Set)
		err := buildFiltersParam(
			params,
			filterList,
			tencentCloudApiDescribeInstanceTypeConfigsParaLimitMaxFiltersNumber,
			tencentCloudApiDescribeInstanceTypeConfigsParaLimitMaxFilterValuessNumber,
		)
		if err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] tencentcloud_instance_types - param: %v", params)
	response, err := client.SendRequest("cvm", params)
	if err != nil {
		return err
	}

	type InstanceTypeConfig struct {
		Zone           string  `json:"Zone"`
		InstanceFamily string  `json:"InstanceFamily"`
		InstanceType   string  `json:"InstanceType"`
		CPU            int     `json:"CPU"`
		GPU            int     `json:"GPU"`
		FPGA           int     `json:"FPGA"`
		Memory         float64 `json:"Memory"`
	}
	var jsonresp struct {
		Response struct {
			Error struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			}
			RequestId             string `json:"RequestId"`
			InstanceTypeConfigSet []InstanceTypeConfig
		}
	}

	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return err
	}
	if jsonresp.Response.Error.Code != "" {
		return fmt.Errorf(
			"tencentcloud_instance_types got error, code:%v, message:%v",
			jsonresp.Response.Error.Code,
			jsonresp.Response.Error.Message,
		)
	}

	var (
		instanceTypeConfigList []InstanceTypeConfig
	)
	instanceConfigList := jsonresp.Response.InstanceTypeConfigSet
	if len(instanceConfigList) == 0 {
		return errors.New("No instance types found")
	}

	for _, instanceConfig := range instanceConfigList {
		cpu := instanceConfig.CPU
		fpga := instanceConfig.FPGA
		gpu := instanceConfig.GPU
		instanceFamily := instanceConfig.InstanceFamily
		instanceType := instanceConfig.InstanceType
		mem := instanceConfig.Memory
		zone := instanceConfig.Zone
		log.Printf(
			"[DEBUG] tencentcloud_instance_types - InstanceType found zone: %v, cpu:% v, mem: %v, fpga: %v, gpu: %v, instanceFamily: %v, instanceType:%v",
			zone,
			cpu,
			mem,
			fpga,
			gpu,
			instanceFamily,
			instanceType,
		)

		if cpuCoreCountOk && memorySizeCountOk {
			if cpu == cpuCoreCount && mem == memorySizeCount {
				instanceTypeConfigList = append(instanceTypeConfigList, instanceConfig)
			}
			continue
		}
		if cpuCoreCountOk {
			if cpu == cpuCoreCount {
				instanceTypeConfigList = append(instanceTypeConfigList, instanceConfig)
			}
			continue
		}
		if memorySizeCountOk {
			if mem == memorySizeCount {
				instanceTypeConfigList = append(instanceTypeConfigList, instanceConfig)
			}
			continue
		}

		instanceTypeConfigList = append(instanceTypeConfigList, instanceConfig)
	}

	if len(instanceTypeConfigList) == 0 {
		return errors.New("No instance types found")
	}

	var (
		result    []map[string]interface{}
		resultIds []string
	)

	for _, instanceTypeConfig := range instanceTypeConfigList {
		m := make(map[string]interface{})
		m["availability_zone"] = instanceTypeConfig.Zone
		m["instance_type"] = instanceTypeConfig.InstanceType
		m["cpu_core_count"] = instanceTypeConfig.CPU
		m["memory_size"] = instanceTypeConfig.Memory
		m["family"] = instanceTypeConfig.InstanceFamily
		result = append(result, m)
		resultIds = append(resultIds, instanceTypeConfig.InstanceType)
	}
	id := dataResourceIdsHash(resultIds)
	d.SetId(id)
	log.Printf("[DEBUG] tencentcloud_instance_types - instances[0]: %#v", result[0])
	if err := d.Set("instance_types", result); err != nil {
		return err
	}
	return nil
}
