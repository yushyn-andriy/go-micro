FROM alpine:latest


RUN mkdir app
WORKDIR /app

COPY mailApp /app/

CMD ["/app/mailApp"]
