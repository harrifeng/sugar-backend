

import tornado.ioloop
import tornado.web

import handlers


def make_server():
    return tornado.web.Application([
        (r"/send_message", handlers.SendMessageHandler),
    ])


if __name__ == '__main__':
    app = make_server()
    port = 7799
    app.listen(port)
    print("server running on 127.0.0.1:%d" % port)
    tornado.ioloop.IOLoop.current().start()

