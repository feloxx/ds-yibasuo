#!/bin/bash

if ansible --version >/dev/null 2>&1; then
  echo "ansible ok"
  exit 0
else
  echo "ansible start offline install"
fi

currentPath=$(cd `dirname $0`; pwd)

otherPath=$currentPath/o
ansiblePath=$currentPath/a

########################
# rpm install base env
echo "other"

sudo rpm -ivh $otherPath/mpfr-3.1.1-4.el7.x86_64.rpm
sudo rpm -ivh $otherPath/libmpc-1.0.1-3.el7.x86_64.rpm
sudo rpm -ivh $otherPath/kernel-headers-3.10.0-1062.el7.x86_64.rpm
sudo rpm -ivh $otherPath/glibc-headers-2.17-292.el7.x86_64.rpm
sudo rpm -ivh $otherPath/glibc-devel-2.17-292.el7.x86_64.rpm
sudo rpm -ivh $otherPath/cpp-4.8.5-39.el7.x86_64.rpm
sudo rpm -ivh $otherPath/gcc-4.8.5-39.el7.x86_64.rpm
sudo rpm -ivh $otherPath/libgcc-4.8.5-39.el7.x86_64.rpm

sudo rpm -ivh $otherPath/sshpass-1.06-2.el7.x86_64.rpm

sudo rpm -ivh $otherPath/libcurl-7.29.0-54.el7.x86_64.rpm
sudo rpm -ivh $otherPath/python-pycurl-7.19.0-19.el7.x86_64.rpm

sudo rpm -ivh $otherPath/bash-completion-2.1-6.el7.noarch.rpm
sudo rpm -ivh $otherPath/cmake-2.8.12.2-2.el7.x86_64.rpm
sudo rpm -ivh $otherPath/libaio-devel-0.3.109-13.el7.x86_64.rpm
sudo rpm -ivh $otherPath/net-tools-2.0-0.25.20131004git.el7.x86_64.rpm
sudo rpm -ivh $otherPath/gd-2.0.35-26.el7.x86_64.rpm
sudo rpm -ivh $otherPath/libxml2-2.9.1-6.el7_2.3.x86_64.rpm
sudo rpm -ivh $otherPath/openssl-1.0.2k-19.el7.x86_64.rpm
sudo rpm -ivh $otherPath/pcre-8.32-17.el7.x86_64.rpm
sudo rpm -ivh $otherPath/pcre-devel-8.32-17.el7.x86_64.rpm
sudo rpm -ivh $otherPath/perl-Data-Dumper-2.145-3.el7.x86_64.rpm
sudo rpm -ivh $otherPath/unzip-6.0-20.el7.x86_64.rpm
sudo rpm -ivh $otherPath/zlib-1.2.7-18.el7.x86_64.rpm

sudo rpm -ivh $otherPath/libuuid-devel-2.23.2-61.el7.x86_64.rpm $otherPath/libXpm-devel-3.5.12-1.el7.x86_64.rpm $otherPath/fontconfig-devel-2.13.0-4.3.el7.x86_64.rpm $otherPath/expat-devel-2.1.0-10.el7_3.x86_64.rpm $otherPath/gd-devel-2.0.35-26.el7.x86_64.rpm
if [ `rpm -qa | grep gd-devel | wc -l` -eq 0 ];then
  echo "error"
  exit 0
fi

sudo rpm -ivh $otherPath/libverto-devel-0.2.5-4.el7.x86_64.rpm $otherPath/libsepol-devel-2.5-10.el7.x86_64.rpm $otherPath/libselinux-devel-2.5-14.1.el7.x86_64.rpm $otherPath/libcom_err-devel-1.42.9-16.el7.x86_64.rpm $otherPath/krb5-devel-1.15.1-37.el7_6.x86_64.rpm $otherPath/keyutils-libs-devel-1.5.8-3.el7.x86_64.rpm $otherPath/openssl-devel-1.0.2k-19.el7.x86_64.rpm
if [ `rpm -qa | grep openssl-devel | wc -l` -eq 0 ];then
  echo "error"
  exit 0
fi

sudo rpm -ivh $otherPath/xz-devel-5.2.2-1.el7.x86_64.rpm $otherPath/libxml2-devel-2.9.1-6.el7_2.3.x86_64.rpm
if [ `rpm -qa | grep libxml2-devel | wc -l` -eq 0 ];then
  echo "error"
  exit 0
fi

# todo 判断包是否都安装全

########################
# src install ansible
echo "ansible"

cd $ansiblePath && tar -zxf pycrypto-2.6.1.tar.gz
cd pycrypto-2.6.1
sudo python setup.py install

cd $ansiblePath && tar -zxf yaml-0.1.5.tar.gz
cd yaml-0.1.5
./configure --prefix=/usr/local
make --jobs=`grep processor /proc/cpuinfo | wc -l`
sudo make install

cd $ansiblePath && tar -zxf PyYAML-3.11.tar.gz
cd PyYAML-3.11
sudo python setup.py install

cd $ansiblePath && tar -zxf MarkupSafe-1.1.1.tar.gz
cd MarkupSafe-1.1.1
sudo python setup.py install

cd $ansiblePath && tar -zxf Jinja2-2.9.6.tar.gz
cd Jinja2-2.9.6
sudo python setup.py install

cd $ansiblePath && tar -zxf jmespath-0.9.4.tar.gz
cd jmespath-0.9.4
sudo python setup.py install

cd $ansiblePath && tar -zxf paramiko-2.7.1.tar.gz
cd paramiko-2.7.1
sudo python setup.py install

cd $ansiblePath && tar -zxf ansible-2.8.7.tar.gz
cd ansible-2.8.7
sudo python setup.py install

########################
# check
echo "check"

a=`ansible --version`
result=$(echo $a | grep "ansible 2.8.7")

if [ "$result" != "" ];then
  echo "ok"
else
  echo "error"
fi