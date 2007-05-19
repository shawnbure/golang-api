DB_NAME=erdsea_db

PSQL_ADMIN_PSW=some_pass
PSQL_USER=user
PSQL_PSW=password

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
  sudo -u postgres createdb -O $PSQL_USER $DB_NAME

  # edit conf

  sudo ufw allow from any to any port 5432 proto tcp
}
