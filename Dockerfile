FROM alpine:latest

WORKDIR /app
COPY app .

RUN chmod +x app

CMD ["./app"]
