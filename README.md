# postrabbit

>

## About

A bridge between Postgres and RabbitMQ writing on GO [![Build Status](https://travis-ci.org/andreychuk/postrabbit.svg?branch=master)](https://travis-ci.org/andreychuk/postrabbit)


## Getting Started

1. Make sure you have [GO](https://golang.org) and [DEP](https://github.com/golang/dep) installed.

2. Install your dependencies
    ```
    dep ensure -vendor-only

    ```
3. Build your app
    ```
    go build

    ```

4. Run your app
    ```
    ./postrabbit
    ```

## Environment Variables

    CHANNEL_LIST - tets,test1,test2
    RABBITMQ_URL - RabbitMQ url
    POSTGRES_URL - Postgres url
    DEFAULT_EXCHANGE_NAME - exchange name for simple messages
    DELAY_EXCHANGE_NAME - exchange name for delay messages

## License

Copyright (c) 2017

Licensed under the [MIT license](LICENSE).