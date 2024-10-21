# GNUPlex

GNUPlex is a lightweight media display and media library database with a
friendly web interface.

# Installation and Updating

```shell
curl -s https://gnuplex.janie.page/install | sh
```

This will prompt you for a directory in your PATH, then install GNUPlex there.
If you want to script this, you can use the `path` query paramter. (There's also
a `systemd` user parameter.)

```shell
curl -s https://gnuplex.janie.page/install?path=%2Fhome%2Fjanie%2F.local%2Fbin%2Fgnuplex | sh
```

To install GNUPlex as a SystemD user service, run:

```shell
gnuplex install-user-service
```

To update GNUPlex:

```shell
gnuplex update
```

# Releases

| Version      | Date                           | Details                         |
| ------------ | ------------------------------ | ------------------------------- |
| 0.99         | 2024-10-20                     | v1.0, beta.                     |
| Alpha stages | ~Christmas 2022 - October 2024 | See commits before 1.0 release. |
