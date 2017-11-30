#!/usr/bin/env bash

# Export variables from .env file
source ./.env
export $(cut -d= -f1 ./.env)

export OLD_CHANNEL_COUNT=$(docker exec -it example_rabbitmq_1 rabbitmqctl list_queues | grep 'my-channel' | awk 'NR >= 1 {print $2}' | sed 's/[^0-9]*//g')
export OLD_CHANNEL_DELAYED_COUNT=$(docker exec -it example_rabbitmq_1 rabbitmqctl list_queues | grep 'my-delayed-channel' | awk 'NR >= 1 {print $2}' | sed 's/[^0-9]*//g')

docker exec -it example_postgresql_1 psql -U postgres -c "SELECT pg_notify('not-existing-channel', '{\"variable\": \"value\"}');"
docker exec -it example_postgresql_1 psql -U postgres -c "SELECT pg_notify('my-channel', '{\"variable\": \"value\"}');"
docker exec -it example_postgresql_1 psql -U postgres -c "SELECT pg_notify('my-channel', '{}');" #Empty payload
docker exec -it example_postgresql_1 psql -U postgres -c "SELECT pg_notify('my-delayed-channel', '{\"variable\": \"value\", \"delay\": 1}');"

export NEW_CHANNEL_COUNT=$(docker exec -it example_rabbitmq_1 rabbitmqctl list_queues | grep 'my-channel' | awk 'NR >= 1 {print $2} '  | sed 's/[^0-9]*//g')
export NEW_CHANNEL_DELAYED_COUNT=$(docker exec -it example_rabbitmq_1 rabbitmqctl list_queues | grep 'my-delayed-channel' | awk 'NR >= 1 {print $2}'  | sed 's/[^0-9]*//g')

export EXPECTED_CHANNEL_COUNT=$(($OLD_CHANNEL_COUNT + 2))
export EXPECTED_CHANNEL_DELAYED_COUNT=$(($OLD_CHANNEL_DELAYED_COUNT + 1))

echo "OLD_CHANNEL_COUNT: $OLD_CHANNEL_COUNT"
echo "EXPECTED_CHANNEL_COUNT: $EXPECTED_CHANNEL_COUNT"
echo "NEW_CHANNEL_COUNT: $NEW_CHANNEL_COUNT"
if [ "$EXPECTED_CHANNEL_COUNT" != "$NEW_CHANNEL_COUNT" ]; then
    exit 1
fi

echo "OLD_CHANNEL_DELAYED_COUNT: $OLD_CHANNEL_DELAYED_COUNT"
echo "EXPECTED_CHANNEL_DELAYED_COUNT: $EXPECTED_CHANNEL_DELAYED_COUNT"
echo "NEW_CHANNEL_DELAYED_COUNT: $NEW_CHANNEL_DELAYED_COUNT"
if [ "$EXPECTED_CHANNEL_DELAYED_COUNT" != "$NEW_CHANNEL_DELAYED_COUNT" ]; then
    docker-compose -f docker-compose.yml logs
    exit 1
fi

exit 0