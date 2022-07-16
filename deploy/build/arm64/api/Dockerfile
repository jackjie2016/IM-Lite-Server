FROM showurl/zerobase

WORKDIR /app

COPY ./bin /app/api

RUN chmod +x /app/api && mkdir /app/etc

VOLUME /app/etc/api.yaml

CMD ["/app/api"]
