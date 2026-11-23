
# Steps to reproduce the PoC 

A small PoC which mimic an attack made by the Darkhotel APT. The PoC is built on an infrastructure using Mininet and Docker containers.

---

## Disclaimer

This PoC is modified code from: https://github.com/K1tt3h/CVE-2025-32463-POC, of the author (K1tt3h). This code is provided for educational and research purposes only.

By using any part of this repository, you agree that:
- You will not use this code or technique to gain unauthorized access to systems you do not own or have explicit permission to test.
- The authors of this repository assumes no responsibility or liability for any misuse, damage, or consequences caused by this PoC or related material.

Use this code at your own risk.

---

## Prerequisites

### Building C2 image

```
cd poc/c2
docker build -f Dockerfile -t cdci-c2:latest .
```

### Building victim image

```
cd poc/victim
# Create a zip with mock content
bash create_zips.sh

# Encrypt the shell script with the exploit and copy the resulting string into the go source code
bash encrypt_exploit.sh exploit.sh

cd go_script
# Update the main.go with the result from the previous bash command

# Compile the malware
CGO_ENABLED=0 go build -ldflags="-s -w" -o unzipper .
# Optionally check that the strings from the script with exploit are not present
# strings unzipper | grep "supervisorctl"

cd ..

# Build the victim image
docker build -f Dockerfile -t cdci-victim:latest .
```

### Prepare mininet setup

```
cd poc

sudo dnf -y install git gcc make python3-devel python3-pip python3-setuptools \
  iptables stress-ng iproute net-tools curl openvswitch socat iperf3 ethtool help2man
sudo systemctl enable --now openvswitch

# Setup containernet
git clone https://github.com/containernet/containernet.git
cd containernet
sudo make -f Makefile install-mnexec install-manpages
python3 -m venv venv && source venv/bin/activate
pip install .
deactivate
```

## Run the PoC

```
sudo su
cd poc
source containernet/venv/bin/activate
cd infra
python network.py
```

### C2 machine

```
# Connect to the c2 container
docker exec -it mn.c2 bash

# Check data from victim
supervisorctl status
supervisorctl tail -f heartbeat-server
```

### Victim machine

```
# Connect to the victim container
docker exec -it -u dev mn.victim bash

# Run the malware
# Check password to put it in the program
cat passwords.txt
# Run the malware
./unzipper password top_secret.zip

# Check persistence
sudo supervisorctl status
sudo supervisorctl tail -f heartbeat-agent
```
