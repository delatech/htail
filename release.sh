#!/bin/bash

function usage {
    echo -e "htail release script\n"
    echo "Usage:"
    echo "  $0 version"
    exit 1
}

version=$1
if [ -z "$version" ]; then
    usage
fi

if  [ ! -d "bin" ]; then
    mkdir bin
fi

function xc {
    echo ">>> Cross compiling htail"
    GOOS=linux GOARCH=amd64 go build -o bin/htail-${version}-linux-amd64
    GOOS=linux GOARCH=386 go build -o bin/htail-${version}-linux-i386
    GOOS=darwin GOARCH=amd64 go build -o bin/htail-${version}-darwin-amd64
}

function deb {
    arches="i386 amd64"
    for arch in $arches; do
        echo -e "\n>>> Creating debian package for ${arch}"
        fpm \
            -f \
            -s dir \
            -t deb \
            --vendor "DeLaTech" \
            --name   "htail" \
            --description "An utility that displays files, to stdout and in the browser" \
            --version $version \
            -a $arch \
            -p ./bin/htail-${version}-${arch}.deb \
            ./bin/htail-${version}-linux-${arch}=/usr/bin/htail
    done
}

function osx {
    echo -e "\n>>> Creating osx package"
    fpm \
        -f \
        -s dir \
        -t tar \
        --name   "htail" \
        -p ./bin/htail-darwin-${version}.tar \
        ./bin/htail-${version}-darwin-amd64=/usr/bin/htail
}

function upload {
    aptly repo create delatech
    aptly repo add delatech bin/htail-${version}-i386.deb
    aptly repo add delatech bin/htail-${version}-amd64.deb
    aptly snapshot create delatech-htail-${version} from repo delatech

    export AWS_ACCESS_KEY_ID=$AWS_DELATECH_S3_APT_KEY
    export AWS_SECRET_ACCESS_KEY=$AWS_DELATECH_S3_APT_SECRET

    # for first tie use
    # aptly publish -distribution=squeeze snapshot delatech-htail-${version} s3:apt.delatech.net:
    aptly publish switch squeeze s3:apt.delatech.net: delatech-htail-${version}

    #echo -e "\n>>> Uploading osx package"
    #curl -T ./bin/htail-darwin-${version}.tar -ufuturecat:$BINTRAY_API_KEY https://api.bintray.com/content/delatech/htail-osx/htail/${version}/htail.tar
}

xc
deb
osx
upload
