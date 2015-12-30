#!/usr/bin/env python3
import struct
import sys
from collections import defaultdict
import random

total_packets = 37607469
sampling_rate = 100

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

if len(sys.argv)>1 and sys.argv[1] == '-c':
  for (iphash, sz) in rec_iter("data/packed.bin"):
    t += 1
    tot_sz += sz
    if sz==0: z += 1
  print(t, z, tot_sz)
  exit()

sampled = list(range(total_packets))
random.shuffle(sampled)
sampled = sampled[:total_packets//10]

uns_pkts = defaultdict(int)
uns_byts = defaultdict(int)
smp_pkts = defaultdict(int)
smp_byts = defaultdict(int)
det_pkts = defaultdict(int)
det_byts = defaultdict(int)
ind_pkts = defaultdict(int)
ind_byts = defaultdict(int)
buf_pkts = defaultdict(int)
buf_byts = defaultdict(int)

for i, (iphash, sz) in enumerate(rec_iter("data/packed.bin")):
  t += 1
  tot_sz += sz
  if sz==0: z += 1

  if i%sampling_rate == 0:
    buf_pkt = random.randrange(sampling_rate)

  uns_pkts[iphash] += 1
  uns_byts[iphash] += sz
  if i in sampled:
    smp_pkts[iphash] += 1
    smp_byts[iphash] += sz
  if i%sampling_rate == 0:
    det_pkts[iphash] += 1
    det_byts[iphash] += sz
  if random.randrange(sampling_rate)==0:
    ind_pkts[iphash] += 1
    ind_byts[iphash] += sz
  if i%sampling_rate==buf_pkt:
    buf_pkts[iphash] += 1
    buf_byts[iphash] += sz

print("Unsampled")
print("Flows:", len(uns_pkts), " Packets:", sum(uns_pkts.values()), " Bytes:", sum(uns_byts.values()))
print("Simple")
print("Flows:", len(smp_pkts), " Packets:", sum(smp_pkts.values()), " Bytes:", sum(smp_byts.values()))
print("Deterministic")
print("Flows:", len(det_pkts), " Packets:", sum(det_pkts.values()), " Bytes:", sum(det_byts.values()))
print("Independent")
print("Flows:", len(ind_pkts), " Packets:", sum(ind_pkts.values()), " Bytes:", sum(ind_byts.values()))
print("Buffered")
print("Flows:", len(buf_pkts), " Packets:", sum(buf_pkts.values()), " Bytes:", sum(buf_byts.values()))

print("\nTable 1")
print("r\tsimple\tbuffered\tdeterministic\tindependent")
for i in range(20):
  print(i+1, end='\t')
  print(tuple(smp_pkts.values()).count(i+1), end='\t')
  print(tuple(buf_pkts.values()).count(i+1), end='\t')
  print(tuple(det_pkts.values()).count(i+1), end='\t')
  print(tuple(ind_pkts.values()).count(i+1))

print("\nTable 2")
print("r\tActual mean")
for i in range(19):
  print(i+1, end='\t')
  t = 0
  c = 0
  for (iphash, _) in filter(lambda x: x[1]==i+1, smp_pkts.values()):
    c += 1
    t += uns_pkts[iphash]
  print(t/c)

print("\nTable 4")
print("Method\tActual")
pkts = 0
for iphash in smp_pkts:
  pkts += uns_pkts[pkts]
print("Simple", 1-pkts/total_packets, sep='\t')
pkts = 0
for iphash in buf_pkts:
  pkts += uns_pkts[pkts]
print("Simple", 1-pkts/total_packets, sep='\t')
pkts = 0
for iphash in det_pkts:
  pkts += uns_pkts[pkts]
print("Simple", 1-pkts/total_packets, sep='\t')
pkts = 0
for iphash in ind_pkts:
  pkts += uns_pkts[pkts]
print("Simple", 1-pkts/total_packets, sep='\t')

print("\nTable 5")
print("r\tActual")
for i in range(5):
  print(i+1, end='\t')
  t = 0
  c = 0
  for (iphash, _) in filter(lambda x: x[1] == i+1, smp_pkts.values()):
    c += 1
    t += uns_byts[iphash]


