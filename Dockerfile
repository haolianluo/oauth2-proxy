 FROM golang
 RUN mkdir /code
 WORKDIR /code
 COPY ./ /code
 RUN echo "173.194.75.141 golang.org golang venus" >> /etc/hosts
 RUN echo "203.208.46.1 golang.org golang venus" >> /etc/hosts
 RUN  curl -s https://raw.githubusercontent.com/pote/gpm/v1.4.0/bin/gpm > gpm && chmod +x gpm && ./gpm install
 RUN go build