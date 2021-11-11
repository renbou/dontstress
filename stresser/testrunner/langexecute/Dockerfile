FROM langbase

WORKDIR /usr/local/app

COPY app.sh .
COPY scripts/ scripts

RUN chmod +x app.sh scripts/*

ENTRYPOINT [ "bash", "app.sh" ]