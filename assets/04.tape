Set PlaybackSpeed 1.0
Set TypingSpeed 200ms 
Set Width 900
Set Height 500

Set Shell zsh
Hide
Type@50ms "alias pbcopy='xclip -selection clipboard'"
Enter
Type@50ms "alias pbpaste='xclip -selection clipboard -o'"
Enter
Type@10ms "cat pokemon.md| pbcopy"
Enter
Type "clear"
Enter
Show
Type@100ms "pbpaste | mdtt | pbcopy"
Sleep 500ms
Enter
Sleep 500ms
Type "jjjdddd"
Sleep 100ms
Type "dd"
Sleep 500ms
Type "q"
Sleep 1s
Type@100ms "pbpaste "
Enter
Sleep 1500ms

