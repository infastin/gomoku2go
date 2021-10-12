project('com.github.infastin.gomoku2go')

golang = find_program('go')
gomoku2go_build_path = join_paths(meson.current_source_dir(), 'cmd/')

gomoku2go = custom_target(
  'gomoku2go',
  output: 'gomoku2go',
  command: [ golang, 'build', '-v', '-o', '@OUTPUT@', gomoku2go_build_path ],
  install: true,
  install_dir: 'bin',
)

subdir('data')
