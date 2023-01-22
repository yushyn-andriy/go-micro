FROM alpine:latest


RUN mkdir app
WORKDIR /app

COPY listenerApp /app/

CMD ["/app/listenerApp"]
