FROM alpine:3.10

COPY ./parse_cobertura.php parse_cobertura.php

ENTRYPOINT ["/parse_cobertura.php"]
