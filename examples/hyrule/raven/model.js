
zelda = {
  'name': 'zelda',
  'image': 'debian-stretch',
  'os': 'linux',
  'level': 1,
  'mounts': [{'source': '/space/walrustf', 'point': '/opt/walrus'}]
}

darunia = {
  'name': 'darunia',
  'image': 'debian-stretch',
  'os': 'linux',
  'level': 1,
  'mounts': [{'source': '/space/walrustf', 'point': '/opt/walrus'}]
}

impa = {
  'name': 'impa',
  'image': 'debian-stretch',
  'os': 'linux',
  'level': 3,
  'mounts': [{'source': '/space/walrustf', 'point': '/opt/walrus'}]
}

link = {
  'name': 'link',
  'image': 'debian-stretch',
  'os': 'linux',
  'level': 3,
  'mounts': [{'source': '/space/walrustf', 'point': '/opt/walrus'}]
}

walrus = {
  'name': 'walrus',
  'image': 'debian-stretch',
  'os': 'linux',
  'level': 2,
  'mounts': [{'source': '/space/walrustf', 'point': '/opt/walrus'}]
}

nimbus = {
  'name': 'nimbus',
  'image': 'cumulus-latest',
  'os': 'linux',
  'level': 2,
}

links = [
  Link('zelda', 'eth0', 'nimbus', 'swp1'),
  Link('darunia', 'eth0', 'nimbus', 'swp2'),
  Link('impa', 'eth0', 'nimbus', 'swp3'),
  Link('link', 'eth0', 'nimbus', 'swp4'),
  Link('walrus', 'eth0', 'nimbus', 'swp5')
]

topo = {
  'name': 'hyrule',
  'nodes': [zelda, darunia, impa, link, walrus],
  'switches': [nimbus],
  'links': links
}

