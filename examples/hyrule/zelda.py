#!/usr/bin/env python3

import socket
import sys
sys.path.append('/opt/walrus/python')
import walrus

IP = "192.168.147.3"
PORT = 4747
MSGLEN = 64

wtf = walrus.Test('192.168.147.100', 'hyrule', 'zelda')

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
sock.bind((IP, PORT))

print("awaiting inquery")
data, addr = sock.recvfrom(MSGLEN)
print("inquery received")
msg = data.decode('utf-8')


if msg != 'do you know the muffin man?':
    wtf.error('unexpected message: %s'%msg)
else:
    wtf.ok('got muffin man request')
    bits = 'why yes i know the muffin man'.encode('utf-8')
    sock.sendto(bits, addr)
