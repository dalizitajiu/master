from flask import Flask
from flask import request
import smtplib
import json
from email.mime.text import MIMEText
app = Flask(__name__)
_user = "********@qq.com"
_pwd  = "bpfoypnbzzqabdbd"
_to   = "lixiaomeng19920528@outlook.com"
def actSend(to,msg):
	global _user,_pwd
	msg=MIMEText(msg)
	msg["Subject"]="确认邮件"
	msg["From"]=_user
	msg["To"]=_to
	try:
		s=smtplib.SMTP_SSL("smtp.qq.com",465)
		s.login(_user,_pwd)
		s.sendmail(_user,_to,msg.as_string())
		s.quit()
		print("Success!")
	except e:
		print("Failed,%s"%e)
	pass
@app.route('/')
def hello_world():
    return 'Hello World!'
@app.route('/send_mail',methods=['POST'])
def sendMail():
	actSend(request.form["to"],request.form["msg"])
	return json.dumps({"errno":0,"errmsg":"","data":""})
if __name__ == '__main__':
    app.run()
