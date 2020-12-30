FROM nginx:latest as base

COPY ./client/assets /www/assets


FROM base as dev

COPY ./env/nginx/dev.conf /etc/nginx/conf.d/default.conf


FROM base as prod

COPY ./env/nginx/prod.conf /etc/nginx/conf.d/default.conf
