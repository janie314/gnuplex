import hashlib
import os
import signal
import subprocess
import sys
from pathlib import Path


def source_hash():
    backend = Path(__file__).parent / "backend"
    files = sorted(backend.rglob("*.go"))
    sha = hashlib.sha256()
    for f in files:
        sha.update(str(f).encode())
        with open(f, "rb") as fh:
            sha.update(fh.read())
    return sha.hexdigest()


def platform():
    os = subprocess.check_output("uname -s", shell=True).decode().strip()
    arch = subprocess.check_output("uname -m", shell=True).decode().strip()
    libc = (
        "musl"
        if "musl" in subprocess.check_output("ldd /bin/ls", shell=True).decode().strip()
        else "glibc"
    )
    return f"{os}-{libc}-{arch}".lower()


def run(cmd, **kwargs):
    print(f"+ {cmd}")
    res = subprocess.run(cmd, shell=True, check=True, **kwargs)
    if res.returncode != 0:
        print("cmd failed. exiting...")
        os.exit(res.returncode)


def dev():
    os.chdir(os.path.dirname(__file__))
    run("bun i --cwd frontend")
    os.makedirs("tmp", exist_ok=True)
    procs = [
        subprocess.Popen("caddy run", shell=True),
        subprocess.Popen("bun run --cwd frontend dev", shell=True),
        subprocess.Popen("go run -C backend . -verbose", shell=True),
    ]

    def cleanup(signum, frame):
        for p in procs:
            p.terminate()
        sys.exit(0)

    signal.signal(signal.SIGTERM, cleanup)
    try:
        for p in procs:
            p.wait()
    except KeyboardInterrupt:
        cleanup(None, None)


def dev_compiled():
    os.chdir(os.path.dirname(__file__))
    os.makedirs("tmp", exist_ok=True)
    procs = [
        subprocess.Popen("caddy run --config Caddyfile-compiled", shell=True),
        subprocess.Popen(
            "./backend/bin/gnuplex -verbose -static_files ./backend/static", shell=True
        ),
    ]

    def cleanup(signum, frame):
        for p in procs:
            p.terminate()
        sys.exit(0)

    signal.signal(signal.SIGTERM, cleanup)
    try:
        for p in procs:
            p.wait()
    except KeyboardInterrupt:
        cleanup(None, None)


def frontend_build():
    os.chdir(os.path.dirname(__file__))
    run("bun i --cwd frontend")
    run("bun run --cwd frontend build")


def go_build():
    os.chdir(os.path.dirname(__file__))
    target = os.environ.get("TARGET", "bin/gnuplex")
    run(
        f'go build -C backend -o {target} -ldflags "-X main.SourceHash={source_hash()} '
        + f'-X main.Platform={platform()}" .'
    )


def go_build_ci():
    run(
        "go build -C backend -o /tmp/gnuplex"
        + f' -ldflags "-X main.SourceHash={source_hash()}'
        + f' -X main.Platform={platform()}" .'
    )


def build():
    frontend_build()
    go_build()


def go_build_current():
    os.chdir(os.path.dirname(__file__))
    exe = Path(__file__).parent / "backend/bin/gnuplex"
    if not exe.exists():
        sys.exit(1)
    build_hash = subprocess.check_output([str(exe), "-source_hash"]).decode().strip()
    if source_hash() != build_hash:
        sys.exit(1)


def go_source_hash():
    print(source_hash())


TASKS = {
    "dev": dev,
    "dev_compiled": dev_compiled,
    "frontend_build": frontend_build,
    "go_build": go_build,
    "go_build_ci": go_build_ci,
    "build": build,
    "go_build_current": go_build_current,
    "go_source_hash": go_source_hash,
    "platform": platform,
}


def main():
    if len(sys.argv) < 2:
        print("Available tasks:")
        for t in TASKS:
            print(f"  {t}")
        sys.exit(1)
    task = sys.argv[1]
    if task not in TASKS:
        print(f"Unknown task: {task}")
        sys.exit(1)
    TASKS[task]()


if __name__ == "__main__":
    main()
