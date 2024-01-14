package middleware

import (
	"context"
	"encoding/json"
	"gin/session/rds"
	"time"
)

const userTokenPrefix = "user_token_"
const sessionPrefix = "user_session_"
const defaultExpireTime = time.Minute * 30

type Session struct {
	AccessToken string `json:"access_token"`
	WorkId      string `json:"work_id"`
	Expire      string `json:"expire"`
	ExpireUnix  int64  `json:"expire_unix"`
}

func NewSession(workId string) *Session {
	expireTime := time.Now().Add(defaultExpireTime)
	var sess = Session{
		AccessToken: RandString(16), // 随机生成16位的token
		WorkId:      workId,
		Expire:      expireTime.Format("2006-01-02 15:04:05"),
		ExpireUnix:  expireTime.Unix(),
	}
	if err := SetSession(&sess); err != nil {
		return nil
	}
	return &sess
}

// 保存session信息到redis中
func SetSession(sess *Session) error {
	rds := rds.GetRedis()

	data, err := json.Marshal(sess)
	if err != nil {
		return err
	}
	// set redis: key = token, value = marshal(session), ttl = 1*minute
	_, err = rds.Set(context.Background(), sessionPrefix+sess.AccessToken,
		string(data), defaultExpireTime).Result()
	if err != nil {
		return err
	}
	// set redis: key = workId, value = token, ttl = 1*minute
	_, err = rds.Set(context.Background(), userTokenPrefix+sess.WorkId,
		sess.AccessToken, defaultExpireTime).Result()
	if err != nil {
		return err
	}
	return nil
}

// 根据token获取session
func GetSession(token string) *Session {
	rds := rds.GetRedis()
	// 获取token对应的session
	data, err := rds.Get(context.Background(), sessionPrefix+token).Result()
	if err != nil {
		return nil
	}
	var sess Session
	// 反序列化session
	if err := json.Unmarshal([]byte(data), &sess); err != nil {
		return nil
	}
	return &sess
}

// 根据workId获取session
func GetUserSession(workId string) *Session {
	rds := rds.GetRedis()
	token, err := rds.Get(context.Background(), userTokenPrefix+workId).Result()
	if err != nil {
		return nil
	}
	return GetSession(token)
}
