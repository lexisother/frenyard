# LAYERS FOR THE LAYER THRONE!
# CONTAINERS FOR THE CONTAINER THRONE!

# Dockerfile to run the Linux build, with VERY SPECIFIC requirements
#  arch = 386
#  platform = linux
#  distribution = as old as possible
# it's the only way to be sure
FROM golang:stretch@sha256:ff71d89f78e4fbc348c783d6978edf2721d7f360a11f70df1db91ce5acc36de1
# Start by declaring where I/O goes
RUN mkdir -p /go/src/github.com/20kdc/CCUpdaterUI
VOLUME /go/src/github.com/20kdc/CCUpdaterUI
WORKDIR /go/src/github.com/20kdc/CCUpdaterUI
# Dependency management goes here
RUN apt-get update
RUN apt-get install -y libsdl2-dev
RUN go get github.com/veandco/go-sdl2/sdl ; go build github.com/veandco/go-sdl2/sdl
RUN go get github.com/Masterminds/semver ; go build github.com/Masterminds/semver
RUN go get github.com/golang/freetype ; go build github.com/golang/freetype
RUN go get golang.org/x/image/font/gofont/goregular ; go get golang.org/x/image/font/gofont/gobold
RUN go get github.com/20kdc/go-vkv
# Final build step: get packages & build packages
# NOTE: Output from the 'CUT HERE' set is eligible to be moved upwards
ENTRYPOINT go get -v . ; echo "CUT HERE" ; go build -v -tags static -ldflags "-s -w"
