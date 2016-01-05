set term png

set output 'fig1.png'
set logscale xy
plot 'fig1.txt' using ($1 + 0.6*(rand(0)-.5)):2 with dots
unset logscale xy

binwidth=40
bin(x,width) = width*floor(x/width) + binwidth/2.
set boxwidth binwidth
set output 'fig2a.png'
plot 'fig2a.txt' using (bin($1,binwidth)):(1.0) smooth freq with boxes
set output 'fig2b.png'
plot 'fig2b.txt' using (bin($1,binwidth)):(1.0) smooth freq with boxes
