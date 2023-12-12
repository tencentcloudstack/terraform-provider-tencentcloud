Provides a resource to create a mariadb dedicatedcluster_db_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_dedicatedcluster_db_instance" "dedicatedcluster_db_instance" {
  goods_num 	= 1
  memory 		= 2
  storage 		= 10
  cluster_id 	= "dbdc-24odnuhr"
  vpc_id 		= "vpc-ii1jfbhl"
  subnet_id 	= "subnet-3ku415by"
  db_version_id = "8.0"
  instance_name = "cluster-mariadb-test-1"
}

```
Import

mariadb dedicatedcluster_db_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_mariadb_dedicatedcluster_db_instance.dedicatedcluster_db_instance tdsql-050g3fmv
```