from mininet.net import Containernet
from mininet.node import Controller, OVSBridge
from mininet.cli import CLI

net = Containernet(controller=None)
s1 = net.addSwitch('s1', cls=OVSBridge)
server = net.addDocker('c2', ip='10.0.0.10', dimage='cdci-c2:latest',
                       dcmd='supervisord -n -c /etc/supervisor/supervisord.conf',
                       environment={'HEARTBEAT_TOKEN': 'secret'})

agent = net.addDocker('victim', ip='10.0.0.21', dimage='cdci-victim:latest',
                      dcmd='bash -c "supervisord -c /etc/supervisor/supervisord.conf && exec runuser -l dev"',
                      environment={'HEARTBEAT_SERVER_URL': 'http://10.0.0.10:8080/report',
                                   'HEARTBEAT_TOKEN': 'secret',
                                   'HEARTBEAT_INTERVAL': '5'})

net.addLink(server, s1); net.addLink(agent, s1)
net.start(); CLI(net); net.stop()
