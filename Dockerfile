FROM mysql:8.0

ENV MYSQL_ROOT_PASSWORD=can_we_still_score_for_dop

RUN mkdir -p /docker-entrypoint-initdb.d
COPY record_db.sql /docker-entrypoint-initdb.d/

EXPOSE 3306
