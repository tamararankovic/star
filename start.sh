export STAR_HOSTNAME="star"
export MAGNETAR_HOSTNAME="magnetar"

export NATS_HOSTNAME="nats"
export NATS_PORT=4222

export ETCD_HOSTNAME="etcd"
export ETCD_PORT=2379

export REGISTRATION_SUBJECT="register"
export REGISTRATION_REQ_TIMEOUT_MILLISECONDS=1000
export MAX_REGISTRATION_RETRIES=5

export NODE_ID_DIR_PATH="/etc/c12s"
export NODE_ID_FILE_NAME="nodeid"

docker-compose up --build