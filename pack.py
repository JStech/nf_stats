#!/usr/bin/env python3
import re
import subprocess
import shlex
import struct

tcpd_cmd = "tcpdump (tcp or udp) and not ip6 -q -n -r {:}"
ip = b"\d{1,3}\."*4
tcpd_re = re.compile(b"(?:\d\d:\d\d:\d\d.\d{6}) IP ("+ip+b")(\d{1,5}) > ("+ip+
    b")(\d{1,5}): (?:tcp|UDP,(?: bad)? length) (\d{1,4})")

t = 0
z = 0
tot_sz = 0

with open("data/packed.bin", "wb") as outfile:
  for infile in ("data/equinix-chicago.dirA.20150219-130000.UTC.anon.pcap",
      "data/equinix-chicago.dirA.20150219-130100.UTC.anon.pcap"):
    with subprocess.Popen(shlex.split(
      tcpd_cmd.format(infile)),
      stdout=subprocess.PIPE, bufsize=1) as tcpd:
      for l in tcpd.stdout:
        if b"ip-proto-" in l: continue
        try: (fm_ip, fm_pt, to_ip, to_pt, sz) = tcpd_re.match(l).groups()
        except:
          print(l, end="")
          raise
        iphash = hash((fm_ip, fm_pt, to_ip, to_pt)) & 0xfffffffe
        sz = int(sz)
        if b"53" in (fm_pt, to_pt):
          iphash += 1
        t += 1
        tot_sz += sz
        if sz==0: z += 1
        outfile.write(struct.pack("II", iphash, sz))

print(t, z, tot_sz)
