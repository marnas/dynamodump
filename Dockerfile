# GitHub:       https://github.com/AltoStack/dynamodump
# Twitter:      https://twitter.com/AltoStack
# Website:      https://altostack.io/

FROM golang:alpine AS build

ENV GOOS=linux

WORKDIR /go/src/github.com/AltoStack/dynamodump
COPY . /go/src/github.com/AltoStack/dynamodump/

RUN apk update \
    && apk add --no-cache gcc g++ git make

RUN make install

# ---

FROM golang:alpine

COPY --from=build /go/bin/dynamodump /usr/bin/dynamodump

ENTRYPOINT ["dynamodump"]
CMD ["--help"]