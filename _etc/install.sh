apt-get -y install unattended-upgrades
dpkg-reconfigure -plow unattended-upgrades

apt-get -y install convmv ca-certificates libcap2-bin

adduser --disabled-password --gecos "" martin
su - martin -c 'mkdir /home/martin/.ssh'
su - martin -c 'echo "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMfCM9I5xPrH5bl1TyjbCaqNPALl3cTU1e4NeJNqHTK3 martin" >> /home/martin/.ssh/authorized_keys'

su - martin -c 'mkdir /home/martin/ahouhpuc'
