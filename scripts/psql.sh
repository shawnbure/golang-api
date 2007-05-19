MACHINE_IP="123.123.123"

PSQL_DB=erdsea_db
PSQL_ADMIN_PSW=some_pass
PSQL_USER=user
PSQL_PSW=password

PSQL_ADDR_REPLACE="#listen_addresses = 'localhost'"
PSQL_ADDR_WITH="listen_addresses = '$MACHINE_IP'"
PSQL_CLIENT_CONNECT_OPT="host        all         all             0.0.0.0/0               md5"

function install() {
  if ! [ -x "$(command psql -V)" ]; then
    echo "postgresql is not installed on your system. installing"

    sudo apt update && sudo apt postgresql postgresql-client

    sudo systemctl start postgresql.service && sudo systemctl enable postgresql.service

    sudo -u postgres bash -c "psql -c \"alter user postgres with password '$PSQL_ADMIN_PSW';\""

  else
    echo "postgresql already installed. skipping"

  fi
}

function setup() {
  echo "creating database with user: '$PSQL_USER' and password: '$PSQL_PSW'"

  sudo -u postgres bash -c "psql -c \"create user '$PSQL_USER' with password '$PSQL_PSW';\""
  sudo -u postgres createdb -O $PSQL_USER $PSQL_DB

  sudo sed -i -e "s|$PSQL_ADDR_REPLACE|$PSQL_ADDR_WITH|g" /etc/postgresql/12/main/postgresql.conf

  sudo systemctl restart postgresql.service

  sudo bash -c "echo -e \"\n$PSQL_CLIENT_CONNECT_OPT\" >>/etc/postgresql/12/main/pg_hba.conf"

  sudo ufw allow from any to any port 5432 proto tcp
}
