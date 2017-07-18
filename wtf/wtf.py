#!/usr/bin/env python3

import redis
import sys

def usage():
    print("wtf <test pattern>")

if len(sys.argv) != 2:
    usage()
    sys.exit(1)

host = 'localhost'
pattern = sys.argv[1]

r = redis.StrictRedis(host=host)


colors = {
        'error': '\033[91m',
        'warning': '\033[93m',
        'ok': '\033[92m',
        'test': '\033[94m',
        'participant': '\033[34m',
        'clear': '\033[0m'
        }

diagnostics = {}

for x in r.scan_iter(match=pattern):
    k = x.decode('utf-8')
    
    if '~seq~' in k:
        continue
    
    test, participant, counter, *_ = k.split(':', 4)

    #TODO show timestamp
    key = '[%s:%s:%s]'%(
        test,
        participant,
        counter)
    color_key = '[%s%s%s:%s%s%s]'%(
        colors['test'],
        test,
        colors['clear'], 
        colors['participant'],
        participant,
        colors['clear'])


    if key not in diagnostics:
        diagnostics[key] = {}

    diagnostics[key]['colorkey'] = color_key

    test, participant, counter = k.split(':', 2)
    v = r.get(x).decode('utf-8')
    level, message = v.split(':::')

    message = '%s%s%s'%(
        colors[level], 
        message, 
        colors['clear'])
    
    diagnostics[key]['message'] = message


#TODO sort by timestamp
#diags = sorted(diagnostics.items(), key=lambda x: x[1]['time'])
diags = diagnostics.items()

for x in diags:
    print("%s %s"%(x[1]['colorkey'], x[1]['message'], ))

