# GNUPlex

Version 0.9 (Currituck)

GNUPlex is a lightweight media display and media library database with a
friendly web interface.

The stack is:

- [Gin](https://github.com/gin-gonic/gin) +
  [SQLite](https://www.sqlite.org/index.html)
- [Vite](https://vitejs.dev/config/)

GNUPlex either runs in **user mode** or **server mode**. In user mode, the web
server operates out of `./gnuplex` (this Git repository).

In server mode, a web server is installed in `/var/gnuplex` and run out of a
SystemD service (running as the `gnuplex` user). This is useful for setting up a
computer as a home media display; the GNUPlex web interface would then typically
be acessed from a laptop or smartphone.

# Prerequisites

- mpv
- Go
- NodeJS / npm

Try e.g. `sudo dnf install mpv go nodejs npm` or
`sudo apt install mpv go nodejs npm`.

# Build

```bash
git clone https://gitlab.com/jane314/gnuplex.git
./build
```

# Running (user mode)

```bash
./gnuplex
```

# Testing

```bash
./dev
```

# Install or update server mode

Ensure port port 40000 is open on your server. The site will be hosted at
`http:hostname:40000/`.

```bash
sudo sh install-server
```

# Release History 

| Release | Date | Changes |
| - | - | - |
| 0.91 Currituck | 2022-12-25 | YouTube casting. |
| 0.9 Currituck | 2022-12-25 | Initial release. |

