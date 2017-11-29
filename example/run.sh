#!/usr/bin/env bash

# Export variables from .env file
source ./.env
export $(cut -d= -f1 ./.env)

docker-compose -f docker-compose.yml stop
docker-compose -f docker-compose.yml rm -f
docker volume ls -q | xargs docker volume rm
docker-compose -f docker-compose.yml up --build -d

sleep 3
docker exec -it example_rabbitmq_1 rabbitmqctl wait /var/lib/rabbitmq/mnesia/rabbit@rabbitmq.pid

docker exec -it example_rabbitmq_1 rabbitmqadmin declare exchange name=${DELAY_EXCHANGE_NAME} type=x-delayed-message arguments='{"x-delayed-type":"direct"}'
docker exec -it example_rabbitmq_1 rabbitmqadmin declare queue name=my-delayed-channel durable=false
docker exec -it example_rabbitmq_1 rabbitmqadmin declare binding routing_key=my-delayed-channel source=${DELAY_EXCHANGE_NAME} destination_type="queue" destination="my-delayed-channel"

[ -z "$DEFAULT_EXCHANGE_NAME" ] || docker exec -it example_rabbitmq_1 rabbitmqadmin declare exchange name=${DEFAULT_EXCHANGE_NAME} type=direct
docker exec -it example_rabbitmq_1 rabbitmqadmin declare queue name=my-channel durable=false
docker exec -it example_rabbitmq_1 rabbitmqadmin declare binding routing_key=my-channel source=${DEFAULT_EXCHANGE_NAME} destination=my-channel

# docker-compose -f docker-compose.yml logs -tf
