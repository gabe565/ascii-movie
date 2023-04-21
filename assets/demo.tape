Output ascii-movie.mp4

Set Framerate 30

Set Theme "Builtin Dark"
Set FontSize 17
Set Width 970
Set Height 720
Set Padding 40
Set Margin 20
Set WindowBar Colorful
Set MarginFill "#7A6F64"
Set BorderRadius 10

# Use a shell alias to hide the cursor after typing
# and set --speed=2 flag.
Hide
Type@0 'ascii-movie() ( sleep 0.1 && exec ./ascii-movie "$@" ); clear'
Enter
Sleep 100ms
Show
Sleep 50ms

Type "ascii-movie play"
Sleep 500ms
Enter

# Total movie time, plus a few seconds, divided by 2.
Sleep 1070s
