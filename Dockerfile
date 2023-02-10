FROM scratch

WORKDIR /app

# copy compiled script
COPY cloudflare-ddns /app

ENTRYPOINT [ "/app/cloudflare-ddns"]
