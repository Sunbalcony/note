FROM centos:centos7.9.2009

WORKDIR /root
RUN mkdir static
RUN mkdir conf

COPY static /root/static
#COPY config /root/conf
COPY note /root
VOLUME /root/conf

EXPOSE 8080

CMD ["/root/note"]





