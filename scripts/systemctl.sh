HOME="/root"
USER="root"
API_WORKDIR=$HOME/erdsea/api
CMD_FLAGS="--log-level=*:DEBUG"

start_api() {
  sudo systemctl start erdsea-api.service
}

stop_api() {
  sudo systemctl stop erdsea-api.service
}

systemd_api()  {
  echo "[Unit]
  Description=Erdsea API
  After=network-online.target
  
  [Service]
  User=$USER

  WorkingDirectory=$API_WORKDIR
  ExecStart=$API_WORKDIR/api/cmd $CMD_FLAGS
  StandardOutput=journal
  StandardError=journal
  Restart=always
  RestartSec=3
  LimitNOFILE=65535
  
  [Install]
  WantedBy=multi-user.target" > erdsea-api.service

  sudo mv erdsea-api.service /etc/systemd/system/

  sudo systemctl enable erdsea-api.service
}

