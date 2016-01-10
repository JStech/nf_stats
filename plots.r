# FIG 1
png(file = 'fig1.png', bg = 'transparent')
fig1 = read.table('fig1.txt')
p = fig1[[1]]
b = fig1[[2]]
plot(jitter(p), b+1, log='xy', pch=20, cex=0.1,
    xlab="Packets (with jitter)",
    ylab="Bytes + 1")
title(main="Flows with 1 packet sampled, packets vs. bytes")
dev.off()

# FIG 2
png(file = 'fig2.png', bg='transparent', width=960, height=480)
par(mfrow=c(1,2))
fig2a = read.table('fig2a.txt')
hist(fig2a[[1]], breaks=seq(0,1400,50),
    main="Flows with 1 packet sampled and 5 packets total",
    xlab="Bytes per packet")
fig2b = read.table('fig2b.txt')
hist(fig2b[[1]], breaks=seq(0,1400,50),
    main="Flows with 1 packet sampled and 25 packets total",
    xlab="Bytes per packet")
dev.off()

# FIG 3
png(file = 'fig3.png', bg='transparent')
fig3 = read.table('fig3.txt')
rownames(fig3) <- c('FALSE', 'TRUE')
colnames(fig3) <- c('FALSE', 'TRUE')
mosaicplot(fig3, main="All flows, proportions of port 53 and single packet",
xlab="Port 53", ylab="Single packet")
dev.off()
