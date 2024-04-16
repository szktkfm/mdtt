Set PlaybackSpeed 1.0
Set TypingSpeed 200ms 
Set Width 500
Set Height 400

Set Shell zsh
Type@100ms "mdtt"
Enter
Sleep 500ms
Type "ki"
Type@100ms "FOO"
Ctrl+C
Type "li"
Type@100ms "BAR"
Ctrl+C
Type "jhi"
Type@100ms "foo"
Ctrl+C
Type "li"
Type@100ms "bar"
Sleep 500ms
Ctrl+C
Sleep 500ms
Type "q"
Sleep 1s

