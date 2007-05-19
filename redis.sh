MACHINE_IP=45.32.237.247
REDIS_PASS=newpasslalala

redis-install() {
    sudo apt-get install redis-server

    replace="bind 127.0.0.1 ::1"
    with="bind $MACHINE_IP ::1"
    sudo sed -i -e "s|$replace|$with|g" /etc/redis/redis.conf

    replace="# requirepass foobared"
    with="requirepass $REDIS_PASS"
    sudo sed -i -e "s|$replace|$with|g" /etc/redis/redis.conf
}

redis-restart() {
    sudo systemctl restart redis
}
