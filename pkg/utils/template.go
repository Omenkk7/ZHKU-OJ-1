package utils

import "zhku-oj-server/pkg/models"

func MergeTemplate(choice int, task interface{}) (code interface{}) {
	lg := GetDefaultLogger()

	t := task.(models.LocalTask)
	lg.Infof("题目id：%v语言为：%v", t.ProblemId, t.Language)

	//TODO 通过ProblemId查询模板,然后合并

	//假设已经合并了完整代码，以下示例为java的两数之和
	code = "import java.util.Scanner; \n class a {\n   public static void main(String[] args){\n Scanner sc=new Scanner(System.in);\n System.out.print(sc.nextInt()+sc.nextInt());\n }\n \n}"
	return code
}
