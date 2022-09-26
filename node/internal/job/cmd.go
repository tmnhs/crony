package job

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/logger"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

type CMDHandler struct {
}

func (c *CMDHandler) Run(job *Job) (result string, err error) {
	var (
		cmd         *exec.Cmd
		proc        *JobProc
		sysProcAttr *syscall.SysProcAttr
	)

	//todo 设置属性
	//sysProcAttr, err = j.CreateCmdAttr()
	//if err != nil {
	//	j.Fail(t, err.Error())
	//	return false
	//}

	// 超时控制
	if job.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(job.Timeout)*time.Second)
		defer cancel()
		cmd = exec.CommandContext(ctx, job.Cmd[0], job.Cmd[1:]...)
	} else {
		cmd = exec.Command(job.Cmd[0], job.Cmd[1:]...)
	}

	cmd.SysProcAttr = sysProcAttr
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b

	err = cmd.Start()
	result = b.String()
	if err != nil {
		logger.Fatal(fmt.Sprintf("%s\n%s", b.String(), err.Error()))
		return
	}
	//todo 正在运行的任务

	proc = &JobProc{
		JobProc: &models.JobProc{
			ID:     strconv.Itoa(cmd.Process.Pid),
			JobID:  job.ID,
			Group:  job.Group,
			NodeID: job.RunOn,
			JobProcVal: models.JobProcVal{
				Time: time.Now(),
			},
		},
	}

	err = proc.Start()
	if err != nil {
		return
	}
	defer proc.Stop()

	if err := cmd.Wait(); err != nil {
		logger.Fatal(fmt.Sprintf("%s\n%s", b.String(), err.Error()))
		return
	}
	// todo 将结果写入数据库
	//j.Success(t, b.String())
	return b.String(), nil
}
