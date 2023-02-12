FROM alpine:3.17

WORKDIR /app

# copy compiled script
COPY cloudflare-ddns /app

ENTRYPOINT [ "/app/cloudflare-ddns"]
