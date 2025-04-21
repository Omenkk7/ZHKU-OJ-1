package dto

type ReqProblem struct {
	ID            string   `json:"_id,omitempty" bson:"_id,omitempty"`
	Creator       string   `json:"creator,omitempty" bson:"creator,omitempty"`
	Description   string   `json:"description,omitempty" bson:"description,omitempty"`
	InputExample  string   `json:"input_example,omitempty" bson:"input_example,omitempty"`
	OutputExample string   `json:"output_example,omitempty" bson:"output_example,omitempty"`
	Status        int32    `json:"status,omitempty" bson:"status,omitempty"`
	SubmitNum     int32    `json:"submit_num,omitempty" bson:"submit_num,omitempty"`
	Difficulty    string   `json:"difficulty,omitempty" bson:"difficulty,omitempty"`
	Labels        []string `json:"labels,omitempty" bson:"labels,omitempty"`
	PassNum       int32    `json:"pass_num,omitempty" bson:"pass_num,omitempty"`
	Title         string   `json:"title,omitempty" bson:"title,omitempty"`
	URL           string   `json:"url,omitempty" bson:"url,omitempty"`
	Ctime         int64    `json:"ctime,omitempty" bson:"ctime,omitempty"`
	Mtime         int64    `json:"mtime,omitempty" bson:"mtime,omitempty"`

	//新增文件上传的字段，为了加速开发，先暂时使用JSON传输测试用例和不同语言的模板
	TestExample    string `json:"testExample,omitempty" bson:"testExample,omitempty"`
	JavaTemplate   string `json:"javaTemplate,omitempty" bson:"javaTemplate,omitempty"`
	GoTemplate     string `json:"goTemplate,omitempty" bson:"goTemplate,omitempty"`
	PythonTemplate string `json:"pythonTemplate,omitempty" bson:"pythonTemplate,omitempty"`
}
