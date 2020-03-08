#!/bin/bash

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