# GNUPLEX

# RULES

1. Be careful with main / know what branch you're currently in.
1. One repo to rule them all. (More precisely, we have the same monorepo cred as Google now.)
1. Installation instructions are: `bundle`.
1. `standardrb` is the standard linter.
1. Follow the Gwen / Torvalds rules about commit messages. (https://yarchive.net/comp/linux/commit_messages.html)
1. Defer to Gwen about dev stuff, but also Jane is not a total noob.
1. Extreme arguments will be resolved by a chess match.

```
Personally, I don't write documentation. I don't much even write comments
in my code. My personal feeling is that as long as functions are small and
readable (and logical), and global variables have good names, that's all I
need to do. Others - who do care about comments and docs - can do that
part.
	- Linus Torvalds
```


# DEV ENVIRONMENT

- Run `./server/gnuplex-server` in a terminal window. 
  - Actually we're messing with the server now. Current dev environment instructions are: `just figure it out`.
- Then use one of the gnuplex clients.
