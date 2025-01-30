/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 20:16
@Name: service_user.go
@Description:
*/

package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/dao"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
	"zhku-oj-server/pkg/utils/myUtils"
)

// 通过id更新数据
func UpdateUser(user *models.User) (status string, mes string, data string) {
	lg := utils.GetDefaultLogger()
	dao := dao.NewDao()
	selector := bson.M{
		"_id": user.ID,
	}
	hashPassword, _ := myUtils.HashPassword(user.Password) //hash加密
	update := bson.D{
		{"$set", bson.D{
			{"password", hashPassword},
			{"email", user.Email},
			{"phone", user.Phone},
			{"username", user.Username},
			{"nickname", user.Nickname},
			{"mtime", time.Now().Unix()},
		}},
	}
	lg.Infof("更新_id:%s\nselector:%s\n", user.ID, selector)
	err := dao.UpdateUser(context.Background(), selector, update)
	if err != nil {
		lg.Info("更新错误：", err)
		return utils.StatusFail, "更新失败", ""
	}
	return utils.StatusSuccess, "更新成功", ""
}

// 通过id删除用户
func DeleteUser(id string) (status string, mes string, data string) {
	lg := utils.GetDefaultLogger()
	//1.删之前先查询是否有该条数据
	dao := dao.NewDao()
	objectId, err := primitive.ObjectIDFromHex(id)
	selector := bson.M{
		"_id": objectId,
	}
	lg.Infof("查询_id: %s", objectId)

	//2.调用Dao查询user
	user, err := dao.GetOneUser(context.Background(), selector)
	if err != nil {
		lg.Info("查询用户异常: %v", err)
		return utils.StatusFail, "不存在该用户", ""
	}
	//3.不能删管理员
	if user.Role == utils.StatusAdmin {
		return utils.StatusFail, "不能删除管理员", ""
	}

	//4.删除
	lg.Infof("删除id%s", id)
	err = dao.DeleteUser(context.Background(), selector)
	if err != nil {
		lg.Info("删除失败", err)
		return utils.StatusFail, "删除失败！", ""
	}
	lg.Info("删除成功")
	return utils.StatusSuccess, "删除成功！", ""
}

// 通过id查用户
func GetUserById(id string) (status string, mes string, data string) {
	lg := utils.GetDefaultLogger()
	dao := dao.NewDao()
	objectId, _ := primitive.ObjectIDFromHex(id)
	selector := bson.M{
		"_id": objectId,
	}
	lg.Infof("查询_id: %s", objectId)
	//调用Dao查询user
	user, err := dao.GetOneUser(context.Background(), selector)
	if err != nil || user == nil {
		lg.Info("查询错误：", err)
		return utils.StatusFail, "用户不存在！", ""
	}
	//把结构体转换为json响应
	jsonData, _ := myUtils.StructToJSON(user)
	return utils.StatusSuccess, "查询成功！", jsonData
}

// 查一堆数据
func (s *Service) GetUserList(comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	items, err = s.dao.GetUserList(context.Background(), comQuery)
	return
}

// 注册
// TODO 校验邮箱，电话，密码的格式
func PostUser(postUser *dto.ReqPostUser) (status string, mes string) {
	lg := utils.GetDefaultLogger()
	//1.判空
	if postUser.Username == "" || postUser.Password == "" {
		return utils.StatusFail, "账号或密码不能为空！"
	}
	//2.检查用户名是否存在
	dao := dao.NewDao()
	query := bson.M{
		"username": postUser.Username,
	}
	user, _ := dao.GetOneUser(context.Background(), query)
	if user != nil {
		lg.Infof("用户名%s已被注册", postUser.Username)
		return utils.StatusFail, "用户名已被注册"
	}
	//3.未被注册，一切正常
	hashPassword, _ := myUtils.HashPassword(postUser.Password) //hash加密
	user = &models.User{
		Username: postUser.Username,
		Password: hashPassword,
		Email:    postUser.Email,
		Phone:    postUser.Phone,
		Status:   utils.StatusNormal,
		Ctime:    time.Now().Unix(),
		Mtime:    time.Now().Unix(),
	}
	_, err := dao.CreateUser(context.Background(), user)
	if err != nil {
		lg.Infof("服务器异常:%v", err)
		return utils.StatusFail, "服务器异常"
	}
	return utils.StatusSuccess, "注册成功！"
}

// 登录
func UserLogin(loginUser *dto.ReqPostLoginUser) (msg string, token string, mes string) {
	lg := utils.GetDefaultLogger()
	//1.判空
	if loginUser.Username == "" || loginUser.Password == "" {
		lg.Info("账号或密码为空")
		return utils.StatusFail, "", "账号或密码不能为空！"
	}
	//2.查询用户是否存在
	dao := dao.NewDao()
	query := bson.M{
		"username": loginUser.Username,
	}
	user, err := dao.GetOneUser(context.Background(), query)
	if err != nil {
		lg.Info("查询用户异常: %v", err)
		return utils.StatusFail, "", "账号错误！"
	}
	//3.判断密码是否正确
	if !myUtils.CheckPasswordHash(loginUser.Password, user.Password) {
		return utils.StatusFail, "", "密码错误！"
	}
	//4.判断账号是否被封禁
	if user.Status == utils.StatusBanned {
		return utils.StatusFail, "", "账号已被封禁！"
	}
	//5.成功查询
	jwt, err := myUtils.GenerateToken(user, utils.JwtTokenSecretKey, time.Hour*24) //jwt令牌
	if err != nil {
		lg.Infof("异常信息：%v", err)
		return utils.StatusFail, "", "生成令牌失败！"
	}
	lg.Infof("jwt:%s", user.Username, jwt)
	return utils.StatusSuccess, jwt, "登录成功！"
}

//原登录逻辑
/*func (s *Service) UserLogin(loginUser *dto.ReqPostLoginUser) (token string, err error) {
	lg := utils.GetDefaultLogger()
	userModel, err := s.GetOneUser(map[string]interface{}{
		"name": loginUser.Username,
	})
	if err != nil {
		return
	}
	if err = utils.Encoder.Match([]byte(userModel.Password), []byte(loginUser.Password)); err != nil {
		lg.Errorf("Match: %v", err)
		return
	}
	token, err = utils.GenerateStringToken(userModel.ID.Hex(), userModel.Name, userModel.Role)
	if err != nil {
		lg.Errorf("GenerateToken: %v", err)
		return
	}
	return
}*/

/*func (s *Service) GetOneUser(query interface{}) (user *models.User, err error) {
	user, err = s.dao.GetOneUser(context.Background(), query)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}*/

/*func (s *Service) checkUsernameExists(username string) (bool, error) {
	var query = map[string]interface{}{
		"name": username,
	}
	m, err := s.dao.GetOneUser(context.Background(), query)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
	}
	if m != nil {
		return true, nil
	}
	return false, nil
}

func (s *Service) CheckPostUserParams(user *dto.ReqPostUser) error {
	// 重名校验
	ok, err := s.checkUsernameExists(user.Username)
	if err != nil {
		return err
	}
	if ok {
		return utils.ErrUsernameExisted
	}
	return nil
}

func (s *Service) CheckUserLoginParams(user *dto.ReqPostLoginUser) error {
	if user.Username == "" || user.Password == "" {
		return errors.New("username or password is empty")
	}
	return nil
}

func (s *Service) BuildPostUser(u *dto.ReqPostUser) (*models.User, error) {
	passwordHash, err := utils.Encoder.Encode([]byte(u.Password))
	return &models.User{
		Nickname: u.Username,
		Email:    u.Email,
		Password: string(passwordHash),
		Phone:    u.Phone,
		Role:     utils.DefaultRole,
		Ctime:    time.Now().Unix(),
		Mtime:    time.Now().Unix(),
	}, err

}

func (s *Service) CreateUser(user *models.User) (string, error) {
	return s.dao.CreateUser(context.Background(), user)
}*/
