package menu

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		td          *TemplateData
		text        string
		output      string
		expectedErr bool
	}{
		{
			name:   "default",
			td:     &TemplateData{},
			text:   "Version: {{.LauncherVersion}}",
			output: "Version: ",
		},
		{
			name:   "single option",
			td:     &TemplateData{ServerHostname: "localhost"},
			text:   "Hostname: {{.ServerHostname}}",
			output: "Hostname: localhost",
		},
		{
			name:   "multiple option",
			td:     &TemplateData{ServerHostname: "localhost", LauncherVersion: "0.0.0"},
			text:   "Hostname: {{.ServerHostname}}, launcher version: {{.LauncherVersion}}",
			output: "Hostname: localhost, launcher version: 0.0.0",
		},
		{
			name:   "unsupported capability",
			td:     &TemplateData{ServerHostname: "localhost", LauncherVersion: "0.0.0"},
			text:   "This capability is {{if hasCapability `bad capability`}}supported{{else}}unsupported{{end}}.",
			output: "This capability is unsupported.",
		},
		{
			name:   "relativeTime 2 hours ago",
			td:     &TemplateData{LastUpdateTime: time.Now().Add(-2 * time.Hour).Unix()},
			text:   fmt.Sprintf("This Menu Was Last Updated {{if hasCapability `relativeTime`}}{{relativeTime .LastUpdateTime}}{{else}}never{{end}}."),
			output: "This Menu Was Last Updated 2 Hours Ago.",
		},
		{
			name:   "relativeTime 15 minutes ago",
			td:     &TemplateData{LastUpdateTime: time.Now().Add(-15*time.Minute - 30*time.Second).Unix()},
			text:   fmt.Sprintf("This Menu Was Last Updated {{if hasCapability `relativeTime`}}{{relativeTime .LastUpdateTime}}{{else}}never{{end}}."),
			output: "This Menu Was Last Updated 15 Minutes Ago.",
		},
		{
			name:   "relativeTime 111 seconds ago",
			td:     &TemplateData{LastUpdateTime: time.Now().Add(-111 * time.Second).Unix()},
			text:   fmt.Sprintf("This Menu Was Last Updated {{if hasCapability `relativeTime`}}{{relativeTime .LastUpdateTime}}{{else}}never{{end}}."),
			output: "This Menu Was Last Updated 111 Seconds Ago.",
		},
		{
			name:   "relativeTime one minute ago",
			td:     &TemplateData{LastUpdateTime: time.Now().Add(-1 * time.Minute).Unix()},
			text:   fmt.Sprintf("This Menu Was Last Updated {{if hasCapability `relativeTime`}}{{relativeTime .LastUpdateTime}}{{else}}never{{end}}."),
			output: "This Menu Was Last Updated One Minute Ago.",
		},
		{
			name:   "relativeTime 7 seconds ago",
			td:     &TemplateData{LastUpdateTime: time.Now().Add(-7 * time.Second).Unix()},
			text:   fmt.Sprintf("This Menu Was Last Updated {{if hasCapability `relativeTime`}}{{relativeTime .LastUpdateTime}}{{else}}never{{end}}."),
			output: "This Menu Was Last Updated 7 Seconds Ago.",
		},
		{
			name:   "relativeTime just now",
			td:     &TemplateData{LastUpdateTime: time.Now().Add(-1 * time.Second).Unix()},
			text:   fmt.Sprintf("This Menu Was Last Updated {{if hasCapability `relativeTime`}}{{relativeTime .LastUpdateTime}}{{else}}never{{end}}."),
			output: "This Menu Was Last Updated Just Now.",
		},
		{
			name:   "relativeTime very soon",
			td:     &TemplateData{ServerHostname: "localhost", LauncherVersion: "0.0.0"},
			text:   fmt.Sprintf("This Is Starting {{if hasCapability `relativeTime`}}{{relativeTime %d}}{{else}}never{{end}}.", time.Now().Add(2*time.Minute+30*time.Second).Unix()),
			output: "This Is Starting Very Soon.",
		},
		{
			name:   "relativeTime 15 minutes",
			td:     &TemplateData{ServerHostname: "localhost", LauncherVersion: "0.0.0"},
			text:   fmt.Sprintf("This Is Starting {{if hasCapability `relativeTime`}}{{relativeTime %d}}{{else}}never{{end}}.", time.Now().Add(15*time.Minute+30*time.Second).Unix()),
			output: "This Is Starting In 15 Minutes.",
		},
		{
			name:   "relativeTime one hour",
			td:     &TemplateData{ServerHostname: "localhost", LauncherVersion: "0.0.0"},
			text:   fmt.Sprintf("This Is Starting {{if hasCapability `relativeTime`}}{{relativeTime %d}}{{else}}never{{end}}.", time.Now().Add(1*time.Hour+30*time.Second).Unix()),
			output: "This Is Starting In About An Hour.",
		},
		{
			name:   "relativeTime 111 minutes",
			td:     &TemplateData{ServerHostname: "localhost", LauncherVersion: "0.0.0"},
			text:   fmt.Sprintf("This Is Starting {{if hasCapability `relativeTime`}}{{relativeTime %d}}{{else}}never{{end}}.", time.Now().Add(111*time.Minute+30*time.Second).Unix()),
			output: "This Is Starting In 111 Minutes.",
		},
		{
			name:   "relativeTime four hours",
			td:     &TemplateData{ServerHostname: "localhost", LauncherVersion: "0.0.0"},
			text:   fmt.Sprintf("This Is Starting {{if hasCapability `relativeTime`}}{{relativeTime %d}}{{else}}never{{end}}.", time.Now().Add(4*time.Hour+30*time.Second).Unix()),
			output: "This Is Starting In 4 Hours.",
		},
		{
			name:   "relativeTime one day",
			td:     &TemplateData{ServerHostname: "localhost", LauncherVersion: "0.0.0"},
			text:   fmt.Sprintf("This Is Starting {{if hasCapability `relativeTime`}}{{relativeTime %d}}{{else}}never{{end}}.", time.Now().Add(24*time.Hour+30*time.Second).Unix()),
			output: "This Is Starting In One Day.",
		},
		{
			name:   "relativeTime three days",
			td:     &TemplateData{ServerHostname: "localhost", LauncherVersion: "0.0.0"},
			text:   fmt.Sprintf("This Is Starting {{if hasCapability `relativeTime`}}{{relativeTime %d}}{{else}}never{{end}}.", time.Now().Add(3*24*time.Hour+30*time.Second).Unix()),
			output: "This Is Starting In 3 Days.",
		},
		{
			name:   "relativeTime seven days",
			td:     &TemplateData{ServerHostname: "localhost", LauncherVersion: "0.0.0"},
			text:   fmt.Sprintf("This Is Starting {{if hasCapability `relativeTime`}}{{relativeTime %d}}{{else}}never{{end}}.", time.Now().Add(7*24*time.Hour+30*time.Second).Unix()),
			output: "This Is Starting In 7 Days.",
		},
		{
			name:   "relativeTime two weeks",
			td:     &TemplateData{ServerHostname: "localhost", LauncherVersion: "0.0.0"},
			text:   fmt.Sprintf("This Is Starting {{if hasCapability `relativeTime`}}{{relativeTime %d}}{{else}}never{{end}}.", time.Now().Add(14*24*time.Hour+30*time.Second).Unix()),
			output: "This Is Starting In 2 Weeks.",
		},
		{
			name:        "undefined",
			td:          &TemplateData{},
			text:        "Version: {{GARBAGE}}",
			output:      "",
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tp := NewTemplateParser(tt.td)
			o, err := tp.Parse(tt.text)
			if !tt.expectedErr {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.output, o)
		})
	}
}
