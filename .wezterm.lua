local wezterm = require 'wezterm'
local config = {}

if wezterm.config_builder then
  config = wezterm.config_builder()
end

-- Set character height to half
config.cell_width = 0.1
config.line_height = 0.025
config.initial_cols = 1000
config.initial_rows = 1000

return config
