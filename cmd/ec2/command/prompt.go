package command

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	prompt "github.com/c-bata/go-prompt"
)

func selectInstancePrompt(instances []*ec2.Instance) string {
	var suggests []prompt.Suggest
	for _, i := range instances {
		var name string
		for _, t := range i.Tags {
			if *t.Key == "Name" {
				name = *t.Value
			}
		}

		s := prompt.Suggest{Text: name,
			Description: *i.State.Name,
		}
		suggests = append(suggests, s)
	}

	return prompt.Input(
		"Choose a ec2 instance >>> ",
		func(in prompt.Document) []prompt.Suggest {
			return prompt.FilterContains(suggests, in.GetWordBeforeCursor(), true)
		},
		prompt.OptionTitle("ec2"),
		prompt.OptionPrefixTextColor(prompt.Green),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.Blue),
	)
}
