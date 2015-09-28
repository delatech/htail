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

function publish_debian {
echo -e ">>> Publishing debian packages"
    aptly repo create delatech
    aptly repo add delatech bin/htail-${version}-i386.deb
    aptly repo add delatech bin/htail-${version}-amd64.deb
    aptly snapshot create delatech-htail-${version} from repo delatech


    # for first tie use
    # aptly publish -distribution=squeeze snapshot delatech-htail-${version} s3:apt.delatech.net:
    aptly publish switch squeeze s3:apt.delatech.net: delatech-htail-${version}

}

function publish_homebrew {
    echo -e "\n>>> Publishing osx package"
    gzip -f ./bin/htail-darwin-${version}.tar

    sha1sum=`sha1sum ./bin/htail-darwin-${version}.tar.gz | awk '{print $1}'`
    aws s3 cp bin/htail-darwin-${version}.tar.gz s3://release.delatech.net/htail/htail-${version}.tar.gz --acl=public-read

    cat <<EOF > $DELATECH_BREWTAP/Formula/htail.rb
#encoding: utf-8

require 'formula'

class Htail < Formula
    homepage 'https://github.com/delatech/htail'
    version '${version}'

    url 'http://release.delatech.net.s3-website-eu-west-1.amazonaws.com/htail/htail-${version}.tar.gz'
    sha1 '${sha1sum}'

    depends_on :arch => :intel

    def install
        bin.install 'bin/htail'
    end
end
EOF
    cd $DELATECH_BREWTAP
    git add Formula/htail.rb
    git ci -m"Update htail to v${version}"
    git push origin master
}

export AWS_ACCESS_KEY_ID=$AWS_DELATECH_S3_APT_KEY
export AWS_SECRET_ACCESS_KEY=$AWS_DELATECH_S3_APT_SECRET

xc
#deb
#publish_debian
osx
publish_homebrew
