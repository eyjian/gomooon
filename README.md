# mooonutils

提供各类工具函数，还文件类、字符串类、时间类和系统类等。

# mooonstr

专注提供字符串类工具函数。

# moooncrypto

提供常用加密解密封装。

# mooongorm

提供 gorm 的辅助函数。

# mooonhttp

提供 http 的辅助函数。

# txcloud

提供腾讯云的身份证、银行卡和手机号核验功能。

# mooonwepay

提供微信支付的下载账单和回执单等功能。

# mooonpdf

提供 jpg 转 png 和 pdf 转 png 的函数。

mooonpdf 依赖开源的 pdf 库 go-fitz，而 go-fitz 又依赖开源的 c 库 mupdf，对 ld 版本要求，2.23 版本的 ld 编译报错：

```
go/pkg/mod/github.com/gen2brain/go-fitz@v1.22.2/libs/libmupdf_linux_amd64.a(colorspace.o): unrecognized relocation (0x2a) in section `.text.fz_find_icc_link'
```

需要将 ld 升级到 2.28 版本。ld 在开源的 binutils 中，下载地址：

```
https://ftp.gnu.org/gnu/binutils/
```

另外种方式从重新编译 mupdf 开始：

```
https://github.com/ArtifexSoftware/mupdf
```

而 mupdf 又依赖 ghostpdl：

```
https://github.com/plangrid/ghostpdl/tree/master
```

所有最好是通过升级 binutils 的方式解决。
