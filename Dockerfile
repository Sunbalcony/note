FROM centos:centos7.9.2009

WORKDIR /root
RUN mkdir static
RUN mkdir config

COPY static /root/static/
COPY config /root/config/
COPY note /root/

EXPOSE 8080

CMD ["/root/note"]





