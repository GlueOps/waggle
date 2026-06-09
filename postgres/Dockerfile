FROM postgres:18@sha256:78481659c47e862334611ccdaf7c369c986b3046da9857112f3b309114a65fb4

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["postgres", "-c", "ssl=on", "-c", "ssl_cert_file=/var/lib/postgresql/server.crt", "-c", "ssl_key_file=/var/lib/postgresql/server.key"]
