import tornado.web
import random

from qcloudsms_py import SmsSingleSender
from qcloudsms_py.httpclient import HTTPError
import config



class SendMessageHandler(tornado.web.RequestHandler):
    def get(self):
        phone_number = self.get_argument('phone_number')
        code = self.get_argument('code')
        ssender = SmsSingleSender(config.appid, config.appkey)
        params = [code, "3"]
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
