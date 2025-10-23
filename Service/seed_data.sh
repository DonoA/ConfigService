curl -X POST http://localhost:8080/configs -H "Content-Type: application/json" -d '{
   "config": {
        "service": "example-service",
        "name": "example-config1",
        "defaultValue": "example-value1",
        "type": "string"
   }
 }'

echo ""

 curl -X POST http://localhost:8080/configs -H "Content-Type: application/json" -d '{
   "config": {
        "service": "example-service",
        "name": "example-config2",
        "defaultValue": "example-value2",
        "type": "string"
   }
 }'

echo ""

 curl -X POST http://localhost:8080/configs/example-service/example-config1/overrides -H "Content-Type: application/json" -d '{
   "override": {
        "entityType": "user",
        "entityId": "user123",
        "value": "user123-value1"
   }
 }'

echo ""

curl -X POST http://localhost:8080/configs -H "Content-Type: application/json" -d '{
   "config": {
        "service": "example-service2",
        "name": "example-config",
        "defaultValue": "example-value",
        "type": "string"
   }
 }'

echo ""
