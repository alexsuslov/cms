FROM debian:latest
RUN apt update && \
  apt -y upgrade && \
  apt install -y  ca-certificates

RUN mkdir /app
WORKDIR /app
COPY bin/cms .
RUN useradd  -ms /bin/bash app
RUN chown -R app:app /app
VOLUME /app/vol

EXPOSE 80
USER app

CMD ["/app/server"]


