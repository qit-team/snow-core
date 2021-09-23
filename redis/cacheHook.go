package redis

import (
	"context"
	"fmt"
	"github.com/SkyAPM/go2sky"
	v3 "github.com/SkyAPM/go2sky/reporter/grpc/language-agent"
	goredis "github.com/go-redis/redis/v8"
)

const (
	ComponentIDRedis int32 = 7
)

type SkyWalkingHook struct {
	tracer *go2sky.Tracer
}

type Go2skyKey interface{}

func NewSkyWalkingHook(tracer *go2sky.Tracer) *SkyWalkingHook {
	return &SkyWalkingHook{tracer: tracer}
}

func (h *SkyWalkingHook) BeforeProcess(ctx context.Context, cmd goredis.Cmder) (context.Context, error) {
	peer := "redis"
	if p, ok := ctx.Value("peer").(string); ok {
		peer = p
	}
	args := fmt.Sprintf("%v", cmd.Args())
	span, err := h.tracer.CreateExitSpan(ctx, fmt.Sprintf("%v %v", cmd.Name(), args), peer, func(header string) error {
		return nil
	})
	if err != nil {
		return nil, err
	}
	span.SetComponent(ComponentIDRedis)
	span.Tag("args", args)
	span.SetSpanLayer(v3.SpanLayer_Cache)
	var key Go2skyKey = fmt.Sprintf("%v %v", cmd.Name(), args)
	c := context.WithValue(ctx, key, span)
	return c, nil
}

func (h *SkyWalkingHook) AfterProcess(c context.Context, cmd goredis.Cmder) error {
	var key Go2skyKey = fmt.Sprintf("%v %v", cmd.Name(), cmd.Args())
	span := c.Value(key).(go2sky.Span)
	span.Tag("cache_results", cmd.String())
	span.End()
	return nil
}

func (h *SkyWalkingHook) BeforeProcessPipeline(ctx context.Context, cmds []goredis.Cmder) (context.Context, error) {
	peer := "redis"
	if p, ok := ctx.Value("peer").(string); ok {
		peer = p
	}
	pipelineInfo := ""
	cmdStr := ""
	for _, cmd := range cmds {
		pipelineInfo += fmt.Sprintf("%v %v", cmd.Name(), cmd.Args())
		cmdStr += " " + cmd.Name()
	}
	span, err := h.tracer.CreateExitSpan(ctx, pipelineInfo, peer, func(header string) error {
		return nil
	})
	if err != nil {
		return nil, err
	}
	span.SetComponent(ComponentIDRedis)
	span.Tag("pipeline", pipelineInfo)
	span.SetSpanLayer(v3.SpanLayer_Cache)
	c := context.WithValue(ctx, cmdStr, span)
	return c, nil
}

func (h *SkyWalkingHook) AfterProcessPipeline(c context.Context, cmds []goredis.Cmder) error {
	pipelineInfo := ""
	cmdStr := ""
	for _, cmd := range cmds {
		pipelineInfo += fmt.Sprintf("%v %v \n", cmd.Name(), cmd.Args())
		cmdStr += " " + cmd.Name()
	}
	span := c.Value(cmdStr).(go2sky.Span)
	span.Tag("cache_results", pipelineInfo)
	span.End()
	return nil
}
