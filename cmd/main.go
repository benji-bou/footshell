package main

import (
	"log"
	"os"

	"github.com/benji-bou/footshell/loader"
	"github.com/c-bata/go-prompt"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/urfave/cli/v2"
)

var (
	prefixAttackSymbolName   = "Attack"
	prefixModifierSymbolName = "Modifier"
	prefixDisplaySymbolName  = "Display"
)

type AttackVectorPlugin interface {
	Attack(payload string) (string, error)
}

type OutputDisplayPlugin interface {
	Display(output interface{}) error
}

type ModifierPlugin interface {
	Transform(input interface{}) (interface{}, error)
}

func main() {
	app := &cli.App{
		Name:  "FootShell",
		Usage: "footshell: simple modular shell for foothold",
		Authors: []*cli.Author{
			{
				Name: "Benjamin Bouachour",
			},
		},
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "attacker",
				Aliases:  []string{"a"},
				Usage:    "Path for attack plugin to load",
				Category: "Attack Plugin load",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     "modifier",
				Aliases:  []string{"m"},
				Usage:    "Path for modifier plugin to load",
				Category: "Modifier Plugin load",
				Required: false,
			},

			&cli.PathFlag{
				Name:     "display",
				Aliases:  []string{"d"},
				Usage:    "Path for display plugin to load",
				Category: "Attack Plugin load",
				Required: false,
			},
		},
		Before: func(cCtx *cli.Context) error {
			pterm.DefaultBigText.WithLetters(
				putils.LettersFromStringWithStyle("FootShell ", pterm.FgMagenta.ToStyle()))
			return nil
		},
		Commands: []*cli.Command{},
		Action: func(ctx *cli.Context) error {
			attack, display, modifiersPlugin, err := LoadAllPlugins(ctx)
			if err != nil {
				pterm.Error.Println(err)
				return err
			}
			return RunShell(attack, display, modifiersPlugin...)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func LoadAllPlugins(ctx *cli.Context) (AttackVectorPlugin, OutputDisplayPlugin, []ModifierPlugin, error) {
	spinner, _ := pterm.DefaultSpinner.Start("🔪 Loading Attack vector plugin...")
	defer func() {
		if spinner.IsActive {
			spinner.Stop()
		}
	}()
	attack, symbolNameLoaded, err := loader.Symbol[AttackVectorPlugin](ctx.Path("attacker"), prefixAttackSymbolName)
	if err != nil {
		return nil, nil, nil, err
	}
	pterm.Success.Printfln("Attack vector %s is loaded", symbolNameLoaded)

	var display OutputDisplayPlugin
	if ctx.IsSet("display") {
		spinner.UpdateText("🖥️ Loading Display plugin...")
		display, symbolNameLoaded, err = loader.Symbol[OutputDisplayPlugin](ctx.Path("display"), prefixDisplaySymbolName)
		if err != nil {
			return nil, nil, nil, err
		}
		pterm.Success.Printfln("Display %s is loaded", symbolNameLoaded)
	}
	spinner.UpdateText("💫 Loading Modifiers plugin...")
	modifiersPath := ctx.StringSlice("modifier")
	modifiers := make([]ModifierPlugin, 0, len(modifiersPath))
	for _, modifierPath := range modifiersPath {
		modifierSymbol, symbolNameLoaded, err := loader.Symbol[ModifierPlugin](modifierPath, prefixModifierSymbolName)
		if err != nil {
			return nil, nil, nil, err
		}
		pterm.Success.Printfln("Modifiers %s is loaded", symbolNameLoaded)
		modifiers = append(modifiers, modifierSymbol)
	}
	spinner.Success("🚀 Ready to attack!")
	return attack, display, modifiers, nil
}

func DisplayOrDefault(display OutputDisplayPlugin, output interface{}) {
	if display != nil {
		err := display.Display(output)
		if err != nil {
			pterm.Error.Println(err)
		}
	} else {
		pterm.DefaultBasicText.Println(pterm.Green(output))
	}
}

func RunShell(attacker AttackVectorPlugin, display OutputDisplayPlugin, modifier ...ModifierPlugin) error {

	p := prompt.New(func(input string) {
		var res interface{}
		rawRes, err := attacker.Attack(input)
		if err != nil {
			pterm.Error.Println(err)
			return
		}
		res = rawRes
		for _, m := range modifier {
			res, err = m.Transform(res)
			if err != nil {
				break
			}
		}
		if err != nil {
			pterm.Error.Println(err)
			DisplayOrDefault(display, rawRes)
		} else {
			DisplayOrDefault(display, res)
		}
	}, func(d prompt.Document) []prompt.Suggest {
		return []prompt.Suggest{}
	},
		prompt.OptionPrefix("$ > "),
		prompt.OptionPrefixTextColor(prompt.DarkGreen),
	)
	p.Run()
	return nil
}
