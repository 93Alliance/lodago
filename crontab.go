package lodago

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// 对 https://github.com/robfig/cron 库的补充，由于这个库不支持一次性事件，
// 并且前端对于 * * * * * 格式的时间设定不友好，所以做了一些简单封装。

// ScheduleType 时间表类型
type ScheduleType int32

// 定时任务的时间表类型
const (
	Yearly        ScheduleType = 1 // 每年
	Monthly       ScheduleType = 2 // 每月
	Weekly        ScheduleType = 3 // 每周
	Daily         ScheduleType = 4 // 每天
	Hourly        ScheduleType = 5 // 每小时
	IntervalMonth ScheduleType = 6 // 每隔几个月
	IntervalDay   ScheduleType = 7 // 每隔几天
	Every         ScheduleType = 8 // 间隔时间，只支持[时][分]
	Once          ScheduleType = 9 // 一次性
)

// Job 任务
type Job func()

// Crontab 定时任务调度器
type Crontab struct {
	cron     *cron.Cron
	entryIDs map[string]cron.EntryID
	locker   sync.RWMutex
}

// NewCrontab 创建定时器
func NewCrontab() *Crontab {
	return &Crontab{
		cron.New(),
		make(map[string]cron.EntryID),
		sync.RWMutex{},
	}
}

// Start 启动
func (c *Crontab) Start() {
	c.cron.Start()
}

// Stop 停止
func (c *Crontab) Stop() {
	c.cron.Stop()
}

// AddJob 添加任务，返回值是job id，可以用于删除任务
func (c *Crontab) AddJob(cronTime *CronTime, job Job) (cron.EntryID, error) {
	spec, err := cronTime.ToSpec()
	if err != nil {
		return 0, err
	}
	cronTime.Key = RandString(12) // 12位的随机数字+大小写字母
	id, err := c.cron.AddFunc(spec, c.jobDecorate(*cronTime, job))
	c.setEntryID(cronTime.Key, id)
	return id, err
}

// RemoveJob 删除一个任务
func (c *Crontab) RemoveJob(id cron.EntryID) {
	c.cron.Remove(id)
}

// GetEntries 获得所有实体
func (c *Crontab) GetEntries() []cron.Entry {
	return c.cron.Entries()
}

// jobDecorate Job任务装饰器，主要用于解决周期性和一次性任务的执行逻辑不一样。
func (c *Crontab) jobDecorate(cronTime CronTime, job Job) Job {
	if cronTime.Type == Once { // 一次性事件单独处理
		return func() {
			// 假如当前时间是2020年3月12日12时23分，设定了一个2021年4月15日2时3分的任务，那么
			// 在2020年4月15日2时3分会执行这个函数，也就是t1时间，但是由于执行需要时间，导致t1会
			// 略比这个时间大，那么t1-t2就是一个负值，而当真正到了2021年4月15日2时3分，t1-t2应该
			// 是一个正值，所以得出结论：
			//  【时间差为正数代表需要执行，负数为不执行】
			if cronTime.isEver() {
				job() // 原先任务正常执行
				id, ok := c.getEntryID(cronTime.Key)
				if ok {
					c.RemoveJob(id)           // 删除这个job
					c.rmEntryID(cronTime.Key) // 删除这个key
				}
			}
		}
	}
	return job
}

// 设置id
func (c *Crontab) setEntryID(key string, id cron.EntryID) {
	c.locker.Lock()
	c.entryIDs[key] = id
	c.locker.Unlock()
}

// 移除id
func (c *Crontab) rmEntryID(key string) {
	c.locker.Lock()
	delete(c.entryIDs, key)
	c.locker.Unlock()
}

// 移除id
func (c *Crontab) getEntryID(key string) (cron.EntryID, bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()
	id, ok := c.entryIDs[key]
	return id, ok
}

// CronTime 时间结构
type CronTime struct {
	Type   ScheduleType `json:"type"`
	Year   string       `json:"year"`
	Month  string       `json:"month"`
	Day    string       `json:"day"`
	Hour   string       `json:"hour"`
	Minute string       `json:"minute"`
	Week   string       `json:"week"`
	Key    string       `json:"key"`
}

// ToSpec 转换成spec函数
// 【周期性】
// 【每年】 输入[月][日][时][分] -- 30 22 1 1 * 每年的1月1日22点30分执行 【已验证】
// 【每月】 输入[日][时][分] -- 30 22 1 * * 每月1日22点30分执行 【已验证】
// 【每天】 输入[时][分] -- 30 22 * * * 每天22点30分执行，【已验证】
// 【每周】 输入[星期][时][分] -- 30 22 * * 1 每周一22点30分执行 【已验证】
// 【每小时】输入[分] -- 15 * * * * 每小时的15分执行 【已验证】
// 【每隔几月】 输入[月][日][时][分] -- 30 22 3 */3 *
// 【每隔几天】 输入[日][时][分] -- 30 22 */3 * * 每隔3天的22点30分执行
// 【每隔小时】 输入[时][分] -- @every 1h30m 每隔1小时30分执行
// 【一次性】 输入[年][月][日][时][分] 由于cron不支持一次性任务，所以只能通过周期性时间删除自身解决。
func (c *CronTime) ToSpec() (string, error) {
	switch c.Type {
	case Yearly:
		if !c.isNums(c.Month, c.Day, c.Hour, c.Minute) {
			return "", errors.New("Time format is error")
		}
		return fmt.Sprintf("%s %s %s %s *", c.Minute, c.Hour, c.Day, c.Month), nil
	case Monthly:
		if !c.isNums(c.Day, c.Hour, c.Minute) {
			return "", errors.New("Time format is error")
		}
		return fmt.Sprintf("%s %s %s * *", c.Minute, c.Hour, c.Day), nil
	case Weekly:
		if !c.isNums(c.Week, c.Hour, c.Minute) {
			return "", errors.New("Time format is error")
		}
		return fmt.Sprintf("%s %s * * %s", c.Minute, c.Hour, c.Week), nil
	case Daily:
		if !c.isNums(c.Hour, c.Minute) {
			return "", errors.New("Time format is error")
		}
		return fmt.Sprintf("%s %s * * *", c.Minute, c.Hour), nil
	case Hourly:
		if !c.isNums(c.Minute) {
			return "", errors.New("Time format is error")
		}
		return fmt.Sprintf("%s * * * *", c.Minute), nil
	case IntervalMonth:
		if !c.isNums(c.Month, c.Day, c.Hour, c.Minute) {
			return "", errors.New("Time format is error")
		}
		return fmt.Sprintf("%s %s %s */%s *", c.Minute, c.Hour, c.Day, c.Month), nil
	case IntervalDay:
		if !c.isNums(c.Day, c.Hour, c.Minute) {
			return "", errors.New("Time format is error")
		}
		return fmt.Sprintf("%s %s */%s * *", c.Minute, c.Hour, c.Day), nil
	case Every:
		if !c.isNums(c.Hour, c.Minute) {
			return "", errors.New("Time format is error")
		}
		return fmt.Sprintf("@every %sh%sm", c.Minute, c.Hour), nil
	case Once:
		if !c.isNums(c.Year, c.Month, c.Day, c.Hour, c.Minute) || c.isEver() {
			return "", errors.New("Time format is error")
		}
		return fmt.Sprintf("%s %s %s %s *", c.Minute, c.Hour, c.Day, c.Month), nil
	default:
		return "", errors.New("Schedule type is error")
	}
}

// 判断一些字符串是否都是整数
func (c *CronTime) isNums(strs ...string) bool {
	for _, str := range strs {
		if !IsNum(str) {
			return false
		}
	}
	return true
}

// 判断时间是否已经过去，例如 CronTime 中的时间比现在的时间要早
func (c *CronTime) isEver() bool {
	t1 := time.Now()
	year, _ := strconv.Atoi(c.Year)
	month, _ := strconv.Atoi(c.Month)
	day, _ := strconv.Atoi(c.Day)
	hour, _ := strconv.Atoi(c.Hour)
	minute, _ := strconv.Atoi(c.Minute)
	t2 := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local)
	sub := t1.Sub(t2)
	if sub.Seconds() >= 0 {
		return true
	}
	return false
}
