FROM langbase

WORKDIR /usr/local/app

COPY app.sh .
COPY tests tests/

RUN chmod +x app.sh tests/*

ENTRYPOINT [ "bash", "app.sh" ]