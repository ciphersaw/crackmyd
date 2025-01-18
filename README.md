# crackmyd

[![release](https://img.shields.io/github/v/release/ciphersaw/crackmyd)](https://github.com/ciphersaw/crackmyd) [![go](https://img.shields.io/badge/go-1.16-blue)](https://golang.org/)

The crackmyd is a lightweight tool to analyse the user.MYD file in MySQL and its congeneric RDBMS like MariaDB.

Specifically, it can extract  the host, user, and password from user.MYD, and find out the plaintext of password by brute force attack.

## Strategy

Up to now, there are three strategies for brute force attack as following:

| Strategy | Name              | Description                                                  |
| -------- | ----------------- | ------------------------------------------------------------ |
| 1        | Same As User Name | Check if the password is equal to the hash of user.          |
| 2        | Simple Guess      | Check if the password is equal to the hash of the weak passwords in default list, or the user-defined passwords in dictionary. |
| 3        | Suffix Combo      | Check if the password is equal to the hash of combination of user with the default suffixes, or the user-defined suffixes in dictionary. |

## Usage

The usages below are demonstrated in Linux, as the same as in Windows.

### Help

Input `./crackmyd --help` to get usages:

```bash
[root@kali ~]# ./crackmyd --help
Usage: ./crackmyd [options] <file>
  -password string
        Assign the user-defined dictionary of passwords for cracking.
  -suffix string
        Assign the user-defined dictionary of suffixes for cracking.
  -version
        Print the version of crackmyd.
```

### Version

Input `./crackmyd --version` to get the current version:

```bash
[root@kali ~]# ./crackmyd --version
v0.0.2
```

### Crack

Input `./crackmyd <file>` to crack user.MYD, and take `./example/mariadb_10.3_user.MYD` as example:

```bash
[root@kali ~]# ./crackmyd ./example/mariadb_10.3_user.MYD
+---------------+----------+------------------------------------------+-----------+
|     HOST      |   USER   |                 PASSWORD                 | PLAINTEXT |
+---------------+----------+------------------------------------------+-----------+
| localhost     | kali     | AA1420F182E88B9E5F874F6FBE7459291E8F4601 | qwerty    |
| 127.0.0.1     | app      | A8D20C5A046510B95A331561757262487D8313FB | app123    |
| 192.168.2.102 | kalinew  | 0971389F8E7F1AEB999104BDA7A0FA145087F348 |           |
| %             | crackmyd | 826C5303305ACCCE9A470BC26E53BC90CB8718B1 | crackmyd  |
| 172.16.3.103  | appnew   | F984F198C3CF24D4F41607E8510D53914BF88B1D |           |
| %             | stronger | 2D6D817D0F2ED9815092D259CF492DE19D4B4CFD |           |
+---------------+----------+------------------------------------------+-----------+
```

According to the result, it cracks the password `crackmyd` of user `crackmyd` by *Strategy 1*, the password `qwerty` of user `kali` by *Strategy 2*, and the password `app123` of user `app` by *Strategy 3*.

Based on the default configs all above, now assign the dictionary of passwords or suffixes for cracking:

```bash
[root@kali ~]# ./crackmyd --password ./example/password.txt --suffix ./example/suffix.txt ./example/mariadb_10.3_user.MYD
+---------------+----------+------------------------------------------+------------------+
|     HOST      |   USER   |                 PASSWORD                 |    PLAINTEXT     |
+---------------+----------+------------------------------------------+------------------+
| localhost     | kali     | AA1420F182E88B9E5F874F6FBE7459291E8F4601 | qwerty           |
| 127.0.0.1     | app      | A8D20C5A046510B95A331561757262487D8313FB | app123           |
| 192.168.2.102 | kalinew  | 0971389F8E7F1AEB999104BDA7A0FA145087F348 | q1w2e3r4         |
| %             | crackmyd | 826C5303305ACCCE9A470BC26E53BC90CB8718B1 | crackmyd         |
| 172.16.3.103  | appnew   | F984F198C3CF24D4F41607E8510D53914BF88B1D | appnew@gmail.com |
| %             | stronger | 2D6D817D0F2ED9815092D259CF492DE19D4B4CFD |                  |
+---------------+----------+------------------------------------------+------------------+
```

With the dictionaries, it eventually cracks the password `q1w2e3r4` of user `kalinew` by *Strategy 2*, and the password `appnew@gmail.com` of user `appnew` by *Strategy 3*.

In addition, it never cracks the password of user `stronger`, who has the strong password `IamStronger@2021`.

## Notice

- Please run crackmyd with a higher permission accessible to user.MYD.
- If it failed to crack the hash of password, the value of plaintext would be empty.
