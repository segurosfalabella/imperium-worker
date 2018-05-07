FROM scratch

ADD /imperium-worker /imperium-worker

ENTRYPOINT ["/imperium-worker"]
