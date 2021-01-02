#!/usr/bin/env bash

SQL_DUMP_FOLDER="/mnt/hgfs/progdev/archive/c/sql/"

docker run -p 3306:3306 -v ${SQL_DUMP_FOLDER}:/docker-entrypoint-initdb.d/ -e MYSQL_ROOT_PASSWORD=mariadbrootpassword -e MYSQL_DATABASE=asagi  mariadb
