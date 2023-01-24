#!/usr/bin/env python3

import select
import socket


def echo(host='10.0.0.1', port=7, timeout=1, msg='hello'):
    ai = socket.getaddrinfo(host, port, type=socket.SOCK_DGRAM)
    family, _, _, _, addr = ai[0]
    sock = socket.socket(family, socket.SOCK_DGRAM)
    sock.sendto(msg.encode(), addr)
    rlist, _, _ = select.select([sock], [], [], timeout)
    if rlist:
        buf = sock.recv(4096)
        print(f'recv {buf.decode()}')
    else:
        print('timeout')


def main():
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument('-H', '--host', default='10.0.0.1')
    parser.add_argument('-p', '--port', type=int, default=7)
    parser.add_argument('-t', '--timeout', type=int, default=1)
    parser.add_argument('-m', '--msg', default='hello')
    args = parser.parse_args()
    echo(args.host, args.port, args.timeout, args.msg)


if __name__ == '__main__':
    main()
