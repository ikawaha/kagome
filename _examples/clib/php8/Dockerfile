FROM php:latest

RUN apt-get update \
    && apt-get install -y libffi-dev

RUN docker-php-ext-install ffi

WORKDIR /app
