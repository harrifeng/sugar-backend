

import tornado.ioloop
import tornado.web

import handlers


def make_server():
    return tornado.web.Application([
        (r"/send_message", handlers.SendMessageHandler),
    ])


if __name__ == '__main__':
    app = make_server()
    app.listen(7799)
    tornado.ioloop.IOLoop.current().start()