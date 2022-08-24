FROM alpine:3.14

WORKDIR /app

COPY ./bin/distapp .

EXPOSE $APP_PORT
EXPOSE $GOSSIP_PORT

CMD if [[ -z "$GOSSIP_LEADER" ]]; then /app/distapp --app-port $APP_PORT --gossip-port $GOSSIP_PORT; else /app/distapp --app-port $APP_PORT --gossip-port $GOSSIP_PORT --gossip-leader $GOSSIP_LEADER; fi