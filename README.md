# Installation

Build binaries into `/opt/maildock`:

~~~
GOPATH=/opt/maildock go install github.com/yinyin/go-maildock/cmd/...
~~~

Setup user `maildock` with home folder `/opt/maildock`:

~~~
sudo useradd -s /usr/sbin/nologin -r -M -d /opt/maildock maildock
~~~

Install and enable services:

~~~
cp maildock-smtpd.service /lib/systemd/system/
sudo systemctl enable maildock-smtpd.service
sudo systemctl start maildock-smtpd.service
~~~

