FROM alpine:latest

COPY target/sheazuzu-linux-amd64 /

ENV mysql_host ""
ENV mysql_port 0
ENV mysql_database_name ""
ENV mysql_username ""
ENV mysql_password ""

EXPOSE 8080

ENTRYPOINT /sheazuzu-linux-amd64 \
    --database.endpoint=$mysql_host \
    --database.username=$mysql_username \
    --database.password=$mysql_password \
    --database.name=$mysql_database_name \
    --database.config=utf8&parseTime=True&loc=Local \
    --database.port=$mysql_port



