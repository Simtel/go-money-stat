package cmd

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/cobra"
	"log"
)

func RunYearChart() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "year-chart",
		Short:     "Показать график доходов и расходов",
		ValidArgs: []string{"current", "last"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {

		if err := ui.Init(); err != nil {
			log.Fatalf("failed to initialize termui: %v", err)
		}
		defer ui.Close()

		bc := widgets.NewBarChart()
		bc.Data = []float64{3, 2, 5, 3, 9, 3}
		bc.Labels = []string{"S0", "S1", "S2", "S3", "S4", "S5"}
		bc.Title = "Bar Chart"
		bc.SetRect(5, 5, 100, 25)
		bc.BarWidth = 5
		bc.BarColors = []ui.Color{ui.ColorRed, ui.ColorGreen}
		bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorBlue)}
		bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorYellow)}

		ui.Render(bc)

		uiEvents := ui.PollEvents()
		for {
			e := <-uiEvents
			switch e.ID {
			case "q", "<C-c>":
				return nil
			}
		}

	}

	return cmd
}
