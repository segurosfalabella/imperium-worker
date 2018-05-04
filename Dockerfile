FROM scratch

ADD /imperium-worker /imperium-worker

CMD ["/imperium-worker"]
