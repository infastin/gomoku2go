schemasDir = join_paths(get_option('datadir'), 'glib-2.0', 'schemas')

gnome = import('gnome')

install_data(
  meson.project_name() + '.gschema.xml',
  install_dir: schemasDir
)

gnome.post_install(glib_compile_schemas: true)
