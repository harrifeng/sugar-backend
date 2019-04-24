# 参数表
# 1 sex 1 女 2 男
# 2 age 年龄
# 3 height 身高
# 4 weight 体重
# 5 yundong 1-5 表示运动强度
# 6 type 糖尿病类型 1 1型 2 2型 3 妊娠 4 特殊
# 7 bing 有无并发症 0 没有 1 有
# 8 xue 血糖控制 0 稳定不变 1 空腹高 2 餐后高 3 经常低血糖
# 9 fuyong 服用胰岛素

import jieba
import sys
import json
import random

def is_number(s):
    try:
        float(s)
        return True
    except ValueError:
        pass
    try:
        import unicodedata
        unicodedata.numeric(s)
        return True
    except (TypeError, ValueError):
        pass
    return False


def only_num(answer):
    s2 = answer.lower()
    m = "abcdefghijklmnopqrstuvwxyz"
    for c in answer:
        if c in m:
            answer = answer.replace(c, "")
    return answer


def jb(answer):
    seg_list = list(jieba.cut_for_search(answer, HMM=True))
    return seg_list


def ask(qu):
    answer = '错误'
    if qu == 1:
        answer = "请问您的性别是什么呢？（男/女）"
    elif qu == 2:
        answer = "请问您的年龄是多少？"
    elif qu == 3:
        answer = "您的身高有多少cm呢？"
    elif qu == 4:
        answer = "您的体重有多少kg呢？"
    elif qu == 5:
        answer = "您平时运动多吗？"
    elif qu == 6:
        answer = "您的糖尿病类型是哪种类型呢？（1型/2型/妊娠型）如果没有确诊请回复没有，糖导将为你判断类型倾向"
    elif qu == 7:
        answer = "您现在是否有并发症状？如果有，已经确诊了哪些呢？"
    elif qu == 8:
        answer = "您现在血糖控制效果是否稳定？如果不稳定，是在哪些时候高/低呢？"
    elif qu == 9:
        answer = "您是否在使用胰岛素辅助治理？胰岛素方案是什么样子的呢？"
    elif qu == 10:
        answer = "现在的视力怎么样？"
    elif qu == 11:
        answer = "您当前的治疗方式是什么样的呢？"
    elif qu == 12:
        answer = "您是在什么时候确诊呢？"
    elif qu == 13:
        answer = "您的发病年龄是多少呢？"
    elif qu == 14:
        answer = "您的家族史是否明显？"
    elif qu == 15:
        answer = "您的胰岛素分泌情况如何？"
    elif qu == 16:
        answer = "您的病发速度如何？"
    elif qu == 17:
        answer = "是不是感觉视力不如从前了？"
    elif qu == 18:
        answer = "是不是经常感觉口渴无力？"
    elif qu == 27:
        answer = "正在根据您的BMI和BMR指数分析预判您的糖尿病类型……"
    elif qu == 28:
        answer = "您的糖尿病类型可能为2型糖尿病"
    elif qu == 29:
        answer = "您的糖尿病类型可能为1型糖尿病"
    elif qu == 31:
        answer = "嗨，糖友！欢迎使用糖导智能诊断助手，让我们开始诊断吧！"
    elif qu == 32:
        answer = "请问您现在的身体具体症状是怎样呢"
    return answer


def question(flag, msg):
    data = {'style': ''}
    pp = 0
    if flag == 1:
        answer = msg
        ci = jb(answer)
        size = len(ci)
        for i in range(0, size):
            if ci[i] == "女":
                pp = 1
            elif ci[i] == "男":
                pp = 2
        data = {'style': 'sex', 'num': pp}
    elif flag == 2:
        answer = msg
        ci = jb(answer)
        size = len(ci)
        for i in range(0, size):
            if ci[i] == "岁" or ci[i] == "年纪" or ci[i] == "年龄":
                sui1 = ci[i - 1]
                if i == size-1:
                    sui2 = ci[i - 1]
                else:
                    sui2 = ci[i + 1]
                if is_number(sui1) is True:
                    pp = int(sui1)
                elif is_number(sui2) is True:
                    pp = int(sui2)
                else:
                    pp = 0
            else:
                pp = int(ci[i])
        data = {'style': 'age', 'num': pp}

    elif flag == 3:
        answer = msg
        ci = jb(answer)
        size = len(ci)
        for i in range(0, size):
            if ci[i] == "身高":
                gao1 = ci[i - 1]
                gao2 = ci[i + 1]
                gao1 = only_num(gao1)
                gao2 = only_num(gao2)
                if is_number(gao2) is True:
                    pp = int(gao2)
                elif is_number(gao1) is True:
                    pp = int(gao1)
                else:
                    pp = 0
            else:
                gao = ci[i]
                gao = only_num(gao)
                pp = int(gao)
        data = {'style': 'height', 'num': pp}

    elif flag == 4:
        answer = msg
        ci = jb(answer)
        size = len(ci)
        for i in range(0, size):
            if ci[i] == "体重":
                ti1 = ci[i - 1]
                ti2 = ci[i + 1]
                ti1 = only_num(ti1)
                ti2 = only_num(ti2)
                if is_number(ti2) is True:
                    pp = int(ti2)
                elif is_number(ti1) is True:
                    pp = int(ti1)
                else:
                    pp = 0
            else:
                ti = ci[i]
                ti = only_num(ti)
                pp = int(ti)
        data = {'style': 'weight', 'num': pp}
    elif flag == 5:
        pp = 1
        answer = msg
        ci = jb(answer)
        size = len(ci)
        for i in range(0, size):
            if ci[i] == "不":
                pp=1
            elif ci[i] == "偶尔":
                pp = 2
            elif ci[i] == "经常":
                pp=3
            elif ci[i] == "每天":
                pp=5
            if ci[i] == "次":
                cici = ci[i - 1]
                cici = only_num(cici)
                if cici >= 3 and pp <= 4:
                    pp = 4
        data = {'style': 'sport', 'num': pp}
    elif flag == 6:
        answer = msg
        ci = jb(answer)
        size = len(ci)
        for i in range(0, size):
            if ci[i] == "1" or ci[i] == "1型" or ci[i] == "1型糖尿病" or ci[i] == "一型糖尿病" or ci[i] == "一型":
                pp = 1
            elif ci[i] == "2" or ci[i] == "2型" or ci[i] == "2型糖尿病" or ci[i] == "二型糖尿病" or ci[i] == "二型":
                pp = 2
            elif ci[i] == "妊娠" or ci[i] == "妊娠糖尿病":
                pp = 3
            elif ci[i] == "没有":
                pp = 4  # 需要判断
        data = {'style': 'type', 'num': pp}

    elif flag == 7:
        pp = 0
        answer = msg
        ci = jb(answer)
        size = len(ci)
        for i in range(0, size):
            if ci[i] == "无" or ci[i] == "没有" or ci[i] == "不知道" or ci[i] == "不":
                pp = 0
            else:
                pp = 1
        data = {'style': 'bing', 'num': pp}

    elif flag == 8:
        pp = 0
        answer = msg
        ci = jb(answer)
        size = len(ci)
        for i in range(0, size):
            if ci[i] == "是" or ci[i] == "是的" or ci[i] == "嗯" or ci[i] == "恩" or ci[i] == "稳定":
                pp = 0
            elif ci[i] == "高" or ci[i] == "餐前高" or ci[i] == "前高" or ci[i] == "空腹高" or ci[i] == "空腹" or ci[i] == "腹高":
                pp = 1
            elif ci[i] == "高" or ci[i] == "餐后高" or ci[i] == "后高":
                pp = 2
            elif ci[i] == "低" or ci[i] == "低血糖":
                pp = 3
        data = {'style': 'xue', 'num': pp}

    elif flag == 9:
        answer = msg
        ci = jb(answer)
        size = len(ci)
        for i in range(0, size):
            if ci[i] == "无" or ci[i] == "没有" or ci[i] == "不知道" or ci[i] == "不在" or ci[i] == "不":
                pp = 0
            else:
                pp = 1
        data = {'style': 'fuyong', 'num': pp}

    elif flag == 10:
        answer = msg
        data = {'style': 'ask', 'num': 1}
    elif flag == 11:
        answer = msg
        data = {'style': 'ask', 'num': 1}
    elif flag == 12:
        answer = msg
        data = {'style': 'ask', 'num': 1}

    elif flag == 13:
        pp = 0
        answer = msg
        ci = jb(answer)
        size = len(ci)
        for i in range(0, size):
            if ci[i] == "岁" or ci[i] == "年纪" or ci[i] == "年龄":
                sui1 = ci[i - 1]
                if i == size - 1:
                    sui2 = ci[i - 1]
                else:
                    sui2 = ci[i + 1]
                if is_number(sui1) is True:
                    pp = int(sui1)
                elif is_number(sui2) is True:
                    pp = int(sui2)
            else:
                pp = int(ci[i])
        if pp <= 20:
            data = {'style': 'judge1', 'num': 1}
        else:
            data = {'style': 'judge2', 'num': 1}
    elif flag == 14:
        pp = 0
        answer = msg
        ci = jb(answer)
        size = len(ci)
        for i in range(0, size):
            if ci[i] == "是" or ci[i] == "是的" or ci[i] == "明显" or ci[i] == "嗯" or ci[i] == "恩":
                pp = 1
            else:
                pp = 2
        if pp == 1:
            data = {'style': 'judge1', 'num': 1}
        else:
            data = {'style': 'judge2', 'num': 1}
    elif flag == 15:
        pp = 0
        answer = msg
        ci = jb(answer)
        size = len(ci)
        for i in range(0, size):
            if ci[i] == "没有" or ci[i] == "几乎" or ci[i] == "零" or ci[i] == "无":
                pp = 1
            else:
                pp = 2
        if pp == 1:
            data = {'style': 'judge1', 'num': 1}
        else:
            data = {'style': 'judge2', 'num': 1}
    elif flag == 16:
        pp = 0
        answer = msg
        ci = jb(answer)
        size = len(ci)
        for i in range(0, size):
            if ci[i] == "快" or ci[i] == "很快" or ci[i] == "迅速" or ci[i] == "速度快":
                pp = 1
            else:
                pp = 2
        if pp == 1:
            data = {'style': 'judge1', 'num': 1}
        else:
            data = {'style': 'judge2', 'num': 1}
    elif flag == 17:
        pp = 0
        answer = msg
        ci = jb(answer)
        size = len(ci)
        for i in range(0, size):
            if ci[i] == "是" or ci[i] == "是的" or ci[i] == "有" or ci[i] == "稍微" or ci[i] == "嗯" or ci[i] == "恩":
                pp = 1
            else:
                pp = 2
        if pp == 1:
            data = {'style': 'judge1', 'num': 1}
        else:
            data = {'style': 'judge2', 'num': 1}
    elif flag == 18:
        pp = 0
        answer = msg
        ci = jb(answer)
        size = len(ci)
        for i in range(0, size):
            if ci[i] == "是" or ci[i] == "是的" or ci[i] == "有" or ci[i] == "稍微" or ci[i] == "嗯" or ci[i] == "恩":
                pp = 1
            else:
                pp = 2
        if pp == 1:
            data = {'style': 'judge1', 'num': 1}
        else:
            data = {'style': 'judge2', 'num': 1}
    return data


if __name__ == '__main__':
    rt = sys.argv[1]
    if rt == 'q':
        ask_id = int(sys.argv[2])
        print(ask(ask_id))
    else:
        question_flag = int(sys.argv[2])
        question_answer = sys.argv[3]
        print(json.dumps(question(question_flag, question_answer)))

