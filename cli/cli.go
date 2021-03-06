package cli

import (
	"encoding/base64"
	"fmt"
	"net/http"
	_ "net/http/pprof" //so we can have debugging on localhost:3000 - See http://godoc.org/net/http/pprof
	"os"
	"strconv"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/vodolaz095/gossha/config"
	"github.com/vodolaz095/gossha/models"
	"github.com/vodolaz095/gossha/ssh"
	"github.com/vodolaz095/gossha/version"
)

// Greet writes some motivating text alongside with application version
func Greet() string {
	g := "ICBfX19fICAgICAgX19fXyBfX19fICBfICAgXyAgICAgICAKIC8gX19ffCBfX18vIF9fXy8gX19ffHwgfCB8IHwgX18gXyAKfCB8ICBfIC8gXyBcX19fIFxfX18gXHwgfF98IHwvIF9gIHwKfCB8X3wgfCAoXykgfF9fKSB8X18pIHwgIF8gIHwgKF98IHwKIFxfX19ffFxfX18vX19fXy9fX19fL3xffCB8X3xcX18sX3wKICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAK"
	gg, _ := base64.StdEncoding.DecodeString(g)
	var ggg []string
	ggg = append(ggg, "")
	ggg = append(ggg, string(gg))
	ggg = append(ggg, "Cross-platform ssh-server based chat program.")
	ggg = append(ggg, fmt.Sprintf("Build: %v", version.Version))
	ggg = append(ggg, "Homepage: https://github.com/vodolaz095/gossha")
	ggg = append(ggg, "Error reporting: https://github.com/vodolaz095/gossha/issues")
	ggg = append(ggg, "API documentation: https://godoc.com/github.com/vodolaz095/gossha")
	ggg = append(ggg, "           ")
	return strings.Join(ggg, "\r\n")
}

//ProcessConsoleCommand is a dispatcher for processing console commands and main entry point for application
func ProcessConsoleCommand() {
	var rootCmd = &cobra.Command{
		Use:   "gossha",
		Short: "GoSSHa is a cross-platform ssh-server based chat program",
		Long:  "GoSSHa is a cross-platform ssh-server based chat program, with data persisted into relational databases of MySQL, PostgreSQL or Sqlite3. Public channel (with persisted messages) and private message (not stored) are supported. Application has serious custom scripting and hacking potential.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Greet())
			fmt.Println()
			fmt.Println("Try `gossha help` for help...")
			fmt.Println()
			if config.RuntimeConfig.Debug {
				fmt.Println("Debug server is listening on http://localhost:3000/debug/pprof!")
				go func() {
					fmt.Println(http.ListenAndServe("localhost:3000", nil))
				}()
			}
			if len(config.RuntimeConfig.BindTo) > 0 {
				for _, v := range config.RuntimeConfig.BindTo[:(len(config.RuntimeConfig.BindTo) - 1)] {
					go ssh.StartSSHD(v)
				}
				ssh.StartSSHD(config.RuntimeConfig.BindTo[len(config.RuntimeConfig.BindTo)-1])
			} else {
				ssh.StartSSHD(fmt.Sprintf("0.0.0.0:%v", config.RuntimeConfig.Port))
			}
		},
	}
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Outputs program version and exits",
		Long:  "Outputs program version and exits",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Greet())
			os.Exit(0)
		},
	}
	var passwdCmd = &cobra.Command{
		Use:   "passwd [username] [password]",
		Short: "Creates user or set new password to existent one",
		Long:  "Creates user or set new password to existent one",
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 2:
				name := args[0]
				password := args[1]
				err := models.CreateUser(name, password, false)
				if err != nil {
					panic(err)
				}
				fmt.Printf("User %v is created and/or new password is set!\n", name)
				os.Exit(0)
				break
			case 1:
				name := args[0]
				fmt.Print("Enter password:")
				password, err := gopass.GetPasswd()
				if err != nil {
					panic(err)
				}
				err = models.CreateUser(name, string(password), false)
				if err != nil {
					panic(err)
				}
				fmt.Printf("User %v is created and/or new password is set!\n", name)
				os.Exit(0)
				break
			default:
				fmt.Printf("Enter user's name!\n")
				os.Exit(1)
			}
		},
	}
	var makeRootUserCmd = &cobra.Command{
		Use:   "root [username] [password]",
		Short: "Creates root user or set new password to existent one",
		Long:  "Creates root user or set new password to existent one",
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 2:
				name := args[0]
				password := args[1]
				err := models.CreateUser(name, password, true)
				if err != nil {
					panic(err)
				}
				fmt.Printf("Admin %v is created and/or new password is set!\n", name)
				os.Exit(0)
				break
			case 1:
				name := args[0]
				fmt.Print("Enter password:")
				password, err := gopass.GetPasswd()
				if err != nil {
					panic(err)
				}
				err = models.CreateUser(name, string(password), true)
				if err != nil {
					panic(err)
				}
				fmt.Printf("Admin %v is created and/or new password is set!\n", name)
				os.Exit(0)
				break
			default:
				fmt.Printf("Enter user's name!\n")
				os.Exit(1)
			}
		},
	}
	var banCmd = &cobra.Command{
		Use:   "ban [username]",
		Short: "Delete user and all his/her messages",
		Long:  "Delete user and all his/her messages",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				name := args[0]
				err := models.BanUser(name)
				if err != nil {
					panic(err)
				}
				fmt.Printf("User %v is banned!\n", name)
				os.Exit(0)
			} else {
				fmt.Printf("Enter user's name!\n")
				os.Exit(1)
			}
		},
	}
	var dumpConfig = &cobra.Command{
		Use:   "dumpcfg",
		Short: "Outputs the configuration in YAML format",
		Long:  "Outputs the configuration in YAML format. Save this config in `$HOME/.gossha/gossha.toml` or `/etc/gossha/gossha.toml`",
		Run: func(cmd *cobra.Command, args []string) {
			json, err := config.RuntimeConfig.Dump()
			if err != nil {
				panic(err)
			}
			fmt.Println(json)
		},
	}

	var listUsers = &cobra.Command{
		Use:   "list",
		Short: "List users",
		Long:  "List users",
		Run: func(cmd *cobra.Command, args []string) {
			var users []models.User
			k := 0
			err := models.DB.Table("user").Order("user.name ASC").Find(&users).Error
			if err != nil {
				panic(err)
			}
			fmt.Println("Users in database:")
			for _, u := range users {
				k++
				if u.Root {
					fmt.Printf("%v) %v (root!) - online on %v \n", k, u.Name, u.LastSeenOnline.Format("2006-1-2 15:04:05"))
				} else {
					fmt.Printf("%v) %v - online on %v \n", k, u.Name, u.LastSeenOnline.Format("2006-1-2 15:04:05"))
				}
			}
		},
	}

	var listMessages = &cobra.Command{
		Use:   "log [limit]",
		Short: "Show last messages, default limit is 10",
		Long:  "Show last messages, default limit is 10",
		Run: func(cmd *cobra.Command, args []string) {
			var ret []models.Notification
			var messages []models.Message
			var limit int
			if len(args) == 1 {
				l, _ := strconv.ParseInt(args[0], 10, 8)
				if l > 0 {
					limit = int(l)
				} else {
					limit = 10
				}
			} else {
				limit = 10
			}
			err := models.DB.Table("message").Preload("User").Limit(limit).Order("message.id desc").Find(&messages).Error
			if err != nil {
				panic(err)
			}
			for _, m := range messages {
				ret = append(ret, models.Notification{User: m.User, Message: m, IsSystem: false, IsChat: true})
			}
			for _, n := range ret {
				var u = n.User
				var m = n.Message
				var online string
				if u.IsOnline() {
					online = "*"
				} else {
					online = "x"
				}
				fmt.Printf("[%v@%v(%v) %v]{%v}:%v\r\n", u.Name, m.Hostname, m.IP, online, m.CreatedAt.Format("2006-1-2 15:04:05"), m.Message)
			}
		},
	}

	//Note! - this flags are actually used in `config.go#InitConfig`.
	//They are copied here to make application more user friendly!
	/*
		rootCmd.PersistentFlags().Uint("port", 27015, "set the port to listen for connections")
		rootCmd.PersistentFlags().Bool("debug", false, "start pprof debugging on http://localhost:3000/debug/pprof/. See `http://godoc.org/net/http/pprof`")
		rootCmd.PersistentFlags().String("driver", "sqlite3", "set the database driver to use, possible values are `sqlite3`,`mysql`,`postgres`")
		rootCmd.PersistentFlags().String("connectionString", GetDatabasePath(), MakeDSNHelp())
		rootCmd.PersistentFlags().String("sshPublicKeyPath", GetPublicKeyPath(), "location of public ssh key to be used with server, usually the $HOME/.ssh/id_rsa.pub")
		rootCmd.PersistentFlags().String("sshPrivateKeyPath", GetPrivateKeyPath(), "location of private ssh key to be used with server, usually the $HOME/.ssh/id_rsa")
		rootCmd.PersistentFlags().String("homedir", GetHomeDir(), "The home directory of module, usually $HOME/.gossha")
		rootCmd.PersistentFlags().String("executeOnMessage", "", "Script to execute on each message")
		rootCmd.PersistentFlags().String("executeOnPrivateMessage", "", "Script to execute on each private message")
	*/
	rootCmd.AddCommand(versionCmd, passwdCmd, makeRootUserCmd, banCmd, dumpConfig, listUsers, listMessages)
	rootCmd.Execute()
}
