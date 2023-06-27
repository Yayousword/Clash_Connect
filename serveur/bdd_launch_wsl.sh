#!/bin/bash

if psql --version &>/dev/null; then
    echo "already installed"
else
    sudo apt-get update
    sudo apt-get install postgresql postgresql-contrib
fi


# Create a new PostgreSQL user and database
sudo -u postgres psql -c "CREATE USER clash_connect WITH PASSWORD '52fdc5a882ad0cc490297a43dce208cc36639f0c5224fc47bc849a978bd16d98';"
sudo -u postgres psql -c "CREATE DATABASE clash_connect;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE clash_connect TO clash_connect;"
sudo -u postgres psql -c "\c clash_connect"

sudo -u postgres psql -d clash_connect -c "CREATE TABLE authentication (
    id VARCHAR(100) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(300) NOT NULL
); GRANT SELECT, INSERT, UPDATE, DELETE ON authentication TO clash_connect;"

sudo -u postgres psql -d clash_connect -c "CREATE TABLE publication (
    id VARCHAR(255),
    time VARCHAR(255),
    rate INTEGER NOT NULL,
    nombre_rate INTEGER NOT NULL,
    PRIMARY KEY (id,time)
); GRANT SELECT, INSERT, UPDATE, DELETE ON publication TO clash_connect;"

sudo -u postgres psql -d clash_connect -c "CREATE TABLE comment (
    id_sender VARCHAR(255) ,
    id_publication VARCHAR(255),
    time_publication VARCHAR(255),
    time_comment TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    comment TEXT NOT NULL,
    PRIMARY KEY (id_sender,id_publication,time_publication,time_comment)
); GRANT SELECT, INSERT, UPDATE, DELETE ON comment TO clash_connect;"

sudo -u postgres psql -d clash_connect -c "CREATE TABLE chat (
    id_sender VARCHAR(255) ,
    time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    message TEXT NOT NULL,
    PRIMARY KEY (id_sender, time)
); GRANT SELECT, INSERT, UPDATE, DELETE ON chat TO clash_connect;"

# Start the PostgreSQL service
sudo /etc/init.d/postgresql start

#sudo systemctl restart postgresql.service
