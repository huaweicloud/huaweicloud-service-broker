#!/bin/bash

./broker --config=config.json

####################################################################################################################################

# Catalog
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X GET "http://username:password@localhost:3000/v2/catalog"

# Provision PostgreSQL
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PUT "http://username:password@localhost:3000/v2/service_instances/testpg?accepts_incomplete=true" -d '{"service_id":"a2c9adda-6511-462c-9934-b3fd8236e9f0","plan_id":"d42fc3cc-1341-4aa3-866e-01bc5243dc3e","organization_guid":"organization_id","space_guid":"space_id"}'

# Last Operation PostgreSQL
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X GET "http://username:password@localhost:3000/v2/service_instances/testpg/last_operation"

# Bind PostgreSQL
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PUT "http://username:password@localhost:3000/v2/service_instances/testpg/service_bindings/pg-1-binding" -d '{"service_id":"a2c9adda-6511-462c-9934-b3fd8236e9f0","plan_id":"d42fc3cc-1341-4aa3-866e-01bc5243dc3e"}'
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PUT "http://username:password@localhost:3000/v2/service_instances/testpg/service_bindings/pg-2-binding" -d '{"service_id":"a2c9adda-6511-462c-9934-b3fd8236e9f0","plan_id":"d42fc3cc-1341-4aa3-866e-01bc5243dc3e","parameters":{"dbname":"cf_manualdb_1"}}'
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PUT "http://username:password@localhost:3000/v2/service_instances/testpg/service_bindings/pg-3-binding" -d '{"service_id":"a2c9adda-6511-462c-9934-b3fd8236e9f0","plan_id":"d42fc3cc-1341-4aa3-866e-01bc5243dc3e","parameters":{"dbname":"cf_manualdb_1"}}'
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PUT "http://username:password@localhost:3000/v2/service_instances/testpg/service_bindings/pg-4-binding" -d '{"service_id":"a2c9adda-6511-462c-9934-b3fd8236e9f0","plan_id":"d42fc3cc-1341-4aa3-866e-01bc5243dc3e","parameters":{"dbname":"cf_manualdb_2"}}'

# Update PostgreSQL
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PATCH "http://username:password@localhost:3000/v2/service_instances/testpg?accepts_incomplete=true" -d '{"service_id":"a2c9adda-6511-462c-9934-b3fd8236e9f0","plan_id":"80768f31-5c2c-40e8-8135-59fe3d710dc3","previous_values":{"service_id":"a2c9adda-6511-462c-9934-b3fd8236e9f0","plan_id":"d42fc3cc-1341-4aa3-866e-01bc5243dc3e","organization_guid":"organization_id","space_guid":"space_id"},"parameters":{"apply_immediately":true}}'

# Last Operation PostgreSQL
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X GET "http://username:password@localhost:3000/v2/service_instances/testpg/last_operation"

# Unbind PostgreSQL
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X DELETE "http://username:password@localhost:3000/v2/service_instances/testpg/service_bindings/pg-4-binding?service_id=a2c9adda-6511-462c-9934-b3fd8236e9f0&plan_id=80768f31-5c2c-40e8-8135-59fe3d710dc3"
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X DELETE "http://username:password@localhost:3000/v2/service_instances/testpg/service_bindings/pg-3-binding?service_id=a2c9adda-6511-462c-9934-b3fd8236e9f0&plan_id=80768f31-5c2c-40e8-8135-59fe3d710dc3"
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X DELETE "http://username:password@localhost:3000/v2/service_instances/testpg/service_bindings/pg-2-binding?service_id=a2c9adda-6511-462c-9934-b3fd8236e9f0&plan_id=80768f31-5c2c-40e8-8135-59fe3d710dc3"
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X DELETE "http://username:password@localhost:3000/v2/service_instances/testpg/service_bindings/pg-1-binding?service_id=a2c9adda-6511-462c-9934-b3fd8236e9f0&plan_id=80768f31-5c2c-40e8-8135-59fe3d710dc3"

# Deprovision PostgreSQL
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X DELETE "http://username:password@localhost:3000/v2/service_instances/testpg?accepts_incomplete=true&service_id=a2c9adda-6511-462c-9934-b3fd8236e9f0&plan_id=80768f31-5c2c-40e8-8135-59fe3d710dc3"

# Last Operation PostgreSQL
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X GET "http://username:password@localhost:3000/v2/service_instances/testpg/last_operation"

####################################################################################################################################

# Provision Errors
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PUT "http://username:password@localhost:3000/v2/service_instances/testmy" -d '{"service_id":"ce71b484-d542-40f7-9dd4-5526e38c81ba","plan_id":"5b8282cf-a669-4ffc-b426-c169a7bbfc71","organization_guid":"organization_id","space_guid":"space_id"}'
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PUT "http://username:password@localhost:3000/v2/service_instances/testmy?accepts_incomplete=true" -d '{"service_id":"ce71b484-d542-40f7-9dd4-5526e38c81ba","plan_id":"unknown","organization_guid":"organization_id","space_guid":"space_id"}'
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PUT "http://username:password@localhost:3000/v2/service_instances/testmy?accepts_incomplete=true" -d '{"service_id":"ce71b484-d542-40f7-9dd4-5526e38c81ba","plan_id":"5b8282cf-a669-4ffc-b426-c169a7bbfc71","organization_guid":"organization_id","space_guid":"space_id","parameters":{"((("}}'

# Update Errors
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PATCH "http://username:password@localhost:3000/v2/service_instances/testmy" -d '{"service_id":"ce71b484-d542-40f7-9dd4-5526e38c81ba","plan_id":"b4943c34-fd33-47a6-8f0b-eb4f462fd746","previous_values":{"service_id":"ce71b484-d542-40f7-9dd4-5526e38c81ba","plan_id":"5b8282cf-a669-4ffc-b426-c169a7bbfc71","organization_guid":"organization_id","space_guid":"space_id"},"parameters":{"apply_immediately":true}}'
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PATCH "http://username:password@localhost:3000/v2/service_instances/testmy?accepts_incomplete=true" -d '{"service_id":"unknown","plan_id":"b4943c34-fd33-47a6-8f0b-eb4f462fd746","previous_values":{"service_id":"ce71b484-d542-40f7-9dd4-5526e38c81ba","plan_id":"5b8282cf-a669-4ffc-b426-c169a7bbfc71","organization_guid":"organization_id","space_guid":"space_id"},"parameters":{"apply_immediately":true}}'
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PATCH "http://username:password@localhost:3000/v2/service_instances/testmy?accepts_incomplete=true" -d '{"service_id":"ce71b484-d542-40f7-9dd4-5526e38c81ba","plan_id":"unknown","previous_values":{"service_id":"ce71b484-d542-40f7-9dd4-5526e38c81ba","plan_id":"5b8282cf-a669-4ffc-b426-c169a7bbfc71","organization_guid":"organization_id","space_guid":"space_id"},"parameters":{"apply_immediately":true}}'
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PATCH "http://username:password@localhost:3000/v2/service_instances/testmy?accepts_incomplete=true" -d '{"service_id":"ce71b484-d542-40f7-9dd4-5526e38c81ba","plan_id":"b4943c34-fd33-47a6-8f0b-eb4f462fd746","previous_values":{"service_id":"ce71b484-d542-40f7-9dd4-5526e38c81ba","plan_id":"5b8282cf-a669-4ffc-b426-c169a7bbfc71","organization_guid":"organization_id","space_guid":"space_id"},"parameters":{"((("}}'
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PATCH "http://username:password@localhost:3000/v2/service_instances/unknown?accepts_incomplete=true" -d '{"service_id":"ce71b484-d542-40f7-9dd4-5526e38c81ba","plan_id":"b4943c34-fd33-47a6-8f0b-eb4f462fd746","previous_values":{"service_id":"ce71b484-d542-40f7-9dd4-5526e38c81ba","plan_id":"5b8282cf-a669-4ffc-b426-c169a7bbfc71","organization_guid":"organization_id","space_guid":"space_id"},"parameters":{"apply_immediately":true}}'

# Deprovision Errors
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X DELETE "http://username:password@localhost:3000/v2/service_instances/testmy?service_id=ce71b484-d542-40f7-9dd4-5526e38c81ba&plan_id=b4943c34-fd33-47a6-8f0b-eb4f462fd746"
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X DELETE "http://username:password@localhost:3000/v2/service_instances/testmy?accepts_incomplete=true&service_id=ce71b484-d542-40f7-9dd4-5526e38c81ba&plan_id=unknown"
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X DELETE "http://username:password@localhost:3000/v2/service_instances/unknown?accepts_incomplete=true&service_id=ce71b484-d542-40f7-9dd4-5526e38c81ba&plan_id=b4943c34-fd33-47a6-8f0b-eb4f462fd746"

# Bind Errors
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PUT "http://username:password@localhost:3000/v2/service_instances/testmy/service_bindings/mysql-1-binding" -d '{"service_id":"unknown","plan_id":"5b8282cf-a669-4ffc-b426-c169a7bbfc71"}'
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PUT "http://username:password@localhost:3000/v2/service_instances/testmy/service_bindings/mysql-1-binding" -d '{"service_id":"ce71b484-d542-40f7-9dd4-5526e38c81ba","plan_id":"unknown"}'
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PUT "http://username:password@localhost:3000/v2/service_instances/testmy/service_bindings/mysql-1-binding" -d '{"service_id":"ce71b484-d542-40f7-9dd4-5526e38c81ba","plan_id":"5b8282cf-a669-4ffc-b426-c169a7bbfc71","parameters":{"((("}}'
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X PUT "http://username:password@localhost:3000/v2/service_instances/unknown/service_bindings/mysql-1-binding" -d '{"service_id":"ce71b484-d542-40f7-9dd4-5526e38c81ba","plan_id":"5b8282cf-a669-4ffc-b426-c169a7bbfc71"}'

# Unbind Errors
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X DELETE "http://username:password@localhost:3000/v2/service_instances/testmy/service_bindings/mysql-1-binding?service_id=ce71b484-d542-40f7-9dd4-5526e38c81ba&plan_id=unknown"
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X DELETE "http://username:password@localhost:3000/v2/service_instances/unknown/service_bindings/mysql-1-binding?service_id=ce71b484-d542-40f7-9dd4-5526e38c81ba&plan_id=b4943c34-fd33-47a6-8f0b-eb4f462fd746"

# Last Operation Errors
curl -H 'Accept: application/json' -H 'Content-Type: application/json' -H 'X-Broker-Api-Version: 2.x' -X GET "http://username:password@localhost:3000/v2/service_instances/unknown/last_operation"

####################################################################################################################################
