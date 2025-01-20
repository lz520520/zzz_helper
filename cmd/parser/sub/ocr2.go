package sub

import (
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"zzz_helper/internal/mylog"
	"zzz_helper/internal/utils/file2"
)

var (
	ocrContent = ""
)

var ocr2Cmd = &cobra.Command{
	Use:   "ocr2",
	Short: "parser ocr content without use ocr",
	RunE: func(cmd *cobra.Command, args []string) error {
		disk, err := ParseDriveInfo(ocrContent)
		if err != nil {
			mylog.CommonLogger.Err(err).Send()
		}
		m := map[string]interface{}{
			uuid.New().String(): disk,
		}
		result, _ := yaml.Marshal(m)
		file2.AppendFile(output, result)
		mylog.CommonLogger.Info().Msgf("parser success, write to %s", output)
		return nil
	},
}

func init() {
	ocr2Cmd.Flags().StringVarP(&ocrContent, "content", "c", "", "ocr content")
	ocr2Cmd.Flags().StringVarP(&output, "output", "o", "out/driver.yml", "output path")

	ocr2Cmd.MarkFlagRequired("content")
	RootCmd.AddCommand(ocr2Cmd)
}
