# This script is not intended to be run automatically. Run each command by hand.

apt-get update
apt-get upgrade

# install required packages
apt-get -y install convmv ca-certificates libcap2-bin iptables-persistent

# install and configure unattended-upgrades
apt-get -y install unattended-upgrades
dpkg-reconfigure -plow unattended-upgrades

# root authorized_keys
echo "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMfCM9I5xPrH5bl1TyjbCaqNPALl3cTU1e4NeJNqHTK3 martin_o" >> /root/.ssh/authorized_keys
echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDFiGIx2ofxTRfkFWJ3VHh94pm5evmUqQYMXKmuqJvnClzFVIAjgJtKiCQmnXDRa3RbSmAFC/MIFiM3noPfqM4fbde32cjH1T3h762HdvVRW6AOyYeYCM3rrdxz2JR9feP95kJq0UVd1xvlik2QZCnHEX6yDNO4MtasQ30CgiQsLGUxQOJjXuu/W+1MtH/djgPY+8eRAycfpOjINRT5okPGm4ThrCGNHk1HzBXK1RLjxpxYz9Sp8gWAUaCKJ9oAr+dd2hKcQICC8VMbtotyEFEqr+8Ihe+pDDhEJwwQlbCRvwDHNj9pI1qmD09c3AaytEj9ily+5wSl2yzxWKRjaWqt martin_s" >> /root/.ssh/authorized_keys
echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsD/LUw9LOT00uQSOZonqbZSK8h1iQZL/am9pa3QeeS15qmmpcB6Ii/9fq+JTYDLOVmpO1tqkObp62bzJAjWqJAG6J0HGZv0Fws2ybwKUFNNg5Q/F5LdwKkSn/if7YJgwaiuBIMIyq8eIxPbOwFPSUG+r4HA8EvhLyta4otLAcF+hSywKkFwbvytkIHs++Vor9Mt8JL8MNYUREyL8HKERzPjs8wN8BJCmlIUSmp27sqNhP8MiP3l13/B9tC9cfwiOEwzzN2hWyGgVFDI9D0NoIWoNlkXRzl4EdrxLI5tTmECAWpIwsNwTUM4Lxp/aL3PrUKses7YgRBdYhIljrBTtt damien_m" >> /root/.ssh/authorized_keys

# create user "martin" and add authorized_keys
adduser --disabled-password --gecos "" martin
su - martin -c 'mkdir /home/martin/.ssh'
su - martin -c 'echo "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMfCM9I5xPrH5bl1TyjbCaqNPALl3cTU1e4NeJNqHTK3 martin_o" >> /home/martin/.ssh/authorized_keys'
su - martin -c 'echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDFiGIx2ofxTRfkFWJ3VHh94pm5evmUqQYMXKmuqJvnClzFVIAjgJtKiCQmnXDRa3RbSmAFC/MIFiM3noPfqM4fbde32cjH1T3h762HdvVRW6AOyYeYCM3rrdxz2JR9feP95kJq0UVd1xvlik2QZCnHEX6yDNO4MtasQ30CgiQsLGUxQOJjXuu/W+1MtH/djgPY+8eRAycfpOjINRT5okPGm4ThrCGNHk1HzBXK1RLjxpxYz9Sp8gWAUaCKJ9oAr+dd2hKcQICC8VMbtotyEFEqr+8Ihe+pDDhEJwwQlbCRvwDHNj9pI1qmD09c3AaytEj9ily+5wSl2yzxWKRjaWqt martin_s" >> /home/martin/.ssh/authorized_keys'
su - martin -c 'echo " ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsD/LUw9LOT00uQSOZonqbZSK8h1iQZL/am9pa3QeeS15qmmpcB6Ii/9fq+JTYDLOVmpO1tqkObp62bzJAjWqJAG6J0HGZv0Fws2ybwKUFNNg5Q/F5LdwKkSn/if7YJgwaiuBIMIyq8eIxPbOwFPSUG+r4HA8EvhLyta4otLAcF+hSywKkFwbvytkIHs++Vor9Mt8JL8MNYUREyL8HKERzPjs8wN8BJCmlIUSmp27sqNhP8MiP3l13/B9tC9cfwiOEwzzN2hWyGgVFDI9D0NoIWoNlkXRzl4EdrxLI5tTmECAWpIwsNwTUM4Lxp/aL3PrUKses7YgRBdYhIljrBTtt damien_m" >> /home/martin/.ssh/authorized_keys'

# disable ssh password authentication
sed -i 's/.*PasswordAuthentication\ \+\(no\|yes\).*/PasswordAuthentication no/' /etc/ssh/sshd_config
systemctl restart sshd.service

# create the application directory
su - martin -c 'mkdir /home/martin/ahouhpuc'

# ahouhpuc.service
# scp _etc/ahouhpuc.service root@HOST:/etc/systemd/system/ahouhpuc.service
systemctl daemon-reload && systemctl enable ahouhpuc.service

# configure iptables

# iptables: clear all
iptables -t filter -F
iptables -t nat -F
iptables -t filter -X
ip6tables -t filter -F
ip6tables -t filter -X

# iptables: don't break already established connections
iptables -A INPUT -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
iptables -A OUTPUT -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
ip6tables -A INPUT -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT
ip6tables -A OUTPUT -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT

# iptables: allow ssh (22)
iptables -t filter -A INPUT -p tcp --dport ssh -j ACCEPT
ip6tables -t filter -A INPUT -p tcp --dport ssh -j ACCEPT

# iptables: allow ping
iptables -A INPUT -p icmp -j ACCEPT
ip6tables -A INPUT -p ipv6-icmp -j ACCEPT

# iptables: disallow any incomming connection
iptables -t filter -P INPUT DROP
iptables -t filter -P FORWARD DROP
ip6tables -t filter -P INPUT DROP
ip6tables -t filter -P FORWARD DROP

# iptables: allow loopback
iptables -A INPUT -i lo -j ACCEPT
ip6tables -A INPUT -i lo -j ACCEPT

# iptables: allow http (80)
iptables -t filter -A INPUT -p tcp --dport http -j ACCEPT
ip6tables -t filter -A INPUT -p tcp --dport http -j ACCEPT

# iptables: allow https (443)
iptables -t filter -A INPUT -p tcp --dport https -j ACCEPT
ip6tables -t filter -A INPUT -p tcp --dport https -j ACCEPT

# iptables: save
iptables-save > /etc/iptables/rules.v4
ip6tables-save > /etc/iptables/rules.v6

# done! let's reboot just in case...
reboot
