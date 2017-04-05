/******************************************************************************
 *
 * Walrus Test Framework
 * Copyright Ryan Goodfellow 2017, all rights reserved
 * GPLv3+
 *
 *****************************************************************************/
#include <stdlib.h>
#include <hiredis/hiredis.h>
#include "walrus.h"


int __WTFxxx(struct WTFTest *t, const char *level, const char *fmt, ...)
{
  redisContext *c = redisConnect(t->collector, REDIS_PORT);
  if(c == NULL || c->err)
  {
    printf("redis error: %s\n", c->errstr);
    return FAILURE;
  }

  va_list vl1, vl2;
  va_start(vl1, fmt);
  va_copy(vl2, vl1);
  size_t sz = vsnprintf(NULL, 0, fmt, vl1);
  va_end(vl1);
  char *buf = malloc(sz+1);
  vsnprintf(buf, sz+1, fmt, vl2);
  va_end(vl2);
  buf[sz] = 0;

  //this is the only mutation we make to t, lets me thread safe about ot
  //TODO this is a moot point right now b/c hiredis does not appear to be
  //thread safe, but i'll look into this soon.
  int counter = __sync_fetch_and_add(&t->counter, 1);

  redisReply *r = 
  (redisReply*)redisCommand(c, "SET %s:%s:%d %s:::%s", 
      t->test, 
      t->participant,
      counter,
      level,
      buf
  );
  if(c->err)
  {
    printf("redis error: %s\n", c->errstr);
    exit(FAILURE);
  }
  if(r->type == REDIS_REPLY_ERROR)
  {
    printf("redis reply error: %s\n", r->str);
    exit(FAILURE);
  }

  r = (redisReply*)redisCommand(c, "TIME");
  if(c->err)
  {
    printf("redis error: %s\n", c->errstr);
    exit(FAILURE);
  }
  if(r->type != REDIS_REPLY_ARRAY)
  {
    printf("redis error: TIME did not return array\n");
    exit(FAILURE);
  }
  redisReply *seconds = r->element[0];
  redisReply *microseconds= r->element[1];
  redisCommand(c, "DEL %s:%s:%d:~time~", 
      t->test, 
      t->participant,
      counter
  );
  
  r = (redisReply*)redisCommand(c, "RPUSH %s:%s:%d:~time~ %s", 
      t->test, 
      t->participant,
      counter,
      seconds->str
  );
  if(c->err)
  {
    printf("redis error (timestamp push): %s\n", c->errstr);
    exit(FAILURE);
  }
  if(r->type == REDIS_REPLY_ERROR)
  {
    printf("redis reply error (timestamp push): %s\n", r->str);
    exit(FAILURE);
  }
  r = (redisReply*)redisCommand(c, "RPUSH %s:%s:%d:~time~ %s", 
      t->test, 
      t->participant,
      counter,
      microseconds->str
  );
  
  free(buf);
  redisFree(c);

  return SUCCESS;
}

