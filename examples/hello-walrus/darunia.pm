#!/usr/bin/env perl

use IO::Socket::INET;
use lib '/opt/walrus/perl';
use Walrus;

my $wtf = Walrus->new('192.168.147.100', 'hyrule', 'darunia');
$wtf->warning("grrrrrr");


my $sock = IO::Socket::INET->new(
  LocalAddr => '192.168.147.4',
  LocalPort => 4747,
  Proto => 'udp') or die "socket: $@";

print("awaiting msg\n");

my $msg = "";
my $sender = $sock->recv($msg, 64);

if ($msg eq "do you know the mUfFiN MaN?")
{
  print("expected message recv'd\n");
  $wtf->ok("expected message recv'd");
  $sock->send("the muffin man is ME!");
}
else
{
  print("strange message recv'd\n");
  $wtf->error("strange message recv'd");
  $sock->send("errr.... wat!?");
}


