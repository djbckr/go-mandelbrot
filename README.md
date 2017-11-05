# A smooth Mandelbrot Set rendering engine.

There were a few things I wanted to accomplish with this project.

1) Learn how to concurrently process items in Go.
2) Make a super-high-resolution Mandelbrot, with no banding.
3) Learn a bit about what the Mandelbrot Set is.

The concurrent piece was not too hard. The idea for me was to run each frame
in its own goroutine, but it was important not to dump 3600 frames to render
all at the same time. Instead, I setup a set of goroutines == number of processors
on the computer. This allows each frame to render in a reasonable amount of time,
and when the routine finishes one, it picks up the next one until the channel is closed.

The next thing I wanted to have was a 4K 60fps video of a Mandelbrot zoom, but it was
important to me that it not have banding; a common issue with naive implementations.
Research led me to the code in the mandelbrot() function, and it works very nicely.
I don't fully understand what it does (math isn't my strong suit), but I'm very happy
with what comes out.

The program will read a ```colors.json``` file with an array of colors, and it will
create a gradient palette based on those colors. If a color file isn't found, it will
make a preset gradient that isn't all that bad. You probably do not want any of your
colors to be close to black, as black is the "inside" of the set. Maybe I'll make
a modification where you can set whatever color you want for that. See below for the format.

The part that actually took the most time to program, and also takes the most processing
time, is the anti-aliasing piece. When the Mandelbrot gets close to the "inside", the
colors swing wildly, and without anti-aliasing and converted to video, it looks
like a bunch of bugs are messing about.

The output is a set of png files, which can be turned into a video using ffmpeg or another program of your choosing. I typically use the free version of DaVinci Resolve.

This program does not rely on any external libraries. Just run ```go build``` then ```./mandelbrot```
___
### The ```colors.json``` format
```
[
  {"R":95,"G":59,"B":45},
  {"R":188,"G":214,"B":240},
  {"R":163,"G":106,"B":54},
  {"R":54,"G":100,"B":133},
  {"R":145,"G":122,"B":41},
  {"R":132,"G":57,"B":26},
  {"R":246,"G":212,"B":180},
  {"R":148,"G":176,"B":213},
  {"R":189,"G":106,"B":63}
]
```
It's recommended to have at least 3 color objects, and you can have as many as you like. I find that 5~10 colors looks pretty good. The values should be between 0 and 255, and the program will interpolate these color gradients to 16-bits per color. As noted above, you want to stay clear of
black so that the final image has good contrast.