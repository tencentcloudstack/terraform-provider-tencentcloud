package pts_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcpts "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/pts"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// go test -i; go test -test.run TestAccTencentCloudPtsScenarioResource_basic -v
func TestAccTencentCloudPtsScenarioResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckPtsScenarioDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsScenario,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPtsScenarioExists("tencentcloud_pts_scenario.scenario"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "name", "pts-js"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "type", "pts-js"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.geo_regions_load_distribution.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.geo_regions_load_distribution.0.percentage", "100"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.geo_regions_load_distribution.0.region", "ap-guangzhou"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.geo_regions_load_distribution.0.region_id", "1"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.load_spec.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.load_spec.0.concurrency.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.load_spec.0.concurrency.0.graceful_stop_seconds", "3"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.load_spec.0.concurrency.0.iteration_count", "0"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.load_spec.0.concurrency.0.max_requests_per_second", "0"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.load_spec.0.concurrency.0.stages.#", "4"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.load_spec.0.concurrency.0.stages.0.duration_seconds", "120"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.load_spec.0.concurrency.0.stages.0.target_virtual_users", "2"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.load_spec.0.concurrency.0.stages.1.duration_seconds", "120"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.load_spec.0.concurrency.0.stages.1.target_virtual_users", "4"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.load_spec.0.concurrency.0.stages.2.duration_seconds", "120"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.load_spec.0.concurrency.0.stages.2.target_virtual_users", "5"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.load_spec.0.concurrency.0.stages.3.duration_seconds", "240"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "load.0.load_spec.0.concurrency.0.stages.3.target_virtual_users", "5"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "test_scripts.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_pts_scenario.scenario", "test_scripts.0.encoded_content"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "test_scripts.0.load_weight", "100"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "test_scripts.0.name", "script.js"),
					resource.TestCheckResourceAttrSet("tencentcloud_pts_scenario.scenario", "test_scripts.0.size"),
					resource.TestCheckResourceAttr("tencentcloud_pts_scenario.scenario", "test_scripts.0.type", "js"),
					resource.TestCheckResourceAttrSet("tencentcloud_pts_scenario.scenario", "test_scripts.0.updated_at"),
				),
			},
			{
				ResourceName:      "tencentcloud_pts_scenario.scenario",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPtsScenarioDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcpts.NewPtsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_pts_scenario" {
			continue
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		projectId := idSplit[0]
		scenarioId := idSplit[1]

		scenario, err := service.DescribePtsScenario(ctx, projectId, scenarioId)
		if scenario != nil {
			return fmt.Errorf("pts scenario %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckPtsScenarioExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		projectId := idSplit[0]
		scenarioId := idSplit[1]

		service := svcpts.NewPtsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		scenario, err := service.DescribePtsScenario(ctx, projectId, scenarioId)
		if scenario == nil {
			return fmt.Errorf("pts scenario %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccPtsScenario = testAccPtsProject + `

resource "tencentcloud_pts_scenario" "scenario" {
    name            = "pts-js"
    project_id      = tencentcloud_pts_project.project.id
    type            = "pts-js"

    domain_name_config {
    }

    load {
        geo_regions_load_distribution {
            percentage = 100
            region     = "ap-guangzhou"
            region_id  = 1
        }

        load_spec {
            concurrency {
                graceful_stop_seconds   = 3
                iteration_count         = 0
                max_requests_per_second = 0

                stages {
                    duration_seconds     = 120
                    target_virtual_users = 2
                }
                stages {
                    duration_seconds     = 120
                    target_virtual_users = 4
                }
                stages {
                    duration_seconds     = 120
                    target_virtual_users = 5
                }
                stages {
                    duration_seconds     = 240
                    target_virtual_users = 5
                }
            }
        }
    }

    sla_policy {
    }

    test_scripts {
        encoded_content = <<-EOT
            // Send a http get request
            import http from 'pts/http';
            import { check, sleep } from 'pts';

            export default function () {
              // simple get request
              const resp1 = http.get('http://httpbin.org/get');
              console.log(resp1.body);
              // if resp1.body is a json string, resp1.json() transfer json format body to a json object
              console.log(resp1.json());
              check('status is 200', () => resp1.statusCode === 200);

              // sleep 1 second
              sleep(1);

              // get request with headers and parameters
              const resp2 = http.get('http://httpbin.org/get', {
                headers: {
                  Connection: 'keep-alive',
                  'User-Agent': 'pts-engine',
                },
                query: {
                  name1: 'value1',
                  name2: 'value2',
                },
              });

              console.log(resp2.json().args.name1); // 'value1'
              check('body.args.name1 equals value1', () => resp2.json().args.name1 === 'value1');
            }
        EOT
        load_weight     = 100
        name            = "script.js"
        size            = 838
        type            = "js"
        updated_at      = "2022-11-11T16:18:37+08:00"
    }
}

`
