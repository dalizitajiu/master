from flask import Flask
from flask import request
import smtplib
import json
from email.mime.text import MIMEText
app = Flask(__name__)
_user = "********@qq.com"
_pwd = "bpfoypnbzzqabdbd"


def actSend(to, msg):
    global _user, _pwd
    msg = MIMEText(msg)
    msg["Subject"] = "确认邮件"
    msg["From"] = _user
    msg["To"] = to
    try:
        stmpobj = smtplib.SMTP_SSL("smtp.qq.com", 465)
        stmpobj.login(_user, _pwd)
        stmpobj.sendmail(_user, to, msg.as_string())
        stmpobj.quit()
        print("Success!")
    except Exception as exception:
        print("Failed,%stmpobj" % exception)


@app.route('/')
def hello_world():
    '''helloworld'''
    return 'Hello World!'


@app.route('/send_mail', methods=['POST'])
def sendMail():
    '''发送邮件'''
    actSend(request.form["to"], request.form["msg"])
    return json.dumps({"errno": 0, "errmsg": "", "data": ""})


if __name__ == '__main__':
    app.run()
