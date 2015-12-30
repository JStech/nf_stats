#!/usr/bin/env python3
import struct

def rec_iter(filename, bufsize=4096):
  bufsize &= 0xffffffffff8
  data = b""
  with open(filename, "rb") as infile:
    while True:
      buff = infile.read(bufsize)
      if not buff:
        return
      for i in range(0, len(buff), 8):
        yield struct.unpack("II", buff[i:i+8])

t = 0
z = 0
tot_sz = 0

for (iphash, sz) in rec_iter("data/packed.bin"):
  t += 1
  tot_sz += sz
  if sz==0: z += 1

print(t, z, tot_sz)
