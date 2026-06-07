package server

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/shabykov/intervals-mcp/client"
)

type Server struct {
	mcpServer *mcp.Server
}

func NewServer(clientConfig client.Config) *Server {

	cl := client.NewClient(clientConfig)

	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "intervals-icu",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "list_activities",
		Description: "Список тренировок за период с ключевыми метриками (нагрузка, мощность, пульс, CTL/ATL). Параметры: oldest, newest (YYYY-MM-DD).",
	}, cl.ListActivities)
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_activity",
		Description: "Полные данные одной тренировки по id. intervals=true добавляет разбивку по интервалам/лапам."}, cl.GetActivity)
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_activity_streams",
		Description: "Посекундные ряды одной тренировки (watts/heartrate/cadence/velocity_smooth/altitude) для детального разбора. Параметры: id, types.",
	}, cl.GetStreams)
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_wellness",
		Description: "Wellness за период: CTL, ATL, форма (TSB), rampRate, HRV, пульс покоя, вес. Параметры: oldest, newest (YYYY-MM-DD).",
	}, cl.GetWellness)
	mcp.AddTool(mcpServer, &mcp.Tool{
		Name:        "get_athlete",
		Description: "Профиль атлета: FTP, зоны, настройки. Без параметров.",
	}, cl.GetAthlete)

	return &Server{mcpServer: mcpServer}
}

func (s *Server) Run(ctx context.Context) error {
	if err := s.mcpServer.Run(ctx, &mcp.StdioTransport{}); err != nil {
		return err
	}
	return nil
}
