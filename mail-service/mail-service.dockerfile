FROM alpine:latest


RUN mkdir app
WORKDIR /app

COPY mailApp /app/
COPY templates/ /app/templates

CMD ["/app/mailApp"]
