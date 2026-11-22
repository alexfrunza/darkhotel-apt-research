# Disclaimer

This PoC is modified code from: https://github.com/K1tt3h/CVE-2025-32463-POC, of the author (K1tt3h). This code is provided for educational and research purposes only.

By using any part of this repository, you agree that:
- You will not use this code or technique to gain unauthorized access to systems you do not own or have explicit permission to test.
- The authors of this repository assumes no responsibility or liability for any misuse, damage, or consequences caused by this PoC or related material.

Use this code at your own risk.

---

For building c2:
```
docker build -f Dockerfile -t cdci-c2:latest .
```

For building the script you will need:

passwords.txt with the password to decrypt zip archive
a zip archive with valuable content

Run:
```
bash encrypt_exploit.sh exploit.sh
```

And then copy the result variables in the go code to have the exploit script as a string encrypted in the malware binary.
Then go to the go directory and compile the program:
```
CGO_ENABLED=0 go build -ldflags="-s -w" -o unzipper .
# Optionally check that the strings from the script with exploit are not present:
# strings unzipper | grep "supervisorctl"
```

Build the dockerimage of the victim
```
docker build -f Dockerfile -t cdci-victim:latest .
# Run for tests:
# docker run -it --rm debian-sudo:1.9.17
```

Run the malware:
```
# Check password to put it in the program
cat passwords.txt
# Run the malware
./unzipper password top_secret.zip
```

Check persistence
```
sudo supervisorctl status
sudo supervisorctl tail -f heartbeat-agent

supervisorctl status
supervisorctl tail -f heartbeat-server
```

To prepare infra:
```
sudo dnf -y install git gcc make python3-devel python3-pip python3-setuptools \
  iptables stress-ng iproute net-tools curl openvswitch socat iperf3 ethtool help2man

sudo systemctl enable --now openvswitch

# Containernet
git clone https://github.com/containernet/containernet.git
cd containernet
sudo make -f Makefile install-mnexec install-manpages
python3 -m venv venv && source venv/bin/activate
pip install .

sudo su
source containernet/venv/bin/activate
cd infra
python network.py

docker exec -it mn.c2 bash
docker exec -it -u dev mn.victim bash
```
