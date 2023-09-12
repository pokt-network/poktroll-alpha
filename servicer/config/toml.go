package config

import serverconfig "github.com/cosmos/cosmos-sdk/server/config"

const ConfigTemplate = serverconfig.DefaultConfigTemplate + `

###############################################################################
###                         Servicer Configuration                          ###
###############################################################################

blocks-per-session = "{{ Servicer.BlocksPerSession }}"
`
