# crackmyd

[![release](https://img.shields.io/github/v/release/ciphersaw/crackmyd)](https://github.com/ciphersaw/crackmyd) [![go](https://img.shields.io/badge/go-1.23-blue)](https://golang.org/)

[English](README.md) | 简体中文

**crackmyd** 是一个轻量级工具，用于分析 MySQL 及其同类 RDBMS（如 MariaDB）中的 `user.MYD` 文件。

具体来说，它可以从 `user.MYD` 文件中提取用户名、口令、允许用户登录所使用的主机IP等信息，并通过暴力破解找到口令明文。

## 破解策略

目前支持三种暴力破解策略，如下所示：

| 策略编号 | 名称              | 描述                                                   |
| -------- | ----------------- | ------------------------------------------------------ |
| 1        | Same As User Name | 检查口令是否等于用户名的哈希值。                       |
| 2        | Simple Guess      | 检查口令是否等于默认或指定字典中弱口令的哈希值。       |
| 3        | Suffix Combo      | 检查口令是否等于用户名拼接默认或指定后缀组合的哈希值。 |

## 用法演示

接下来在 Linux 系统中演示相关用法，在 Windows 系统中用法相同。

### 使用说明

输入 `./crackmyd --help` 获取使用说明：

```bash
[root@kali ~]# ./crackmyd --help
Usage: ./crackmyd [options] <file>
  -h, --help              Print the usage of crackmyd.
  -p, --password string   Assign the user-defined dictionary of passwords for cracking.
  -s, --suffix string     Assign the user-defined dictionary of suffixes for cracking.
  -v, --version           Print the version of crackmyd.
```

### 获取版本

输入 `./crackmyd --version` 获取当前版本：

```bash
[root@kali ~]# ./crackmyd --version
v0.1.0
```

### 破解口令

输入 `./crackmyd <file>` 来破解 `user.MYD` 文件，以 `./example/mariadb_10.3_user.MYD` 为例：

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

根据以上结果，它通过 *策略 1* 破解了用户 `crackmyd` 的口令 `crackmyd`，通过 *策略 2* 破解了用户 `kali` 的口令 `qwerty`，并通过 *策略 3* 破解了用户 `app` 的口令 `app123`。

除了使用默认字典，还可以指定弱口令字典与拼接后缀进行破解：

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

使用指定字典后，它最终通过 *策略 2* 破解了用户 `kalinew` 的口令 `q1w2e3r4`，并通过 *策略 3* 破解了用户 `appnew` 的口令 `appnew@gmail.com`。

此外，它未能破解用户 `stronger` 的口令，该用户拥有强口令 `IamStronger@2021`。

## 注意事项

- 请确保以足够高的权限运行 **crackmyd** 工具，以便访问 `user.MYD` 文件。
- 若未能破解口令哈希值，则将明文口令置为空值。
