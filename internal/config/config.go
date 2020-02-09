/*
Package config has global variable for configuration.
*/
package config

const (
	// Logo is the logo of Argocui.
	Logo = `
    _____                                 _________  ____ ___.___ 
   /  _  \_______  ____   ____            \_   ___ \|    |   \   |
  /  /_\  \_  __ \/ ___\ /  _ \   ____    /    \  \/|    |   /   | 
 /    |    \  | \/ /_/  >  <_> ) _______  \     \___|    |  /|   |
 \____|__  /__|  \___  / \____/   ___      \______  /______/ |___|
         \/     /_____/                           \/              
`
	// Version is the version of Argocui.
	Version = "v0.0.2"
	// ArgoVersion is the version of the Argo package.
	ArgoVersion = "v2.4.1"
	// HomePage is the url of Argocui repository.
	HomePage = "github.com/hanjunlee/argocui"
)

var (
	// ReadOnly is the read-only mode.
	ReadOnly = false
)
