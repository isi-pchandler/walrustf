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
        t = r.time()
        key = '%s:%s:%s:%s'%(
            s.test,
            s.participant,
            t[0],
            t[1])
        value = '%s:::%s'%(
            level,
            msg)

        r.set(key, value)

