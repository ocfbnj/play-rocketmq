#!/bin/bash

RMQ_CONTAINER=$(docker ps -a|awk '/rmq-/ {print $1}')
if [[ -n "$RMQ_CONTAINER" ]]; then
   echo "Removing RocketMQ Container..."
   docker rm -fv $RMQ_CONTAINER
   # Wait till the existing containers are removed
   sleep 5
fi

prepare_dir()
{
    dirs=("logs/namesrv" "logs/broker0" "logs/broker1" "store/broker0" "store/broker1")

    for dir in ${dirs[@]}
    do
        if [ ! -d "`pwd`/${dir}" ]; then
            mkdir -p "`pwd`/${dir}"
            chmod a+rw "`pwd`/${dir}"
        fi
    done
}

prepare_dir

# Run nameserver and broker
docker-compose -f ./docker-compose.yml up -d
