FROM python:alpine

RUN pip3 install localstripe

EXPOSE 8420

ENTRYPOINT ["localstripe"]