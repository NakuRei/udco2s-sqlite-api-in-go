FROM golang:1.22-alpine3.19 as builder

ARG USERNAME=user

RUN apk add --no-cache build-base sqlite-dev

RUN addgroup $USERNAME \
  && adduser -D -G $USERNAME -h /home/$USERNAME $USERNAME

WORKDIR /home/$USERNAME/workspace

RUN mkdir -p db \
  && chown -R $USERNAME:$USERNAME db

USER $USERNAME

ENTRYPOINT ["bash"]
