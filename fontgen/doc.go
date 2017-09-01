/*
Fontgen generates font files to be used by github.com/ZiRo/captcha

To install, simply run:
	github.com/ZiRo/captcha

This progam also requires ImageMagick's convert binary.

Fontgen has two options
	usage: fontgen [flags] filename
	  -font string
			Font name (default "Monospace")
	  -size string
			Font size (default "30")

Size is the desired font size in point. Values around 25 have proven to work the best.
Font is the name of the font you want to generate. To find a list of all available fonts,
just run
	convert -list font

For example:
	fontgen -font DejaVu-Sans-Mono -size 25 DejaVu-Sans-Mono.gob
*/
package main
