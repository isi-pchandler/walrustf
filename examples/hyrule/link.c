#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <wtf/walrus.h>
#include <sys/time.h>
#include <errno.h>
#include <stdlib.h>

#define N 64

int init();
char* ask(
    struct WTFTest *t, 
    int sock,
    const char *who, 
    const char *addr, 
    const char *question,
    ssize_t *len);

int main()
{
  int sock = init();

  struct WTFTest t = {
    .collector = "192.168.147.100", 
    .test = "hyrule", 
    .participant = "link",
    .counter = 0
  };

  /// prepare my question
  const char *question = "do you know the muffin man?";

  /// ask zelda
  const char *expected = "why yes i know the muffin man";
  ssize_t len = 0;
  char *response = ask(&t, sock, "zelda", "192.168.147.3", question, &len);
  if(response == NULL || len == 0)
  {

  }
  else if(strncmp(expected, response, len) != 0)
  {
    WTFerror(&t, "unexpected zelda response `%s`", response);
  }
  else
  {
    WTFok(&t, "zelda response is good");
  }
  free(response);

  /// ask darunia
  expected = "the muffin man is ME!";
  len = 0;
  response = ask(&t, sock, "darunia", "192.168.147.4", question, &len);
  if(response == NULL || len == 0)
  {

  }
  else if(strncmp(expected, response, len) != 0)
  {
    WTFerror(&t, "unexpected darunia response `%s`", response);
  }
  else
  {
    WTFok(&t, "darunia response is good");
  }
  free(response);
  
  return EXIT_SUCCESS;
}

int init()
{
  /// prepare to listen for answers
  int sock = socket(AF_INET, SOCK_DGRAM, 0);
  struct timeval tv;

  // only wait 3 seconds for replies
  tv.tv_sec = 3;
  tv.tv_usec = 0;
  setsockopt(sock, SOL_SOCKET, SO_RCVTIMEO, &tv, sizeof(tv));

  // listen on port 4747
  struct sockaddr_in listen = {
    .sin_family = AF_INET,
    .sin_port = htons(4747),
    .sin_addr = INADDR_ANY
  };
  bind(sock, (const struct sockaddr*)&listen, sizeof(listen));

  return sock;
}


char* ask(
    struct WTFTest *t,
    int sock,
    const char *who, 
    const char *addr, 
    const char *question, 
    ssize_t *len)
{

  // fill in destination struct
  struct sockaddr_in saddr = {
    .sin_family = AF_INET,
    .sin_port = htons(4747),
    .sin_addr = inet_addr(addr)
  };
  socklen_t sz = sizeof(saddr);

  // allocate memory for response
  char *answer = malloc(N);
  memset(answer, 0, N);
  
  errno = 0;

  // send & receive
  sendto(sock, question, strlen(question), 0, (const struct sockaddr*)&saddr, sz);
  *len = recvfrom(sock, answer, N, 0, (struct sockaddr*)&saddr, &sz);

  // check for errors
  if(errno == EAGAIN || errno == EWOULDBLOCK)
  {
    WTFerror(t, "%s did not respond", who);
    free(answer);
    *len = 0;
    return NULL;
  }
  else if(errno)
  {
    WTFerror(t, "error talking to %s code=%d", who, errno);
    free(answer);
    *len = 0;
    return NULL;
  }

  // OK
  return answer;
}
