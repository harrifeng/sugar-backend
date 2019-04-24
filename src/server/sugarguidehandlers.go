package server

import (
	"bytes"
	"db"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type Answer struct {
	Style string `json:"style"`
	Num   int    `json:"num"`
}

type SugarGuider struct {
	Pos    int
	Flag   int
	Ws     *websocket.Conn
	UserId int
	Result map[string]int
	Li     []int
	LeftQs int
}

func NewSugarGuider(ws *websocket.Conn, userId int) *SugarGuider {
	return &SugarGuider{
		Pos:    31,
		Flag:   0,
		Ws:     ws,
		UserId: userId,
		Result: map[string]int{
			"sex":    -1,
			"age":    -1,
			"height": -1,
			"weight": -1,
			"sport":  -1,
			"type":   -1,
			"bing":   -1,
			"xue":    -1,
			"fuyong": -1,
			"ask":    0,
			"judge1": 0,
			"judge2": 0,
		},
		Li: []int{-1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1},
	}
}

func (sg *SugarGuider) nowPos() string {
	return strconv.Itoa(sg.Pos)
}

func (sg *SugarGuider) setPos(pos int) {
	sg.Pos = pos
}

func (sg *SugarGuider) randNextQuestion() bool {
	flag := false
	for i := 0; i < len(sg.Li); i++ {
		if sg.Li[i] == 0 {
			flag = true
		}
	}
	if !flag {
		return flag
	}
	for sg.Pos = rand.Intn(18) + 1; sg.Li[sg.Pos] == 1; sg.Pos = rand.Intn(18) + 1 {
	}
	return true
}

func (sg *SugarGuider) isEnd() bool {
	return sg.Result["sex"] != -1 && sg.Result["age"] != -1 && sg.Result["height"] != -1 && sg.Result["weight"] != -1 &&
		sg.Result["sport"] != -1 && sg.Result["type"] != -1 && sg.Result["bing"] != -1 && sg.Result["xue"] != -1 &&
		sg.Result["fuyong"] != -1
}

func (sg *SugarGuider) query() error {
	question, err := sg.getQuestion()
	if err != nil {
		return err
	}
	return sg.Ws.WriteJSON(gin.H{
		"questionId": sg.nowPos(),
		"msg":        string(question),
	})
}

func (sg *SugarGuider) answer() error {
	_, message, err := sg.Ws.ReadMessage()
	if err != nil {
		return err
	}
	ans, err := sg.parseAnswer(string(message))
	if err != nil {
		return err
	}

	if ans.Style == "judge1" || ans.Style == "judge2" {
		sg.Result[ans.Style]++
	} else {
		sg.Result[ans.Style] = ans.Num
	}
	sg.Li[sg.Pos] = 1

	if sg.Result["type"] == 4 && sg.Flag == 0 {
		for i := 13; i <= 18; i++ {
			sg.Li[i] = 0
			sg.Flag = 1
		}
	}
	return nil
}

func (sg *SugarGuider) getQuestion() (string, error) {
	sugarDir, err := sugarGuideDir()
	if err != nil {
		return "", err
	}
	cmd := exec.Command(filepath.Join(sugarDir, "venv", "bin", "python"), filepath.Join(sugarDir, "guide.py"), "q", strconv.Itoa(sg.Pos))
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func (sg *SugarGuider) parseAnswer(ans string) (Answer, error) {
	sugarDir, err := sugarGuideDir()
	if err != nil {
		return Answer{}, err
	}
	cmd := exec.Command(filepath.Join(sugarDir, "venv", "bin", "python"), filepath.Join(sugarDir, "guide.py"), "a", strconv.Itoa(sg.Pos), ans)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return Answer{}, err
	}
	var result Answer
	err = json.Unmarshal(out.Bytes(), &result)
	if err != nil {
		return Answer{}, err
	}
	return result, nil
}

func (sg *SugarGuider) calcPlan(t string) ([]byte, error) {
	sugarDir, err := sugarGuideDir()
	if err != nil {
		return nil, err
	}
	var cmd *exec.Cmd
	switch t {
	case "yinshi":
		cmd = exec.Command(filepath.Join(sugarDir, "venv", "bin", "python"), filepath.Join(sugarDir, "exp.py"),
			"yinshi", strconv.Itoa(sg.Result["sex"]), strconv.Itoa(sg.Result["age"]), strconv.Itoa(sg.Result["height"]),
			strconv.Itoa(sg.Result["weight"]), strconv.Itoa(sg.Result["sport"]))
	case "yundong":
		cmd = exec.Command(filepath.Join(sugarDir, "venv", "bin", "python"), filepath.Join(sugarDir, "exp.py"),
			"yundong", strconv.Itoa(sg.Result["sex"]), strconv.Itoa(sg.Result["age"]), strconv.Itoa(sg.Result["type"]),
			strconv.Itoa(sg.Result["bing"]))
		fmt.Println(cmd.Args)
	case "kongtang":
		cmd = exec.Command(filepath.Join(sugarDir, "venv", "bin", "python"), filepath.Join(sugarDir, "exp.py"),
			"kongtang", strconv.Itoa(sg.Result["sex"]), strconv.Itoa(sg.Result["age"]), strconv.Itoa(sg.Result["type"]),
			strconv.Itoa(sg.Result["fuyong"]), strconv.Itoa(sg.Result["bing"]), strconv.Itoa(sg.Result["xue"]))
	default:
		return nil, errors.New("no type of plan")
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	fmt.Println(t, out.String())
	return out.Bytes(), nil
}

func sugarGuideDir() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "sugarguide"), nil
}

func (sg *SugarGuider) analysePlan() error {
	err := sg.Ws.WriteJSON(gin.H{
		"questionId": "9998",
		"msg":        "糖导正在分析您的情况，请稍等....",
	})
	if err != nil {
		return err
	}
	dietPlan, err := sg.calcPlan("yinshi")
	if err != nil {
		return err
	}
	sportPlan, err := sg.calcPlan("yundong")
	if err != nil {
		return err
	}
	controlPlan, err := sg.calcPlan("kongtang")
	if err != nil {
		return err
	}
	err = db.SaveSugarGuidePlan(sg.UserId, dietPlan, sportPlan, controlPlan)
	if err != nil {
		return err
	}
	return sg.Ws.WriteJSON(gin.H{
		"questionId": "9999",
		"msg":        "分析完成，已经给你发送给你一份健康周报，记得在糖家查收哦",
	})
}

func sugarGuideWebsocket(c *gin.Context) {
	userId, _ := c.Get("user_id")
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()
	//初始化糖导
	sugarGuider := NewSugarGuider(ws, userId.(int))
	if err := sugarGuider.query(); err != nil {
		fmt.Println(err)
		return
	}
	sugarGuider.randNextQuestion()
	if err := sugarGuider.query(); err != nil {
		fmt.Println(err)
		return
	}
	for {
		if err := sugarGuider.answer(); err != nil {
			fmt.Println(err)
			return
		}
		if sugarGuider.isEnd() {
			if sugarGuider.Result["type"] == 4 {
				if sugarGuider.Result["judge2"] >= sugarGuider.Result["judge1"] {
					sugarGuider.Result["type"] = 2
					sugarGuider.setPos(27)
					if err := sugarGuider.query(); err != nil {
						fmt.Println(err)
						return
					}
					sugarGuider.setPos(28)
					if err := sugarGuider.query(); err != nil {
						fmt.Println(err)
						return
					}
				} else {
					sugarGuider.Result["type"] = 1
					sugarGuider.setPos(27)
					if err := sugarGuider.query(); err != nil {
						fmt.Println(err)
						return
					}
					sugarGuider.setPos(29)
					if err := sugarGuider.query(); err != nil {
						fmt.Println(err)
						return
					}
				}
			}
			if err := sugarGuider.analysePlan(); err != nil {
				fmt.Println(err)
				return
			}
			break
		} else {
			sugarGuider.randNextQuestion()
			if err := sugarGuider.query(); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func getWeeklyNewspaper(userId int) responseBody {
	exist, err := db.CheckWeeklyNewspaper(userId)
	if err != nil {
		return responseInternalServerError(err)
	}
	if !exist {
		return responseNormalError("还没做糖导，尚未生成健康周报")
	}
	dietPlan, sportPlan, controlPlan, err := db.GetWeeklyNewspaper(userId)
	if err != nil {
		return responseInternalServerError(err)
	}
	bloodRecord, exist, err := db.GetBloodSugarRecordFromRecordDate(userId, time.Now())
	if err != nil {
		return responseInternalServerError(err)
	}
	bloodMap := make(gin.H)
	if !exist {
		bloodMap["0"] = "0"
		bloodMap["1"] = "0"
		bloodMap["2"] = "0"
		bloodMap["3"] = "0"
		bloodMap["4"] = "0"
		bloodMap["5"] = "0"
		bloodMap["6"] = "0"
	} else {
		err = json.Unmarshal([]byte(bloodRecord.Level), &bloodMap)
		if err != nil {
			return responseInternalServerError(err)
		}
	}
	return responseOKWithData(gin.H{
		"diet": gin.H{
			"change":     dietPlan.Change,
			"cereals":    dietPlan.Cereals,
			"fruit":      dietPlan.Fruit,
			"meat":       dietPlan.Meat,
			"milk":       dietPlan.Milk,
			"fat":        dietPlan.Fat,
			"vegetables": dietPlan.Vegetables,
		},
		"sport": gin.H{
			"sport1": sportPlan.Sport1,
			"sport2": sportPlan.Sport2,
			"sport3": sportPlan.Sport3,
			"sport4": sportPlan.Sport4,
			"time1":  sportPlan.Time1,
			"time2":  sportPlan.Time2,
			"time3":  sportPlan.Time3,
			"time4":  sportPlan.Time4,
			"week1":  sportPlan.Week1,
			"week2":  sportPlan.Week2,
			"week3":  sportPlan.Week3,
			"week4":  sportPlan.Week4,
		},
		"control": gin.H{
			"min1":   controlPlan.Min1,
			"min2":   controlPlan.Min2,
			"max1":   controlPlan.Max1,
			"max2":   controlPlan.Max2,
			"sleep1": controlPlan.Sleep1,
			"sleep2": controlPlan.Sleep2,
		},
		"level": bloodMap,
	})
}
