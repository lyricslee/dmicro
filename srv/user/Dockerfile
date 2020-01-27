#FROM 192.168.56.190/library/scratch:latest
FROM scratch
WORKDIR /
COPY user .
COPY conf conf

ENTRYPOINT ["/user"]