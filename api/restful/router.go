package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/flow/internal/flow"
	"github.com/quanxiang-cloud/flow/internal/mq"
	"github.com/quanxiang-cloud/flow/internal/server/options"
	"github.com/quanxiang-cloud/flow/pkg/config"
	"github.com/quanxiang-cloud/flow/pkg/misc/logger"
	"github.com/quanxiang-cloud/flow/pkg/misc/mysql2"
)

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
	// ServerPath server path
	ServerPath = "/api/v1/flow"
)

// Router 路由
type Router struct {
	c *config.Configs

	engine *gin.Engine
}

// NewRouter 开启路由
func NewRouter(c *config.Configs) (*Router, error) {
	engine, err := newRouter(c)
	if err != nil {
		return nil, err
	}
	db, err := mysql2.New(c.Mysql, logger.Logger)
	if err != nil {
		return nil, err
	}
	optDB := options.WithDB(db)
	triggerRule, err := flow.NewTriggerRule(c, optDB)
	if err != nil {
		return nil, err
	}
	engine.Any("/send", mq.Subscription(triggerRule))
	//Flow router
	flow, err := NewFlow(c, optDB)
	if err != nil {
		return nil, err
	}
	v1 := engine.Group(ServerPath)
	{
		v1.POST("/deleteFlow/:id", flow.deleteFlow)
		v1.POST("/flowInfo/:ID", flow.info)
		v1.POST("/flowList", flow.flowList)
		v1.POST("/saveFlow", flow.saveFlow)
		v1.POST("/getNodes/:ID", flow.getNodes)

		v1.POST("/getVariableList", flow.getVariableList)
		v1.POST("/saveFlowVariable", flow.saveFlowVariable)
		v1.POST("/deleteFlowVariable/:ID", flow.deleteFlowVariable)

		v1.POST("/updateFlowStatus", flow.updateFlowStatus)

		v1.POST("/correlationFlowList", flow.correlationFlowList)
		v1.POST("/deleteApp", flow.deleteApp)
		v1.POST("/copyFlow/:ID", flow.copyFlow)

		v1.POST("/appReplicationExport", flow.appReplicationExport)
		v1.POST("/appReplicationImport", flow.appReplicationImport)
	}

	// Instance router
	instance, err := NewInstance(c, optDB)
	if err != nil {
		return nil, err
	}
	v2 := engine.Group(ServerPath + "/instance")
	{
		v2.POST("/myApplyList", instance.myApplyList)
		v2.POST("/waitReviewList", instance.waitReviewList)
		v2.POST("/reviewedList", instance.reviewedList)
		v2.POST("/ccToMeList", instance.ccToMeList)
		v2.POST("/allList", instance.allList)

		v2.POST("/flowInfo/:processInstanceID", instance.flowInfo)
		v2.POST("/cancel/:processInstanceID", instance.cancel)
		v2.POST("/sendBack/:processInstanceID/:taskID", instance.sendBack)
		v2.POST("/resubmit/:processInstanceID", instance.resubmit)
		v2.POST("/stepBack/:processInstanceID/:taskID", instance.stepBack)
		v2.POST("/stepBackActivityList/:processInstanceID", instance.stepBackNodes)
		v2.POST("/ccFlow/:processInstanceID/:taskID", instance.ccFlow)
		v2.POST("/readFlow/:processInstanceID/:taskID", instance.readFlow)
		v2.POST("/handleCc", instance.handleCc)
		v2.POST("/handleRead/:processInstanceID/:taskID", instance.handleRead)
		v2.POST("/addSign/:processInstanceID/:taskID", instance.addSign)
		v2.POST("/deliverTask/:processInstanceID/:taskID", instance.deliverTask)
		v2.POST("/getFlowInstanceCount", instance.flowInstanceCount)
		v2.POST("/reviewTask/:processInstanceID/:taskID", instance.reviewTask)
		v2.POST("/getFlowInstanceForm/:processInstanceID", instance.getFlowInstanceForm)
		v2.POST("/getFormData/:processInstanceID/:taskID", instance.getFormData)
		v2.POST("/processHistories/:processInstanceID", instance.processHistories)
	}

	// comment router
	comment, err := NewComment(c, optDB)
	if err != nil {
		return nil, err
	}
	v3 := engine.Group(ServerPath + "/comment")
	{
		v3.POST("/addComment", comment.addComment)
		v3.POST("/getComments/:flowInstanceID/:taskID", comment.getComments)
	}

	// abnormal router
	abnormal, err := NewAbnormalTask(c, optDB)
	if err != nil {
		return nil, err
	}
	v4 := engine.Group(ServerPath + "/abnormalTask")
	{
		v4.POST("/list", abnormal.list)
		v4.POST("/adminStepBack/:processInstanceID/:taskID", abnormal.adminStepBack)
		v4.POST("/adminSendBack/:processInstanceID/:taskID", abnormal.adminSendBack)
		v4.POST("/adminAbandon/:processInstanceID/:taskID", abnormal.adminAbandon)
		v4.POST("/adminDeliverTask/:processInstanceID/:taskID", abnormal.adminDeliverTask)
		v4.POST("/adminGetTaskForm/:processInstanceID/:taskID", abnormal.adminGetTaskForm)
	}

	// handout router
	handOut, err := NewHandOut(c, optDB)
	if err != nil {
		return nil, err
	}
	v5 := engine.Group("/api/v1")
	{
		v5.POST("/handout/:code", handOut.handOut)
	}

	// urge router
	urge, err := NewUrge(c, optDB)
	if err != nil {
		return nil, err
	}
	v6 := engine.Group(ServerPath + "/urge")
	{
		v6.POST("/taskUrge", urge.taskUrge)
	}

	formula, err := NewFormula(c, optDB)
	if err != nil {
		return nil, err
	}
	v7 := engine.Group(ServerPath + "/formula")
	{
		v7.POST("/calculation", formula.Calculation)
	}

	return &Router{
		c:      c,
		engine: engine,
	}, nil
}

func newRouter(c *config.Configs) (*gin.Engine, error) {
	if c.Model == "" || (c.Model != ReleaseMode && c.Model != DebugMode) {
		c.Model = ReleaseMode
	}
	gin.SetMode(c.Model)
	engine := gin.New()
	// engine.Use(logger.GinLogger(), logger.GinRecovery())
	engine.Use(Recover)
	return engine, nil
}

// Run 启动服务
func (r *Router) Run() {
	r.engine.Run(r.c.Port)
}

// Close 关闭服务
func (r *Router) Close() {
}
