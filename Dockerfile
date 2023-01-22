FROM golang:1.19

WORKDIR /flightowl

COPY . .

RUN ["go", "build"]

ENV PORT=8000

EXPOSE 8000

CMD ["./flightowl-api"]
