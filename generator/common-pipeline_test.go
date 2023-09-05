package generator

import (
	"strings"
	"testing"
)

type testCase struct {
	name         string
	tpc          TrivyPluginConfig
	templatePath string
	hasError     bool
	errMsg       string
	expVal       string
}

func TestGenerateTrivyStep(t *testing.T) {
	tpc := TrivyPluginConfig{
		Severity:       "CRITICAL,HIGH",
		SecurityChecks: "config,secret,vuln",
		IgnoreUnfixed:  true,
		SkipFiles:      "cosign.key",
	}
	expected := `
steps:
  - command: ls
    plugins:
      - equinixmetal-buildkite/trivy#:
          severity: "CRITICAL,HIGH"
          ignore-unfixed: true
          security-checks: "config,secret,vuln"
          skip-files: "cosign.key"
`
	cases := []testCase{
		{"success", tpc, "../templates/*", false, "", expected},
		{"wrong_template_path", tpc, "templates/*", true, "template: pattern matches no files", ""},
	}

	g := Generator{
		TrivyPluginEnabled: true,
		TPConfig:           tpc,
	}

	for _, tc := range cases {
		t.Run(
			tc.name,
			func(t *testing.T) {
				var sb strings.Builder
				err := GenerateTrivyStep(g, &sb, tc.templatePath)
				if tc.hasError {
					if !strings.Contains(err.Error(), tc.errMsg) {
						t.Fatalf("Error %s does not contain %s\n", err.Error(), tc.errMsg)
					}
				} else {
					if strings.TrimSpace(tc.expVal) != strings.TrimSpace(sb.String()) {
						t.Fatalf("Not equal: \n"+
							"expected: %s\n"+
							"actual  : %s\n", strings.TrimSpace(tc.expVal), strings.TrimSpace(sb.String()))
					}
				}
			},
		)
	}
}
