package root

import (
	"github.com/henrywhitaker3/boiler"
	"github.com/henrywhitaker3/shorturl/cmd/consume"
	"github.com/henrywhitaker3/shorturl/cmd/migrate"
	"github.com/henrywhitaker3/shorturl/cmd/routes"
	"github.com/henrywhitaker3/shorturl/cmd/secrets"
	"github.com/henrywhitaker3/shorturl/cmd/seed"
	"github.com/henrywhitaker3/shorturl/cmd/serve"
	"github.com/spf13/cobra"
)

func New(b *boiler.Boiler) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "api",
		Short:   "Golang template API",
		Version: string(b.Version()),
		GroupID: "app",
	}

	cmd.AddGroup(&cobra.Group{
		ID:    "app",
		Title: "Application commands",
	})
	cmd.AddGroup(&cobra.Group{
		ID:    "unconf",
		Title: "Configuration commands",
	})

	cmd.AddCommand(serve.New(b))
	cmd.AddCommand(migrate.New(b))
	cmd.AddCommand(routes.New(b))
	cmd.AddCommand(consume.New(b))
	cmd.AddCommand(seed.New(b))
	cmd.AddCommand(secrets.New())

	cmd.PersistentFlags().
		StringP("config", "c", "shorturl.yaml", "The path to the api config file")

	return cmd
}
