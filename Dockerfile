 FROM golang
 #RUN mkdir /work/src/github.com/bitly/oauth2_proxy
 WORKDIR /work/src/github.com/bitly/oauth2_proxy
 COPY ./ /work/src/github.com/bitly/oauth2_proxy
 RUN  curl -s https://raw.githubusercontent.com/pote/gpm/v1.4.0/bin/gpm > gpm && chmod +x gpm && ./gpm install
 RUN go build