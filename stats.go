package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
)

type Record struct {
	IPhash uint32
	Size   uint32
}

func count(m map[uint32]uint32, t uint32) uint32 {
	c := uint32(0)
	for _, v := range m {
		if v == t {
			c++
		}
	}
	return c
}

func sum(m map[uint32]uint32) uint32 {
	s := uint32(0)
	for _, v := range m {
		s += v
	}
	return s
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	total_packets := uint32(37607469)
	sampling_rate := uint32(100)

	var record Record

	sampled := make(map[uint32]bool)
	uns_pkts := make(map[uint32]uint32)
	uns_byts := make(map[uint32]uint32)
	smp_pkts := make(map[uint32]uint32)
	smp_byts := make(map[uint32]uint32)
	det_pkts := make(map[uint32]uint32)
	det_byts := make(map[uint32]uint32)
	ind_pkts := make(map[uint32]uint32)
	ind_byts := make(map[uint32]uint32)
	buf_pkts := make(map[uint32]uint32)
	buf_byts := make(map[uint32]uint32)

	for uint32(len(sampled)) < total_packets/sampling_rate {
		sampled[uint32(rand.Intn(int(total_packets)))] = true
	}

	var t uint32
	var z uint32
	var tot_sz uint64
	var dns uint32
	var udp uint32

	infile, err := os.Open("data/packed.bin")
	check(err)

	var buf_pkt uint32
	for i := uint32(0); ; i++ {
		if err = binary.Read(infile, binary.LittleEndian, &record); err != nil {
			break
		}

		t += 1
		tot_sz += uint64(record.Size)
		if record.Size == 0 {
			z += 1
			//continue
		}
		if record.IPhash&1 > 0 {
			dns += 1
		}
		if record.IPhash&2 > 0 {
			udp += 1
		}

		if i%uint32(sampling_rate) == 0 {
			buf_pkt = uint32(rand.Intn(int(sampling_rate)))
		}

		uns_pkts[record.IPhash] += 1
		uns_byts[record.IPhash] += record.Size
		if sampled[i] {
			smp_pkts[record.IPhash] += 1
			smp_byts[record.IPhash] += record.Size
		}
		if i%uint32(sampling_rate) == 0 {
			det_pkts[record.IPhash] += 1
			det_byts[record.IPhash] += record.Size
		}
		if rand.Intn(int(sampling_rate)) == 0 {
			ind_pkts[record.IPhash] += 1
			ind_byts[record.IPhash] += record.Size
		}
		if i%uint32(sampling_rate) == buf_pkt {
			buf_pkts[record.IPhash] += 1
			buf_byts[record.IPhash] += record.Size
		}
	}

	if err != io.EOF {
		check(err)
	}

	fmt.Println(t, z, tot_sz, dns, udp)
	fmt.Printf("Unsampled      Flows: %10d  Packets %10d  Bytes %10d\n", len(uns_pkts), sum(uns_pkts), sum(uns_byts))
	fmt.Printf("Simple         Flows: %10d  Packets %10d  Bytes %10d\n", len(smp_pkts), sum(smp_pkts), sum(smp_byts))
	fmt.Printf("Deterministic  Flows: %10d  Packets %10d  Bytes %10d\n", len(det_pkts), sum(det_pkts), sum(det_byts))
	fmt.Printf("Independent    Flows: %10d  Packets %10d  Bytes %10d\n", len(ind_pkts), sum(ind_pkts), sum(ind_byts))
	fmt.Printf("Buffered       Flows: %10d  Packets %10d  Bytes %10d\n", len(buf_pkts), sum(buf_pkts), sum(buf_byts))

	fmt.Println("\nTable 1")
	fmt.Println("r\tsimple\tbuffered\tdeterministic\tindependent")
	for i := uint32(0); i < 20; i++ {
		fmt.Printf("%d\t", i+1)
		fmt.Printf("%d\t", count(smp_pkts, i+1))
		fmt.Printf("%d\t", count(buf_pkts, i+1))
		fmt.Printf("%d\t", count(det_pkts, i+1))
		fmt.Printf("%d\n", count(ind_pkts, i+1))
	}

	fmt.Println("\nTable 2")
	fmt.Println("r\tActual mean\tNew estimate")
	for i := uint32(0); i < 19; i++ {
		fmt.Printf("%d\t", i+1)
		t = uint32(0)
		c := 0
		for IPhash, v := range smp_pkts {
			if v != i+1 {
				continue
			}
			c += 1
			t += uns_pkts[IPhash]
		}
		fmt.Printf("%f\t%f\n", float32(t)/float32(c),
			float32((i+2)*sampling_rate*count(smp_pkts, i+2))/float32(count(smp_pkts,
				i+1)))
	}

	// TODO: why isn't this working?
	fmt.Println("\nTable 4")
	fmt.Println("Method\tEstimate\tActual")
	pkts := uint32(0)
	for IPhash := range smp_pkts {
		pkts += uns_pkts[IPhash]
	}
	fmt.Printf("Simple\t%f\t%f\n",
		1-float32(count(smp_pkts, 1))/float32(len(smp_pkts)),
		1-float32(pkts)/float32(total_packets))
	pkts = 0
	for IPhash := range buf_pkts {
		pkts += uns_pkts[IPhash]
	}
	fmt.Printf("Buffered\t%f\t%f\n",
		1-float32(count(buf_pkts, 1))/float32(len(buf_pkts)),
		1-float32(pkts)/float32(total_packets))
	pkts = 0
	for IPhash := range det_pkts {
		pkts += uns_pkts[IPhash]
	}
	fmt.Printf("Deterministic\t%f\t%f\n",
		1-float32(count(det_pkts, 1))/float32(len(det_pkts)),
		1-float32(pkts)/float32(total_packets))
	pkts = 0
	for IPhash := range ind_pkts {
		pkts += uns_pkts[IPhash]
	}
	fmt.Printf("Independent\t%f\t%f\n",
		1-float32(count(det_pkts, 1))/float32(len(det_pkts)),
		1-float32(pkts)/float32(total_packets))

	// TODO: why isn't this working?
	fmt.Println("\nRepeat rate")
	reprate := uint32(0)
	for _, v := range uns_pkts {
		reprate += v * v
	}
	fmt.Println("True:", float64(reprate)/(float64(total_packets)*float64(total_packets)))
	reprate = 0
	for i := uint32(0); i < 1000; i++ {
		reprate += (i + 1) * i * count(smp_pkts, i+1)
	}
	fmt.Println("Est: ", float64(reprate)/(float64(total_packets)*float64(total_packets-1)))

	fmt.Println("\nTable 5")
	fmt.Println("r\tEstimate\tActual")
	for i := uint32(0); i < 5; i++ {
		fmt.Printf("%d\t", i+1)
		est := uint32(0)
		act := uint32(0)
		for IPhash, p := range smp_pkts {
			if p == i+1 {
				est += smp_byts[IPhash] * sampling_rate
				act += uns_byts[IPhash]
			}
		}
		fmt.Printf("%f\t%f\n",
			float32(est)/float32(1000*count(smp_pkts, i+1)),
			float32(act)/float32(1000*count(smp_pkts, i+1)))
	}

	fmt.Println("\nTable 6")
	fmt.Println("k\tTrue sum\tTrue mean\tProxy sum\tProxy mean")
	for k := uint32(0); k < 5; k++ {
		fmt.Printf("%d\t", k+1)
		tsum := uint32(0)
		psum := uint32(0)
		for IPhash, p := range smp_pkts {
			if IPhash&2 > 0 || p != k+1 {
				continue
			}
			tsum += uns_byts[IPhash]
			psum += smp_byts[IPhash]
		}
		fmt.Printf("%d\t%f\t%d\t%f\n", tsum, float32(tsum)/float32(len(smp_pkts)),
			psum, float32(psum)/float32(len(smp_pkts)))
	}

	fmt.Println("\nFigure 1 written to file")
	fig1, err := os.Create("fig1.txt")
	defer fig1.Close()
	check(err)
	for IPhash := range smp_pkts {
		fmt.Fprintf(fig1, "%d %d\n", smp_pkts[IPhash], smp_byts[IPhash])
	}

	fmt.Println("\nFigure 2 written to files")
	fig2a, err := os.Create("fig2a.txt")
	defer fig2a.Close()
	fig2b, err := os.Create("fig2b.txt")
	defer fig2b.Close()
	for IPhash, p := range smp_pkts {
		if p == 1 && uns_pkts[IPhash] == 5 {
			fmt.Fprintf(fig2a, "%f\n", float32(uns_byts[IPhash])/5.)
		} else if p == 1 && uns_pkts[IPhash] == 25 {
			fmt.Fprintf(fig2b, "%f\n", float32(uns_byts[IPhash])/25.)
		}
	}

	fmt.Println("\nDNS")

	single_dns := 0
	for IPhash, p := range uns_pkts {
		if p == 1 && (IPhash&1 > 0) {
			single_dns++
		}
	}
	fmt.Printf("%d flows, %d single-packet flows, %d DNS flows, %d single-packet DNS flows\n",
		len(uns_pkts), count(uns_pkts, 1), dns, single_dns)
}
