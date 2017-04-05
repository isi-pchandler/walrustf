package Walrus;

use Redis;

use warnings;
use strict;

sub new {
  my $class = shift;

  my $self = {
    collector => shift,
    test => shift,
    participant => shift,
    counter => 0,
  };
  bless($self, $class);
  return($self);
}

sub report {
  my ($self, $level, $msg) = @_;
  my $redis = Redis->new(server => "$self->{collector}:6379");
  my $key = "$self->{test}:$self->{participant}:$self->{counter}";
  my $value = "${level}:::$msg";
  $redis->set($key => $value);
  my @time = $redis->time();
  $redis->del("$key:~time~");
  $redis->rpush("$key:~time~" => $time[0]);
  $redis->rpush("$key:~time~" => $time[1]);
  $redis->quit();
  $self->{counter}++;
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
