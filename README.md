# crackmyd

[![release](https://img.shields.io/github/v/release/ciphersaw/crackmyd)](https://github.com/ciphersaw/crackmyd) [![go](https://img.shields.io/badge/go-1.16-blue)](https://golang.org/)

The crackmyd is a lightweight tool to analyse the user.MYD file in MySQL and its congeneric RDBMS like MariaDB.

Specifically, it can extract  the host, user, and password from user.MYD, and find out the plaintext of password by brute force attack.

## Usage

The usages below are demonstrated in Linux, as the same as in Windows.

### Help

Input `./crackmyd --help` to get usages:

```bash
[root@kali ~]# ./crackmyd --help
Usage: ./crackmyd [options] <file>
  -version
        Print the version of crackmyd.
```

### Version

Input `./crackmyd --version` to get the current version:

```bash
[root@kali ~]# ./crackmyd --version
v0.0.1
```

### Crack

Input `./crackmyd <file>` to crack user.MYD, and now take `./example/user.MYD` as example:

```bash
[root@kali ~]# ./crackmyd ./example/user.MYD
+---------------+----------+------------------------------------------+-----------+
|     HOST      |   USER   |                 PASSWORD                 | PLAINTEXT |
+---------------+----------+------------------------------------------+-----------+
| %             | kali     | 7CFA671F67454DD0660587E2177286063B7B458F | kali123   |
| 192.168.2.104 | app      | AA1420F182E88B9E5F874F6FBE7459291E8F4601 | qwerty    |
| 127.0.0.1     | crackmyd | 826C5303305ACCCE9A470BC26E53BC90CB8718B1 | crackmyd  |
| 0.0.0.0       | strong   | BCC78A72D0A2929990A5036AF58ED265399E7A72 |           |
+---------------+----------+------------------------------------------+-----------+
```

## Notice

- Please run crackmyd with a higher permission accessible to user.MYD.
- If it failed to crack the hash of password, the value of plaintext would be empty.
