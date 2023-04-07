Output ascii-movie.mp4

# Full video is almost 18 min, so record at 2x speed, then output at 1/2 speed
# so that recordings only take ~9 min.
Set Framerate 30
Set PlaybackSpeed 0.5
Set TypingSpeed 20ms

Set FontSize 17
Set Width 980
Set Height 620
Set Padding 40
Set Margin 20
Set WindowBar Colorful
Set MarginFill "#FFDDBD"
Set BorderRadius 10

# Use a shell alias to hide the cursor after typing
# and set --speed=2 flag.
Hide
Type@0 "alias ascii-movie='sleep 0.1 && tput civis && ./ascii-movie --speed=2' && clear"
Enter
Sleep 100ms
Show
Sleep 50ms

Type "ascii-movie play"
Sleep 250ms
Enter

# Total movie time, plus a few seconds, divided by 2.
Sleep 536s
