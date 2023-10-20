# koghi - a rly simple fetch in go

![11ms execution time](preview.png)
> dasd aa desd
> asdfasdfasf
> bambam
> asdfasdf
> edad balsl
> AAAAAAAAA
> ababababa
> asdfasdf

## how to use
compile is ``go build``
install is ``go build && sudo cp ./koghi /usr/bin/``

## configuration
for all the config look into main.go!!  
you can change your wm name (incase youre on wayland) and toggle arch package fetching  
also you obviously can change ascii art  
file paths are in util.go if you need to change anything  

```reqs: just go, works on p much all linux machines with no external libraries (except for getting ur window manager from X11)```

ldd output: 
```
linux-vdso.so.1
libc.so.6 => /usr/lib/libc.so.6
/lib64/ld-linux-x86-64.so.2 => /usr/lib64/ld-linux-x86-64.so.2
```
