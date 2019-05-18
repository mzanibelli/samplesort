FROM alpine:latest
RUN apk add curl
RUN curl -O 'https://essentia.upf.edu/documentation/extractors/essentia-extractors-v2.1_beta2-linux-i686.tar.gz'
RUN tar xvzf essentia-extractors-v2.1_beta2-linux-i686.tar.gz
RUN install -t /usr/local/bin essentia-extractors-v2.1_beta2/streaming_extractor_music
ENTRYPOINT ["/usr/local/bin/streaming_extractor_music"]
