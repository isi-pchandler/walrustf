package Walrus;

use Redis;

use warnings;
use strict;

sub new {
  my $class = shift;

  my $self = {
    collector => shift,
    test => shift,
    participant => shift
  };
  bless($self, $class);
  return($self);
}

sub report {
  my ($self, $level, $msg) = @_;

  # connect to collector
  my $redis = Redis->new(server => "$self->{collector}:6379");

  # get the time according to collector
  my @time = $redis->time();

  # create a key/value pair for this diagnostic
  my $key = "$self->{test}:$self->{participant}:$time[0]:$time[1]";
  my $value = "${level}:::$msg";

  # send and close connection
  $redis->set($key => $value);
  $redis->quit();
}

sub error {
  my ($self, $msg) = @_;
  $self->report('error', $msg);
}

sub warning {
  my ($self, $msg) = @_;
  $self->report('warning', $msg);
}

sub ok {
  my ($self, $msg) = @_;
  $self->report('ok', $msg);
}

1;
