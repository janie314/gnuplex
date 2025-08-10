import os
import subprocess
import sys
import signal
import shutil
from pathlib import Path
import hashlib


def source_hash():
    backend = Path(__file__).parent / "backend"
    files = sorted(backend.rglob("*.go"))
    sha = hashlib.sha256()
    for f in files:
        sha.update(str(f).encode())
        with open(f, "rb") as fh:
            sha.update(fh.read())
    return sha.hexdigest()


def run(cmd, **kwargs):
    print(f"Running: {cmd}")
    subprocess.run(cmd, shell=True, check=True, **kwargs)


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
        subprocess.Popen("caddy run", shell=True),
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
        f'go build -C backend -o {target} -ldflags "-X main.SourceHash={source_hash()}" .'
    )


def go_build_ci():
    backend = Path(__file__).parent / "backend"
    os.chdir(backend)
    run(f'go build -o /tmp/gnuplex -ldflags "-X main.SourceHash={source_hash()}" .')


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
