#!/usr/bin/env python3
import re
import subprocess
import shlex

tcpd_cmd = "tcpdump (tcp or udp) and not ip6 -q -n -r {:}"
ip = b"\d{1,3}\."*4
tcpd_re = re.compile(b"(?:\d\d:\d\d:\d\d.\d{6}) IP ("+ip+b")(\d{1,5}) > ("+ip+
    b")(\d{1,5}): (?:tcp|UDP, length) (\d{1,4})")

with subprocess.Popen(shlex.split(
  tcpd_cmd.format("data/equinix-chicago.dirA.20150219-130000.UTC.anon.pcap")),
  stdout=subprocess.PIPE, bufsize=1) as tcpd:
  t = 0
  z = 0
  for l in tcpd.stdout:
    try: (fm_ip, fm_pt, to_ip, to_pt, sz) = tcpd_re.match(l).groups()
    except:
      print(l, end="")
      raise
    iphash = hash((fm_ip, fm_pt, to_ip, to_pt)) & 0xfffffffe
    sz = int(sz)
    if "53" in (fm_pt, to_pt):
      iphash += 1
    t += 1
    if sz==0: z += 1
    print((fm_ip, fm_pt, to_ip, to_pt), hex(iphash), sz, t, z)
    if iphash & 1: break
