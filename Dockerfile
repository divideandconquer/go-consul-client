FROM scratch

ADD /bin/app /app

ENTRYPOINT ["/app"]