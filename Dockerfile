FROM php:8.1-cli-alpine

RUN apk add --no-cache bash

COPY parse_cobertura.php /app/parse_cobertura.php

WORKDIR /app

RUN php /app/parse_cobertura.php
