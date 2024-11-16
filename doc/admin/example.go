package admin

type Example struct {
}

type ExampleExampleFormat struct {
	Code int32             `json:"code"`
	Msg  string            `json:"msg"`
	Data ExampleExampleRes `json:"data"`
}

type ExampleExampleMsgsExample struct {
	Example_1 string   `json:"example_1"` //oss accessKeyId
	Example_2 int64    `json:"example_2"` //oss accessKeyId
	Example_3 []string `json:"example_3"` //oss accessKeyId
}

type ExampleExampleMsgsExample2 struct {
	Example_1 string                       `json:"example_1"` //oss accessKeyId
	Example_2 []string                     `json:"example_2"` //oss accessKeyId
	Example_3 []*ExampleExampleMsgsExample `json:"example_3"` //oss accessKeyId
	Example_4 *ExampleExampleMsgsExample   `json:"example_4"` //oss accessKeyId
}

type ExampleExampleReq struct {
	Id string `json:"id"  validate:"required"` //id
}

type ExampleExampleRes struct {
	MsgsExample []*ExampleExampleMsgsExample2 `json:"msgs_example"` //msgs_example
	Name        string                        `json:"name"`         //name
	Age         float64                       `json:"age"`          //age
	Games       []*ExampleExampleGame         `json:"games"`        //games
	NextPlay    *ExampleExampleGame           `json:"next_play"`    //games
	Play        ExampleExampleGame            `json:"play"`         //games
}

type ExampleExampleGame struct {
	Name string `json:"name"` //name
	Time int32  `json:"time"` //time
}

// Example
// @Tags 示例接口组
// @Summary 示例接口
// @Security apiKey
// @accept application/json
// @Produce application/json
// @Param data query ExampleExampleReq true "数据"
// @Success 200 {object} ExampleExampleFormat
// @Router /admin/example/example [GET]
func (Example) Example() {

}
