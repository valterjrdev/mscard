#!/bin/sh

docker-compose exec api curl -X 'POST' 'http://127.0.0.1:8000/operation-types' -H 'accept: application/json' \
-H 'Content-Type: application/json' \
-d '{
"description": "COMPRA A VISTA",
"negative": "true"
}'

docker-compose exec api curl -X 'POST' 'http://127.0.0.1:8000/operation-types' -H 'accept: application/json' \
-H 'Content-Type: application/json' \
-d '{
"description": "COMPRA PARCELADA",
"negative": "true"
}'

docker-compose exec api curl -X 'POST' 'http://127.0.0.1:8000/operation-types' -H 'accept: application/json' \
-H 'Content-Type: application/json' \
-d '{
"description": "SAQUE",
"negative": "true"
}'

docker-compose exec api curl -X 'POST' 'http://127.0.0.1:8000/operation-types' -H 'accept: application/json' \
-H 'Content-Type: application/json' \
-d '{
"description": "PAGAMENTO",
"negative": "false"
}'