#!/bin/bash

huaweicloud-service-broker -config=config.json -port=3000

####################################################################################################################################

# Catalog
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X GET "http://username:password@localhost:3000/v2/catalog"

# Provision RDS MySQL 5.6.39
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PUT "http://username:password@localhost:3000/v2/service_instances/testrdsmysql?accepts_incomplete=true" -d '{"service_id":"275f3e0b-86fd-4303-946c-171374d29150","plan_id":"fc1b0ebf-aabf-4ab1-87e9-83544a5902e8","organization_guid":"organization_id","space_guid":"space_id","parameters":{"name":"testrdsmysql", "database_password": "Mysql1234!"}}}'

# Last Operation RDS MySQL 5.6.39
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X GET "http://username:password@localhost:3000/v2/service_instances/testrdsmysql/last_operation"

# Update RDS MySQL 5.6.39
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PATCH "http://username:password@localhost:3000/v2/service_instances/testrdsmysql?accepts_incomplete=true" -d '{"service_id":"275f3e0b-86fd-4303-946c-171374d29150","plan_id":"fc1b0ebf-aabf-4ab1-87e9-83544a5902e8","previous_values":{"service_id":"275f3e0b-86fd-4303-946c-171374d29150","plan_id":"fc1b0ebf-aabf-4ab1-87e9-83544a5902e8","organization_guid":"organization_id","space_guid":"space_id"},"parameters":{"volume_size":400}}'

# Last Operation RDS MySQL 5.6.39
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X GET "http://username:password@localhost:3000/v2/service_instances/testrdsmysql/last_operation"

# Bind RDS MySQL 5.6.39
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PUT "http://username:password@localhost:3000/v2/service_instances/testrdsmysql/service_bindings/testrdsmysqlbinding" -d '{"service_id":"275f3e0b-86fd-4303-946c-171374d29150","plan_id":"fc1b0ebf-aabf-4ab1-87e9-83544a5902e8"}'

# Unbind RDS MySQL 5.6.39
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X DELETE "http://username:password@localhost:3000/v2/service_instances/testrdsmysql/service_bindings/testrdsmysqlbinding?service_id=275f3e0b-86fd-4303-946c-171374d29150&plan_id=fc1b0ebf-aabf-4ab1-87e9-83544a5902e8"

# Deprovision RDS MySQL 5.6.39
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X DELETE "http://username:password@localhost:3000/v2/service_instances/testrdsmysql?accepts_incomplete=true&service_id=275f3e0b-86fd-4303-946c-171374d29150&plan_id=fc1b0ebf-aabf-4ab1-87e9-83544a5902e8"

# Last Operation RDS MySQL 5.6.39
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X GET "http://username:password@localhost:3000/v2/service_instances/testrdsmysql/last_operation"
