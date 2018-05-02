FROM python:3-slim

RUN mkdir /usr/src/app
WORKDIR /usr/src/app
ADD requirements.txt .
RUN pip install -r requirements.txt
ADD . .
