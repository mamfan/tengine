local wezterm = require 'wezterm'
local config = {}
local gpus = wezterm.gui.enumerate_gpus()

if wezterm.config_builder then
  config = wezterm.config_builder()
end

config.webgpu_preferred_adapter = gpus[1]
config.front_end = 'WebGpu'

-- Set character height to half
local scale = 8
config.cell_width = 0.1 * scale / 2
config.line_height = 0.025 * scale * 0.9
config.initial_cols = 360
config.initial_rows = 200

max_fps = 24

config.use_fancy_tab_bar = false
config.show_tabs_in_tab_bar = false
config.show_new_tab_button_in_tab_bar = false


config.skip_close_confirmation_for_processes_named = {
  'cmd.exe',
  'pwsh.exe',
  'powershell.exe',
  'go.exe',
  'core.exe',
  'asm.exe',
}


return config
