apt-get -y install unattended-upgrades
dpkg-reconfigure -plow unattended-upgrades

echo "deb http://http.debian.net/debian wheezy-backports main" >> /etc/apt/sources.list
apt-get update
apt-get -t wheezy-backports install nginx
/etc/init.d/nginx start

apt-get -y install convmv ca-certificates

adduser --disabled-password --gecos "" martin
su - martin -c 'mkdir /home/martin/.ssh'
su - martin -c 'echo "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAw9R1/dJA3wrhZ5fCTWBg5gVZhfleWQJ6bDMfplTZ7TjCYuq0/KkjYGGAxB4IplR0NMeVAjfrs2RWMuUSwDmI3Fr+y1xVrHWdwpESciOvx7k0YnVhETIxbLmnVCSkcTzyYCjdmQvxNwElkr55TEt+1zVpWMNTx9d5bNjcgXoaZyqAM4PTF2O9KCOiUOVsiklygCM6GY4dVAC/Z3+Xhsp4/q/wojGlNEzjtKQAD6OXD3ogmQl9TPAURo7QdOtGhIYo6sp7eq4XtsdidSHCNPaXsS4d6MM9+LTXtVlxzoBwRFiw4k/625BCLj4RnlDbu+vBvY6ZtCnj5I/rRm7MQsIC+w== martin" >> /home/martin/.ssh/authorized_keys'

su - martin -c 'mkdir /home/martin/ahouhpuc'
