import tornado.web
import random

from qcloudsms_py import SmsSingleSender
from qcloudsms_py.httpclient import HTTPError
import config


def rand_code():
    s = ""
    for i in range(0, 6):
        s += random.choice('1234567890')
    return s


class SendMessageHandler(tornado.web.RequestHandler):
    def get(self):
        phone_number = self.get_argument('phone_number')
        ssender = SmsSingleSender(config.appid, config.appkey)
        params = [rand_code(), "3"]
        try:
            result = ssender.send_with_param(86, phone_number,
                                             config.template_id, params, sign=config.sms_sign, extend="",
                                             ext="")
            if result['result'] == 0:
                self.write('ok')
            else:
                self.write('error')
        except HTTPError as e:
            print(e)
            self.write('error')
        except Exception as e:
            print(e)
            self.write('error')
