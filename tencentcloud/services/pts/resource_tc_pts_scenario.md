Provides a resource to create a pts scenario

Example Usage

```hcl
resource "tencentcloud_pts_scenario" "scenario" {
    name            = "pts-js"
    project_id      = "project-45vw7v82"
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

```
Import

pts scenario can be imported using the id, e.g.
```
$ terraform import tencentcloud_pts_scenario.scenario scenario_id
```