import redis

class Test:
    def __init__(s, collector, test, participant):
        s.collector = collector
        s.test = test
        s.participant = participant
        s.counter = 0

    def error(s, msg):
        s.__report('error', msg)

    def warning(s, msg):
        s.__report('warning', msg)

    def ok(s, msg):
        s.__report('ok', msg)

    def __report(s, level, msg):
        r = redis.StrictRedis(host=s.collector)
        key = '%s:%s:%d'%(
            s.test,
            s.participant,
            s.counter,
            )
        value = '%s:::%s'%(
            level,
            msg)

        r.set(key, value)
        t = r.time()
        r.delete("%s:~time~"%key)
        r.rpush("%s:~time~"%key, t[0])
        r.rpush("%s:~time~"%key, t[1])

        s.counter += 1
