#!/bin/sh

docker-compose exec api curl -X 'POST' 'http://127.0.0.1:8000/operations' -H 'accept: application/json' \
-H 'Content-Type: application/json' \
-d '{
"description": "COMPRA A VISTA",
"debit": "true"
}'

docker-compose exec api curl -X 'POST' 'http://127.0.0.1:8000/operations' -H 'accept: application/json' \
-H 'Content-Type: application/json' \
-d '{
"description": "COMPRA PARCELADA",
"debit": "true"
}'

docker-compose exec api curl -X 'POST' 'http://127.0.0.1:8000/operations' -H 'accept: application/json' \
-H 'Content-Type: application/json' \
-d '{
"description": "SAQUE",
"debit": "true"
}'

docker-compose exec api curl -X 'POST' 'http://127.0.0.1:8000/operations' -H 'accept: application/json' \
-H 'Content-Type: application/json' \
-d '{
"description": "PAGAMENTO",
"debit": "false"
}'