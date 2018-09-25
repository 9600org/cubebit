# Golang Pi Zero library for the 4tronix Cube:Bit

The [Cube:Bit)[https://4tronix.co.uk/cubebit] is a nifty 5x5x5 LED cube
stack for use with various micros, including the Raspberry PI Zero. 

This library builds on the work of [libws281x](https://github.com/jgarff/rpi_ws281x)
and [Máximo Cuadros' nifty Golang wrapper of the same](https://github.com/mcuadros/go-rpi-ws281x)
to provide an API which more closes matches the 3D cube of the physical
device.

While there's a bunch of support for using the *Cube:Bit* with the *BBC
MicroBit*, I wasn't able to find much else for the Raspberry PI, so I knocked
this up mainly just to play with it on the PI.

This is a video of what it looks like:<br/>
[![Video of cubebit demo running on 5x5x5 cube](http://img.youtube.com/vi/DZc7rKozVuI/0.jpg)](http://www.youtube.com/watch?v=DZc7rKozVuI "Cube:bit demo")



# The hardware

It's really just a string of WS281x *NeoPixels* formed into a cube shape,
although cunningly the path the *string* takes is slightly unexpected:

The *odd* layers are laid out like so:
<pre>
,-21--22--23--24--25
'-20--19--18--17--16,
,-11--12--13--14--15'
'-10---9---8---7---6,
   1---2---3---4---5'
</pre>

and the *even* layers are laid out like this:
<pre>
   .__.    .__.
  21  20  11  10   1
  22  19  12   9   2
  23  18  13   8   3
  24  17  14   7   4
  25  16  15   6   5
       '--'    '---'
</pre>

Still, that's easily dealt with.

Unfortunately the LEDs themselves seem not to have many levels of brightness
available, so the types of effects you can make with it *slightly* limited,
but it's still a brilliant bit of fun.

# The software

It's all a bit hacky and very much work-in-progress, but the bits that are
there dor *something* at least :)

There are some "demos" under the [cmd](cmd) directory, have a look at these to
see some usage examples. 
The `cubebit` package provides a direct cube-like API for setting the colours
of the LEDs.

The `renderer` package is a simple layer which sits above `cubebit` and
provides a mechanism for rendering simple shapes into the LED space.
