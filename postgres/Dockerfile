FROM postgres:18@sha256:4aabea78cf39b90e834caf3af7d602a18565f6fe2508705c8d01aa63245c2e20

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["postgres", "-c", "ssl=on", "-c", "ssl_cert_file=/var/lib/postgresql/server.crt", "-c", "ssl_key_file=/var/lib/postgresql/server.key"]
