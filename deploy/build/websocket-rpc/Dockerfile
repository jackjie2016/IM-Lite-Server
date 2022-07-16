FROM showurl/zerobase

WORKDIR /app

COPY ./bin /app/websocket

RUN chmod +x /app/websocket && mkdir /app/etc

VOLUME /app/etc/websocket.yaml

CMD ["/app/websocket"]
