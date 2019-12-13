#!/bin/sh
export GOOS="$1"
if [ "$GOOS" = linux ]; then
 rm CCUpdaterUI.linux CCUpdaterUI.linux.zip
 # Linux: standard non-static compilation (to prevent linking with deps we haven't worked out the licensing of)
 #  from within a container (to ensure that any libc linkage is as portable as possible)
 # IF IT FAILS HERE BECAUSE OF A MISSING CONTAINER, RUN THIS:
 # docker build --tag ccupdaterui-env . -f linux.Dockerfile
 docker run -v "`pwd`/..:/go/src/github.com/20kdc/CCUpdaterUI:rw" ccupdaterui-env
 mv ../CCUpdaterUI CCUpdaterUI.linux
 zip CCUpdaterUI.linux.zip CCUpdaterUI.linux
elif [ "$GOOS" = windows ]; then
 rm CCUpdaterUI.windows.exe CCUpdaterUI.windows.zip
 # Windows: https://github.com/veandco/go-sdl2#cross-compiling
 #  this is honestly just the same solution as with Linux, just with a different method
 #  it's for the same reasons, too - we need to link with msvcrt.dll and/or no CRT at all for version compat.
 #  what we absolutely DON'T want to do is link with, say, the UCRT
 # two additional gotchas we have to deal with:
 #  + a manual adjustment for GetDoubleClickTime
 #  + -x -v to ensure you can see what happens
 cd ..
 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc GOOS=windows GOARCH=386 go build -v -tags static -ldflags "-s -w"
 cd build
 echo "If an issue occurred with GetDoubleClickTime:"
 echo " + Get the MinGW devlibs from https://www.libsdl.org/download-2.0.php"
 echo " + Overwrite libSDL2_windows_386.a with 32-bit libSDL2.a"
 echo " + Overwrite libSDL2main_windows_386.a with 32-bit libSDL2main.a"
 mv ../CCUpdaterUI.exe CCUpdaterUI.windows.exe
 zip CCUpdaterUI.windows.zip CCUpdaterUI.windows.exe
elif [ "$GOOS" = darwin ]; then
 echo "Sorry!"
 # Mac OS X:
 #if [ -z "$DARLING" ]; then
 # echo "Please install Darling ( https://www.darlinghq.org/ ) and put it in env. variable DARLING"
 #else
 # echo "Sorry, TODO! You'll have to do it yourself"
 #fi
else
 echo "Target OS $1 not supported"
fi
