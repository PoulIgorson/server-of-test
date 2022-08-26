FROM node:lts as node-builder

# Install OS libraries
RUN apt-get update \
  && apt-get install -y gcc

WORKDIR /app
COPY . /app
WORKDIR /app/svelte_components

RUN npm install
RUN npm rebuild node-sass
RUN npm run build

FROM python:latest as python-builder

# Keeps Python from generating .pyc files in the container
ENV PYTHONDONTWRITEBYTECODE=1

# Turns off buffering for easier container logging
ENV PYTHONUNBUFFERED=1

COPY --from=node-builder /app /app

WORKDIR /app

# Install pip requirements
COPY requirements.txt .
RUN python -m pip install --upgrade pip && pip install -r requirements.txt

EXPOSE 80

COPY release-tasks.sh /app/release-tasks.sh

ENTRYPOINT ["/app/release-tasks.sh"]