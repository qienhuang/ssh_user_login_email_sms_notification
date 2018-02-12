**Send you an email or SMS when someone login to your Raspberry Pi/Linux system**

Gmail will reject connection if you didn't "Allow less secure apps."  https://myaccount.google.com/lesssecureapps

```
go run main.go
#Before you build it, don't forget to set disableLogDisplay = true
go build main.go
cp main /bin/sshlog
nano /etc/profile.d/log.sh
```

#Add the following lines in log.sh
```
#!/bin/bash
/bin/sshlog
```

![example](https://github.com/qienhuang/ssh_user_login_email_sms_notification/raw/master/email_example.jpg)
)
